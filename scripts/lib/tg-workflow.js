const fs = require('fs');
const path = require('path');
const { spawnSync } = require('child_process');

const META_REQUIRED_FIELDS = [
  'schema_version',
  'id',
  'slug',
  'title',
  'status',
  'created_at',
  'updated_at',
  'source_exploration',
  'owner',
  'task_backend',
  'archive_targets',
];

const STATUS_VALUES = new Set([
  'draft',
  'proposed',
  'approved',
  'active',
  'implemented',
  'archived',
  'rejected',
]);

const TASK_BACKENDS = new Set(['beads', 'github', 'linear', 'none']);
const ARCHIVE_TARGET_TYPES = new Set(['module', 'architecture', 'decision']);
const ARCHIVE_TARGET_ROLES = new Set(['primary', 'secondary']);
const TASK_PRIORITIES = new Set(['P0', 'P1', 'P2']);
const PROPOSAL_ID_RE = /^CR\d{8}$/;
const DEFAULT_CONFIG = {
  project: {
    id: 'unknown-project',
    name: 'Unknown Project',
    slug: 'unknown-project',
  },
  paths: {
    docs_root: 'docs',
    state_root: '.workflow/state',
  },
  task_backend: {
    type: 'beads',
  },
  memory_provider: {
    enabled: false,
  },
};

const STATUS_TRANSITIONS = {
  draft: new Set(['draft', 'proposed', 'rejected']),
  proposed: new Set(['proposed', 'approved', 'rejected', 'draft']),
  approved: new Set(['approved', 'active', 'rejected']),
  active: new Set(['active', 'implemented', 'rejected']),
  implemented: new Set(['implemented', 'active', 'archived']),
  archived: new Set(['archived']),
  rejected: new Set(['rejected', 'draft']),
};

function parseScalar(raw) {
  const value = raw.trim();
  if (value === 'null') return null;
  if (value === 'true') return true;
  if (value === 'false') return false;
  if (/^-?\d+$/.test(value)) return Number(value);
  if ((value.startsWith('"') && value.endsWith('"')) || (value.startsWith("'") && value.endsWith("'"))) {
    return value.slice(1, -1);
  }
  return value;
}

function nextMeaningfulLine(lines, startIndex) {
  for (let index = startIndex; index < lines.length; index += 1) {
    const trimmed = lines[index].trim();
    if (!trimmed || trimmed.startsWith('#')) continue;
    return index;
  }
  return lines.length;
}

function lineIndent(line) {
  return line.length - line.trimStart().length;
}

function parseYamlBlock(lines, startIndex, indent) {
  const firstIndex = nextMeaningfulLine(lines, startIndex);
  if (firstIndex >= lines.length) return { value: null, nextIndex: firstIndex };

  const firstLine = lines[firstIndex];
  if (lineIndent(firstLine) < indent) {
    return { value: null, nextIndex: firstIndex };
  }

  if (firstLine.trimStart().startsWith('- ')) {
    return parseYamlArray(lines, firstIndex, indent);
  }

  return parseYamlObject(lines, firstIndex, indent);
}

function parseYamlArray(lines, startIndex, indent) {
  const result = [];
  let index = startIndex;

  while (index < lines.length) {
    index = nextMeaningfulLine(lines, index);
    if (index >= lines.length) break;

    const line = lines[index];
    const currentIndent = lineIndent(line);
    if (currentIndent < indent) break;
    if (currentIndent > indent) {
      throw new Error(`Unexpected indentation at line ${index + 1}`);
    }

    const trimmed = line.trimStart();
    if (!trimmed.startsWith('- ')) break;

    const itemText = trimmed.slice(2).trim();
    if (!itemText) {
      const nested = parseYamlBlock(lines, index + 1, indent + 2);
      result.push(nested.value);
      index = nested.nextIndex;
      continue;
    }

    const inlineObjectMatch = itemText.match(/^([^:]+):\s*(.*)$/);
    if (inlineObjectMatch) {
      const item = {};
      item[inlineObjectMatch[1].trim()] = parseScalar(inlineObjectMatch[2]);

      const nestedStart = nextMeaningfulLine(lines, index + 1);
      if (nestedStart < lines.length && lineIndent(lines[nestedStart]) > indent) {
        const nested = parseYamlObject(lines, nestedStart, indent + 2);
        Object.assign(item, nested.value || {});
        index = nested.nextIndex;
      } else {
        index += 1;
      }

      result.push(item);
      continue;
    }

    result.push(parseScalar(itemText));
    index += 1;
  }

  return { value: result, nextIndex: index };
}

function parseYamlObject(lines, startIndex, indent) {
  const result = {};
  let index = startIndex;

  while (index < lines.length) {
    index = nextMeaningfulLine(lines, index);
    if (index >= lines.length) break;

    const line = lines[index];
    const currentIndent = lineIndent(line);
    if (currentIndent < indent) break;
    if (currentIndent > indent) {
      throw new Error(`Unexpected indentation at line ${index + 1}`);
    }

    const trimmed = line.trim();
    const match = trimmed.match(/^([^:]+):\s*(.*)$/);
    if (!match) {
      throw new Error(`Invalid mapping at line ${index + 1}`);
    }

    const key = match[1].trim();
    const rawValue = match[2];
    if (rawValue) {
      result[key] = parseScalar(rawValue);
      index += 1;
      continue;
    }

    const nested = parseYamlBlock(lines, index + 1, indent + 2);
    result[key] = nested.value;
    index = nested.nextIndex;
  }

  return { value: result, nextIndex: index };
}

function parseSimpleYaml(text) {
  const lines = text.replace(/\r\n/g, '\n').split('\n');
  return parseYamlBlock(lines, 0, 0).value || {};
}

function parseTaskMap(text) {
  const lines = text.replace(/\r\n/g, '\n').split('\n');
  const result = {
    proposal_id: null,
    backend: null,
    tasks: [],
  };

  let currentTask = null;
  let collectingCompletion = false;

  for (const line of lines) {
    const trimmed = line.trim();
    if (!trimmed) continue;

    const backendMatch = trimmed.match(/^- Backend:\s*(.+)$/);
    if (backendMatch) {
      result.backend = backendMatch[1].trim();
      continue;
    }

    const proposalMatch = trimmed.match(/^- Proposal ID:\s*(.+)$/);
    if (proposalMatch) {
      result.proposal_id = proposalMatch[1].trim();
      continue;
    }

    const taskMatch = trimmed.match(/^###\s+(.+)$/);
    if (taskMatch) {
      currentTask = {
        task_id: taskMatch[1].trim(),
        depends_on: [],
        capability_refs: [],
        decision_refs: [],
        archive_target_refs: [],
        completion_definition: [],
      };
      result.tasks.push(currentTask);
      collectingCompletion = false;
      continue;
    }

    if (!currentTask) continue;

    const fieldMatch = trimmed.match(/^- ([^:]+):\s*(.*)$/);
    if (fieldMatch) {
      const field = fieldMatch[1].trim().toLowerCase();
      const value = fieldMatch[2].trim();
      collectingCompletion = field === 'completion definition';

      if (field === 'title') currentTask.title = value;
      else if (field === 'outcome') currentTask.outcome = value;
      else if (field === 'priority') currentTask.priority = value;
      else if (field === 'depends on' && value) currentTask.depends_on = value.split(',').map((item) => item.trim()).filter(Boolean);
      else if (field === 'capability refs' && value) currentTask.capability_refs = value.split(',').map((item) => item.trim()).filter(Boolean);
      else if (field === 'decision refs' && value) currentTask.decision_refs = value.split(',').map((item) => item.trim()).filter(Boolean);
      else if (field === 'archive target refs' && value) currentTask.archive_target_refs = value.split(',').map((item) => item.trim()).filter(Boolean);
      continue;
    }

    if (collectingCompletion) {
      const completionMatch = trimmed.match(/^- (.+)$/);
      if (completionMatch) {
        currentTask.completion_definition.push(completionMatch[1].trim());
      }
    }
  }

  return result;
}

function parseCliArgs(argv) {
  const args = { _: [] };

  for (let index = 0; index < argv.length; index += 1) {
    const token = argv[index];
    if (!token.startsWith('--')) {
      args._.push(token);
      continue;
    }

    const key = token.slice(2);
    const next = argv[index + 1];
    const value = !next || next.startsWith('--') ? true : next;
    if (value !== true) {
      index += 1;
    }

    if (args[key] === undefined) {
      args[key] = value;
    } else if (Array.isArray(args[key])) {
      args[key].push(value);
    } else {
      args[key] = [args[key], value];
    }
  }

  return args;
}

function readFileRequired(filePath) {
  return fs.readFileSync(filePath, 'utf8');
}

function ensureDir(dirPath) {
  fs.mkdirSync(dirPath, { recursive: true });
}

function fileExists(filePath) {
  return fs.existsSync(filePath);
}

function getWorkflowConfig(cwd = process.cwd()) {
  const configPath = path.join(cwd, 'workflow', 'config.json');
  if (!fileExists(configPath)) {
    return { ...DEFAULT_CONFIG, configPath };
  }

  const parsed = JSON.parse(readFileRequired(configPath));
  return {
    project: { ...DEFAULT_CONFIG.project, ...(parsed.project || {}) },
    paths: { ...DEFAULT_CONFIG.paths, ...(parsed.paths || {}) },
    task_backend: { ...DEFAULT_CONFIG.task_backend, ...(parsed.task_backend || {}) },
    memory_provider: { ...DEFAULT_CONFIG.memory_provider, ...(parsed.memory_provider || {}) },
    ...parsed,
    configPath,
  };
}

function getDocsRoot(cwd = process.cwd()) {
  return path.join(cwd, getWorkflowConfig(cwd).paths.docs_root || 'docs');
}

function getProposalsRoot(cwd = process.cwd()) {
  return path.join(getDocsRoot(cwd), 'proposals');
}

function getTemplateRoot(cwd = process.cwd()) {
  return path.join(cwd, 'workflow', 'templates', 'docs');
}

function slugify(value) {
  return String(value || '')
    .normalize('NFKD')
    .replace(/[^\w\s-]/g, '')
    .trim()
    .toLowerCase()
    .replace(/[_\s-]+/g, '-')
    .replace(/^-+|-+$/g, '') || 'proposal';
}

function nowIso() {
  return new Date().toISOString().replace(/\.\d{3}Z$/, 'Z');
}

function formatDateCode(date = new Date()) {
  const year = String(date.getFullYear()).slice(-2);
  const month = String(date.getMonth() + 1).padStart(2, '0');
  const day = String(date.getDate()).padStart(2, '0');
  return `${year}${month}${day}`;
}

function isProposalId(value) {
  return PROPOSAL_ID_RE.test(String(value || ''));
}

function yamlScalar(value) {
  if (value === null || value === undefined) return 'null';
  if (typeof value === 'boolean') return value ? 'true' : 'false';
  if (typeof value === 'number') return String(value);
  return JSON.stringify(String(value));
}

function serializeYaml(value, indent = 0) {
  const prefix = ' '.repeat(indent);

  if (Array.isArray(value)) {
    return value
      .map((item) => {
        if (item && typeof item === 'object' && !Array.isArray(item)) {
          const entries = Object.entries(item);
          if (entries.length === 0) return `${prefix}- {}`;
          const [firstKey, firstValue] = entries[0];
          const firstLine = `${prefix}- ${firstKey}: ${isPlainScalar(firstValue) ? yamlScalar(firstValue) : ''}`.trimEnd();
          const restLines = [];
          if (!isPlainScalar(firstValue)) {
            restLines.push(serializeYaml(firstValue, indent + 4));
          }
          for (const [key, nestedValue] of entries.slice(1)) {
            if (isPlainScalar(nestedValue)) {
              restLines.push(`${' '.repeat(indent + 2)}${key}: ${yamlScalar(nestedValue)}`);
            } else {
              restLines.push(`${' '.repeat(indent + 2)}${key}:`);
              restLines.push(serializeYaml(nestedValue, indent + 4));
            }
          }
          return [firstLine, ...restLines].join('\n');
        }

        return `${prefix}- ${yamlScalar(item)}`;
      })
      .join('\n');
  }

  if (value && typeof value === 'object') {
    return Object.entries(value)
      .map(([key, nestedValue]) => {
        if (isPlainScalar(nestedValue)) {
          return `${prefix}${key}: ${yamlScalar(nestedValue)}`;
        }
        return `${prefix}${key}:\n${serializeYaml(nestedValue, indent + 2)}`;
      })
      .join('\n');
  }

  return `${prefix}${yamlScalar(value)}`;
}

function isPlainScalar(value) {
  return value === null || ['string', 'number', 'boolean'].includes(typeof value);
}

function renderTemplate(templateText, replacements) {
  let text = templateText;
  for (const [key, value] of Object.entries(replacements)) {
    const pattern = new RegExp(key.replace(/[.*+?^${}()|[\]\\]/g, '\\$&'), 'g');
    text = text.replace(pattern, value);
  }
  return text;
}

function findNextProposalId(cwd = process.cwd(), date = new Date()) {
  const proposalsRoot = getProposalsRoot(cwd);
  ensureDir(proposalsRoot);

  const prefix = `CR${formatDateCode(date)}`;
  let maxSequence = 0;

  for (const entry of fs.readdirSync(proposalsRoot, { withFileTypes: true })) {
    if (!entry.isDirectory()) continue;
    const match = entry.name.match(new RegExp(`^${prefix}(\\d{2})(?:-|$)`));
    if (!match) continue;
    maxSequence = Math.max(maxSequence, Number(match[1]));
  }

  const nextSequence = String(maxSequence + 1).padStart(2, '0');
  return `${prefix}${nextSequence}`;
}

function resolveProposalDir(target, cwd = process.cwd()) {
  if (!target) {
    throw new Error('Missing proposal id or path');
  }

  const absoluteTarget = path.isAbsolute(target) ? target : path.join(cwd, target);
  if (fileExists(absoluteTarget) && fs.statSync(absoluteTarget).isDirectory()) {
    return absoluteTarget;
  }

  if (isProposalId(target)) {
    const proposalsRoot = getProposalsRoot(cwd);
    if (!fileExists(proposalsRoot)) {
      throw new Error(`Proposal root not found: ${proposalsRoot}`);
    }

    const children = fs.readdirSync(proposalsRoot, { withFileTypes: true });
    const match = children.find((entry) => entry.isDirectory() && entry.name.startsWith(`${target}-`));

    if (!match) {
      throw new Error(`Proposal not found for id ${target}`);
    }

    return path.join(proposalsRoot, match.name);
  }

  throw new Error(`Proposal directory not found: ${target}`);
}

function loadProposalContext(target, cwd = process.cwd()) {
  const proposalDir = resolveProposalDir(target, cwd);
  const metaPath = path.join(proposalDir, 'meta.yaml');
  const taskMapPath = path.join(proposalDir, 'task-map.md');
  const designPath = path.join(proposalDir, 'design.md');
  const notesPath = path.join(proposalDir, 'notes.md');

  const meta = parseSimpleYaml(readFileRequired(metaPath));
  const taskMap = parseTaskMap(readFileRequired(taskMapPath));

  return {
    proposalDir,
    metaPath,
    taskMapPath,
    designPath,
    notesPath,
    meta,
    taskMap,
  };
}

function validateProposalContext(context, cwd = process.cwd()) {
  const errors = [];
  const warnings = [];
  const { meta, taskMap, proposalDir, designPath, notesPath } = context;

  for (const field of META_REQUIRED_FIELDS) {
    if (meta[field] === undefined || meta[field] === null || meta[field] === '') {
      errors.push(`meta.yaml missing required field: ${field}`);
    }
  }

  if (meta.id && !isProposalId(meta.id)) {
    errors.push(`meta.yaml id must match CRYYMMDDNN: ${meta.id}`);
  }

  if (meta.schema_version && meta.schema_version !== 'v1') {
    errors.push(`unsupported schema_version: ${meta.schema_version}`);
  }

  if (meta.status && !STATUS_VALUES.has(meta.status)) {
    errors.push(`invalid proposal status: ${meta.status}`);
  }

  if (meta.task_backend && !TASK_BACKENDS.has(meta.task_backend)) {
    errors.push(`invalid task backend: ${meta.task_backend}`);
  }

  if (!Array.isArray(meta.archive_targets) || meta.archive_targets.length === 0) {
    errors.push('meta.yaml must define at least one archive target');
  } else {
    const primaryTargets = meta.archive_targets.filter((target) => target?.role === 'primary');
    if (primaryTargets.length !== 1) {
      errors.push(`meta.yaml must define exactly one primary archive target, found ${primaryTargets.length}`);
    }

    for (const target of meta.archive_targets) {
      if (!ARCHIVE_TARGET_TYPES.has(target?.type)) {
        errors.push(`invalid archive target type: ${target?.type}`);
      }
      if (!ARCHIVE_TARGET_ROLES.has(target?.role)) {
        errors.push(`invalid archive target role: ${target?.role}`);
      }
      if (!target?.path) {
        errors.push('archive target path is required');
      }
    }
  }

  if (!meta.links || typeof meta.links !== 'object') {
    errors.push('meta.yaml links block is required');
  } else {
    for (const key of ['design', 'task_map', 'notes']) {
      if (!meta.links[key]) {
        errors.push(`meta.yaml links.${key} is required`);
      }
    }
  }

  if (!fileExists(designPath)) {
    errors.push(`design.md not found: ${designPath}`);
  }
  if (!fileExists(notesPath)) {
    warnings.push(`notes.md not found: ${notesPath}`);
  }

  if (!taskMap.proposal_id) {
    errors.push('task-map.md missing Proposal ID');
  } else if (meta.id && taskMap.proposal_id !== meta.id) {
    errors.push(`task-map proposal id ${taskMap.proposal_id} does not match meta id ${meta.id}`);
  }

  if (!taskMap.backend) {
    errors.push('task-map.md missing Backend');
  } else if (meta.task_backend && taskMap.backend !== meta.task_backend) {
    errors.push(`task-map backend ${taskMap.backend} does not match meta task_backend ${meta.task_backend}`);
  }

  if (meta.task_backend !== 'none' && taskMap.tasks.length === 0) {
    errors.push('task-map.md must define tasks when task_backend is not none');
  }

  const taskIds = new Set();
  for (const task of taskMap.tasks) {
    if (!task.task_id) {
      errors.push('task missing task_id');
      continue;
    }
    if (taskIds.has(task.task_id)) {
      errors.push(`duplicate task id: ${task.task_id}`);
    }
    taskIds.add(task.task_id);

    for (const field of ['title', 'outcome']) {
      if (!task[field]) {
        errors.push(`${task.task_id} missing ${field}`);
      }
    }

    if (task.priority && !TASK_PRIORITIES.has(task.priority)) {
      errors.push(`${task.task_id} has invalid priority ${task.priority}`);
    }
    if (!Array.isArray(task.capability_refs) || task.capability_refs.length === 0) {
      errors.push(`${task.task_id} must reference at least one capability`);
    }
    if (!Array.isArray(task.completion_definition) || task.completion_definition.length === 0) {
      errors.push(`${task.task_id} must define completion criteria`);
    }
  }

  for (const task of taskMap.tasks) {
    for (const dependency of task.depends_on || []) {
      if (!taskIds.has(dependency)) {
        errors.push(`${task.task_id} depends on unknown task ${dependency}`);
      }
    }
  }

  const archiveTargetRefs = new Set(
    (meta.archive_targets || []).map((target) => `${target.type}:${path.basename(target.path)}`)
  );
  for (const task of taskMap.tasks) {
    for (const ref of task.archive_target_refs || []) {
      if (!archiveTargetRefs.has(ref)) {
        warnings.push(`${task.task_id} archive target ref not found in meta.yaml: ${ref}`);
      }
    }
  }

  if (meta.source_exploration) {
    const sourcePath = path.isAbsolute(meta.source_exploration)
      ? meta.source_exploration
      : path.join(cwd, meta.source_exploration);
    if (!fileExists(sourcePath)) {
      warnings.push(`source exploration path does not exist locally: ${meta.source_exploration}`);
    }
  }

  if (!proposalDir.startsWith(getProposalsRoot(cwd))) {
    warnings.push(`proposal dir is outside configured proposal root: ${proposalDir}`);
  }

  return { errors, warnings };
}

function runCommand(command, args, cwd = process.cwd(), timeout = 20000) {
  const result = spawnSync(command, args, {
    cwd,
    encoding: 'utf8',
    timeout,
  });

  if (result.error) {
    return {
      ok: false,
      error: result.error.message,
      stdout: result.stdout || '',
      stderr: result.stderr || '',
      status: result.status,
    };
  }

  if (result.status !== 0) {
    return {
      ok: false,
      error: (result.stderr || result.stdout || '').trim() || `${command} exited with ${result.status}`,
      stdout: result.stdout || '',
      stderr: result.stderr || '',
      status: result.status,
    };
  }

  return {
    ok: true,
    stdout: result.stdout || '',
    stderr: result.stderr || '',
    status: result.status,
  };
}

function beadTaskSummary(proposalId, cwd = process.cwd()) {
  const command = runCommand('bd', ['query', `spec=${proposalId}`, '--json'], cwd);

  if (!command.ok) {
    return {
      available: false,
      error: command.error,
      tasks: [],
    };
  }

  try {
    const tasks = JSON.parse(command.stdout || '[]');
    const openTasks = tasks.filter((task) => !['closed', 'done', 'completed'].includes(String(task.status || '').toLowerCase()));
    return {
      available: true,
      tasks,
      openTasks,
    };
  } catch (error) {
    return {
      available: false,
      error: `Failed to parse bd output: ${error.message}`,
      tasks: [],
    };
  }
}

function loadTemplate(cwd, relativePath) {
  return readFileRequired(path.join(getTemplateRoot(cwd), relativePath));
}

function archiveTargetRef(target) {
  return `${target.type}:${path.basename(target.path)}`;
}

function getDefaultOwner() {
  const fromGit = runCommand('git', ['config', 'user.name'], process.cwd(), 5000);
  if (fromGit.ok && fromGit.stdout.trim()) return fromGit.stdout.trim();
  return process.env.USER || 'unknown-owner';
}

function createProposalSkeleton(options, cwd = process.cwd()) {
  const docsRoot = getDocsRoot(cwd);
  const proposalsRoot = path.join(docsRoot, 'proposals');
  ensureDir(proposalsRoot);

  const id = options.id || findNextProposalId(cwd);
  if (!isProposalId(id)) {
    throw new Error(`proposal id must match CRYYMMDDNN: ${id}`);
  }

  const slug = slugify(options.slug || options.title);
  const proposalDir = path.join(proposalsRoot, `${id}-${slug}`);
  if (fileExists(proposalDir)) {
    throw new Error(`proposal directory already exists: ${proposalDir}`);
  }

  const createdAt = nowIso();
  const title = options.title;
  const taskBackend = options.taskBackend || getWorkflowConfig(cwd).task_backend?.type || 'beads';
  const archiveTargets = options.archiveTargets;
  const primaryTarget = archiveTargets.find((target) => target.role === 'primary');
  const primaryRef = archiveTargetRef(primaryTarget);
  const owner = options.owner || getDefaultOwner();

  ensureDir(proposalDir);

  const meta = {
    schema_version: 'v1',
    id,
    slug,
    title,
    status: options.status || 'proposed',
    created_at: createdAt,
    updated_at: createdAt,
    source_exploration: options.sourceExploration,
    owner,
    task_backend: taskBackend,
    task_epic_id: null,
    archive_targets: archiveTargets,
    tags: options.tags || [],
    links: {
      design: 'design.md',
      task_map: 'task-map.md',
      notes: 'notes.md',
    },
  };

  const proposalTemplate = loadTemplate(cwd, path.join('proposals', 'proposal.md'));
  const designTemplate = loadTemplate(cwd, path.join('proposals', 'design.md'));
  const taskMapTemplate = loadTemplate(cwd, path.join('proposals', 'task-map.md'));
  const notesTemplate = loadTemplate(cwd, path.join('proposals', 'notes.md'));

  const replacements = {
    '<Proposal Title>': title,
    'CR20260520': id,
    'CR26052001': id,
    'module:example-module': primaryRef,
    '2026-05-20': createdAt.slice(0, 10),
  };

  fs.writeFileSync(path.join(proposalDir, 'meta.yaml'), `${serializeYaml(meta)}\n`, 'utf8');
  fs.writeFileSync(path.join(proposalDir, 'proposal.md'), renderTemplate(proposalTemplate, replacements), 'utf8');
  fs.writeFileSync(path.join(proposalDir, 'design.md'), renderTemplate(designTemplate, replacements), 'utf8');
  fs.writeFileSync(path.join(proposalDir, 'task-map.md'), renderTemplate(taskMapTemplate, replacements), 'utf8');
  fs.writeFileSync(path.join(proposalDir, 'notes.md'), renderTemplate(notesTemplate, replacements), 'utf8');

  return {
    id,
    slug,
    proposalDir,
    meta,
  };
}

function writeProposalMeta(metaPath, meta) {
  fs.writeFileSync(metaPath, `${serializeYaml(meta)}\n`, 'utf8');
}

function transitionProposalStatus(context, nextStatus) {
  const currentStatus = context.meta.status;
  if (!STATUS_VALUES.has(nextStatus)) {
    throw new Error(`invalid target status: ${nextStatus}`);
  }

  if (!STATUS_TRANSITIONS[currentStatus] || !STATUS_TRANSITIONS[currentStatus].has(nextStatus)) {
    throw new Error(`invalid status transition: ${currentStatus} -> ${nextStatus}`);
  }

  context.meta.status = nextStatus;
  context.meta.updated_at = nowIso();
  writeProposalMeta(context.metaPath, context.meta);
  return context.meta;
}

function ensureNotesFile(context, cwd = process.cwd()) {
  if (fileExists(context.notesPath)) return false;
  const template = loadTemplate(cwd, path.join('proposals', 'notes.md'));
  const content = renderTemplate(template, {
    '<Proposal Title>': context.meta.title,
    '2026-05-20': nowIso().slice(0, 10),
  });
  fs.writeFileSync(context.notesPath, content, 'utf8');
  return true;
}

function beadCreateIssue(args, cwd = process.cwd()) {
  const command = runCommand('bd', ['create', ...args, '--silent'], cwd);
  if (!command.ok) {
    throw new Error(`bd create failed: ${command.error}`);
  }

  const issueId = command.stdout.trim();
  if (!issueId) {
    throw new Error('bd create did not return an issue id');
  }
  return issueId;
}

function ensureBeadsTasks(context, cwd = process.cwd()) {
  if (context.meta.task_epic_id) {
    return {
      reused: true,
      epic_id: context.meta.task_epic_id,
      created_task_ids: [],
    };
  }

  const summary = beadTaskSummary(context.meta.id, cwd);
  if (summary.available && summary.tasks.length > 0) {
    throw new Error(`Beads already has tasks for ${context.meta.id}; set meta.task_epic_id or clean up existing tasks first`);
  }
  if (!summary.available) {
    const probe = runCommand('bd', ['context', '--json'], cwd);
    if (!probe.ok) {
      throw new Error(`Beads backend unavailable: ${probe.error}`);
    }
  }

  const epicDescription = `Proposal ${context.meta.id}\n\n${context.meta.proposalDir || context.proposalDir}`;
  const epicLabels = [
    `proposal:${context.meta.id}`,
    `workflow:proposal`,
    `archive:${path.basename((context.meta.archive_targets || [])[0]?.path || 'unknown')}`,
  ];

  const epicId = beadCreateIssue([
    '--title',
    `${context.meta.id} ${context.meta.title}`,
    '--type',
    'epic',
    '--description',
    epicDescription,
    '--spec-id',
    context.meta.id,
    '--labels',
    epicLabels.join(','),
  ], cwd);

  const taskIdMap = new Map();
  for (const task of context.taskMap.tasks) {
    const labels = [
      `proposal:${context.meta.id}`,
      `task-map:${task.task_id}`,
      ...task.capability_refs.map((ref) => `capability:${ref}`),
      ...task.archive_target_refs.map((ref) => `archive-ref:${ref}`),
    ];

    const descriptionLines = [
      task.outcome ? `Outcome: ${task.outcome}` : '',
      task.decision_refs.length > 0 ? `Decision refs: ${task.decision_refs.join(', ')}` : '',
      task.archive_target_refs.length > 0 ? `Archive refs: ${task.archive_target_refs.join(', ')}` : '',
    ].filter(Boolean);

    const beadId = beadCreateIssue([
      '--title',
      task.title,
      '--type',
      'task',
      '--description',
      descriptionLines.join('\n'),
      '--acceptance',
      task.completion_definition.join('\n'),
      '--priority',
      task.priority || 'P1',
      '--parent',
      epicId,
      '--spec-id',
      context.meta.id,
      '--labels',
      labels.join(','),
    ], cwd);

    taskIdMap.set(task.task_id, beadId);
  }

  for (const task of context.taskMap.tasks) {
    for (const dependency of task.depends_on) {
      const currentId = taskIdMap.get(task.task_id);
      const dependencyId = taskIdMap.get(dependency);
      const command = runCommand('bd', ['link', currentId, dependencyId], cwd);
      if (!command.ok) {
        throw new Error(`bd link failed for ${task.task_id} -> ${dependency}: ${command.error}`);
      }
    }
  }

  context.meta.task_epic_id = epicId;
  context.meta.updated_at = nowIso();
  writeProposalMeta(context.metaPath, context.meta);

  return {
    reused: false,
    epic_id: epicId,
    created_task_ids: Array.from(taskIdMap.values()),
  };
}

module.exports = {
  archiveTargetRef,
  beadTaskSummary,
  createProposalSkeleton,
  ensureBeadsTasks,
  ensureNotesFile,
  findNextProposalId,
  formatDateCode,
  getDocsRoot,
  getWorkflowConfig,
  isProposalId,
  loadProposalContext,
  parseCliArgs,
  parseSimpleYaml,
  parseTaskMap,
  resolveProposalDir,
  serializeYaml,
  slugify,
  transitionProposalStatus,
  validateProposalContext,
  writeProposalMeta,
};
