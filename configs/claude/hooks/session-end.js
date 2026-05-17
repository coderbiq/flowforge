const fs = require('fs').promises;
const path = require('path');
const { spawn } = require('child_process');

const DEFAULT_CONFIG = {
  memoryService: {
    endpoint: process.env.TG_MEMORY_ENDPOINT || 'http://127.0.0.1:8000',
    apiKey: process.env.TG_MEMORY_API_KEY || '',
  },
  session: {
    minContentLength: 100,
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

async function ensureMemoryDir(directory) {
  const memoryDir = path.join(directory, '.memory');
  const sessionsDir = path.join(memoryDir, 'sessions');

  try {
    await fs.mkdir(sessionsDir, { recursive: true });
  } catch {
  }

  return { memoryDir, sessionsDir };
}

async function getSessionFilePath(directory, sessionId) {
  const { sessionsDir } = await ensureMemoryDir(directory);
  return path.join(sessionsDir, `${sessionId}.json`);
}

async function loadSessionState(directory, sessionId) {
  const filePath = await getSessionFilePath(directory, sessionId);
  try {
    const raw = await fs.readFile(filePath, 'utf8');
    return JSON.parse(raw);
  } catch {
    return null;
  }
}

async function saveSessionState(directory, sessionId, state) {
  const filePath = await getSessionFilePath(directory, sessionId);
  state.updated_at = new Date().toISOString();
  await fs.writeFile(filePath, JSON.stringify(state, null, 2));
}

async function updateActiveSession(directory, sessionId) {
  const { memoryDir } = await ensureMemoryDir(directory);
  const activePath = path.join(memoryDir, 'active.json');

  const activeData = {
    active_session: sessionId,
    last_switched: new Date().toISOString(),
  };

  await fs.writeFile(activePath, JSON.stringify(activeData, null, 2));
}

function checkUserMarkers(messages) {
  if (!Array.isArray(messages)) return { remember: false, skip: false };

  const recentMessages = messages.slice(-5).map(m => m?.content || '').join(' ').toLowerCase();

  return {
    remember: recentMessages.includes('#remember'),
    skip: recentMessages.includes('#skip'),
  };
}

async function parseTranscript(transcriptPath) {
  try {
    const content = await fs.readFile(transcriptPath, 'utf8');
    const lines = content.trim().split('\n');
    const messages = [];

    for (const line of lines) {
      if (!line.trim()) continue;

      try {
        const entry = JSON.parse(line);

        if (entry.type === 'user' || entry.type === 'assistant') {
          const msg = entry.message;
          if (msg && msg.role && msg.content) {
            let contentText = '';
            if (typeof msg.content === 'string') {
              contentText = msg.content;
            } else if (Array.isArray(msg.content)) {
              contentText = msg.content
                .filter(block => block.type === 'text')
                .map(block => block.text)
                .join('\n');
            }

            if (contentText) {
              messages.push({
                role: msg.role,
                content: contentText
              });
            }
          }
        }
      } catch {
        continue;
      }
    }

    return { messages };
  } catch (error) {
    console.error('[Session Hook] Failed to parse transcript:', error.message);
    return { messages: [] };
  }
}

function analyzeConversation(messages) {
  const analysis = {
    working_files: [],
    completed_tasks: [],
    current_task: '',
    next_steps: [],
  };

  if (!Array.isArray(messages)) return analysis;

  const allContent = messages.map(m => m.content || '').join('\n');

  const filePattern = /(?:file|path|in|at|edit|write|read|update|modify):\s*([^\s,]+\.(?:js|ts|jsx|tsx|py|rs|go|java|cpp|c|md|json|yaml|yml))/gi;
  const fileMatches = allContent.matchAll(filePattern);
  for (const match of fileMatches) {
    if (match[1] && !analysis.working_files.includes(match[1])) {
      analysis.working_files.push(match[1]);
    }
  }

  analysis.current_task = 'Session completed';
  analysis.next_steps.push('Review session state file for details');

  return analysis;
}

async function onSessionEnd(context) {
  try {
    const directory = context.workingDirectory || process.cwd();
    const sessionId = context.sessionId || `session_${Date.now()}`;

    console.log(`[Session Hook] Session ending: ${sessionId}`);

    if (context.conversation && context.conversation.messages) {
      const markers = checkUserMarkers(context.conversation.messages);
      if (markers.skip) {
        console.log('[Session Hook] Session state update skipped by #skip marker');
        return;
      }
    }

    let sessionState = await loadSessionState(directory, sessionId);
    if (!sessionState) {
      sessionState = {
        id: sessionId,
        created_at: new Date().toISOString(),
        status: 'active',
        metadata: {
          title: 'New Session',
          description: '',
        },
        state: {
          working_files: [],
          completed_tasks: [],
          current_task: '',
          next_steps: [],
        },
        summary: {
          short: '',
          long: '',
        },
      };
    }

    const contentLength = context.conversation?.messages
      ? context.conversation.messages.reduce((sum, m) => sum + (m?.content?.length || 0), 0)
      : 0;

    const markers = checkUserMarkers(context.conversation?.messages || []);

    if (contentLength < DEFAULT_CONFIG.session.minContentLength && !markers.remember) {
      console.log(`[Session Hook] Session content too short (${contentLength} chars), skipping update`);
      return;
    }

    if (context.conversation?.messages) {
      const analysis = analyzeConversation(context.conversation.messages);
      sessionState.state = {
        ...sessionState.state,
        ...analysis,
      };
    }

    sessionState.status = 'idle';
    sessionState.updated_at = new Date().toISOString();

    await saveSessionState(directory, sessionId, sessionState);
    await updateActiveSession(directory, sessionId);

    console.log(`[Session Hook] Session state updated: ${sessionId}`);

  } catch (error) {
    console.error('[Session Hook] Error in session end:', error.message);
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

async function runInBackground(context) {
  const scriptPath = __filename;
  const contextJson = JSON.stringify(context);

  const child = spawn(process.execPath, [scriptPath, '--background'], {
    detached: true,
    stdio: 'ignore',
    env: {
      ...process.env,
      TG_SESSION_CONTEXT: contextJson,
    },
  });

  child.unref();
}

async function main() {
  try {
    if (process.argv.includes('--background')) {
      const contextJson = process.env.TG_SESSION_CONTEXT;
      if (!contextJson) {
        console.error('[Session Hook] No context provided for background processing');
        process.exit(1);
      }

      const context = JSON.parse(contextJson);
      await onSessionEnd(context);
      return;
    }

    const stdinContext = await readStdinContext();

    let context;

    if (stdinContext && stdinContext.transcript_path) {
      console.log(`[Session Hook] Reading transcript: ${stdinContext.transcript_path}`);
      console.log(`[Session Hook] Session end reason: ${stdinContext.reason || 'unknown'}`);

      const conversation = await parseTranscript(stdinContext.transcript_path);

      context = {
        workingDirectory: stdinContext.cwd || process.cwd(),
        sessionId: stdinContext.session_id || `session_${Date.now()}`,
        reason: stdinContext.reason,
        conversation: conversation
      };

      console.log(`[Session Hook] Parsed ${conversation.messages.length} messages from transcript`);
    } else {
      console.log('[Session Hook] No stdin context - skipping session state update');
      return;
    }

    await runInBackground(context);
    console.log('[Session Hook] Session state update started in background');

  } catch (error) {
    console.error('[Session Hook] Session end hook failed:', error);
    process.exit(1);
  }
}

module.exports = {
  name: 'tg-session-end',
  version: '1.0.0',
  description: 'Update Tangram V2 session state on session end',
  trigger: 'session-end',
  handler: onSessionEnd,
  config: {
    async: true,
    timeout: 5000,
    priority: 'normal'
  },
};

if (require.main === module) {
  main();
}
