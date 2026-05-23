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
const ARCHIVE_TARGET_TYPES = new Set(['module', 'architecture', 'decision', 'convention']);
const ARCHIVE_TARGET_ROLES = new Set(['primary', 'secondary']);
const TASK_PRIORITIES = new Set(['P0', 'P1', 'P2']);
const SIZE_CLASSES = new Set(['small', 'medium', 'large']);
const OWNERSHIP_TYPES = new Set(['module', 'system', 'cross-module', 'convention']);
const DESIGN_LAYOUTS = new Set(['single', 'split']);
const DOCUMENT_TYPES = new Set([
  'exploration',
  'proposal',
  'design',
  'model',
  'finding',
  'decision',
  'journal',
  'note',
  'task-map',
  'convention',
  'module',
  'architecture',
  'adr',
]);
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

const OWNERSHIP_TO_ARCHIVE_TYPE = {
  module: 'module',
  system: 'architecture',
  'cross-module': 'architecture',
  convention: 'convention',
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
  if (value === '[]') return [];
  if (value === '{}') return {};
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
  const parsed = parseFrontmatterDocument(text);
  const lines = text.replace(/\r\n/g, '\n').split('\n');
  const result = {
    proposal_id: parsed.frontmatter?.proposal_id || null,
    backend: parsed.frontmatter?.task_backend || null,
    tasks: [],
  };

  const sourceLines = parsed.frontmatter ? parsed.body.split('\n') : lines;

  let currentTask = null;
  let collectingCompletion = false;

  for (const line of sourceLines) {
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
        model_refs: [],
        convention_refs: [],
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
      else if (field === 'model refs' && value) currentTask.model_refs = value.split(',').map((item) => item.trim()).filter(Boolean);
      else if (field === 'convention refs' && value) currentTask.convention_refs = value.split(',').map((item) => item.trim()).filter(Boolean);
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
  const parsed = parseFrontmatterDocument(text);
  if (parsed.frontmatter) {
    const frontmatter = parsed.frontmatter;
    return {
      title: frontmatter.title || null,
      status: frontmatter.status || null,
      created: frontmatter.created || null,
      updated: frontmatter.updated || null,
      ownership: Array.isArray(frontmatter.ownership) ? frontmatter.ownership : [],
      reusable_rules: Array.isArray(frontmatter.reusable_rules) ? frontmatter.reusable_rules : [],
      expected_size_class: frontmatter.expected_size_class || null,
      question: frontmatter.question || null,
      exploration_slug: frontmatter.exploration_slug || null,
    };
  }

  const lines = text.replace(/\r\n/g, '\n').split('\n');
  const result = {
    title: null,
    status: null,
    created: null,
    updated: null,
    ownership: [],
    reusable_rules: [],
    expected_size_class: null,
  };

  for (const line of lines) {
    const trimmed = line.trim();
    if (!trimmed) continue;

    const titleMatch = trimmed.match(/^#\s+(.+)$/);
    if (titleMatch && !result.title) {
      result.title = titleMatch[1].trim();
      continue;
    }

    const statusMatch = trimmed.match(/^-\s*Status:\s*(.+)$/i);
    if (statusMatch && !result.status) {
      result.status = statusMatch[1].trim();
      continue;
    }

    const createdMatch = trimmed.match(/^-\s*Created:\s*(.+)$/i);
    if (createdMatch && !result.created) {
      result.created = createdMatch[1].trim();
      continue;
    }

    const updatedMatch = trimmed.match(/^-\s*Updated:\s*(.+)$/i);
    if (updatedMatch && !result.updated) {
      result.updated = updatedMatch[1].trim();
      continue;
    }

    const sizeMatch = trimmed.match(/^-\s*Expected size class for the resulting proposal:\s*(.+)$/i);
    if (sizeMatch && !result.expected_size_class) {
      result.expected_size_class = sizeMatch[1].trim();
      continue;
    }

    const primaryOwnershipMatch = trimmed.match(/^-\s*primary:\s*(.+)$/i);
    if (primaryOwnershipMatch) {
      try {
        result.ownership.push(parseOwnershipEntry(`${primaryOwnershipMatch[1].trim()}:primary`, 0));
      } catch (error) {
        result.ownership.push({
          type: null,
          target: primaryOwnershipMatch[1].trim(),
          role: 'primary',
        });
      }
      continue;
    }

    const secondaryOwnershipMatch = trimmed.match(/^-\s*secondary:\s*(.+)$/i);
    if (secondaryOwnershipMatch) {
      try {
        result.ownership.push(parseOwnershipEntry(`${secondaryOwnershipMatch[1].trim()}:secondary`, 1));
      } catch (error) {
        result.ownership.push({
          type: null,
          target: secondaryOwnershipMatch[1].trim(),
          role: 'secondary',
        });
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

function getProjectRulesRoot(cwd = process.cwd(), workspaceName = null) {
  return path.join(getProjectRoot(cwd), 'docs', 'flowforge', '_rules');
}

function getIntakeRoot(cwd = process.cwd(), workspaceName = null) {
  return path.join(getWorkspaceDocsRoot(workspaceName, cwd), 'intake');
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
    if (value.length === 0) {
      return `${prefix}[]`;
    }
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
        if (Array.isArray(nestedValue) && nestedValue.length === 0) {
          return `${prefix}${key}: []`;
        }
        if (nestedValue && typeof nestedValue === 'object' && !Array.isArray(nestedValue) && Object.keys(nestedValue).length === 0) {
          return `${prefix}${key}: {}`;
        }
        if (isPlainScalar(nestedValue)) {
          return `${prefix}${key}: ${yamlScalar(nestedValue)}`;
        }
        return `${prefix}${key}:\n${serializeYaml(nestedValue, indent + 2)}`;
      })
      .join('\n');
  }

  return `${prefix}${yamlScalar(value)}`;
}

function parseFrontmatterDocument(text) {
  const normalized = String(text || '').replace(/\r\n/g, '\n');
  const lines = normalized.split('\n');
  if (lines[0]?.trim() !== '---') {
    return {
      frontmatter: null,
      body: normalized,
    };
  }

  let closingIndex = -1;
  for (let index = 1; index < lines.length; index += 1) {
    if (lines[index].trim() === '---') {
      closingIndex = index;
      break;
    }
  }

  if (closingIndex === -1) {
    return {
      frontmatter: null,
      body: normalized,
    };
  }

  const frontmatterText = lines.slice(1, closingIndex).join('\n');
  const body = lines.slice(closingIndex + 1).join('\n').replace(/^\n+/, '');

  return {
    frontmatter: parseSimpleYaml(frontmatterText),
    body,
  };
}

function renderFrontmatter(frontmatter, body = '') {
  const normalizedBody = String(body || '').replace(/^\n+/, '');
  const renderedFrontmatter = serializeYaml(frontmatter);
  return normalizedBody
    ? `---\n${renderedFrontmatter}\n---\n\n${normalizedBody}`
    : `---\n${renderedFrontmatter}\n---`;
}

function workspaceDocRef(workspaceName, ref) {
  return `${workspaceName}:${ref}`;
}

function toPosixPath(filePath) {
  return String(filePath || '').split(path.sep).join('/');
}

function normalizeRefList(values) {
  return Array.from(new Set((values || [])
    .filter((value) => value !== null && value !== undefined)
    .map((value) => String(value).trim())
    .filter(Boolean)));
}

function buildDocumentFrontmatter({
  doc_type,
  title,
  status,
  workspace,
  module_scope = [],
  system_scope = [],
  convention_scope = [],
  ownership = [],
  information_class,
  topics = [],
  related_docs = [],
  archive_target = 'none',
  created,
  updated,
  ...rest
}) {
  const frontmatter = {
    doc_type,
    title,
    status,
    workspace,
    module_scope: normalizeRefList(module_scope),
    system_scope: normalizeRefList(system_scope),
    convention_scope: normalizeRefList(convention_scope),
    ownership,
    information_class,
    topics: normalizeRefList(topics),
    related_docs: normalizeRefList(related_docs),
    archive_target,
    created,
    updated,
    ...rest,
  };

  for (const key of Object.keys(frontmatter)) {
    if (frontmatter[key] === undefined) {
      delete frontmatter[key];
    }
  }

  return frontmatter;
}

function renderDocumentWithFrontmatter(frontmatter, body) {
  const normalizedBody = String(body || '').replace(/^\n+/, '');
  return `${renderFrontmatter(frontmatter, normalizedBody)}\n`;
}

function renderDocumentTemplate(templateText, frontmatter, replacements) {
  return renderDocumentWithFrontmatter(frontmatter, renderTemplate(templateText, replacements));
}

function updateDocumentFrontmatter(filePath, patch = {}) {
  if (!fileExists(filePath)) return false;
  const { frontmatter, body } = parseFrontmatterDocument(readFileRequired(filePath));
  if (!frontmatter) return false;
  writeText(filePath, renderDocumentWithFrontmatter({ ...frontmatter, ...patch }, body));
  return true;
}

function validateDocumentFrontmatter(filePath, expectedDocType, errors, warnings, options = {}) {
  if (!fileExists(filePath)) {
    errors.push(`${path.basename(filePath)} not found`);
    return null;
  }

  const label = options.label || path.basename(filePath);
  const { frontmatter, body } = parseFrontmatterDocument(readFileRequired(filePath));
  if (!frontmatter) {
    errors.push(`${label} missing YAML frontmatter`);
    return null;
  }

  for (const field of ['doc_type', 'title', 'status', 'workspace', 'module_scope', 'system_scope', 'convention_scope', 'ownership', 'information_class', 'topics', 'related_docs', 'archive_target', 'created', 'updated']) {
    if (frontmatter[field] === undefined || frontmatter[field] === null || frontmatter[field] === '') {
      errors.push(`${label} missing ${field}`);
    }
  }

  if (expectedDocType && frontmatter.doc_type !== expectedDocType) {
    errors.push(`${label} doc_type must be ${expectedDocType}, got ${frontmatter.doc_type}`);
  }

  for (const field of ['module_scope', 'system_scope', 'convention_scope', 'topics', 'related_docs']) {
    if (frontmatter[field] !== undefined && !Array.isArray(frontmatter[field])) {
      errors.push(`${label} ${field} must be an array`);
    }
  }

  if (frontmatter.ownership !== undefined && !Array.isArray(frontmatter.ownership)) {
    errors.push(`${label} ownership must be an array`);
  }

  if (frontmatter.archive_target === undefined) {
    errors.push(`${label} missing archive_target`);
  }

  if (options.allowedStatuses && frontmatter.status && !options.allowedStatuses.has(frontmatter.status)) {
    errors.push(`${label} has invalid status ${frontmatter.status}`);
  }

  if (options.requireBodyHeading && body && !body.includes(options.requireBodyHeading)) {
    warnings.push(`${label} should include body heading: ${options.requireBodyHeading}`);
  }

  return frontmatter;
}

function frontmatterDocRef(workspaceName, relativePath) {
  if (!relativePath) return 'none';
  return workspaceDocRef(workspaceName, relativePath);
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
  const meta = parseSimpleYaml(readFileRequired(metaPath));
  const taskMapPath = meta?.links?.task_map
    ? path.join(proposalDir, meta.links.task_map)
    : path.join(proposalDir, 'task-map.md');
  const designPath = meta?.links?.design
    ? path.join(proposalDir, meta.links.design)
    : path.join(proposalDir, 'design.md');
  const modelPath = meta?.links?.model
    ? path.join(proposalDir, meta.links.model)
    : path.join(proposalDir, 'model', 'README.md');
  const notesPath = meta?.links?.notes
    ? path.join(proposalDir, meta.links.notes)
    : path.join(proposalDir, 'notes.md');
  const taskMap = parseTaskMap(readFileRequired(taskMapPath));

  return {
    proposalDir,
    metaPath,
    taskMapPath,
    designPath,
    modelPath,
    notesPath,
    meta,
    taskMap,
  };
}


function isMilestoneTaskId(taskId) {
  return /^MILESTONE-/i.test(String(taskId || ''));
}

function isTemplatePlaceholder(value) {
  if (value === null || value === undefined) return false;
  const text = String(value).trim();
  if (!text) return true;
  return text.startsWith('<') && text.endsWith('>');
}

function validateProposalContext(context, cwd = process.cwd()) {
  const errors = [];
  const warnings = [];
  const { meta, taskMap, proposalDir, taskMapPath, designPath, modelPath, notesPath } = context;

  const proposalRoot = path.dirname(proposalDir);

  const schemaVersion = meta.schema_version || 'v1';
  const workspaceConfig = getWorkflowConfig(cwd).docs || {};
  const workspaceNames = new Set(Object.keys(workspaceConfig.workspaces || {}));
  const designLayout = getProposalDesignLayout(proposalDir, meta);
  const modelLayout = getProposalModelLayout(proposalDir, meta);

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
      'size_class',
      'ownership',
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
    if (meta.size_class && !SIZE_CLASSES.has(meta.size_class)) {
      errors.push(`invalid proposal size_class: ${meta.size_class}`);
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

    if (!Array.isArray(meta.ownership) || meta.ownership.length === 0) {
      errors.push('meta.yaml must define at least one ownership entry');
    } else {
      const primaryOwnership = meta.ownership.filter((entry) => entry?.role === 'primary');
      if (primaryOwnership.length !== 1) {
        errors.push(`meta.yaml must define exactly one primary ownership entry, found ${primaryOwnership.length}`);
      }

      const ownershipTargets = new Set();
      for (const entry of meta.ownership) {
        if (!entry?.type) {
          errors.push('ownership type is required');
          continue;
        }
        if (!OWNERSHIP_TYPES.has(entry.type)) {
          errors.push(`invalid ownership type: ${entry.type}`);
        }
        if (!entry?.target) {
          errors.push('ownership target is required');
        }
        if (!ARCHIVE_TARGET_ROLES.has(entry?.role)) {
          errors.push(`invalid ownership role: ${entry?.role}`);
        }
        if (entry?.workspace && workspaceNames.size > 0 && !workspaceNames.has(entry.workspace)) {
          errors.push(`unknown ownership workspace: ${entry.workspace}`);
        }
        if (entry?.target) {
          const key = `${entry.type}:${entry.target}`;
          if (ownershipTargets.has(key)) {
            warnings.push(`duplicate ownership entry detected: ${key}`);
          }
          ownershipTargets.add(key);
        }
      }
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

    if (meta.size_class === 'small') {
      if (!fileExists(path.join(proposalDir, 'design.md'))) {
        errors.push(`design.md not found for small proposal: ${path.join(proposalDir, 'design.md')}`);
      }
      if (fileExists(path.join(proposalDir, 'design'))) {
        warnings.push('small proposal contains a design/ directory; consider using the single-file design.md layout');
      }
    }

    if (meta.size_class === 'medium') {
      if (!fileExists(path.join(proposalDir, 'design.md')) && !fileExists(path.join(proposalDir, 'design', 'README.md'))) {
        errors.push('medium proposal must provide either design.md or design/README.md');
      }
    }

    if (meta.size_class === 'large') {
      if (!fileExists(path.join(proposalDir, 'design', 'README.md'))) {
        errors.push('large proposal must provide design/README.md');
      }
      if (fileExists(path.join(proposalDir, 'design.md'))) {
        errors.push('large proposal must not use root design.md');
      }
      if (!fileExists(path.join(proposalDir, 'model', 'README.md'))) {
        errors.push('large proposal must provide model/README.md');
      }
    }

    if (designLayout === 'split' && !meta.links.design.includes('design/README.md')) {
      warnings.push('proposal contains a split design/ layout but meta.links.design does not point to design/README.md');
    }

    if (modelLayout === 'split' && !meta.links.model) {
      errors.push('proposal contains a model/ directory but meta.yaml links.model is missing');
    }

    if (modelLayout === 'split' && meta.size_class === 'small') {
      errors.push('small proposals should not use a split model/ layout');
    }

    const primaryOwnership = primaryOwnershipEntries(meta.ownership);
    const archiveTargetsByTypeRef = new Set(
      (meta.archive_targets || []).map((target) => `${target.type}:${target.ref || target.path}`)
    );
    for (const ownership of meta.ownership || []) {
      const expectedType = OWNERSHIP_TO_ARCHIVE_TYPE[ownership.type];
      if (!expectedType) continue;
      const ownershipKey = `${expectedType}:${ownership.target}`;
      if (!archiveTargetsByTypeRef.has(ownershipKey)) {
        errors.push(`ownership entry has no matching archive target: ${ownership.type}:${ownership.target}`);
      }
    }
    if (primaryOwnership.length === 1) {
      const primary = primaryOwnership[0];
      const expectedType = OWNERSHIP_TO_ARCHIVE_TYPE[primary.type];
      if (expectedType) {
        const primaryKey = `${expectedType}:${primary.target}`;
        const primaryArchiveTarget = meta.archive_targets.find((target) => target.role === 'primary');
        if (!primaryArchiveTarget || `${primaryArchiveTarget.type}:${primaryArchiveTarget.ref || primaryArchiveTarget.path}` !== primaryKey) {
          errors.push('primary ownership must correspond to the primary archive target');
        }
      }
    }

    const proposalFrontmatter = validateDocumentFrontmatter(path.join(proposalDir, 'proposal.md'), 'proposal', errors, warnings, {
      label: 'proposal.md',
      allowedStatuses: STATUS_VALUES,
    });
    if (proposalFrontmatter) {
      if (proposalFrontmatter.proposal_id && proposalFrontmatter.proposal_id !== meta.id) {
        errors.push(`proposal.md proposal_id must match meta.yaml id: ${proposalFrontmatter.proposal_id}`);
      }
      if (proposalFrontmatter.size_class && proposalFrontmatter.size_class !== meta.size_class) {
        errors.push(`proposal.md size_class must match meta.yaml size_class: ${proposalFrontmatter.size_class}`);
      }
      const primaryArchiveTarget = meta.archive_targets.find((target) => target.role === 'primary');
      if (primaryArchiveTarget && proposalFrontmatter.archive_target && proposalFrontmatter.archive_target !== workspaceDocRef(primaryArchiveTarget.workspace, primaryArchiveTarget.ref)) {
        warnings.push(`proposal.md archive_target does not match the primary archive target: ${proposalFrontmatter.archive_target}`);
      }
      if (proposalFrontmatter.ownership_primary && primaryOwnership[0]) {
        const expectedPrimary = `${primaryOwnership[0].type}:${primaryOwnership[0].target}`;
        if (proposalFrontmatter.ownership_primary !== expectedPrimary) {
          warnings.push(`proposal.md ownership_primary does not match meta.yaml primary ownership: ${proposalFrontmatter.ownership_primary}`);
        }
      }
    }

    const taskMapFrontmatter = validateDocumentFrontmatter(taskMapPath, 'task-map', errors, warnings, {
      label: 'task-map.md',
      allowedStatuses: STATUS_VALUES,
    });
    if (taskMapFrontmatter && taskMapFrontmatter.proposal_id && taskMapFrontmatter.proposal_id !== meta.id) {
      errors.push(`task-map.md proposal_id must match meta.yaml id: ${taskMapFrontmatter.proposal_id}`);
    }
    if (taskMapFrontmatter && taskMapFrontmatter.task_backend && taskMapFrontmatter.task_backend !== meta.task_backend) {
      errors.push(`task-map.md task_backend must match meta.yaml task_backend: ${taskMapFrontmatter.task_backend}`);
    }

    const notesFrontmatter = validateDocumentFrontmatter(notesPath, 'note', errors, warnings, {
      label: 'notes.md',
      allowedStatuses: STATUS_VALUES,
    });
    if (notesFrontmatter && notesFrontmatter.proposal_id && notesFrontmatter.proposal_id !== meta.id) {
      errors.push(`notes.md proposal_id must match meta.yaml id: ${notesFrontmatter.proposal_id}`);
    }

    if (designLayout === 'single') {
      const designDoc = validateDocumentFrontmatter(path.join(proposalDir, 'design.md'), 'design', errors, warnings, {
        label: 'design.md',
        allowedStatuses: STATUS_VALUES,
      });
      if (designDoc && designDoc.proposal_id && designDoc.proposal_id !== meta.id) {
        errors.push(`design.md proposal_id must match meta.yaml id: ${designDoc.proposal_id}`);
      }
    } else {
      const designReadme = validateDocumentFrontmatter(path.join(proposalDir, 'design', 'README.md'), 'design', errors, warnings, {
        label: 'design/README.md',
        allowedStatuses: STATUS_VALUES,
      });
      if (designReadme && designReadme.proposal_id && designReadme.proposal_id !== meta.id) {
        errors.push(`design/README.md proposal_id must match meta.yaml id: ${designReadme.proposal_id}`);
      }
      for (const filePath of listMarkdownFiles(path.join(proposalDir, 'design'), { recursive: false })) {
        if (path.basename(filePath) === 'README.md') continue;
        const sectionDoc = validateDocumentFrontmatter(filePath, 'design', errors, warnings, {
          label: path.relative(proposalDir, filePath).split(path.sep).join('/'),
          allowedStatuses: STATUS_VALUES,
        });
        if (sectionDoc && sectionDoc.proposal_id && sectionDoc.proposal_id !== meta.id) {
          errors.push(`${path.relative(proposalDir, filePath).split(path.sep).join('/')} proposal_id must match meta.yaml id: ${sectionDoc.proposal_id}`);
        }
      }
    }

    if (modelLayout === 'split') {
      const modelReadme = validateDocumentFrontmatter(path.join(proposalDir, 'model', 'README.md'), 'model', errors, warnings, {
        label: 'model/README.md',
        allowedStatuses: STATUS_VALUES,
      });
      if (modelReadme && modelReadme.proposal_id && modelReadme.proposal_id !== meta.id) {
        errors.push(`model/README.md proposal_id must match meta.yaml id: ${modelReadme.proposal_id}`);
      }
      for (const filePath of listMarkdownFiles(path.join(proposalDir, 'model'), { recursive: true, ignoreDirs: ['parts'] })) {
        if (path.basename(filePath) === 'README.md') continue;
        const modelDoc = validateDocumentFrontmatter(filePath, 'model', errors, warnings, {
          label: path.relative(proposalDir, filePath).split(path.sep).join('/'),
          allowedStatuses: STATUS_VALUES,
        });
        if (modelDoc && modelDoc.proposal_id && modelDoc.proposal_id !== meta.id) {
          errors.push(`${path.relative(proposalDir, filePath).split(path.sep).join('/')} proposal_id must match meta.yaml id: ${modelDoc.proposal_id}`);
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
    errors.push(`design entry not found: ${designPath}`);
  }
  if (meta?.links?.model && !fileExists(modelPath)) {
    errors.push(`model README not found: ${modelPath}`);
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
    if (!isMilestoneTaskId(task.task_id) && (!Array.isArray(task.capability_refs) || task.capability_refs.length === 0)) {
      errors.push(`${task.task_id} must reference at least one capability`);
    }
    if (!Array.isArray(task.completion_definition) || task.completion_definition.length === 0) {
      errors.push(`${task.task_id} must define completion criteria`);
    }
  }

  for (const task of taskMap.tasks) {
    for (const dependency of task.depends_on || []) {
      if (isTemplatePlaceholder(dependency)) continue;
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
      if (isTemplatePlaceholder(ref)) continue;
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

function parseOwnershipEntry(raw, index) {
  const parts = String(raw).split(':');
  if (parts.length < 2) {
    throw new Error(`invalid --ownership value at item ${index + 1}: ${raw}`);
  }

  const [type, target, maybeRole] = parts;
  if (!type || !target) {
    throw new Error(`invalid --ownership value at item ${index + 1}: ${raw}`);
  }

  const role = maybeRole === 'primary' || maybeRole === 'secondary'
    ? maybeRole
    : index === 0
      ? 'primary'
      : 'secondary';

  return {
    type,
    target,
    role,
  };
}

function archiveTargetRef(target) {
  if (target?.key) return target.key;
  if (target?.ref) return `${target.type}:${path.basename(target.ref)}`;
  return `${target?.type || 'unknown'}:${path.basename(target?.path || 'unknown')}`;
}

function normalizeOwnershipEntry(entry, workspaceName) {
  if (!entry) return null;
  const type = entry.type;
  const target = entry.target || entry.ref || entry.path;
  if (!type || !target) return null;
  return {
    type,
    target,
    role: entry.role || 'secondary',
    workspace: entry.workspace || workspaceName,
  };
}

function appendOwnershipEntries(target, entries, workspaceName) {
  const seen = new Set(target.map((entry) => `${entry.type}:${entry.target}:${entry.role}`));
  for (const rawEntry of entries || []) {
    const entry = normalizeOwnershipEntry(rawEntry, workspaceName);
    if (!entry) continue;
    const key = `${entry.type}:${entry.target}:${entry.role}`;
    if (seen.has(key)) continue;
    seen.add(key);
    target.push(entry);
  }
  return target;
}

function inferOwnershipFromArchiveTargets(archiveTargets) {
  return (archiveTargets || [])
    .map((target) => {
      const type = OWNERSHIP_TO_ARCHIVE_TYPE[target?.type];
      if (!type) return null;
      return {
        type,
        target: target.ref || target.path || target.key,
        role: target.role || 'secondary',
      };
    })
    .filter(Boolean);
}

function primaryOwnershipEntries(ownership) {
  return (ownership || []).filter((entry) => entry?.role === 'primary');
}

function getProposalDesignLayout(proposalDir, meta = {}) {
  const splitDesign = fileExists(path.join(proposalDir, 'design', 'README.md'));
  const singleDesign = fileExists(path.join(proposalDir, 'design.md'));
  if (splitDesign) return 'split';
  if (singleDesign) return 'single';
  const declaredDesign = meta?.links?.design;
  if (declaredDesign && declaredDesign.includes('design/README.md')) return 'split';
  return 'single';
}

function getProposalModelLayout(proposalDir, meta = {}) {
  const splitModel = fileExists(path.join(proposalDir, 'model', 'README.md'));
  if (splitModel) return 'split';
  const declaredModel = meta?.links?.model;
  if (declaredModel && declaredModel.includes('model/README.md')) return 'split';
  return 'none';
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

  const conventionsRoot = path.join(workspaceRoot, 'conventions');
  if (typeAllowed('convention') && fileExists(conventionsRoot)) {
    for (const entry of fs.readdirSync(conventionsRoot, { withFileTypes: true })) {
      if (!entry.isFile() || !entry.name.endsWith('.md')) continue;
      corpus.push({
        workspace: workspaceName,
        ref: path.join('conventions', entry.name),
        type: 'convention',
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

function renderOwnershipTargetList(ownership, types) {
  const matches = (ownership || []).filter((entry) => types.includes(entry.type));
  if (matches.length === 0) {
    return 'none';
  }

  return matches
    .map((entry) => `${entry.target}${entry.role === 'primary' ? ' (primary)' : ''}`)
    .join(', ');
}

function renderScopeListFromOwnership(ownership, types) {
  return normalizeRefList((ownership || [])
    .filter((entry) => types.includes(entry.type))
    .map((entry) => entry.target));
}

function ownershipEntryToFrontmatter(entry) {
  if (!entry) return null;
  return {
    type: entry.type,
    target: entry.target,
    role: entry.role || 'secondary',
  };
}

function normalizeOwnershipForFrontmatter(ownership) {
  return (ownership || [])
    .map((entry) => ownershipEntryToFrontmatter(entry))
    .filter(Boolean);
}

function ownershipEntriesForTypes(ownership, types) {
  return (ownership || []).filter((entry) => types.includes(entry.type));
}

function buildArchiveDocumentOwnership(targetType, context, targetRef) {
  const proposalOwnership = context.meta.ownership || [];
  const target = targetRef;
  if (targetType === 'module') {
    return [{ type: 'module', target, role: 'primary' }];
  }

  if (targetType === 'architecture') {
    const preferred = primaryOwnershipEntries(ownershipEntriesForTypes(proposalOwnership, ['system', 'cross-module']))[0];
    return [preferred ? { type: preferred.type, target: preferred.target, role: 'primary' } : { type: 'system', target, role: 'primary' }];
  }

  if (targetType === 'convention') {
    return [{ type: 'convention', target, role: 'primary' }];
  }

  if (targetType === 'decision') {
    const preferred = primaryOwnershipEntries(proposalOwnership)[0];
    return [preferred ? { type: preferred.type, target: preferred.target, role: 'primary' } : { type: 'system', target, role: 'primary' }];
  }

  return normalizeOwnershipForFrontmatter(proposalOwnership);
}

function makeDocFrontmatter({
  docType,
  title,
  status,
  workspace,
  fileRef,
  ownership = [],
  informationClass,
  topics = [],
  relatedDocs = [],
  archiveTarget = null,
  extra = {},
  created,
  updated,
}) {
  const timestamp = nowIso();
  const moduleScope = renderScopeListFromOwnership(ownership, ['module']);
  const systemScope = renderScopeListFromOwnership(ownership, ['system', 'cross-module']);
  const conventionScope = renderScopeListFromOwnership(ownership, ['convention']);

  return buildDocumentFrontmatter({
    doc_type: docType,
    title,
    status,
    workspace,
    module_scope: moduleScope,
    system_scope: systemScope,
    convention_scope: conventionScope,
    ownership: normalizeOwnershipForFrontmatter(ownership),
    information_class: informationClass || docType,
    topics,
    related_docs: relatedDocs,
    archive_target: archiveTarget || workspaceDocRef(workspace, fileRef),
    created: created || timestamp,
    updated: updated || timestamp,
    ...extra,
  });
}

function listMarkdownFiles(dirPath, options = {}) {
  if (!fileExists(dirPath)) return [];
  const recursive = options.recursive || false;
  const ignoreDirs = new Set(options.ignoreDirs || []);
  const results = [];

  for (const entry of fs.readdirSync(dirPath, { withFileTypes: true })) {
    if (entry.isDirectory()) {
      if (!recursive || ignoreDirs.has(entry.name)) continue;
      results.push(...listMarkdownFiles(path.join(dirPath, entry.name), options));
      continue;
    }

    if (!entry.isFile() || !entry.name.endsWith('.md')) continue;
    results.push(path.join(dirPath, entry.name));
  }

  return results;
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
  const sizeClass = options.sizeClass || 'small';
  if (!SIZE_CLASSES.has(sizeClass)) {
    throw new Error(`invalid size class: ${sizeClass}`);
  }

  const designLayout = options.designLayout || (sizeClass === 'large' ? 'split' : 'single');
  if (!DESIGN_LAYOUTS.has(designLayout)) {
    throw new Error(`invalid design layout: ${designLayout}`);
  }
  const modelLayout = designLayout === 'split' ? 'split' : 'none';
  if (sizeClass === 'small' && designLayout !== 'single') {
    throw new Error('small proposals must use the single-file design layout');
  }
  if (sizeClass === 'large' && designLayout !== 'split') {
    throw new Error('large proposals must use the split design layout');
  }

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
  const ownership = (options.ownership && options.ownership.length > 0)
    ? options.ownership
    : inferOwnershipFromArchiveTargets(archiveTargets);
  const primaryOwnership = primaryOwnershipEntries(ownership)[0] || null;
  const reusableRules = options.reusableRules || [];
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
    size_class: sizeClass,
    ownership,
    source_explorations: sourceExplorations,
    canonical_corpus: canonicalCorpus,
    owner,
    task_backend: taskBackend,
    task_epic_id: null,
    archive_targets: archiveTargets,
    tags: options.tags || [],
    links: {
      design: designLayout === 'split' ? 'design/README.md' : 'design.md',
      ...(designLayout === 'split' ? { model: 'model/README.md' } : {}),
      task_map: 'task-map.md',
      notes: 'notes.md',
    },
  };

  const proposalTemplate = loadTemplate(cwd, path.join('proposals', 'proposal.md'));
  const taskMapTemplate = loadTemplate(cwd, path.join('proposals', 'task-map.md'));
  const notesTemplate = loadTemplate(cwd, path.join('proposals', 'notes.md'));
  const singleDesignTemplate = loadTemplate(cwd, path.join('proposals', 'design.md'));
  const designReadmeTemplate = loadTemplate(cwd, path.join('proposals', 'design', 'README.md'));
  const designSectionTemplates = {
    'architecture.md': loadTemplate(cwd, path.join('proposals', 'design', 'architecture.md')),
    'model.md': loadTemplate(cwd, path.join('proposals', 'design', 'model.md')),
    'lifecycle.md': loadTemplate(cwd, path.join('proposals', 'design', 'lifecycle.md')),
    'flow.md': loadTemplate(cwd, path.join('proposals', 'design', 'flow.md')),
    'api.md': loadTemplate(cwd, path.join('proposals', 'design', 'api.md')),
    'constraints.md': loadTemplate(cwd, path.join('proposals', 'design', 'constraints.md')),
    'tradeoffs.md': loadTemplate(cwd, path.join('proposals', 'design', 'tradeoffs.md')),
  };
  const modelReadmeTemplate = loadTemplate(cwd, path.join('proposals', 'model', 'README.md'));
  const proposalFileRef = toPosixPath(path.relative(docsRoot, path.join(proposalDir, 'proposal.md')));
  const taskMapFileRef = toPosixPath(path.relative(docsRoot, path.join(proposalDir, 'task-map.md')));
  const notesFileRef = toPosixPath(path.relative(docsRoot, path.join(proposalDir, 'notes.md')));
  const designReadmeRef = toPosixPath(path.relative(docsRoot, path.join(proposalDir, designLayout === 'split' ? 'design/README.md' : 'design.md')));
  const modelReadmeRef = toPosixPath(path.relative(docsRoot, path.join(proposalDir, 'model/README.md')));
  const relatedDocRefs = normalizeRefList([
    ...sourceExplorations.map((source) => workspaceDocRef(source.workspace || workspaceName, source.ref)),
    ...canonicalCorpus.map((entry) => workspaceDocRef(entry.workspace, entry.ref)),
  ]);
  const proposalFrontmatter = makeDocFrontmatter({
    docType: 'proposal',
    title,
    status: meta.status,
    workspace: workspaceName,
    fileRef: proposalFileRef,
    ownership,
    informationClass: 'proposal',
    topics: options.tags || [],
    relatedDocs: relatedDocRefs,
    archiveTarget: primaryTarget ? workspaceDocRef(primaryTarget.workspace, primaryTarget.ref) : proposalFileRef,
    extra: {
      proposal_id: id,
      size_class: sizeClass,
      ownership_primary: primaryOwnership ? `${primaryOwnership.type}:${primaryOwnership.target}` : undefined,
      design_layout: designLayout,
    },
    created: createdAt,
    updated: createdAt,
  });
  const taskMapFrontmatter = makeDocFrontmatter({
    docType: 'task-map',
    title: `${title} Task Map`,
    status: meta.status,
    workspace: workspaceName,
    fileRef: taskMapFileRef,
    ownership,
    informationClass: 'task-map',
    topics: options.tags || [],
    relatedDocs: [proposalFileRef],
    archiveTarget: taskMapFileRef,
    extra: {
      proposal_id: id,
      task_backend: taskBackend,
    },
    created: createdAt,
    updated: createdAt,
  });
  const notesFrontmatter = makeDocFrontmatter({
    docType: 'note',
    title: `${title} Notes`,
    status: meta.status,
    workspace: workspaceName,
    fileRef: notesFileRef,
    ownership,
    informationClass: 'note',
    topics: options.tags || [],
    relatedDocs: [proposalFileRef],
    archiveTarget: notesFileRef,
    extra: {
      proposal_id: id,
      note_kind: 'progress',
    },
    created: createdAt,
    updated: createdAt,
  });
  const designFrontmatter = makeDocFrontmatter({
    docType: 'design',
    title: `${title} Design`,
    status: meta.status,
    workspace: workspaceName,
    fileRef: designReadmeRef,
    ownership,
    informationClass: 'design',
    topics: options.tags || [],
    relatedDocs: [proposalFileRef, modelLayout === 'split' ? modelReadmeRef : null],
    archiveTarget: primaryTarget ? workspaceDocRef(primaryTarget.workspace, primaryTarget.ref) : designReadmeRef,
    extra: {
      proposal_id: id,
      design_section: 'entry',
      canonical_entry_point: proposalFileRef,
    },
    created: createdAt,
    updated: createdAt,
  });
  const modelReadmeFrontmatter = makeDocFrontmatter({
    docType: 'model',
    title: `${title} Models`,
    status: meta.status,
    workspace: workspaceName,
    fileRef: modelReadmeRef,
    ownership,
    informationClass: 'model',
    topics: options.tags || [],
    relatedDocs: [proposalFileRef, designReadmeRef],
    archiveTarget: primaryTarget ? workspaceDocRef(primaryTarget.workspace, primaryTarget.ref) : modelReadmeRef,
    extra: {
      proposal_id: id,
      model_name: 'index',
      model_role: 'shared',
      data_scope: 'derived',
      model_status_in_proposal: 'retained',
    },
    created: createdAt,
    updated: createdAt,
  });

  const replacements = {
    '<Proposal Title>': title,
    '<Size class>': sizeClass,
    '<Primary ownership>': primaryOwnership ? `${primaryOwnership.type}:${primaryOwnership.target}` : '<Primary ownership>',
    '<Secondary ownership>': ownership.filter((entry) => entry.role === 'secondary').map((entry) => `${entry.type}:${entry.target}`).join(', ') || 'none',
    '<Owning modules>': renderOwnershipTargetList(ownership, ['module']),
    '<Owning systems>': renderOwnershipTargetList(ownership, ['system', 'cross-module']),
    '<Owning conventions>': renderOwnershipTargetList(ownership, ['convention']),
    '<Promotes reusable rules>': reusableRules.length > 0 ? 'yes' : 'no',
    '<Document layout>': designLayout === 'split'
      ? 'design/README.md plus model/README.md and supporting section files'
      : 'design.md',
    '<Primary archive target key>': primaryTarget?.key || 'primary-archive-target',
    '<Reusable rules block>': reusableRules.length > 0
      ? reusableRules.map((rule) => `- ${rule.title}: ${rule.summary || ''}`.trimEnd()).join('\n')
      : '- none',
    'CR20260520': id,
    'CR26052001': id,
    '- Backend: beads': `- Backend: ${taskBackend}`,
    '<Canonical corpus reviewed>': renderCanonicalCorpusList(canonicalCorpus, proposalDir, cwd),
    '<Knowledge impact>': renderKnowledgeImpact(canonicalCorpus),
    'module:example-module': primaryRef,
    '2026-05-20': createdAt.slice(0, 10),
  };

  writeText(path.join(proposalDir, 'meta.yaml'), `${serializeYaml(meta)}\n`);
  writeText(path.join(proposalDir, 'proposal.md'), renderDocumentTemplate(proposalTemplate, proposalFrontmatter, replacements));
  if (designLayout === 'single') {
    writeText(path.join(proposalDir, 'design.md'), renderDocumentTemplate(singleDesignTemplate, designFrontmatter, replacements));
  } else {
    const designDir = path.join(proposalDir, 'design');
    const modelDir = path.join(proposalDir, 'model');
    ensureDir(designDir);
    ensureDir(modelDir);
    writeText(path.join(designDir, 'README.md'), renderDocumentTemplate(designReadmeTemplate, designFrontmatter, replacements));
    for (const [fileName, templateText] of Object.entries(designSectionTemplates)) {
      const sectionFrontmatter = makeDocFrontmatter({
        docType: 'design',
        title: `${title} ${fileName.replace(/\.md$/, '')}`,
        status: meta.status,
        workspace: workspaceName,
        fileRef: toPosixPath(path.relative(docsRoot, path.join(designDir, fileName))),
        ownership,
        informationClass: 'design',
        topics: options.tags || [],
        relatedDocs: [proposalFileRef, designReadmeRef],
        archiveTarget: primaryTarget ? workspaceDocRef(primaryTarget.workspace, primaryTarget.ref) : toPosixPath(path.relative(docsRoot, path.join(designDir, fileName))),
        extra: {
          proposal_id: id,
          design_section: fileName.replace(/\.md$/, ''),
          canonical_entry_point: designReadmeRef,
        },
        created: createdAt,
        updated: createdAt,
      });
      writeText(path.join(designDir, fileName), renderDocumentTemplate(templateText, sectionFrontmatter, replacements));
    }
    writeText(path.join(modelDir, 'README.md'), renderDocumentTemplate(modelReadmeTemplate, modelReadmeFrontmatter, replacements));
  }
  writeText(path.join(proposalDir, 'task-map.md'), renderDocumentTemplate(taskMapTemplate, taskMapFrontmatter, replacements));
  writeText(path.join(proposalDir, 'notes.md'), renderDocumentTemplate(notesTemplate, notesFrontmatter, replacements));

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
  const body = renderTemplate(template, {
    '<Proposal Title>': context.meta.title,
    '2026-05-20': todayDate(),
  });
  const relativeRef = toPosixPath(path.relative(getDocsRoot(cwd, context.meta.workspace), context.notesPath));
  const frontmatter = makeDocFrontmatter({
    docType: 'note',
    title: `${context.meta.title} Notes`,
    status: context.meta.status,
    workspace: context.meta.workspace,
    fileRef: relativeRef,
    ownership: context.meta.ownership || [],
    informationClass: 'note',
    topics: context.meta.tags || [],
    relatedDocs: [workspaceDocRef(context.meta.workspace, toPosixPath(path.relative(getDocsRoot(cwd, context.meta.workspace), context.proposalDir)) + '/proposal.md')],
    archiveTarget: relativeRef,
    extra: {
      proposal_id: context.meta.id,
      note_kind: 'progress',
    },
    created: context.meta.created_at || nowIso(),
    updated: nowIso(),
  });
  writeText(context.notesPath, renderDocumentWithFrontmatter(frontmatter, body));
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
  updateDocumentFrontmatter(context.notesPath, { updated: timestamp });
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
      ...task.model_refs.map((ref) => `model-ref:${ref}`),
      ...task.convention_refs.map((ref) => `convention-ref:${ref}`),
    ];

    const descriptionLines = [
      task.outcome ? `Outcome: ${task.outcome}` : '',
      task.decision_refs.length > 0 ? `Decision refs: ${task.decision_refs.join(', ')}` : '',
      task.archive_target_refs.length > 0 ? `Archive refs: ${task.archive_target_refs.join(', ')}` : '',
      task.model_refs.length > 0 ? `Model refs: ${task.model_refs.join(', ')}` : '',
      task.convention_refs.length > 0 ? `Convention refs: ${task.convention_refs.join(', ')}` : '',
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
  updateDocumentFrontmatter(filePath, { updated: nowIso() });
  return true;
}

function ensureModuleArchiveTarget(targetPath, context, cwd = process.cwd()) {
  const moduleDir = path.extname(targetPath) === '.md' ? path.dirname(targetPath) : targetPath;
  ensureDir(moduleDir);
  const docsRoot = getDocsRoot(cwd, context.meta.workspace);
  const moduleDirRef = toPosixPath(path.relative(docsRoot, moduleDir));
  const readmeRef = toPosixPath(path.relative(docsRoot, path.join(moduleDir, 'README.md')));
  const designRef = toPosixPath(path.relative(docsRoot, path.join(moduleDir, 'design.md')));
  const apiRef = toPosixPath(path.relative(docsRoot, path.join(moduleDir, 'api.md')));
  const historyRef = toPosixPath(path.relative(docsRoot, path.join(moduleDir, 'history.md')));
  const ownership = buildArchiveDocumentOwnership('module', context, moduleDirRef);
  const commonRelatedDocs = normalizeRefList([
    workspaceDocRef(context.meta.workspace, toPosixPath(path.relative(docsRoot, context.proposalDir)) + '/proposal.md'),
  ]);

  const templateFiles = [
    ['README.md', path.join('modules', 'README.md')],
    ['design.md', path.join('modules', 'design.md')],
    ['api.md', path.join('modules', 'api.md')],
    ['history.md', path.join('modules', 'history.md')],
  ];

  for (const [fileName, templatePath] of templateFiles) {
    const absolutePath = path.join(moduleDir, fileName);
    if (fileExists(absolutePath)) continue;
    let rendered;
    if (fileName === 'history.md') {
      const frontmatter = makeDocFrontmatter({
        docType: 'module',
        title: `${path.basename(moduleDir)} History`,
        status: 'active',
        workspace: context.meta.workspace,
        fileRef: historyRef,
        ownership,
        informationClass: 'module',
        topics: context.meta.tags || [],
        relatedDocs: commonRelatedDocs,
        archiveTarget: historyRef,
        extra: {
          module_name: path.basename(moduleDir),
          module_status: 'active',
          primary_proposal: context.meta.id,
        },
        created: todayDate(),
        updated: todayDate(),
      });
      rendered = renderDocumentWithFrontmatter(frontmatter, `# ${path.basename(moduleDir)} History\n`);
    } else {
      const template = loadTemplate(cwd, templatePath);
      const frontmatter = makeDocFrontmatter({
        docType: 'module',
        title: path.basename(moduleDir),
        status: 'active',
        workspace: context.meta.workspace,
        fileRef: fileName === 'README.md' ? readmeRef : fileName === 'api.md' ? apiRef : designRef,
        ownership,
        informationClass: 'module',
        topics: context.meta.tags || [],
        relatedDocs: commonRelatedDocs,
        archiveTarget: fileName === 'README.md' ? readmeRef : fileName === 'api.md' ? apiRef : designRef,
        extra: {
          module_name: path.basename(moduleDir),
          module_status: 'active',
          primary_proposal: context.meta.id,
        },
        created: todayDate(),
        updated: todayDate(),
      });
      rendered = renderDocumentTemplate(template, frontmatter, {
        '<Module Name>': path.basename(moduleDir),
        '<proposal id>': context.meta.id,
        'CR26052001': context.meta.id,
        '2026-05-20': todayDate(),
        'What changed in the module': context.meta.title,
      });
    }
    writeText(absolutePath, rendered);
  }

  const historyPath = path.join(moduleDir, 'history.md');
  const block = [
    `## ${todayDate()}`,
    '',
    `- Proposal: ${context.meta.id}`,
    `- Summary: ${context.meta.title}`,
    `- Source: ${path.relative(path.dirname(historyPath), context.proposalDir)}`,
    '',
  ].join('\n');
  appendArchiveBlock(historyPath, context.meta.id, block);

  const readmePath = path.join(moduleDir, 'README.md');
  appendArchiveBlock(readmePath, context.meta.id, [
    '## Archived proposals',
    '',
    `- ${context.meta.id}: ${context.meta.title}`,
    '',
  ].join('\n'));

  return {
    type: 'module',
    path: moduleDir,
  };
}

function ensureArchitectureArchiveTarget(targetPath, context, cwd = process.cwd()) {
  const docsRoot = getDocsRoot(cwd, context.meta.workspace);
  const fileRef = toPosixPath(path.relative(docsRoot, targetPath));
  const ownership = buildArchiveDocumentOwnership('architecture', context, fileRef);
  if (!fileExists(targetPath)) {
    const template = loadTemplate(cwd, path.join('architecture', 'system.md'));
    const rendered = renderDocumentTemplate(template, makeDocFrontmatter({
      docType: 'architecture',
      title: path.basename(targetPath, path.extname(targetPath)),
      status: 'active',
      workspace: context.meta.workspace,
      fileRef,
      ownership,
      informationClass: 'architecture',
      topics: context.meta.tags || [],
      relatedDocs: [workspaceDocRef(context.meta.workspace, toPosixPath(path.relative(docsRoot, context.proposalDir)) + '/proposal.md')],
      archiveTarget: fileRef,
      extra: {
        architecture_topic: path.basename(targetPath, path.extname(targetPath)),
        architecture_status: 'active',
        primary_proposal: context.meta.id,
      },
      created: todayDate(),
      updated: todayDate(),
    }), {
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
  const docsRoot = getDocsRoot(cwd, context.meta.workspace);
  const fileRef = toPosixPath(path.relative(docsRoot, targetPath));
  const ownership = buildArchiveDocumentOwnership('decision', context, fileRef);
  if (!fileExists(targetPath)) {
    const template = loadTemplate(cwd, path.join('decisions', 'ADR-template.md'));
    const rendered = renderDocumentTemplate(template, makeDocFrontmatter({
      docType: 'adr',
      title: path.basename(targetPath, path.extname(targetPath)),
      status: 'proposed',
      workspace: context.meta.workspace,
      fileRef,
      ownership,
      informationClass: 'adr',
      topics: context.meta.tags || [],
      relatedDocs: [workspaceDocRef(context.meta.workspace, toPosixPath(path.relative(docsRoot, context.proposalDir)) + '/proposal.md')],
      archiveTarget: fileRef,
      extra: {
        adr_id: path.basename(targetPath, path.extname(targetPath)),
        adr_status: 'proposed',
        primary_proposal: context.meta.id,
      },
      created: todayDate(),
      updated: todayDate(),
    }), {
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

function ensureConventionArchiveTarget(targetPath, context, cwd = process.cwd()) {
  const docsRoot = getDocsRoot(cwd, context.meta.workspace);
  const fileRef = toPosixPath(path.relative(docsRoot, targetPath));
  const ownership = buildArchiveDocumentOwnership('convention', context, fileRef);
  if (!fileExists(targetPath)) {
    const template = loadTemplate(cwd, path.join('conventions', 'convention.md'));
    const rendered = renderDocumentTemplate(template, makeDocFrontmatter({
      docType: 'convention',
      title: path.basename(targetPath, path.extname(targetPath)),
      status: 'active',
      workspace: context.meta.workspace,
      fileRef,
      ownership,
      informationClass: 'convention',
      topics: context.meta.tags || [],
      relatedDocs: [workspaceDocRef(context.meta.workspace, toPosixPath(path.relative(docsRoot, context.proposalDir)) + '/proposal.md')],
      archiveTarget: fileRef,
      extra: {
        convention_status: 'active',
        enforcement: 'must',
        applies_to: [path.basename(targetPath, path.extname(targetPath))],
        origin_proposal: context.meta.id,
      },
      created: nowIso(),
      updated: nowIso(),
    }), {
      '<Convention Title>': path.basename(targetPath, path.extname(targetPath)),
      '<proposal id>': context.meta.id,
      '<name or team>': context.meta.owner || 'unknown-owner',
      '<ISO-8601 timestamp>': nowIso(),
      '<ISO date>': todayDate(),
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
    type: 'convention',
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
  if (target.type === 'convention') {
    return ensureConventionArchiveTarget(absolutePath, context, cwd);
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
      size_class: context.meta.size_class || null,
      ownership: context.meta.ownership || [],
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
      created: parsed.created || null,
      updated: parsed.updated || null,
      expected_size_class: parsed.expected_size_class || null,
      ownership: parsed.ownership || [],
    });
  }

  return summaries.sort((a, b) => a.title.localeCompare(b.title, 'zh-Hans-CN'));
}

function resolveExplorationDir(target, cwd = process.cwd()) {
  const absoluteTarget = path.resolve(cwd, target);
  if (fileExists(path.join(absoluteTarget, 'index.md'))) {
    return absoluteTarget;
  }

  const matches = [];
  for (const workspace of listDocsWorkspaces(cwd)) {
    const explorationsRoot = path.join(getWorkspaceDocsRoot(workspace.name, cwd), 'explorations');
    if (!fileExists(explorationsRoot)) continue;

    for (const entry of fs.readdirSync(explorationsRoot, { withFileTypes: true })) {
      if (!entry.isDirectory()) continue;
      const explorationDir = path.join(explorationsRoot, entry.name);
      if (!fileExists(path.join(explorationDir, 'index.md'))) continue;
      if (entry.name === target) {
        matches.push(explorationDir);
      }
    }
  }

  if (matches.length === 1) {
    return matches[0];
  }
  if (matches.length > 1) {
    throw new Error(`Exploration ${target} resolved to multiple directories: ${matches.join(', ')}`);
  }

  throw new Error(`Exploration directory not found: ${target}`);
}

function loadExplorationContext(target, cwd = process.cwd()) {
  const explorationDir = resolveExplorationDir(target, cwd);
  const indexPath = path.join(explorationDir, 'index.md');
  return {
    explorationDir,
    indexPath,
    parsed: parseExplorationIndex(readFileRequired(indexPath)),
    text: readFileRequired(indexPath),
  };
}

function loadProjectRuleBundle(cwd = process.cwd(), workspaceName = null) {
  const workspace = workspaceName ? getDocsWorkspace(workspaceName, cwd) : resolveWorkspaceForCwd(cwd);
  const rulesRoot = getProjectRulesRoot(cwd, workspace.name);
  const orderedFiles = [
    'README.md',
    'workflow.md',
    'classification.md',
    'intake.md',
    'explore.md',
    'propose.md',
    'archive.md',
  ];

  const files = orderedFiles.map((fileName) => {
    const filePath = path.join(rulesRoot, fileName);
    const exists = fileExists(filePath);
    return {
      file_name: fileName,
      path: filePath,
      exists,
      content: exists ? readFileRequired(filePath) : null,
    };
  });

  return {
    workspace,
    rulesRoot,
    files,
    missing_files: files.filter((file) => !file.exists).map((file) => file.file_name),
    available: fileExists(rulesRoot),
  };
}

function resolveIntakeDir(target, cwd = process.cwd()) {
  const absoluteTarget = path.isAbsolute(target)
    ? target
    : path.resolve(cwd, target);

  if (fileExists(path.join(absoluteTarget, 'index.md'))) {
    return absoluteTarget;
  }

  const workspaces = listDocsWorkspaces(cwd);
  const matches = [];
  for (const workspace of workspaces) {
    const intakeRoot = path.join(getWorkspaceDocsRoot(workspace.name, cwd), 'intake');
    if (!fileExists(intakeRoot)) continue;
    const candidate = path.join(intakeRoot, target);
    if (fileExists(path.join(candidate, 'index.md'))) {
      matches.push(candidate);
    }
  }

  if (matches.length === 1) return matches[0];
  if (matches.length > 1) {
    throw new Error(`Intake package ${target} resolved to multiple directories: ${matches.join(', ')}`);
  }

  throw new Error(`Intake package directory not found: ${target}`);
}

function loadIntakeContext(target, cwd = process.cwd()) {
  const intakeDir = resolveIntakeDir(target, cwd);
  const indexPath = path.join(intakeDir, 'index.md');
  const parsed = parseFrontmatterDocument(readFileRequired(indexPath));
  const markdownFiles = listMarkdownFiles(intakeDir)
    .filter((filePath) => path.basename(filePath) !== 'index.md')
    .sort((a, b) => a.localeCompare(b));

  const files = [
    {
      file_name: 'index.md',
      path: indexPath,
      exists: true,
      content: readFileRequired(indexPath),
      frontmatter: parsed.frontmatter || null,
    },
    ...markdownFiles.map((filePath) => ({
      file_name: path.relative(intakeDir, filePath).split(path.sep).join('/'),
      path: filePath,
      exists: true,
      content: readFileRequired(filePath),
      frontmatter: parseFrontmatterDocument(readFileRequired(filePath)).frontmatter || null,
    })),
  ];

  return {
    intakeDir,
    indexPath,
    parsed: parsed.frontmatter || {},
    text: readFileRequired(indexPath),
    files,
    assets: fileExists(path.join(intakeDir, 'assets'))
      ? fs.readdirSync(path.join(intakeDir, 'assets')).sort((a, b) => a.localeCompare(b))
      : [],
  };
}

function validateExplorationContext(context, cwd = process.cwd()) {
  const errors = [];
  const warnings = [];
  const { explorationDir } = context;
  const frontmatter = validateDocumentFrontmatter(context.indexPath, 'exploration', errors, warnings, {
    label: 'index.md',
    allowedStatuses: new Set(['active', 'validated', 'parked', 'rejected', 'archived']),
  });

  if (frontmatter) {
    if (!frontmatter.expected_size_class) errors.push('index.md missing expected_size_class');
    if (frontmatter.expected_size_class && !SIZE_CLASSES.has(frontmatter.expected_size_class)) {
      errors.push(`invalid exploration expected size class: ${frontmatter.expected_size_class}`);
    }

    if (frontmatter.reusable_rules !== undefined) {
      if (!Array.isArray(frontmatter.reusable_rules)) {
        errors.push('index.md reusable_rules must be an array');
      } else {
        for (const [index, rule] of frontmatter.reusable_rules.entries()) {
          if (!rule || typeof rule !== 'object' || Array.isArray(rule)) {
            errors.push(`index.md reusable_rules[${index}] must be an object`);
            continue;
          }
          if (!rule.title) {
            errors.push(`index.md reusable_rules[${index}] missing title`);
          }
        }
      }
    }

    if (!Array.isArray(frontmatter.ownership) || frontmatter.ownership.length === 0) {
      errors.push('index.md must declare at least one ownership entry');
    } else {
      const primaryCount = frontmatter.ownership.filter((entry) => entry.role === 'primary').length;
      if (primaryCount !== 1) {
        errors.push(`index.md must declare exactly one primary ownership entry, found ${primaryCount}`);
      }
      for (const entry of frontmatter.ownership) {
        if (!entry.type) {
          errors.push(`ownership entry could not be parsed: ${entry.target}`);
          continue;
        }
        if (!OWNERSHIP_TYPES.has(entry.type)) {
          errors.push(`invalid exploration ownership type: ${entry.type}`);
        }
      }
    }
  }

  for (const requiredDir of ['journal', 'findings', 'decisions', 'artifacts']) {
    const absolutePath = path.join(explorationDir, requiredDir);
    if (!fileExists(absolutePath)) {
      warnings.push(`exploration missing directory: ${absolutePath}`);
    }
  }

  for (const filePath of listMarkdownFiles(path.join(explorationDir, 'findings'))) {
    const finding = validateDocumentFrontmatter(filePath, 'finding', errors, warnings, {
      label: path.relative(explorationDir, filePath).split(path.sep).join('/'),
    });
    if (finding && finding.exploration_slug && finding.exploration_slug !== path.basename(explorationDir)) {
      errors.push(`${path.relative(explorationDir, filePath).split(path.sep).join('/')} exploration_slug must match the directory slug`);
    }
  }

  for (const filePath of listMarkdownFiles(path.join(explorationDir, 'decisions'))) {
    const decision = validateDocumentFrontmatter(filePath, 'decision', errors, warnings, {
      label: path.relative(explorationDir, filePath).split(path.sep).join('/'),
    });
    if (decision && decision.exploration_slug && decision.exploration_slug !== path.basename(explorationDir)) {
      errors.push(`${path.relative(explorationDir, filePath).split(path.sep).join('/')} exploration_slug must match the directory slug`);
    }
  }

  for (const filePath of listMarkdownFiles(path.join(explorationDir, 'journal'))) {
    const journal = validateDocumentFrontmatter(filePath, 'journal', errors, warnings, {
      label: path.relative(explorationDir, filePath).split(path.sep).join('/'),
    });
    if (journal && journal.exploration_slug && journal.exploration_slug !== path.basename(explorationDir)) {
      errors.push(`${path.relative(explorationDir, filePath).split(path.sep).join('/')} exploration_slug must match the directory slug`);
    }
  }

  return { errors, warnings };
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
  getProjectRulesRoot,
  getIntakeRoot,
  getToolRoot,
  getWorkflowConfigPath,
  getWorkflowConfig,
  getWorkspaceDocsRoot,
  isProposalId,
  listProposalSummaries,
  listExplorationSummaries,
  listDocsWorkspaces,
  loadExplorationContext,
  loadIntakeContext,
  loadProjectRuleBundle,
  loadProposalContext,
  nowIso,
  parseCliArgs,
  parseOwnershipEntry,
  parseSimpleYaml,
  parseTaskMap,
  resolveProposalDir,
  resolveIntakeDir,
  resolveExplorationDir,
  resolveWorkspaceForCwd,
  serializeYaml,
  slugify,
  todayDate,
  transitionProposalStatus,
  validateExplorationContext,
  validateProposalContext,
  writeProposalMeta,
};
