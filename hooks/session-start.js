const fs = require('fs').promises;
const path = require('path');

const DEFAULT_CONFIG = {
  memoryService: {
    endpoint: process.env.TG_MEMORY_ENDPOINT || 'http://127.0.0.1:8000',
    apiKey: process.env.TG_MEMORY_API_KEY || '',
  },
};

function buildUrl(baseUrl, pathname) {
  const normalizedBase = baseUrl.endsWith('/') ? baseUrl : `${baseUrl}/`;
  const normalizedPath = pathname.startsWith('/') ? pathname.slice(1) : pathname;
  return new URL(normalizedPath, normalizedBase).toString();
}

async function requestJson(config, pathname, init = {}) {
  const controller = new AbortController();
  const timeout = setTimeout(() => controller.abort(), 5000);

  try {
    const headers = {
      Accept: 'application/json',
      'Content-Type': 'application/json',
      ...init.headers,
    };

    if (config.memoryService.apiKey) {
      headers.Authorization = `Bearer ${config.memoryService.apiKey}`;
    }

    const response = await fetch(buildUrl(config.memoryService.endpoint, pathname), {
      ...init,
      headers,
      signal: controller.signal,
    });

    const text = await response.text();
    let body = null;
    if (text) {
      try {
        body = JSON.parse(text);
      } catch {
        body = { detail: text };
      }
    }

    if (!response.ok) {
      const detail = body?.detail || body?.error || response.statusText;
      throw new Error(`${response.status} ${detail}`);
    }

    return body;
  } finally {
    clearTimeout(timeout);
  }
}

async function checkDelayedReview(config, projectTag) {
  try {
    const today = new Date().toISOString().slice(0, 10);

    const result = await requestJson(config, '/api/memories/search', {
      method: 'POST',
      body: JSON.stringify({
        query: 'review pending decisions',
        tags: [projectTag, 'review-pending'],
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

    const dueReviews = memories.filter((m) => {
      const reviewAt = m?.metadata?.review_at || m?.memory?.metadata?.review_at;
      if (!reviewAt) return false;
      return reviewAt <= today;
    });

    return dueReviews;
  } catch (error) {
    console.error('[Session Hook] Delayed review check failed:', error.message);
    return [];
  }
}

function formatReviewReminder(reviews) {
  if (!reviews.length) return '';

  const lines = [
    '## ⚠️ 待回顾的决策',
    '',
    '以下决策已到回顾日期：',
    '',
  ];

  for (const review of reviews) {
    const content = review.content || review.memory?.content || 'Unknown';
    const reason = review.metadata?.review_reason || review.memory?.metadata?.review_reason || '无原因';
    const reviewAt = review.metadata?.review_at || review.memory?.metadata?.review_at || 'Unknown';

    lines.push(`- **${content.slice(0, 100)}...`);
    lines.push(`  - 回顾日期: ${reviewAt}`);
    lines.push(`  - 原因: ${reason}`);
    lines.push('');
  }

  lines.push('请回顾这些决策，确认是否需要调整。');
  lines.push('使用 `#skip` 标记跳过，或使用 Memory MCP 更新决策状态。');

  return lines.join('\n');
}

async function onSessionStart(context) {
  try {
    const directory = context.workingDirectory || process.cwd();
    const projectName = path.basename(directory) || 'project';
    const projectTag = `project:${projectName}`;

    console.log(`[Session Hook] Session starting for project: ${projectName}`);

    const dueReviews = await checkDelayedReview(DEFAULT_CONFIG, projectTag);

    if (dueReviews.length > 0) {
      const reminder = formatReviewReminder(dueReviews);
      console.log('\n' + reminder + '\n');
    } else {
      console.log('[Session Hook] No pending reviews found');
    }

  } catch (error) {
    console.error('[Session Hook] Error in session start:', error.message);
  }
}

async function readStdinContext() {
  return new Promise((resolve, reject) => {
    let data = '';

    const timeout = setTimeout(() => {
      resolve(null);
    }, 100);

    process.stdin.setEncoding('utf8');
    process.stdin.on('readable', () => {
      let chunk;
      while ((chunk = process.stdin.read()) !== null) {
        data += chunk;
      }
    });

    process.stdin.on('end', () => {
      clearTimeout(timeout);
      if (data.trim()) {
        try {
          resolve(JSON.parse(data));
        } catch (error) {
          console.error('[Session Hook] Failed to parse stdin JSON:', error.message);
          reject(error);
        }
      } else {
        resolve(null);
      }
    });

    process.stdin.on('error', (error) => {
      clearTimeout(timeout);
      console.error('[Session Hook] Stdin error:', error.message);
      reject(error);
    });
  });
}

async function main() {
  try {
    const stdinContext = await readStdinContext();

    let context;

    if (stdinContext) {
      context = {
        workingDirectory: stdinContext.cwd || process.cwd(),
        sessionId: stdinContext.session_id || `session_${Date.now()}`,
      };
    } else {
      context = {
        workingDirectory: process.cwd(),
        sessionId: `session_${Date.now()}`,
      };
    }

    await onSessionStart(context);

  } catch (error) {
    console.error('[Session Hook] Session start hook failed:', error);
    process.exit(1);
  }
}

module.exports = {
  name: 'tg-session-start',
  version: '1.0.0',
  description: 'Check delayed review on session start for Tangram V2',
  trigger: 'session-start',
  handler: onSessionStart,
  config: {
    async: true,
    timeout: 5000,
    priority: 'normal'
  },
};

if (require.main === module) {
  main();
}
