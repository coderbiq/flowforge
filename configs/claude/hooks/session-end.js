const fs = require('fs').promises;
const path = require('path');
const { spawn } = require('child_process');

const DEFAULT_CONFIG = {
  paths: {
    state_root: '.workflow/state',
  },
  session: {
    minContentLength: 100,
  },
};

function mergeConfig(base, overrides = {}) {
  return {
    ...base,
    ...overrides,
    paths: {
      ...base.paths,
      ...(overrides.paths || {}),
    },
    session: {
      ...base.session,
      ...(overrides.session || {}),
    },
  };
}

async function loadConfig(directory) {
  let config = DEFAULT_CONFIG;
  const configPath = path.join(directory, 'workflow', 'config.json');

  try {
    const raw = await fs.readFile(configPath, 'utf8');
    config = mergeConfig(config, JSON.parse(raw));
  } catch {
    // use defaults
  }

  return config;
}

async function ensureStateDir(directory, config) {
  const stateRoot = path.join(directory, config.paths.state_root);
  const sessionsDir = path.join(stateRoot, 'sessions');

  await fs.mkdir(sessionsDir, { recursive: true });
  return { stateRoot, sessionsDir };
}

async function getSessionFilePath(directory, config, sessionId) {
  const { sessionsDir } = await ensureStateDir(directory, config);
  return path.join(sessionsDir, `${sessionId}.json`);
}

async function loadSessionState(directory, config, sessionId) {
  const filePath = await getSessionFilePath(directory, config, sessionId);
  try {
    return JSON.parse(await fs.readFile(filePath, 'utf8'));
  } catch {
    return null;
  }
}

async function saveSessionState(directory, config, sessionId, state) {
  const filePath = await getSessionFilePath(directory, config, sessionId);
  state.updated_at = new Date().toISOString();
  await fs.writeFile(filePath, JSON.stringify(state, null, 2));
}

async function updateActiveSession(directory, config, sessionId) {
  const { stateRoot } = await ensureStateDir(directory, config);
  const activePath = path.join(stateRoot, 'active-session.json');
  const activeData = {
    session_id: sessionId,
    updated_at: new Date().toISOString(),
  };
  await fs.writeFile(activePath, JSON.stringify(activeData, null, 2));
}

function checkUserMarkers(messages) {
  if (!Array.isArray(messages)) return { remember: false, skip: false };

  const recentMessages = messages.slice(-5).map((message) => message?.content || '').join(' ').toLowerCase();
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
        if (entry.type !== 'user' && entry.type !== 'assistant') continue;

        const msg = entry.message;
        if (!msg?.role || !msg?.content) continue;

        let contentText = '';
        if (typeof msg.content === 'string') {
          contentText = msg.content;
        } else if (Array.isArray(msg.content)) {
          contentText = msg.content
            .filter((block) => block.type === 'text')
            .map((block) => block.text)
            .join('\n');
        }

        if (contentText) {
          messages.push({ role: msg.role, content: contentText });
        }
      } catch {
        // ignore invalid lines
      }
    }

    return { messages };
  } catch (error) {
    console.error('[tg-memory] Failed to parse transcript:', error.message);
    return { messages: [] };
  }
}

function extractProposalId(text) {
  const match = text.match(/\bCR\d{8}\b/);
  return match ? match[0] : null;
}

function analyzeConversation(messages) {
  const analysis = {
    proposal_id: null,
    current_focus: 'Session completed',
    active_files: [],
    completed_items: [],
    next_actions: ['Review session notes and continue the active proposal'],
    notes: '',
  };

  if (!Array.isArray(messages) || messages.length === 0) {
    return analysis;
  }

  const allContent = messages.map((message) => message.content || '').join('\n');
  const proposalId = extractProposalId(allContent);
  if (proposalId) {
    analysis.proposal_id = proposalId;
  }

  const filePattern = /([A-Za-z0-9_./-]+\.(?:js|ts|jsx|tsx|py|rs|go|java|cpp|c|md|json|yaml|yml))/g;
  const matches = allContent.match(filePattern) || [];
  analysis.active_files = [...new Set(matches)].slice(0, 20);

  const lastUserMessage = [...messages].reverse().find((message) => message.role === 'user');
  if (lastUserMessage?.content) {
    analysis.current_focus = lastUserMessage.content.slice(0, 160);
  }

  analysis.notes = messages
    .slice(-3)
    .map((message) => `${message.role}: ${message.content.slice(0, 120)}`)
    .join('\n');
  return analysis;
}

async function onSessionEnd(context) {
  try {
    const directory = context.workingDirectory || process.cwd();
    const config = await loadConfig(directory);
    const sessionId = context.sessionId || `session_${Date.now()}`;

    if (context.conversation?.messages) {
      const markers = checkUserMarkers(context.conversation.messages);
      if (markers.skip) {
        console.log('[tg-memory] Session state update skipped by #skip marker');
        return;
      }
    }

    let sessionState = await loadSessionState(directory, config, sessionId);
    if (!sessionState) {
      sessionState = {
        session_id: sessionId,
        updated_at: new Date().toISOString(),
        proposal_id: null,
        current_focus: '',
        active_files: [],
        completed_items: [],
        next_actions: [],
        notes: '',
      };
    }

    const contentLength = context.conversation?.messages
      ? context.conversation.messages.reduce((sum, message) => sum + (message?.content?.length || 0), 0)
      : 0;
    const markers = checkUserMarkers(context.conversation?.messages || []);

    if (contentLength < config.session.minContentLength && !markers.remember) {
      console.log(`[tg-memory] Session content too short (${contentLength} chars), skipping update`);
      return;
    }

    if (context.conversation?.messages) {
      sessionState = {
        ...sessionState,
        ...analyzeConversation(context.conversation.messages),
      };
    }

    await saveSessionState(directory, config, sessionId, sessionState);
    await updateActiveSession(directory, config, sessionId);
    console.log(`[tg-memory] Session state updated: ${sessionId}`);
  } catch (error) {
    console.error('[tg-memory] Session end failed:', error.message);
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

async function runInBackground(context) {
  const child = spawn(process.execPath, [__filename, '--background'], {
    detached: true,
    stdio: 'ignore',
    env: {
      ...process.env,
      TG_SESSION_CONTEXT: JSON.stringify(context),
    },
  });

  child.unref();
}

async function main() {
  try {
    if (process.argv.includes('--background')) {
      const contextJson = process.env.TG_SESSION_CONTEXT;
      if (!contextJson) {
        console.error('[tg-memory] No context provided for background processing');
        process.exit(1);
      }

      await onSessionEnd(JSON.parse(contextJson));
      return;
    }

    const stdinContext = await readStdinContext();
    if (!stdinContext?.transcript_path) {
      console.log('[tg-memory] No transcript context - skipping session state update');
      return;
    }

    const conversation = await parseTranscript(stdinContext.transcript_path);
    const context = {
      workingDirectory: stdinContext.cwd || process.cwd(),
      sessionId: stdinContext.session_id || `session_${Date.now()}`,
      reason: stdinContext.reason,
      conversation,
    };

    await runInBackground(context);
    console.log('[tg-memory] Session state update started in background');
  } catch (error) {
    console.error('[tg-memory] Session end hook failed:', error);
    process.exit(1);
  }
}

module.exports = {
  name: 'tg-session-end',
  version: '2.0.0',
  description: 'tg-workflow session-end hook',
  trigger: 'session-end',
  handler: onSessionEnd,
  config: {
    async: true,
    timeout: 5000,
    priority: 'normal',
  },
};

if (require.main === module) {
  main();
}
