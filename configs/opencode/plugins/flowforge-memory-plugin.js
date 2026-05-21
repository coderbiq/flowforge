import { readFile, writeFile, mkdir } from "node:fs/promises"
import { homedir } from "node:os"
import path from "node:path"

function resolveEnv(...names) {
  for (const name of names) {
    const value = process.env[name]
    if (value) return value
  }
  return undefined
}

const DEFAULT_CONFIG = {
  project: {
    id: null,
    slug: null,
  },
  paths: {
    state_root: ".flowforge/state",
  },
  memory_provider: {
    type: "memory-mcp",
    enabled: false,
    endpoint: "http://127.0.0.1:8000",
    apiKey: "",
    timeoutMs: 5000,
    tags: [],
  },
  session: {
    minContentLength: 100,
  },
}

function mergeConfig(base, overrides = {}) {
  return {
    ...base,
    ...overrides,
    project: {
      ...base.project,
      ...(overrides.project || {}),
    },
    paths: {
      ...base.paths,
      ...(overrides.paths || {}),
    },
    memory_provider: {
      ...base.memory_provider,
      ...(overrides.memory_provider || {}),
    },
    session: {
      ...base.session,
      ...(overrides.session || {}),
    },
  }
}

function environmentOverrides() {
  const endpoint = resolveEnv(
    "FLOWFORGE_MEMORY_ENDPOINT",
    "OPENCODE_FLOWFORGE_MEMORY_ENDPOINT",
    "OPENCODE_MEMORY_ENDPOINT",
  )
  const apiKey = resolveEnv(
    "FLOWFORGE_MEMORY_API_KEY",
    "OPENCODE_FLOWFORGE_MEMORY_API_KEY",
    "OPENCODE_MEMORY_API_KEY",
  )

  const memory_provider = {}
  if (endpoint !== undefined) {
    memory_provider.endpoint = endpoint
  }
  if (apiKey !== undefined) {
    memory_provider.apiKey = apiKey
  }

  return {
    memory_provider,
  }
}

async function loadConfig(directory) {
  let config = DEFAULT_CONFIG
  const configPaths = [
    path.join(directory, ".flowforge", "config.json"),
    path.join(directory, ".opencode", "flowforge-memory-plugin.json"),
    path.join(homedir(), ".config", "opencode", "flowforge-memory-plugin.json"),
  ]

  for (const configPath of configPaths) {
    try {
      const raw = await readFile(configPath, "utf8")
      config = mergeConfig(config, JSON.parse(raw))
    } catch {
      // continue
    }
  }

  return mergeConfig(config, environmentOverrides())
}

function getProjectTags(directory, config) {
  const fallbackSlug = path.basename(directory) || "project"
  const slug = config.project.slug || config.project.id || fallbackSlug
  return config.memory_provider.tags?.length ? config.memory_provider.tags : [`project:${slug}`]
}

function buildUrl(baseUrl, pathname) {
  const normalizedBase = baseUrl.endsWith("/") ? baseUrl : `${baseUrl}/`
  const normalizedPath = pathname.startsWith("/") ? pathname.slice(1) : pathname
  return new URL(normalizedPath, normalizedBase).toString()
}

async function requestJson(config, pathname, init = {}) {
  const controller = new AbortController()
  const timeout = setTimeout(() => controller.abort(), config.memory_provider.timeoutMs)

  try {
    const headers = {
      Accept: "application/json",
      "Content-Type": "application/json",
      ...(init.headers || {}),
    }

    if (config.memory_provider.apiKey) {
      headers.Authorization = `Bearer ${config.memory_provider.apiKey}`
    }

    const response = await fetch(buildUrl(config.memory_provider.endpoint, pathname), {
      ...init,
      headers,
      signal: controller.signal,
    })

    const text = await response.text()
    const body = text ? JSON.parse(text) : null

    if (!response.ok) {
      const detail = body?.detail || body?.error || response.statusText
      throw new Error(`${response.status} ${detail}`)
    }

    return body
  } finally {
    clearTimeout(timeout)
  }
}

async function ensureStateDir(directory, config) {
  const stateRoot = path.join(directory, config.paths.state_root)
  const sessionsDir = path.join(stateRoot, "sessions")
  await mkdir(sessionsDir, { recursive: true })
  return { stateRoot, sessionsDir }
}

async function getSessionFilePath(directory, config, sessionId) {
  const { sessionsDir } = await ensureStateDir(directory, config)
  return path.join(sessionsDir, `${sessionId}.json`)
}

async function loadSessionState(directory, config, sessionId) {
  const filePath = await getSessionFilePath(directory, config, sessionId)
  try {
    return JSON.parse(await readFile(filePath, "utf8"))
  } catch {
    return null
  }
}

async function saveSessionState(directory, config, sessionId, state) {
  const filePath = await getSessionFilePath(directory, config, sessionId)
  state.updated_at = new Date().toISOString()
  await writeFile(filePath, JSON.stringify(state, null, 2))
}

async function updateActiveSession(directory, config, sessionId) {
  const { stateRoot } = await ensureStateDir(directory, config)
  const activePath = path.join(stateRoot, "active-session.json")
  await writeFile(activePath, JSON.stringify({
    session_id: sessionId,
    updated_at: new Date().toISOString(),
  }, null, 2))
}

async function checkDelayedReview(config, tags) {
  if (!config.memory_provider.enabled) {
    return []
  }

  try {
    const today = new Date().toISOString().slice(0, 10)
    const result = await requestJson(config, "/api/memories/search", {
      method: "POST",
      body: JSON.stringify({
        query: "review pending decisions",
        tags: [...tags, "review-pending"],
        limit: 10,
      }),
    })

    const memories = Array.isArray(result)
      ? result
      : Array.isArray(result?.memories)
        ? result.memories
        : Array.isArray(result?.results)
          ? result.results
          : []

    return memories.filter((memory) => {
      const reviewAt = memory?.metadata?.review_at || memory?.memory?.metadata?.review_at
      return reviewAt && reviewAt <= today
    })
  } catch (error) {
    console.error("[flowforge-memory] Delayed review check failed:", error.message)
    return []
  }
}

function formatReviewReminder(reviews) {
  if (!reviews.length) return ""

  return [
    "## Pending Review Decisions",
    "",
    ...reviews.flatMap((review) => {
      const content = review.content || review.memory?.content || "Unknown"
      const reason = review.metadata?.review_reason || review.memory?.metadata?.review_reason || "No reason"
      const reviewAt = review.metadata?.review_at || review.memory?.metadata?.review_at || "Unknown"
      return [
        `- ${content.slice(0, 100)}...`,
        `  review_at: ${reviewAt}`,
        `  reason: ${reason}`,
      ]
    }),
  ].join("\n")
}

function checkUserMarkers(messages) {
  if (!Array.isArray(messages)) return { remember: false, skip: false }

  const recentMessages = messages.slice(-5).map((message) => message?.content || "").join(" ").toLowerCase()
  return {
    remember: recentMessages.includes("#remember"),
    skip: recentMessages.includes("#skip"),
  }
}

function extractProposalId(text) {
  const match = text.match(/\bCR\d{8}\b/)
  return match ? match[0] : null
}

function analyzeMessages(messages) {
  const allContent = messages.map((message) => message?.content || "").join("\n")
  const filePattern = /([A-Za-z0-9_./-]+\.(?:js|ts|jsx|tsx|py|rs|go|java|cpp|c|md|json|yaml|yml))/g
  const activeFiles = [...new Set(allContent.match(filePattern) || [])].slice(0, 20)
  const proposalId = extractProposalId(allContent)

  return {
    proposal_id: proposalId,
    current_focus: messages.at(-1)?.content?.slice(0, 160) || "Session idle",
    active_files: activeFiles,
    completed_items: [],
    next_actions: ["Review session notes and continue the active proposal"],
    notes: messages.slice(-3).map((message) => `${message.role}: ${(message.content || "").slice(0, 120)}`).join("\n"),
  }
}

export const OpenCodeMemoryPlugin = async ({ client, directory }) => {
  const config = await loadConfig(directory)
  const appLog = client.app.log.bind(client.app)
  const projectTags = getProjectTags(directory, config)

  const logInfo = async (message) => {
    await appLog({ body: { service: "flowforge-memory", level: "info", message } }).catch(() => {})
  }

  return {
    event: async ({ event }) => {
      if (event.type === "session.created") {
        const sessionId = event.properties?.info?.id
        if (!sessionId) return

        await updateActiveSession(directory, config, sessionId)
        const dueReviews = await checkDelayedReview(config, projectTags)
        if (dueReviews.length > 0) {
          await logInfo(`Delayed review check: ${dueReviews.length} decisions due`)
          console.log(`\n${formatReviewReminder(dueReviews)}\n`)
        }
      }

      if (event.type === "session.idle") {
        const sessionId = event.properties?.info?.id
        const messages = event.properties?.messages || []
        if (!sessionId) return

        const markers = checkUserMarkers(messages)
        if (markers.skip) {
          await logInfo("Session state update skipped by #skip marker")
          return
        }

        const contentLength = messages.reduce((sum, message) => sum + (message?.content?.length || 0), 0)
        if (contentLength < config.session.minContentLength && !markers.remember) {
          await logInfo(`Session content too short (${contentLength} chars), skipping update`)
          return
        }

        const previous = await loadSessionState(directory, config, sessionId)
        const nextState = {
          session_id: sessionId,
          updated_at: new Date().toISOString(),
          proposal_id: null,
          current_focus: "",
          active_files: [],
          completed_items: [],
          next_actions: [],
          notes: "",
          ...(previous || {}),
          ...analyzeMessages(messages),
        }

        await saveSessionState(directory, config, sessionId, nextState)
        await updateActiveSession(directory, config, sessionId)
        await logInfo(`Session state updated: ${sessionId}`)
      }
    },
  }
}
