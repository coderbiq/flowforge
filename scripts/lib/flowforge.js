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
const HISTORY_MARKER_PREFIX = '<!-- flowforge:proposal:';
const DEFAULT_TOOL_ROOT = '.flowforge';
const LEGACY_TOOL_ROOT = 'workflow';
const DEFAULT_WORKSPACE_NAME = 'default';
const DEFAULT_CONFIG = {
  project: {
    id: 'unknown-project',
    name: 'Unknown Project',
    slug: 'unknown-project',
  },
  paths: {
    tool_root: DEFAULT_TOOL_ROOT,
    state_root: `${DEFAULT_TOOL_ROOT}/state`,
  },
  docs: {
    default_workspace: DEFAULT_WORKSPACE_NAME,
    workspaces: {
      [DEFAULT_WORKSPACE_NAME]: {
        root: 'docs',
        scope: '.',
        kind: 'repository',
      },
    },
  },
  task_backend: {
    type: 'beads',
  },
  memory_provider: {
    enabled: false,
  },
};

const STATUS_TRANSITIONS = {
  draft: new Set(['draft', 'proposed', 'approved', 'rejected']),
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
      else if (field === 'workspace') currentTask.workspace = value;
      else if (field === 'code scope') currentTask.code_scope = value;
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

function parseExplorationIndex(text) {
  const lines = text.replace(/\r\n/g, '\n').split('\n');
  const result = {
    title: null,
    status: null,
    date: null,
  };

  for (const line of lines) {
    const trimmed = line.trim();
    if (!trimmed) continue;

    const titleMatch = trimmed.match(/^#\s+(.+)$/);
    if (titleMatch && !result.title) {
      result.title = titleMatch[1].trim();
      continue;
    }

    const dateMatch = trimmed.match(/^\*\*日期\*\*:\s*(.+)$/);
    if (dateMatch && !result.date) {
      result.date = dateMatch[1].trim();
      continue;
    }

    const statusMatch = trimmed.match(/^\*\*状态\*\*:\s*(.+)$/);
    if (statusMatch && !result.status) {
      result.status = statusMatch[1].trim();
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

function ensureParentDir(filePath) {
  ensureDir(path.dirname(filePath));
}

function fileExists(filePath) {
  return fs.existsSync(filePath);
}

function writeText(filePath, content) {
  ensureParentDir(filePath);
  fs.writeFileSync(filePath, content, 'utf8');
}

function appendText(filePath, content) {
  ensureParentDir(filePath);
  fs.appendFileSync(filePath, content, 'utf8');
}

function findAncestor(startDir, predicate) {
  let current = path.resolve(startDir);
  while (true) {
    if (predicate(current)) {
      return current;
    }
    const parent = path.dirname(current);
    if (parent === current) {
      return null;
    }
    current = parent;
  }
}

function getProjectRoot(cwd = process.cwd()) {
  const start = path.resolve(cwd);
  const configured = findAncestor(start, (dir) => fileExists(path.join(dir, DEFAULT_TOOL_ROOT, 'config.json')) || fileExists(path.join(dir, LEGACY_TOOL_ROOT, 'config.json')));
  if (configured) {
    return configured;
  }

  const sourceLayout = findAncestor(start, (dir) => fileExists(path.join(dir, 'workflow')) && fileExists(path.join(dir, 'agents')) && fileExists(path.join(dir, 'scripts')));
  if (sourceLayout) {
    return sourceLayout;
  }

  return start;
}

function getToolRoot(cwd = process.cwd()) {
  const projectRoot = getProjectRoot(cwd);
  const flowforgeRoot = path.join(projectRoot, DEFAULT_TOOL_ROOT);

  if (fileExists(path.join(flowforgeRoot, 'config.json')) || fileExists(path.join(flowforgeRoot, 'workflow'))) {
    return flowforgeRoot;
  }

  if (fileExists(path.join(projectRoot, 'workflow')) || fileExists(path.join(projectRoot, 'agents')) || fileExists(path.join(projectRoot, 'scripts'))) {
    return projectRoot;
  }

  return flowforgeRoot;
}

function getWorkflowConfigPath(cwd = process.cwd()) {
  const projectRoot = getProjectRoot(cwd);
  const flowforgeConfig = path.join(projectRoot, DEFAULT_TOOL_ROOT, 'config.json');
  if (fileExists(flowforgeConfig)) {
    return flowforgeConfig;
  }

  const legacy = path.join(projectRoot, LEGACY_TOOL_ROOT, 'config.json');
  if (fileExists(legacy)) {
    return legacy;
  }

  return flowforgeConfig;
}

function normalizeDocsConfig(parsed, paths) {
  const rawDocs = parsed.docs && typeof parsed.docs === 'object' ? parsed.docs : null;

  if (!rawDocs || !rawDocs.workspaces || Object.keys(rawDocs.workspaces).length === 0) {
    const docsRoot = parsed.paths?.docs_root || paths.docs_root || 'docs';
    return {
      default_workspace: DEFAULT_WORKSPACE_NAME,
      workspaces: {
        [DEFAULT_WORKSPACE_NAME]: {
          root: docsRoot,
          scope: '.',
          kind: 'repository',
        },
      },
    };
  }

  const workspaces = {};
  for (const [name, workspace] of Object.entries(rawDocs.workspaces)) {
    workspaces[name] = {
      root: workspace?.root || paths.docs_root || 'docs',
      scope: workspace?.scope || '.',
      kind: workspace?.kind || 'project',
      label: workspace?.label,
      owners: workspace?.owners,
    };
  }

  const defaultWorkspace = rawDocs.default_workspace && workspaces[rawDocs.default_workspace]
    ? rawDocs.default_workspace
    : Object.keys(workspaces)[0] || DEFAULT_WORKSPACE_NAME;

  if (!workspaces[defaultWorkspace]) {
    workspaces[defaultWorkspace] = {
      root: parsed.paths?.docs_root || paths.docs_root || 'docs',
      scope: '.',
      kind: 'repository',
    };
  }

  return {
    default_workspace: defaultWorkspace,
    workspaces,
  };
}

function getWorkflowConfig(cwd = process.cwd()) {
  const configPath = getWorkflowConfigPath(cwd);
  const parsed = fileExists(configPath) ? JSON.parse(readFileRequired(configPath)) : {};
  const paths = {
    ...DEFAULT_CONFIG.paths,
    ...(parsed.paths || {}),
  };

  if (!paths.tool_root) {
    paths.tool_root = DEFAULT_TOOL_ROOT;
  }

  if (!paths.state_root) {
    paths.state_root = path.join(paths.tool_root, 'state');
  }

  return {
    ...parsed,
    project: { ...DEFAULT_CONFIG.project, ...(parsed.project || {}) },
    paths,
    docs: normalizeDocsConfig(parsed, paths),
    task_backend: { ...DEFAULT_CONFIG.task_backend, ...(parsed.task_backend || {}) },
    memory_provider: { ...DEFAULT_CONFIG.memory_provider, ...(parsed.memory_provider || {}) },
    configPath,
  };
}

function listDocsWorkspaces(cwd = process.cwd()) {
  const config = getWorkflowConfig(cwd);
  return Object.entries(config.docs?.workspaces || {}).map(([name, workspace]) => ({
    name,
    ...workspace,
  }));
}

function getDocsWorkspace(name = null, cwd = process.cwd()) {
  const config = getWorkflowConfig(cwd);
  const workspaces = config.docs?.workspaces || {};
  const resolvedName = name || config.docs?.default_workspace || Object.keys(workspaces)[0] || DEFAULT_WORKSPACE_NAME;
  const workspace = workspaces[resolvedName];
  if (!workspace) {
    throw new Error(`Workspace not found: ${resolvedName}`);
  }
  return { name: resolvedName, ...workspace };
}

function resolveWorkspaceForCwd(cwd = process.cwd()) {
  const config = getWorkflowConfig(cwd);
  const workspaces = listDocsWorkspaces(cwd);
  if (workspaces.length === 0) {
    return getDocsWorkspace(DEFAULT_WORKSPACE_NAME, cwd);
  }

  const projectRoot = getProjectRoot(cwd);
  const relativeCwd = path.relative(projectRoot, cwd) || '.';
  const normalizedRelative = relativeCwd === '' ? '.' : relativeCwd.split(path.sep).join('/');

  const matches = workspaces
    .map((workspace) => {
      const scope = String(workspace.scope || '.').split(path.sep).join('/');
      if (scope === '.' || scope === '') {
        return { workspace, depth: 0 };
      }
      if (normalizedRelative === scope || normalizedRelative.startsWith(`${scope}/`)) {
        return { workspace, depth: scope.split('/').filter(Boolean).length };
      }
      return null;
    })
    .filter(Boolean)
    .sort((a, b) => b.depth - a.depth);

  if (matches.length === 0) {
    return getDocsWorkspace(config.docs?.default_workspace, cwd);
  }

  const topDepth = matches[0].depth;
  const best = matches.filter((match) => match.depth === topDepth);
  if (best.length > 1) {
    throw new Error(`Ambiguous workspace for cwd ${cwd}: ${best.map((item) => item.workspace.name).join(', ')}`);
  }

  return best[0].workspace;
}

function getWorkspaceDocsRoot(workspaceName = null, cwd = process.cwd()) {
  const workspace = workspaceName ? getDocsWorkspace(workspaceName, cwd) : resolveWorkspaceForCwd(cwd);
  return path.join(getProjectRoot(cwd), workspace.root || 'docs');
}

function getDocsRoot(cwd = process.cwd(), workspaceName = null) {
  return getWorkspaceDocsRoot(workspaceName, cwd);
}

function getProposalsRoot(cwd = process.cwd(), workspaceName = null) {
  return path.join(getWorkspaceDocsRoot(workspaceName, cwd), 'proposals');
}

function getTemplateRoot(cwd = process.cwd()) {
  return path.join(getToolRoot(cwd), 'workflow', 'templates', 'docs');
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

function todayDate() {
  return nowIso().slice(0, 10);
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

function isPlainScalar(value) {
  return value === null || ['string', 'number', 'boolean'].includes(typeof value);
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

function renderTemplate(templateText, replacements) {
  let text = templateText;
  for (const [key, value] of Object.entries(replacements)) {
    const pattern = new RegExp(key.replace(/[.*+?^${}()|[\]\\]/g, '\\$&'), 'g');
    text = text.replace(pattern, value);
  }
  return text;
}

function renderTaskMapTemplate(templateText, replacements, taskBackend) {
  if (taskBackend === 'none') {
    return [
      `# Task Map: ${replacements['<Proposal Title>']}`,
      '',
      '- Backend: none',
      `- Proposal ID: ${replacements.CR26052001 || replacements.CR20260520}`,
      '',
      '## Tasks',
      '',
      'No external task backend is configured for this proposal.',
      '',
    ].join('\n');
  }

  return renderTemplate(templateText, replacements);
}

function findNextProposalId(cwd = process.cwd(), date = new Date(), workspaceName = null) {
  const proposalsRoot = getProposalsRoot(cwd, workspaceName);
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

  const projectRoot = getProjectRoot(cwd);
  const absoluteTarget = path.isAbsolute(target) ? target : path.join(projectRoot, target);
  if (fileExists(absoluteTarget) && fs.statSync(absoluteTarget).isDirectory()) {
    return absoluteTarget;
  }

  if (isProposalId(target)) {
    const workspaces = listDocsWorkspaces(cwd);
    const matches = [];

    for (const workspace of workspaces) {
      const proposalsRoot = getProposalsRoot(cwd, workspace.name);
      if (!fileExists(proposalsRoot)) continue;
      const children = fs.readdirSync(proposalsRoot, { withFileTypes: true });
      for (const entry of children) {
        if (entry.isDirectory() && entry.name.startsWith(`${target}-`)) {
          matches.push(path.join(proposalsRoot, entry.name));
        }
      }
    }

    if (matches.length === 0) {
      throw new Error(`Proposal not found for id ${target}`);
    }

    if (matches.length > 1) {
      throw new Error(`Proposal id ${target} resolved to multiple directories: ${matches.join(', ')}`);
    }

    return matches[0];
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

  const proposalRoot = path.dirname(proposalDir);

  const schemaVersion = meta.schema_version || 'v1';
  const workspaceConfig = getWorkflowConfig(cwd).docs || {};
  const workspaceNames = new Set(Object.keys(workspaceConfig.workspaces || {}));

  if (schemaVersion === 'v2') {
    const v2Required = [
      'schema_version',
      'id',
      'slug',
      'title',
      'status',
      'created_at',
      'updated_at',
      'workspace',
      'scope',
      'source_explorations',
      'owner',
      'task_backend',
      'archive_targets',
      'links',
    ];
    for (const field of v2Required) {
      if (meta[field] === undefined || meta[field] === null || meta[field] === '') {
        errors.push(`meta.yaml missing required field: ${field}`);
      }
    }

    if (meta.id && !isProposalId(meta.id)) {
      errors.push(`meta.yaml id must match CRYYMMDDNN: ${meta.id}`);
    }
    if (meta.status && !STATUS_VALUES.has(meta.status)) {
      errors.push(`invalid proposal status: ${meta.status}`);
    }
    if (meta.task_backend && !TASK_BACKENDS.has(meta.task_backend)) {
      errors.push(`invalid task backend: ${meta.task_backend}`);
    }
    if (meta.workspace && workspaceNames.size > 0 && !workspaceNames.has(meta.workspace)) {
      errors.push(`unknown proposal workspace: ${meta.workspace}`);
    }
    if (meta.scope && !['workspace', 'cross-workspace', 'monorepo'].includes(meta.scope)) {
      errors.push(`invalid proposal scope: ${meta.scope}`);
    }

    if (!Array.isArray(meta.source_explorations) || meta.source_explorations.length === 0) {
      errors.push('meta.yaml must define at least one source exploration');
    } else {
      for (const source of meta.source_explorations) {
        if (!source?.workspace) errors.push('source exploration workspace is required');
        if (!source?.ref) errors.push('source exploration ref is required');
        if (source?.workspace && workspaceNames.size > 0 && !workspaceNames.has(source.workspace)) {
          errors.push(`unknown source exploration workspace: ${source.workspace}`);
        }
        if (source?.ref && path.isAbsolute(source.ref)) {
          errors.push(`source exploration ref must be relative: ${source.ref}`);
        }
      }
    }

    const targetTypeWorkspacePairs = new Set();
    for (const target of meta.archive_targets || []) {
      const targetType = target?.type;
      const targetWorkspace = target?.workspace;
      if (!targetType || !targetWorkspace || !ARCHIVE_TARGET_TYPES.has(targetType)) continue;
      const key = `${targetWorkspace}:${targetType}`;
      if (targetTypeWorkspacePairs.has(key)) continue;
      targetTypeWorkspacePairs.add(key);
      const existingDocs = scanWorkspaceDocsForCanonicalCorpus(targetWorkspace, cwd, new Set([targetType]));
      if (existingDocs.length === 0) {
        warnings.push(`no existing canonical corpus docs found for ${targetWorkspace}:${targetType}; this proposal will establish the baseline`);
      }
    }

    if (meta.canonical_corpus !== undefined) {
      if (!Array.isArray(meta.canonical_corpus)) {
        errors.push('meta.yaml canonical_corpus must be an array when present');
      } else {
        for (const entry of meta.canonical_corpus) {
          if (!entry?.workspace) errors.push('canonical corpus workspace is required');
          if (!entry?.ref) errors.push('canonical corpus ref is required');
          if (!entry?.type) errors.push('canonical corpus type is required');
          if (entry?.type && !ARCHIVE_TARGET_TYPES.has(entry.type)) {
            errors.push(`invalid canonical corpus type: ${entry.type}`);
          }
          if (entry?.workspace && workspaceNames.size > 0 && !workspaceNames.has(entry.workspace)) {
            errors.push(`unknown canonical corpus workspace: ${entry.workspace}`);
          }
          if (entry?.ref && path.isAbsolute(entry.ref)) {
            errors.push(`canonical corpus ref must be relative: ${entry.ref}`);
          }
          if (entry?.workspace && entry?.ref) {
            const corpusRoot = getWorkspaceDocsRoot(entry.workspace, cwd);
            const corpusPath = path.join(corpusRoot, entry.ref);
            if (!fileExists(corpusPath)) {
              errors.push(`canonical corpus path does not exist locally: ${entry.workspace}:${entry.ref}`);
            }
          }
        }
      }
    }

    if (!Array.isArray(meta.archive_targets) || meta.archive_targets.length === 0) {
      errors.push('meta.yaml must define at least one archive target');
    } else {
      const primaryTargets = meta.archive_targets.filter((target) => target?.role === 'primary');
      if (primaryTargets.length !== 1) {
        errors.push(`meta.yaml must define exactly one primary archive target, found ${primaryTargets.length}`);
      }

      const keys = new Set();
      for (const target of meta.archive_targets) {
        if (!target?.key) errors.push('archive target key is required');
        if (target?.key && keys.has(target.key)) errors.push(`duplicate archive target key: ${target.key}`);
        if (target?.key) keys.add(target.key);
        if (!ARCHIVE_TARGET_TYPES.has(target?.type)) {
          errors.push(`invalid archive target type: ${target?.type}`);
        }
        if (!ARCHIVE_TARGET_ROLES.has(target?.role)) {
          errors.push(`invalid archive target role: ${target?.role}`);
        }
        if (!target?.workspace) {
          errors.push('archive target workspace is required');
        } else if (workspaceNames.size > 0 && !workspaceNames.has(target.workspace)) {
          errors.push(`unknown archive target workspace: ${target.workspace}`);
        }
        if (!target?.ref) {
          errors.push('archive target ref is required');
        } else if (path.isAbsolute(target.ref)) {
          errors.push(`archive target ref must be relative: ${target.ref}`);
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
  } else {
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

    if (task.workspace && workspaceNames.size > 0 && !workspaceNames.has(task.workspace)) {
      errors.push(`${task.task_id} has unknown workspace ${task.workspace}`);
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
    (meta.archive_targets || []).map((target) => {
      if (meta.schema_version === 'v2') {
        return target.key;
      }
      return `${target.type}:${path.basename(target.path)}`;
    }).filter(Boolean)
  );
  for (const task of taskMap.tasks) {
    for (const ref of task.archive_target_refs || []) {
      if (!archiveTargetRefs.has(ref)) {
        warnings.push(`${task.task_id} archive target ref not found in meta.yaml: ${ref}`);
      }
    }
  }

  if (meta.schema_version === 'v2') {
    for (const source of meta.source_explorations || []) {
      const sourceRoot = source?.workspace ? getWorkspaceDocsRoot(source.workspace, cwd) : getDocsRoot(cwd);
      const sourcePath = path.join(sourceRoot, source.ref || '');
      if (!fileExists(sourcePath)) {
        warnings.push(`source exploration path does not exist locally: ${source.workspace}:${source.ref}`);
      }
    }
  } else if (meta.source_exploration) {
    const sourcePath = path.isAbsolute(meta.source_exploration)
      ? meta.source_exploration
      : path.join(cwd, meta.source_exploration);
    if (!fileExists(sourcePath)) {
      warnings.push(`source exploration path does not exist locally: ${meta.source_exploration}`);
    }
  }

  const proposalRoots = listDocsWorkspaces(cwd).map((workspace) => getProposalsRoot(cwd, workspace.name));
  if (!proposalRoots.some((root) => proposalDir.startsWith(root))) {
    warnings.push(`proposal dir is outside configured proposal roots: ${proposalDir}`);
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
      epics: [],
      workItems: [],
      openTasks: [],
      openWorkItems: [],
    };
  }

  try {
    const tasks = JSON.parse(command.stdout || '[]');
    const epics = tasks.filter((task) => String(task.issue_type || '').toLowerCase() === 'epic');
    const workItems = tasks.filter((task) => String(task.issue_type || '').toLowerCase() !== 'epic');
    const openTasks = tasks.filter((task) => !['closed', 'done', 'completed'].includes(String(task.status || '').toLowerCase()));
    const openWorkItems = workItems.filter((task) => !['closed', 'done', 'completed'].includes(String(task.status || '').toLowerCase()));
    return {
      available: true,
      tasks,
      epics,
      workItems,
      openTasks,
      openWorkItems,
    };
  } catch (error) {
    return {
      available: false,
      error: `Failed to parse bd output: ${error.message}`,
      tasks: [],
      epics: [],
      workItems: [],
      openTasks: [],
      openWorkItems: [],
    };
  }
}

function loadTemplate(cwd, relativePath) {
  return readFileRequired(path.join(getTemplateRoot(cwd), relativePath));
}

function archiveTargetRef(target) {
  if (target?.key) return target.key;
  if (target?.ref) return `${target.type}:${path.basename(target.ref)}`;
  return `${target?.type || 'unknown'}:${path.basename(target?.path || 'unknown')}`;
}

function normalizeCanonicalCorpusEntry(entry, workspaceName) {
  if (!entry) return null;
  const workspace = entry.workspace || workspaceName;
  const ref = entry.ref || entry.path;
  const type = entry.type;
  if (!workspace || !ref || !type) return null;
  return {
    workspace,
    ref,
    type,
    role: entry.role || 'secondary',
  };
}

function appendCanonicalCorpusEntries(target, entries, workspaceName) {
  const seen = new Set(target.map((entry) => `${entry.workspace}:${entry.type}:${entry.ref}:${entry.role}`));
  for (const rawEntry of entries || []) {
    const entry = normalizeCanonicalCorpusEntry(rawEntry, workspaceName);
    if (!entry) continue;
    const key = `${entry.workspace}:${entry.type}:${entry.ref}:${entry.role}`;
    if (seen.has(key)) continue;
    seen.add(key);
    target.push(entry);
  }
  return target;
}

function getCanonicalCorpusTypesForArchiveTargets(archiveTargets) {
  return new Set((archiveTargets || [])
    .map((target) => target?.type)
    .filter((type) => ARCHIVE_TARGET_TYPES.has(type)));
}

function scanWorkspaceDocsForCanonicalCorpus(workspaceName, cwd = process.cwd(), allowedTypes = null) {
  const workspaceRoot = getWorkspaceDocsRoot(workspaceName, cwd);
  const corpus = [];
  const typeAllowed = (type) => !allowedTypes || allowedTypes.has(type);

  const modulesRoot = path.join(workspaceRoot, 'modules');
  if (typeAllowed('module') && fileExists(modulesRoot)) {
    for (const entry of fs.readdirSync(modulesRoot, { withFileTypes: true })) {
      if (!entry.isDirectory()) continue;
      corpus.push({
        workspace: workspaceName,
        ref: path.join('modules', entry.name),
        type: 'module',
        role: 'secondary',
      });
    }
  }

  const architectureRoot = path.join(workspaceRoot, 'architecture');
  if (typeAllowed('architecture') && fileExists(architectureRoot)) {
    for (const entry of fs.readdirSync(architectureRoot, { withFileTypes: true })) {
      if (!entry.isFile() || !entry.name.endsWith('.md')) continue;
      corpus.push({
        workspace: workspaceName,
        ref: path.join('architecture', entry.name),
        type: 'architecture',
        role: 'secondary',
      });
    }
  }

  const decisionsRoot = path.join(workspaceRoot, 'decisions');
  if (typeAllowed('decision') && fileExists(decisionsRoot)) {
    for (const entry of fs.readdirSync(decisionsRoot, { withFileTypes: true })) {
      if (!entry.isFile() || !entry.name.endsWith('.md')) continue;
      corpus.push({
        workspace: workspaceName,
        ref: path.join('decisions', entry.name),
        type: 'decision',
        role: 'secondary',
      });
    }
  }

  return corpus;
}

function sortCanonicalCorpusEntries(entries) {
  return [...entries].sort((a, b) => {
    const aKey = `${a.workspace}:${a.type}:${a.ref}:${a.role}`;
    const bKey = `${b.workspace}:${b.type}:${b.ref}:${b.role}`;
    return aKey.localeCompare(bKey);
  });
}

function renderCanonicalCorpusList(entries, proposalDir, cwd = process.cwd()) {
  if (!entries || entries.length === 0) {
    return '- <none>';
  }

  const seen = new Set();
  const lines = [];

  for (const entry of sortCanonicalCorpusEntries(entries)) {
    const key = `${entry.workspace}:${entry.type}:${entry.ref}:${entry.role || 'secondary'}`;
    if (seen.has(key)) continue;
    seen.add(key);
    const docsRoot = getWorkspaceDocsRoot(entry.workspace, cwd);
    const absolutePath = path.join(docsRoot, entry.ref);
    const relativeLink = path.relative(proposalDir, absolutePath).split(path.sep).join('/');
    const label = `${entry.type}${entry.role ? `/${entry.role}` : ''}`;
    lines.push(`- ${label} @ ${entry.workspace}: [${entry.ref}](${relativeLink})`);
  }

  return lines.join('\n');
}

function renderKnowledgeImpact(entries) {
  if (!entries || entries.length === 0) {
    return [
      '- What is reused from the canonical corpus: 0 baseline docs reviewed; this proposal establishes the initial baseline',
      '- What changes in the canonical corpus: <proposal-specific edits to formal docs>',
      '- What new material must be added: <new module, architecture, or ADR content>',
    ].join('\n');
  }

  const modules = entries.filter((entry) => entry.type === 'module').length;
  const architecture = entries.filter((entry) => entry.type === 'architecture').length;
  const decisions = entries.filter((entry) => entry.type === 'decision').length;
  return [
    `- What is reused from the canonical corpus: ${modules + architecture + decisions} baseline doc(s) reviewed`,
    '- What changes in the canonical corpus: <proposal-specific edits to formal docs>',
    '- What new material must be added: <new module, architecture, or ADR content>',
  ].join('\n');
}


function getDefaultOwner() {
  const fromGit = runCommand('git', ['config', 'user.name'], process.cwd(), 5000);
  if (fromGit.ok && fromGit.stdout.trim()) return fromGit.stdout.trim();
  return process.env.USER || 'unknown-owner';
}

function createProposalSkeleton(options, cwd = process.cwd()) {
  const workspaceName = options.workspace || resolveWorkspaceForCwd(cwd).name;
  const docsRoot = getDocsRoot(cwd, workspaceName);
  const proposalsRoot = path.join(docsRoot, 'proposals');
  ensureDir(proposalsRoot);

  const id = options.id || findNextProposalId(cwd, new Date(), workspaceName);
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
  const sourceExplorations = (options.sourceExplorations && options.sourceExplorations.length > 0)
    ? options.sourceExplorations
    : [{ workspace: workspaceName, ref: options.sourceExploration }];
  const archiveTargets = (options.archiveTargets || []).map((target, index) => {
    const targetWorkspace = target.workspace || workspaceName;
    const ref = target.ref || target.path;
    const key = target.key || `${target.type}-${slugify(path.basename(ref || `target-${index + 1}`))}`;
    return {
      key,
      type: target.type,
      workspace: targetWorkspace,
      ref,
      role: target.role || (index === 0 ? 'primary' : 'secondary'),
    };
  });
  const primaryTarget = archiveTargets.find((target) => target.role === 'primary');
  const primaryRef = archiveTargetRef(primaryTarget);
  const owner = options.owner || getDefaultOwner();

  ensureDir(proposalDir);

  const canonicalCorpus = [];
  appendCanonicalCorpusEntries(canonicalCorpus, options.canonicalCorpus || [], workspaceName);
  const archiveTypes = getCanonicalCorpusTypesForArchiveTargets(archiveTargets);
  appendCanonicalCorpusEntries(canonicalCorpus, scanWorkspaceDocsForCanonicalCorpus(workspaceName, cwd, archiveTypes), workspaceName);

  const meta = {
    schema_version: 'v2',
    id,
    slug,
    title,
    status: options.status || 'proposed',
    created_at: createdAt,
    updated_at: createdAt,
    workspace: workspaceName,
    scope: options.scope || 'workspace',
    source_explorations: sourceExplorations,
    canonical_corpus: canonicalCorpus,
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
    '- Backend: beads': `- Backend: ${taskBackend}`,
    '<Canonical corpus reviewed>': renderCanonicalCorpusList(canonicalCorpus, proposalDir, cwd),
    '<Knowledge impact>': renderKnowledgeImpact(canonicalCorpus),
    'module:example-module': primaryRef,
    '2026-05-20': createdAt.slice(0, 10),
  };

  writeText(path.join(proposalDir, 'meta.yaml'), `${serializeYaml(meta)}\n`);
  writeText(path.join(proposalDir, 'proposal.md'), renderTemplate(proposalTemplate, replacements));
  writeText(path.join(proposalDir, 'design.md'), renderTemplate(designTemplate, replacements));
  writeText(path.join(proposalDir, 'task-map.md'), renderTaskMapTemplate(taskMapTemplate, replacements, taskBackend));
  writeText(path.join(proposalDir, 'notes.md'), renderTemplate(notesTemplate, replacements));

  return {
    id,
    slug,
    proposalDir,
    meta,
  };
}

function writeProposalMeta(metaPath, meta) {
  writeText(metaPath, `${serializeYaml(meta)}\n`);
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
    '2026-05-20': todayDate(),
  });
  writeText(context.notesPath, content);
  return true;
}

function appendImplementationNote(context, note, cwd = process.cwd()) {
  ensureNotesFile(context, cwd);
  const timestamp = nowIso();
  const date = timestamp.slice(0, 10);
  const block = [
    '',
    `## ${date}`,
    '',
    `### ${timestamp}`,
    '',
    '#### Progress',
    '',
    `- ${note.trim()}`,
    '',
  ].join('\n');
  appendText(context.notesPath, block);
  context.meta.updated_at = nowIso();
  writeProposalMeta(context.metaPath, context.meta);
  return {
    notes_path: context.notesPath,
    timestamp,
  };
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

  const epicDescription = `Proposal ${context.meta.id}\n\n${context.proposalDir}`;
  const epicLabels = [
    `proposal:${context.meta.id}`,
    'workflow:proposal',
    `archive:${archiveTargetRef((context.meta.archive_targets || [])[0])}`,
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
      ...(task.workspace ? [`workspace:${task.workspace}`] : []),
      ...(task.code_scope ? [`code-scope:${task.code_scope}`] : []),
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

function closeBeadsEpic(epicId, proposalId, cwd = process.cwd()) {
  if (!epicId) return false;
  const command = runCommand('bd', ['close', epicId, '--reason', `Archived ${proposalId}`], cwd);
  if (!command.ok) {
    throw new Error(`bd close failed for ${epicId}: ${command.error}`);
  }
  return true;
}

function getArchiveMarker(proposalId) {
  return `${HISTORY_MARKER_PREFIX}${proposalId} -->`;
}

function appendArchiveBlock(filePath, proposalId, block) {
  const marker = getArchiveMarker(proposalId);
  if (fileExists(filePath) && readFileRequired(filePath).includes(marker)) {
    return false;
  }
  appendText(filePath, `\n${marker}\n${block}`);
  return true;
}

function ensureModuleArchiveTarget(targetPath, context, cwd = process.cwd()) {
  ensureDir(targetPath);

  const templateFiles = [
    ['README.md', path.join('modules', 'README.md')],
    ['design.md', path.join('modules', 'design.md')],
    ['api.md', path.join('modules', 'api.md')],
    ['history.md', path.join('modules', 'history.md')],
  ];

  for (const [fileName, templatePath] of templateFiles) {
    const absolutePath = path.join(targetPath, fileName);
    if (fileExists(absolutePath)) continue;
    let rendered;
    if (fileName === 'history.md') {
      rendered = `# ${path.basename(targetPath)} History\n`;
    } else {
      const template = loadTemplate(cwd, templatePath);
      rendered = renderTemplate(template, {
        '<Module Name>': path.basename(targetPath),
        '<proposal id>': context.meta.id,
        'CR26052001': context.meta.id,
        '2026-05-20': todayDate(),
        'What changed in the module': context.meta.title,
      });
    }
    writeText(absolutePath, rendered);
  }

  const historyPath = path.join(targetPath, 'history.md');
  const block = [
    `## ${todayDate()}`,
    '',
    `- Proposal: ${context.meta.id}`,
    `- Summary: ${context.meta.title}`,
    `- Source: ${path.relative(path.dirname(historyPath), context.proposalDir)}`,
    '',
  ].join('\n');
  appendArchiveBlock(historyPath, context.meta.id, block);

  const readmePath = path.join(targetPath, 'README.md');
  appendArchiveBlock(readmePath, context.meta.id, [
    '## Archived proposals',
    '',
    `- ${context.meta.id}: ${context.meta.title}`,
    '',
  ].join('\n'));

  return {
    type: 'module',
    path: targetPath,
  };
}

function ensureArchitectureArchiveTarget(targetPath, context, cwd = process.cwd()) {
  if (!fileExists(targetPath)) {
    const template = loadTemplate(cwd, path.join('architecture', 'system.md'));
    const rendered = renderTemplate(template, {
      '<System Topic>': path.basename(targetPath, path.extname(targetPath)),
      '<proposal id>': context.meta.id,
    });
    writeText(targetPath, rendered);
  }

  const block = [
    `## ${todayDate()} ${context.meta.id}`,
    '',
    `- Status: archived from proposal ${context.meta.id}`,
    `- Summary: ${context.meta.title}`,
    `- Source: ${path.relative(path.dirname(targetPath), context.proposalDir)}`,
    '',
    '### Required follow-through',
    '',
    '- Update the relevant system view and cross-cutting relationships.',
    '',
  ].join('\n');
  appendArchiveBlock(targetPath, context.meta.id, block);

  return {
    type: 'architecture',
    path: targetPath,
  };
}

function ensureDecisionArchiveTarget(targetPath, context, cwd = process.cwd()) {
  if (!fileExists(targetPath)) {
    const template = loadTemplate(cwd, path.join('decisions', 'ADR-template.md'));
    const rendered = renderTemplate(template, {
      'ADR-001: <Title>': `${path.basename(targetPath, path.extname(targetPath))}: ${context.meta.title}`,
      'CR26052001': context.meta.id,
      '2026-05-20': todayDate(),
      '<Title>': context.meta.title,
    });
    writeText(targetPath, rendered);
  }

  const block = [
    `## Update ${todayDate()}`,
    '',
    `- Proposal: ${context.meta.id}`,
    `- Summary: ${context.meta.title}`,
    `- Source: ${path.relative(path.dirname(targetPath), context.proposalDir)}`,
    '',
  ].join('\n');
  appendArchiveBlock(targetPath, context.meta.id, block);

  return {
    type: 'decision',
    path: targetPath,
  };
}

function ensureArchiveTarget(target, context, cwd = process.cwd()) {
  const absolutePath = (() => {
    if (target?.ref && target?.workspace) {
      return path.join(getWorkspaceDocsRoot(target.workspace, cwd), target.ref);
    }
    if (target?.path) {
      return path.isAbsolute(target.path) ? target.path : path.join(cwd, target.path);
    }
    if (target?.ref) {
      return path.isAbsolute(target.ref) ? target.ref : path.join(cwd, target.ref);
    }
    throw new Error('unsupported archive target location');
  })();

  if (target.type === 'module') {
    return ensureModuleArchiveTarget(absolutePath, context, cwd);
  }
  if (target.type === 'architecture') {
    return ensureArchitectureArchiveTarget(absolutePath, context, cwd);
  }
  if (target.type === 'decision') {
    return ensureDecisionArchiveTarget(absolutePath, context, cwd);
  }

  throw new Error(`unsupported archive target type: ${target.type}`);
}

function getArchiveReadiness(context, cwd = process.cwd()) {
  const validation = validateProposalContext(context, cwd);
  const failures = [...validation.errors];
  const warnings = [...validation.warnings];

  if (context.meta.status !== 'implemented') {
    failures.push(`proposal status must be implemented before archive, got ${context.meta.status}`);
  }

  let beadSummary = null;
  if (context.meta.task_backend === 'beads') {
    beadSummary = beadTaskSummary(context.meta.id, cwd);
    if (!beadSummary.available) {
      failures.push(`cannot verify Beads tasks: ${beadSummary.error}`);
    } else if (beadSummary.openWorkItems.length > 0) {
      failures.push(`proposal still has ${beadSummary.openWorkItems.length} open Beads work items`);
    }
  }

  const primaryTarget = (context.meta.archive_targets || []).find((target) => target.role === 'primary');
  if (!primaryTarget) {
    failures.push('proposal must define a primary archive target');
  }

  return {
    failures,
    warnings,
    beadSummary,
  };
}

function archiveProposal(context, cwd = process.cwd()) {
  const readiness = getArchiveReadiness(context, cwd);
  if (readiness.failures.length > 0) {
    const error = new Error(`archive readiness failed for ${context.meta.id}`);
    error.readiness = readiness;
    throw error;
  }

  const updatedTargets = [];
  for (const target of context.meta.archive_targets || []) {
    updatedTargets.push(ensureArchiveTarget(target, context, cwd));
  }

  let epicClosed = false;
  if (context.meta.task_backend === 'beads' && context.meta.task_epic_id) {
    epicClosed = closeBeadsEpic(context.meta.task_epic_id, context.meta.id, cwd);
  }

  transitionProposalStatus(context, 'archived');

  return {
    id: context.meta.id,
    updated_targets: updatedTargets,
    task_epic_closed: epicClosed,
    status: context.meta.status,
  };
}

function listProposalDirs(cwd = process.cwd()) {
  const resolvedWorkspace = resolveWorkspaceForCwd(cwd);
  const proposalsRoot = getProposalsRoot(cwd, resolvedWorkspace.name);
  if (!fileExists(proposalsRoot)) return [];

  return fs.readdirSync(proposalsRoot, { withFileTypes: true })
    .filter((entry) => entry.isDirectory())
    .map((entry) => path.join(proposalsRoot, entry.name))
    .sort();
}

function listProposalSummaries(cwd = process.cwd(), workspaceName = null) {
  const resolvedWorkspace = workspaceName ? getDocsWorkspace(workspaceName, cwd) : resolveWorkspaceForCwd(cwd);
  const proposalsRoot = getProposalsRoot(cwd, resolvedWorkspace.name);
  if (!fileExists(proposalsRoot)) return [];

  return fs.readdirSync(proposalsRoot, { withFileTypes: true })
    .filter((entry) => entry.isDirectory())
    .map((entry) => path.join(proposalsRoot, entry.name))
    .sort()
    .map((proposalDir) => {
      const context = loadProposalContext(proposalDir, cwd);
      return {
        kind: 'proposal',
        id: context.meta.id,
        title: context.meta.title,
        status: context.meta.status,
        task_backend: context.meta.task_backend,
        workspace: context.meta.workspace || resolvedWorkspace.name,
        proposal_dir: proposalDir,
        updated_at: context.meta.updated_at,
        archive_targets: context.meta.archive_targets || [],
      };
    });
}

function listExplorationSummaries(cwd = process.cwd(), workspaceName = null) {
  const resolvedWorkspace = workspaceName ? getDocsWorkspace(workspaceName, cwd) : resolveWorkspaceForCwd(cwd);
  const explorationsRoot = path.join(getWorkspaceDocsRoot(resolvedWorkspace.name, cwd), 'explorations');
  if (!fileExists(explorationsRoot)) return [];

  const summaries = [];
  for (const entry of fs.readdirSync(explorationsRoot, { withFileTypes: true })) {
    if (!entry.isDirectory()) continue;
    const explorationDir = path.join(explorationsRoot, entry.name);
    const indexPath = path.join(explorationDir, 'index.md');
    if (!fileExists(indexPath)) continue;

    const parsed = parseExplorationIndex(readFileRequired(indexPath));
    summaries.push({
      kind: 'exploration',
      title: parsed.title || entry.name,
      status: parsed.status || '进行中',
      workspace: resolvedWorkspace.name,
      exploration_dir: explorationDir,
      index_path: indexPath,
      date: parsed.date || null,
    });
  }

  return summaries.sort((a, b) => a.title.localeCompare(b.title, 'zh-Hans-CN'));
}

module.exports = {
  appendImplementationNote,
  archiveProposal,
  archiveTargetRef,
  beadTaskSummary,
  createProposalSkeleton,
  ensureArchiveTarget,
  ensureBeadsTasks,
  ensureNotesFile,
  fileExists,
  findNextProposalId,
  formatDateCode,
  getArchiveReadiness,
  getDocsWorkspace,
  getDocsRoot,
  getProjectRoot,
  getToolRoot,
  getWorkflowConfigPath,
  getWorkflowConfig,
  getWorkspaceDocsRoot,
  isProposalId,
  listProposalSummaries,
  listExplorationSummaries,
  listDocsWorkspaces,
  loadProposalContext,
  nowIso,
  parseCliArgs,
  parseSimpleYaml,
  parseTaskMap,
  resolveProposalDir,
  resolveWorkspaceForCwd,
  serializeYaml,
  slugify,
  todayDate,
  transitionProposalStatus,
  validateProposalContext,
  writeProposalMeta,
};
