const fs = require('fs').promises;
const path = require('path');
const { homedir } = require('os');

function resolveEnv(...names) {
  for (const name of names) {
    const value = process.env[name];
    if (value) return value;
  }
  return undefined;
}

const USER_CONFIG_PATH = path.join(homedir(), '.config', 'flowforge', 'memory.json');

async function readJsonIfExists(filePath) {
  try {
    const raw = await fs.readFile(filePath, 'utf8');
    return JSON.parse(raw);
  } catch {
    return null;
  }
}

function getProjectKeys(project, fallbackSlug) {
  return [...new Set([
    project?.slug,
    project?.id,
    fallbackSlug,
  ].filter(Boolean))];
}

function mergeMemoryProvider(base, overrides = {}) {
  return {
    ...base,
    ...overrides,
  };
}

function applyUserConfig(config, userConfig, projectKeys) {
  if (!userConfig || typeof userConfig !== 'object') {
    return config;
  }

  const baseMemoryProvider = userConfig.memory_provider && typeof userConfig.memory_provider === 'object'
    ? userConfig.memory_provider
    : {};
  const projectOverrides = userConfig.projects && typeof userConfig.projects === 'object'
    ? projectKeys.map((key) => userConfig.projects[key]).find((value) => value && typeof value === 'object')
    : null;
  const projectMemoryProvider = projectOverrides?.memory_provider && typeof projectOverrides.memory_provider === 'object'
    ? projectOverrides.memory_provider
    : projectOverrides && typeof projectOverrides === 'object'
      ? projectOverrides
      : {};

  return mergeConfig(config, {
    memory_provider: mergeMemoryProvider(baseMemoryProvider, projectMemoryProvider),
  });
}

const DEFAULT_CONFIG = {
  project: {
    id: null,
    slug: null,
  },
  memory_provider: {
    type: 'memory-mcp',
    enabled: false,
    endpoint: resolveEnv(
      'FLOWFORGE_MEMORY_ENDPOINT',
      'CLAUDE_FLOWFORGE_MEMORY_ENDPOINT',
    ) || 'http://127.0.0.1:8000',
    apiKey: resolveEnv(
      'FLOWFORGE_MEMORY_API_KEY',
      'CLAUDE_FLOWFORGE_MEMORY_API_KEY',
    ),
    tags: [],
    timeoutMs: 5000,
  },
};

function buildUrl(baseUrl, pathname) {
  const normalizedBase = baseUrl.endsWith('/') ? baseUrl : `${baseUrl}/`;
  const normalizedPath = pathname.startsWith('/') ? pathname.slice(1) : pathname;
  return new URL(normalizedPath, normalizedBase).toString();
}

function mergeConfig(base, overrides = {}) {
  return {
    ...base,
    ...overrides,
    project: {
      ...base.project,
      ...(overrides.project || {}),
    },
    memory_provider: {
      ...base.memory_provider,
      ...(overrides.memory_provider || {}),
    },
  };
}

async function loadConfig(directory) {
  let config = DEFAULT_CONFIG;
  const configPaths = [
    path.join(directory, '.flowforge', 'config.json'),
    path.join(directory, 'workflow', 'config.json'),
  ];

  for (const configPath of configPaths) {
    try {
      const raw = await fs.readFile(configPath, 'utf8');
      config = mergeConfig(config, JSON.parse(raw));
      break;
    } catch {
      // continue
    }
  }

  const fallbackSlug = path.basename(directory) || 'project';
  const projectKeys = getProjectKeys(config.project, fallbackSlug);
  const userConfig = await readJsonIfExists(USER_CONFIG_PATH);
  config = applyUserConfig(config, userConfig, projectKeys);

  return mergeConfig(config, {
    memory_provider: {
      endpoint: resolveEnv(
        'FLOWFORGE_MEMORY_ENDPOINT',
        'CLAUDE_FLOWFORGE_MEMORY_ENDPOINT',
      ) || config.memory_provider.endpoint,
      apiKey: resolveEnv(
        'FLOWFORGE_MEMORY_API_KEY',
        'CLAUDE_FLOWFORGE_MEMORY_API_KEY',
      ) || config.memory_provider.apiKey,
    },
  });
}

function getProjectContext(directory, config) {
  const fallbackSlug = path.basename(directory) || 'project';
  const slug = config.project.slug || config.project.id || fallbackSlug;
  const tags = config.memory_provider.tags?.length
    ? config.memory_provider.tags
    : [`project:${slug}`];

  return { slug, tags };
}

async function requestJson(config, pathname, init = {}) {
  const controller = new AbortController();
  const timeout = setTimeout(() => controller.abort(), config.memory_provider.timeoutMs);

  try {
    const headers = {
      Accept: 'application/json',
      'Content-Type': 'application/json',
      ...init.headers,
    };

    if (config.memory_provider.apiKey) {
      headers.Authorization = `Bearer ${config.memory_provider.apiKey}`;
    }

    const response = await fetch(buildUrl(config.memory_provider.endpoint, pathname), {
      ...init,
      headers,
      signal: controller.signal,
    });

    const text = await response.text();
    const body = text ? JSON.parse(text) : null;

    if (!response.ok) {
      const detail = body?.detail || body?.error || response.statusText;
      throw new Error(`${response.status} ${detail}`);
    }

    return body;
  } finally {
    clearTimeout(timeout);
  }
}

async function checkDelayedReview(config, tags) {
  if (!config.memory_provider.enabled) {
    return [];
  }

  try {
    const today = new Date().toISOString().slice(0, 10);
    const result = await requestJson(config, '/api/memories/search', {
      method: 'POST',
      body: JSON.stringify({
        query: 'review pending decisions',
        tags: [...tags, 'review-pending'],
        limit: 10,
      }),
    });

    const memories = Array.isArray(result)
      ? result
      : Array.isArray(result?.memories)
        ? result.memories
        : Array.isArray(result?.results)
          ? result.results
          : [];

    return memories.filter((memory) => {
      const reviewAt = memory?.metadata?.review_at || memory?.memory?.metadata?.review_at;
      return reviewAt && reviewAt <= today;
    });
  } catch (error) {
    console.error('[flowforge-memory] Delayed review check failed:', error.message);
    return [];
  }
}

function formatReviewReminder(reviews) {
  if (!reviews.length) return '';

  const lines = [
    '## Pending Review Decisions',
    '',
  ];

  for (const review of reviews) {
    const content = review.content || review.memory?.content || 'Unknown';
    const reason = review.metadata?.review_reason || review.memory?.metadata?.review_reason || 'No reason';
    const reviewAt = review.metadata?.review_at || review.memory?.metadata?.review_at || 'Unknown';

    lines.push(`- ${content.slice(0, 100)}...`);
    lines.push(`  review_at: ${reviewAt}`);
    lines.push(`  reason: ${reason}`);
  }

  return lines.join('\n');
}

async function onSessionStart(context) {
  try {
    const directory = context.workingDirectory || process.cwd();
    const config = await loadConfig(directory);
    const project = getProjectContext(directory, config);

    console.log(`[flowforge-memory] Session starting for project: ${project.slug}`);
    const dueReviews = await checkDelayedReview(config, project.tags);
    if (dueReviews.length) {
      console.log(`\n${formatReviewReminder(dueReviews)}\n`);
    }
  } catch (error) {
    console.error('[flowforge-memory] Session start failed:', error.message);
  }
}

async function readStdinContext() {
  return new Promise((resolve, reject) => {
    let data = '';
    const timeout = setTimeout(() => resolve(null), 100);

    process.stdin.setEncoding('utf8');
    process.stdin.on('readable', () => {
      let chunk;
      while ((chunk = process.stdin.read()) !== null) {
        data += chunk;
      }
    });

    process.stdin.on('end', () => {
      clearTimeout(timeout);
      if (!data.trim()) {
        resolve(null);
        return;
      }

      try {
        resolve(JSON.parse(data));
      } catch (error) {
        reject(error);
      }
    });

    process.stdin.on('error', (error) => {
      clearTimeout(timeout);
      reject(error);
    });
  });
}

async function main() {
  try {
    const stdinContext = await readStdinContext();
    const context = stdinContext
      ? {
          workingDirectory: stdinContext.cwd || process.cwd(),
          sessionId: stdinContext.session_id || `session_${Date.now()}`,
        }
      : {
          workingDirectory: process.cwd(),
          sessionId: `session_${Date.now()}`,
        };

    await onSessionStart(context);
  } catch (error) {
    console.error('[flowforge-memory] Session start hook failed:', error);
    process.exit(1);
  }
}

module.exports = {
  name: 'tg-session-start',
  version: '2.0.0',
  description: 'FlowForge session-start hook',
  trigger: 'session-start',
  handler: onSessionStart,
  config: {
    async: true,
    timeout: 5000,
    priority: 'normal',
  },
};

if (require.main === module) {
  main();
}
