/**
 * Tangram V2 Memory Plugin for OpenCode
 *
 * Features:
 * - session.idle: Update session state file (.memory/sessions/<id>.json)
 * - session.created: Check delayed review (review-pending tags)
 */

import { readFile, writeFile, mkdir } from "node:fs/promises"
import { homedir } from "node:os"
import path from "node:path"

const DEFAULT_CONFIG = {
  memoryService: {
    endpoint: "http://127.0.0.1:8000",
    apiKey: "",
    timeoutMs: 5000,
  },
  session: {
    minContentLength: 100, // Minimum characters to trigger session state update
  },
}

// ============ Config Loading ============

function environmentOverrides() {
  const overrides = { memoryService: {} }

  const endpoint = process.env.TG_MEMORY_ENDPOINT || process.env.OPENCODE_MEMORY_ENDPOINT
  if (endpoint) {
    overrides.memoryService.endpoint = endpoint
  }

  const apiKey = process.env.TG_MEMORY_API_KEY || process.env.OPENCODE_MEMORY_API_KEY
  if (apiKey) {
    overrides.memoryService.apiKey = apiKey
  }

  return overrides
}

function mergeConfig(base, overrides = {}) {
  return {
    ...base,
    ...overrides,
    memoryService: {
      ...base.memoryService,
      ...(overrides.memoryService || {}),
    },
    session: {
      ...base.session,
      ...(overrides.session || {}),
    },
  }
}

async function loadConfig(directory, options = {}) {
  let config = DEFAULT_CONFIG

  const configPaths = [
    path.join(directory, ".opencode", "tg-memory-plugin.json"),
    path.join(homedir(), ".config", "opencode", "tg-memory-plugin.json"),
  ]

  for (const configPath of configPaths) {
    try {
      const raw = await readFile(configPath, "utf8")
      const parsed = JSON.parse(raw)
      config = mergeConfig(config, parsed)
      break
    } catch {
      // Continue to next config path
    }
  }

  return mergeConfig(config, environmentOverrides())
}

// ============ Memory MCP API ============

function buildUrl(baseUrl, pathname) {
  const normalizedBase = baseUrl.endsWith("/") ? baseUrl : `${baseUrl}/`
  const normalizedPath = pathname.startsWith("/") ? pathname.slice(1) : pathname
  return new URL(normalizedPath, normalizedBase).toString()
}

function buildHeaders(config) {
  const headers = {
    Accept: "application/json",
    "Content-Type": "application/json",
  }

  if (config.memoryService.apiKey) {
    headers.Authorization = `Bearer ${config.memoryService.apiKey}`
  }

  return headers
}

async function requestJson(config, pathname, init = {}) {
  const controller = new AbortController()
  const timeout = setTimeout(() => controller.abort(), config.memoryService.timeoutMs)

  try {
    const response = await fetch(buildUrl(config.memoryService.endpoint, pathname), {
      ...init,
      headers: buildHeaders(config),
      signal: controller.signal,
    })

    const text = await response.text()
    let body = null
    if (text) {
      try {
        body = JSON.parse(text)
      } catch {
        body = { detail: text }
      }
    }

    if (!response.ok) {
      const detail = body?.detail || body?.error || response.statusText
      throw new Error(`${response.status} ${detail}`)
    }

    return body
  } finally {
    clearTimeout(timeout)
  }
}

// ============ Session State Management ============

async function ensureMemoryDir(directory) {
  const memoryDir = path.join(directory, ".memory")
  const sessionsDir = path.join(memoryDir, "sessions")

  try {
    await mkdir(sessionsDir, { recursive: true })
  } catch {
    // Directory exists
  }

  return { memoryDir, sessionsDir }
}

async function getSessionFilePath(directory, sessionId) {
  const { sessionsDir } = await ensureMemoryDir(directory)
  return path.join(sessionsDir, `${sessionId}.json`)
}

async function loadSessionState(directory, sessionId) {
  const filePath = await getSessionFilePath(directory, sessionId)
  try {
    const raw = await readFile(filePath, "utf8")
    return JSON.parse(raw)
  } catch {
    return null
  }
}

async function saveSessionState(directory, sessionId, state) {
  const filePath = await getSessionFilePath(directory, sessionId)
  state.updated_at = new Date().toISOString()
  await writeFile(filePath, JSON.stringify(state, null, 2))
}

async function updateActiveSession(directory, sessionId) {
  const { memoryDir } = await ensureMemoryDir(directory)
  const activePath = path.join(memoryDir, "active.json")

  const activeData = {
    active_session: sessionId,
    last_switched: new Date().toISOString(),
  }

  await writeFile(activePath, JSON.stringify(activeData, null, 2))
}

// ============ Harvest API ============

async function harvestSessionContent(config, transcript) {
  try {
    const result = await requestJson(config, "/api/harvest", {
      method: "POST",
      body: JSON.stringify({
        transcript,
        extraction_type: "session_state",
      }),
    })
    return result
  } catch (error) {
    console.error("Harvest API failed:", error.message)
    return null
  }
}

// ============ Delayed Review Check ============

async function checkDelayedReview(config, projectTag) {
  try {
    const today = new Date().toISOString().slice(0, 10)

    const result = await requestJson(config, "/api/memories/search", {
      method: "POST",
      body: JSON.stringify({
        query: "review pending decisions",
        tags: [projectTag, "review-pending"],
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

    // Filter by review_at date
    const dueReviews = memories.filter((m) => {
      const reviewAt = m?.metadata?.review_at || m?.memory?.metadata?.review_at
      if (!reviewAt) return false
      return reviewAt <= today
    })

    return dueReviews
  } catch (error) {
    console.error("Delayed review check failed:", error.message)
    return []
  }
}

function formatReviewReminder(reviews) {
  if (!reviews.length) return ""

  const lines = [
    "## ⚠️ 待回顾的决策",
    "",
    "以下决策已到回顾日期：",
    "",
  ]

  for (const review of reviews) {
    const content = review.content || review.memory?.content || "Unknown"
    const reason = review.metadata?.review_reason || review.memory?.metadata?.review_reason || "无原因"
    const reviewAt = review.metadata?.review_at || review.memory?.metadata?.review_at || "Unknown"

    lines.push(`- **${content.slice(0, 100)}...`)
    lines.push(`  - 回顾日期: ${reviewAt}`)
    lines.push(`  - 原因: ${reason}`)
    lines.push("")
  }

  lines.push("请回顾这些决策，确认是否需要调整。")
  lines.push("使用 `#skip` 标记跳过，或使用 Memory MCP 更新决策状态。")

  return lines.join("\n")
}

// ============ User Control Markers ============

function checkUserMarkers(messages) {
  if (!Array.isArray(messages)) return { remember: false, skip: false }

  const recentMessages = messages.slice(-5).map((m) => m?.content || "").join(" ").toLowerCase()

  return {
    remember: recentMessages.includes("#remember"),
    skip: recentMessages.includes("#skip"),
  }
}

// ============ Plugin Export ============

export const OpenCodeMemoryPlugin = async ({ client, directory }, options = {}) => {
  const config = await loadConfig(directory, options)
  const appLog = client.app.log.bind(client.app)

  const logInfo = async (message) => {
    await appLog({ body: { service: "tg-memory", level: "info", message } }).catch(() => {})
  }

  const logWarn = async (message) => {
    await appLog({ body: { service: "tg-memory", level: "warn", message } }).catch(() => {})
  }

  const projectName = path.basename(directory) || "project"
  const projectTag = `project:${projectName}`

  return {
    event: async ({ event }) => {
      // session.created: Check delayed review
      if (event.type === "session.created") {
        const sessionId = event.properties?.info?.id
        if (!sessionId) return

        // Update active session
        await updateActiveSession(directory, sessionId)

        // Check delayed review
        const dueReviews = await checkDelayedReview(config, projectTag)
        if (dueReviews.length > 0) {
          const reminder = formatReviewReminder(dueReviews)
          // Note: OpenCode doesn't have a direct way to inject messages
          // Log the reminder for visibility
          await logInfo(`Delayed review check: ${dueReviews.length} decisions due`)
          console.log("\n" + reminder + "\n")
        }
      }

      // session.idle: Update session state
      if (event.type === "session.idle") {
        const sessionId = event.properties?.info?.id
        const messages = event.properties?.messages || []

        if (!sessionId) return

        // Check user markers
        const markers = checkUserMarkers(messages)
        if (markers.skip) {
          await logInfo("Session state update skipped by #skip marker")
          return
        }

        // Load existing session state
        let sessionState = await loadSessionState(directory, sessionId)
        if (!sessionState) {
          sessionState = {
            id: sessionId,
            created_at: new Date().toISOString(),
            status: "active",
            metadata: {
              title: "New Session",
              description: "",
            },
            state: {
              working_files: [],
              completed_tasks: [],
              current_task: "",
              next_steps: [],
            },
            summary: {
              short: "",
              long: "",
            },
          }
        }

        // Check content length threshold
        const contentLength = messages.reduce((sum, m) => sum + (m?.content?.length || 0), 0)
        if (contentLength < config.session.minContentLength && !markers.remember) {
          await logInfo(`Session content too short (${contentLength} chars), skipping update`)
          return
        }

        // Call Harvest API to extract session content
        // Note: transcript access depends on OpenCode API
        // For now, just update timestamp
        sessionState.status = "idle"
        sessionState.updated_at = new Date().toISOString()

        await saveSessionState(directory, sessionId, sessionState)
        await logInfo(`Session state updated: ${sessionId}`)
      }

      // session.deleted: Clean up session state
      if (event.type === "session.deleted") {
        const sessionId = event.properties?.info?.id
        if (sessionId) {
          const filePath = await getSessionFilePath(directory, sessionId)
          try {
            const { unlink } = await import("node:fs/promises")
            await unlink(filePath)
            await logInfo(`Session state deleted: ${sessionId}`)
          } catch {
            // File doesn't exist
          }
        }
      }
    },
  }
}
