/**
 * FlowForge project-level pi extension.
 *
 * Deployed by `flowforge agents deploy` (pi host enabled) to
 * `.pi/extensions/flowforge.ts`. Source of truth for the default
 * write-protection glob set is `defaultTestFileGlobs` in
 * `internal/command/agents.go`; the copy below must stay identical
 * item-for-item (drift is a defect).
 *
 * Capabilities:
 * 1. Test-file write guard: blocks `write`/`edit` on paths matching
 *    `agents.test_file_globs` from `.flowforge/config.yaml` (fallback:
 *    FlowForge defaults). Escape valve: `agents.disable_test_guard: true`.
 *    `bash` is never intercepted (command content cannot be matched by
 *    path globs reliably).
 * 2. Licensed test-author entry: a read-only `flowforge_test_authority` tool
 *    plus per-session, per-instance, exact-file grant records under
 *    `.flowforge/test-author-grants/`. Only a background native runner child
 *    session may consume a grant, and only the main session may issue one.
 *    See the "Licensed test-author entry" note above the factory.
 * 3. Native LLM tools `flowforge_frontier` / `flowforge_check` wrapping the
 *    FlowForge CLI (PATH first, then <project>/bin/flowforge).
 * 4. Raw-script delegation model guard: blocks `subagent` dispatches made
 *    through `workflowScript`/`workflowScriptPath` without an explicit
 *    `model` — such runs silently inherit the orchestrating session's
 *    model and bypass the per-role model pins deployed by `flowforge
 *    agents deploy`. Escape valves: dispatch a named `agent:` (model pin
 *    applies), pass `model` explicitly (subagent level and on each
 *    `runs.run()`), or set `agents.disable_raw_script_model_guard: true`.
 */

import { execFile } from "node:child_process";
import { randomUUID } from "node:crypto";
import { existsSync, lstatSync, readFileSync, readlinkSync, statSync } from "node:fs";
import * as path from "node:path";
import { Type } from "typebox";
import type { ExtensionAPI } from "@earendil-works/pi-coding-agent";

/** Mirrors defaultTestFileGlobs in internal/command/agents.go — keep in sync. */
const DEFAULT_TEST_FILE_GLOBS: string[] = [
	"**/*_test.go",
	"**/src/test/**",
	"**/src/integrationTest/**",
	"**/__tests__/**",
	"**/*.test.ts",
	"**/*.test.tsx",
	"**/*.spec.ts",
];

/** Project-relative directory holding licensed test-author grant records. */
const GRANT_DIR_REL = ".flowforge/test-author-grants";
/** A session id may only be used as a file name when it is a safe id. */
const SAFE_SESSION_ID = /^[A-Za-z0-9-]+$/;
/** Roles the main session may write into a grant record. */
const TEST_AUTHOR_ROLES = ["plan", "refine"] as const;
/** Exact grant record shape; any other field invalidates the record. */
const GRANT_FIELDS = ["version", "sessionId", "instanceId", "projectRoot", "role", "paths", "expiresAt"];

interface FlowForgeProjectConfig {
	docsDir: string;
	testFileGlobs: string[];
	disableTestGuard: boolean;
	disableRawScriptModelGuard: boolean;
}

/**
 * Minimal line-level reader for the flat `.flowforge/config.yaml` subset
 * FlowForge needs: top-level `docs_dir`, and inside the `agents:` block
 * `test_file_globs` (block or inline list) and `disable_test_guard`.
 * Missing file, missing `agents:` key, or missing keys fall back to
 * defaults without error.
 */
export function readFlowForgeConfig(configPath: string): FlowForgeProjectConfig {
	const result: FlowForgeProjectConfig = { docsDir: "docs", testFileGlobs: [], disableTestGuard: false, disableRawScriptModelGuard: false };
	let text: string;
	try {
		text = readFileSync(configPath, "utf8");
	} catch {
		return result;
	}
	let inAgents = false;
	let lastKey = "";
	let lastKeyIndent = -1;
	for (const rawLine of text.split("\n")) {
		const line = rawLine.replace(/\r$/, "");
		const trimmed = line.trim();
		if (!trimmed || trimmed.startsWith("#")) continue;
		const indent = line.length - line.trimStart().length;

		const topKey = /^([A-Za-z_][\w-]*):/.exec(trimmed);
		if (topKey && indent === 0) {
			inAgents = topKey[1] === "agents";
			if (!inAgents && topKey[1] === "docs_dir") {
				result.docsDir = unquote(trimmed.slice(topKey[0].length)) || "docs";
			}
			lastKey = "";
			continue;
		}
		if (!inAgents) continue;

		const listItem = /^-\s*(.*)$/.exec(trimmed);
		if (listItem && lastKey === "test_file_globs" && indent > lastKeyIndent) {
			const value = unquote(listItem[1]);
			if (value) result.testFileGlobs.push(value);
			continue;
		}
		const kv = /^([A-Za-z_][\w-]*):\s*(.*)$/.exec(trimmed);
		if (kv) {
			lastKey = kv[1];
			lastKeyIndent = indent;
			if (kv[1] === "disable_test_guard") {
				result.disableTestGuard = kv[2].trim() === "true";
			}
			if (kv[1] === "disable_raw_script_model_guard") {
				result.disableRawScriptModelGuard = kv[2].trim() === "true";
			}
			if (kv[1] === "test_file_globs") {
				const inline = /^\[(.*)\]$/.exec(kv[2].trim());
				if (inline) {
					for (const part of inline[1].split(",")) {
						const value = unquote(part);
						if (value) result.testFileGlobs.push(value);
					}
				}
			}
		}
	}
	return result;
}

function unquote(value: string): string {
	const v = value.trim();
	if (v.length >= 2 && ((v.startsWith('"') && v.endsWith('"')) || (v.startsWith("'") && v.endsWith("'")))) {
		return v.slice(1, -1);
	}
	return v.split(/\s+#/)[0].trim();
}

/**
 * Glob → RegExp: `**` crosses directory separators, `*` matches within
 * one segment. A leading double-star also matches zero directories, so the
 * `*_test.go` default pattern matches root-level files as well as nested
 * ones.
 */
export function globToRegExp(glob: string): RegExp {
	const segments = glob.split("/");
	let source = "^";
	for (let i = 0; i < segments.length; i++) {
		const segment = segments[i];
		if (segment === "**") {
			if (i === segments.length - 1) {
				source += "(?:.*)";
			} else {
				source += "(?:[^/]*(?:/|$))*";
			}
			continue;
		}
		source += segment
			.replace(/[.+^${}()|[\]\\]/g, "\\$&")
			.replace(/\*/g, "[^/]*")
			.replace(/\?/g, ".");
		if (i < segments.length - 1) source += "/";
	}
	return new RegExp(source + "$");
}

/** True when a write/edit target path is protected by one of the globs. */
export function matchesAnyGlob(target: string, globs: string[], projectRoot: string): string | undefined {
	const candidates = pathCandidates(target, projectRoot);
	for (const glob of globs) {
		const pattern = globToRegExp(glob);
		for (const candidate of candidates) {
			if (pattern.test(candidate)) return glob;
		}
	}
	return undefined;
}

/**
 * A tool `path` may arrive absolute, relative to the project root, or with a
 * `./` prefix. Test the raw normalized form plus the project-root-relative
 * form so both shapes match the same globs.
 */
function pathCandidates(target: string, projectRoot: string): string[] {
	const normalized = target.replace(/\\/g, "/").replace(/^\.\//, "");
	const candidates = [normalized];
	const root = path.resolve(projectRoot).replace(/\\/g, "/") + "/";
	if (normalized.startsWith(root)) candidates.push(normalized.slice(root.length));
	return candidates;
}

/** A symlink chain longer than this cannot be trusted to resolve. */
const MAX_LINK_HOPS = 40;

interface PathAudit {
	/** Absolute path after replacing every walked symlink component by its destination. */
	resolved: string;
	/** Absolute paths of the symlink components walked through, dangling ones included. */
	linkComponents: string[];
	/** Set when a component could not be inspected; every caller must deny then. */
	error?: string;
}

/**
 * Walk an absolute path component by component with `lstat`, so symbolic links
 * are seen even when they dangle. `realpathSync` only succeeds while the whole
 * chain exists and silently drops a dangling link's destination, which let a
 * link stand in for an uncreated ordinary file. A link is followed with
 * `readlink` (relative destinations resolve against the link's directory), so
 * `resolved` is the place a write would land whether or not it exists yet. A
 * plain missing component (ENOENT) ends the walk and the remaining suffix is
 * appended verbatim: nothing below it exists, so nothing below it can be a
 * symlink. Any other `lstat`/`readlink` failure is reported as `error` so the
 * caller refuses the write instead of guessing a safe target.
 */
export function auditPath(target: string): PathAudit {
	const requested = path.resolve(target);
	const root = path.parse(requested).root;
	const segments = (absolute: string): string[] => {
		const base = path.parse(absolute).root;
		return absolute.slice(base.length).split(path.sep).filter((segment) => segment.length > 0);
	};
	let current = root;
	let pending = segments(requested);
	const linkComponents: string[] = [];
	while (pending.length > 0) {
		const segment = pending.shift() as string;
		const next = path.join(current, segment);
		let entry: ReturnType<typeof lstatSync>;
		try {
			entry = lstatSync(next);
		} catch (error) {
			const code = (error as NodeJS.ErrnoException).code ?? "";
			const resolved = path.join(next, ...pending);
			if (code === "ENOENT") return { resolved, linkComponents };
			return { resolved, linkComponents, error: `cannot inspect path component ${next}${code ? ` (${code})` : ""}` };
		}
		if (!entry.isSymbolicLink()) {
			current = next;
			continue;
		}
		linkComponents.push(next);
		if (linkComponents.length > MAX_LINK_HOPS) {
			return { resolved: requested, linkComponents, error: `symbolic link chain under ${requested} is too long` };
		}
		let destination: string;
		try {
			destination = readlinkSync(next);
		} catch (error) {
			const code = (error as NodeJS.ErrnoException).code ?? "";
			return { resolved: requested, linkComponents, error: `cannot read symbolic link ${next}${code ? ` (${code})` : ""}` };
		}
		const absoluteDestination = path.resolve(current, destination);
		current = path.parse(absoluteDestination).root;
		pending = segments(absoluteDestination).concat(pending);
	}
	return { resolved: current, linkComponents };
}

/**
 * Canonicalize an absolute path through the walked components, so a target that
 * does not exist yet still resolves every symlink that does exist — including a
 * dangling one, whose destination is kept instead of being collapsed back onto
 * the link path.
 */
export function resolveThroughLinks(target: string): string {
	return auditPath(target).resolved;
}

/**
 * True when `candidate` is `root` or lives inside it. Containment is decided on
 * whole path segments: only an explicit `..` segment escapes the root, so an
 * ordinary child whose name merely starts with two dots still counts as inside.
 */
function isWithinPath(root: string, candidate: string): boolean {
	const relative = path.relative(root, candidate);
	if (relative === "") return true;
	if (path.isAbsolute(relative)) return false;
	return relative !== ".." && !relative.startsWith(`..${path.sep}`);
}

/**
 * True when a write target touches the grant-record zone: inside the directory,
 * the directory itself, or any of its ancestors (so a child cannot replace the
 * directory wholesale). The resolved form carries the same decision through
 * symlinked aliases of the directory.
 */
function touchesGrantRecords(candidate: string, grantDir: string): boolean {
	return isWithinPath(grantDir, candidate) || isWithinPath(candidate, grantDir);
}

/**
 * A grant may only list concrete project-relative files: no globs, no absolute
 * paths, no `.`/`..` segments, no Windows separators, no empty segment.
 */
function isPlainGrantPath(value: unknown): value is string {
	if (typeof value !== "string" || value.length === 0 || value.includes("\\")) return false;
	if (/[*?[\]{}]/.test(value)) return false;
	if (path.isAbsolute(value)) return false;
	return value.split("/").every((segment) => segment !== "" && segment !== "." && segment !== "..");
}

/**
 * Strict expiry contract: full calendar date, seconds, optional milliseconds,
 * and an explicit `Z`. `Date.parse` alone accepts year-only/date-only input,
 * missing seconds and non-`Z` offsets, and normalizes impossible values
 * (Feb 30, hour 24, minute 60), so a malformed expiry could still read as a
 * future time. The pattern rejects the loose shapes; the `Date.UTC` round-trip
 * below rejects values the calendar silently normalizes, and returns the exact
 * instant only for a timestamp that means what it says.
 */
const GRANT_EXPIRY = /^(\d{4})-(\d{2})-(\d{2})T(\d{2}):(\d{2}):(\d{2})(?:\.(\d{1,3}))?Z$/;

export function parseGrantExpiry(value: string): number | undefined {
	const match = GRANT_EXPIRY.exec(value);
	if (!match) return undefined;
	const [, yearText, monthText, dayText, hourText, minuteText, secondText, millisecondText] = match;
	const year = Number(yearText);
	const month = Number(monthText);
	const day = Number(dayText);
	const hour = Number(hourText);
	const minute = Number(minuteText);
	const second = Number(secondText);
	const millisecond = millisecondText === undefined ? 0 : Number(millisecondText.padEnd(3, "0"));
	if (month < 1 || month > 12 || hour > 23 || minute > 59 || second > 59) return undefined;
	const epoch = Date.UTC(year, month - 1, day, hour, minute, second, millisecond);
	const roundTrip = new Date(epoch);
	if (
		roundTrip.getUTCFullYear() !== year ||
		roundTrip.getUTCMonth() !== month - 1 ||
		roundTrip.getUTCDate() !== day ||
		roundTrip.getUTCHours() !== hour ||
		roundTrip.getUTCMinutes() !== minute ||
		roundTrip.getUTCSeconds() !== second ||
		roundTrip.getUTCMilliseconds() !== millisecond
	) {
		return undefined;
	}
	return epoch;
}

interface TestAuthorGrant {
	role: (typeof TEST_AUTHOR_ROLES)[number];
	paths: string[];
	expiresAt: string;
}

type GrantState =
	| { status: "unsupported"; reason: string }
	| { status: "missing"; reason: string }
	| { status: "invalid"; reason: string }
	| { status: "valid"; grant: TestAuthorGrant };

interface GrantExpectation {
	sessionId: string;
	instanceId: string;
	projectRoot: string;
}

/**
 * Strictly parse one grant record against the live session identity. Every
 * rejection path is a denial: unknown fields, wrong types, empty or glob range,
 * traversal, wrong project/role/session/instance, and past or unparsable
 * expiry all invalidate the record instead of narrowing its meaning.
 */
export function evaluateTestAuthorGrant(raw: string, expected: GrantExpectation): GrantState {
	let parsed: unknown;
	try {
		parsed = JSON.parse(raw);
	} catch {
		return { status: "invalid", reason: "grant record is not valid JSON" };
	}
	if (typeof parsed !== "object" || parsed === null || Array.isArray(parsed)) {
		return { status: "invalid", reason: "grant record is not an object" };
	}
	const record = parsed as Record<string, unknown>;
	const keys = Object.keys(record).sort();
	const expectedKeys = [...GRANT_FIELDS].sort();
	if (keys.length !== expectedKeys.length || keys.some((key, index) => key !== expectedKeys[index])) {
		return { status: "invalid", reason: `grant record fields must be exactly ${GRANT_FIELDS.join(", ")}` };
	}
	if (record.version !== 1) {
		return { status: "invalid", reason: "unsupported grant version" };
	}
	if (typeof record.sessionId !== "string" || record.sessionId !== expected.sessionId) {
		return { status: "invalid", reason: "grant sessionId does not match this session" };
	}
	if (typeof record.instanceId !== "string" || record.instanceId.length === 0 || record.instanceId !== expected.instanceId) {
		return { status: "invalid", reason: "grant instanceId does not match this session instance" };
	}
	if (
		typeof record.projectRoot !== "string" ||
		!path.isAbsolute(record.projectRoot) ||
		resolveThroughLinks(record.projectRoot) !== expected.projectRoot
	) {
		return { status: "invalid", reason: "grant projectRoot does not match this project" };
	}
	if (typeof record.role !== "string" || !(TEST_AUTHOR_ROLES as readonly string[]).includes(record.role)) {
		return { status: "invalid", reason: "grant role must be plan or refine" };
	}
	if (!Array.isArray(record.paths) || record.paths.length === 0) {
		return { status: "invalid", reason: "grant paths must be a non-empty array" };
	}
	const paths: string[] = [];
	for (const entry of record.paths) {
		if (!isPlainGrantPath(entry)) {
			return { status: "invalid", reason: "grant paths must be concrete project-relative files (no globs, directories, absolute paths, or traversal)" };
		}
		paths.push(entry);
	}
	if (typeof record.expiresAt !== "string") {
		return { status: "invalid", reason: "grant expiresAt must be a UTC timestamp string" };
	}
	const expiry = parseGrantExpiry(record.expiresAt);
	if (expiry === undefined) {
		return { status: "invalid", reason: "grant expiresAt must be a full UTC timestamp (YYYY-MM-DDTHH:mm:ss[.sss]Z)" };
	}
	if (expiry <= Date.now()) {
		return { status: "invalid", reason: "grant has expired" };
	}
	return { status: "valid", grant: { role: record.role as TestAuthorGrant["role"], paths, expiresAt: record.expiresAt } };
}

interface TestWriteContext {
	projectRoot: string;
	grantDir: string;
	globs: string[];
	globSource: string;
	eligible: boolean;
	readGrantState: () => GrantState;
}

/**
 * Decide whether one `write`/`edit` target must be denied. `undefined` allows
 * the call; a string is the block reason. Grant records are never writable by
 * an eligible child, and a valid grant only ever covers its listed exact files.
 */
export function testWriteDenial(target: string, context: TestWriteContext): string | undefined {
	const { projectRoot, grantDir, globs, globSource, eligible } = context;
	const requested = path.resolve(projectRoot, target);
	const audit = auditPath(requested);
	const resolved = audit.resolved;
	if (audit.error) {
		return `cannot safely resolve the write target (${audit.error}): refusing the write`;
	}

	if (eligible) {
		// Both the lexical `.flowforge/test-author-grants` location and the place
		// that path actually lands (through an existing or dangling `.flowforge`
		// link) are protected, so an alias, an ancestor, or a record that does not
		// exist yet cannot be used to reach the zone. This denial is independent of
		// whether a grant is currently valid.
		const zoneAudit = auditPath(grantDir);
		if (zoneAudit.error) {
			return `cannot safely resolve the test-authority record zone (${zoneAudit.error}): refusing the write`;
		}
		for (const zone of new Set([path.resolve(grantDir), zoneAudit.resolved])) {
			if (touchesGrantRecords(requested, zone) || touchesGrantRecords(resolved, zone)) {
				return `protected test-authority record (${GRANT_DIR_REL}): grant records are issued by the main session and a background session cannot write them, their directory, or its parents`;
			}
		}
	}

	const matched =
		matchesAnyGlob(target, globs, projectRoot) ??
		matchesAnyGlob(requested, globs, projectRoot) ??
		(resolved === requested ? undefined : matchesAnyGlob(resolved, globs, projectRoot));
	if (!matched) return undefined;

	const blocked = (detail: string) => `protected test file (matched ${matched}, source: ${globSource}): ${detail}`;

	if (!eligible) {
		return blocked("preset acceptance tests are authored by Plan; if the task requires editing it, report STATUS: BLOCKED instead");
	}

	const state = context.readGrantState();
	if (state.status === "missing") {
		return blocked(`no test-author grant for this session instance (${state.reason})`);
	}
	if (state.status === "invalid") {
		return blocked(`test-author grant rejected: ${state.reason}`);
	}
	if (state.status === "unsupported") {
		return blocked("test-author grants are only consumable by a background native runner child session");
	}

	if (!isWithinPath(projectRoot, requested)) {
		return blocked("test-author grants only cover paths inside the project root");
	}
	if (!state.grant.paths.some((entry) => path.resolve(projectRoot, entry) === requested)) {
		return blocked("the test-author grant does not list this file");
	}
	try {
		if (statSync(requested).isDirectory()) {
			return blocked("test-author grants list files, not directories");
		}
	} catch {
		// Not created yet: the exact-file match above already decides.
	}
	if (audit.linkComponents.length > 0 || resolved !== requested) {
		return blocked("test-author grants do not extend through symbolic links or aliases");
	}
	return undefined;
}

/** PATH first (executable presence check), then <project>/bin/flowforge. */
export function resolveFlowForgeBinary(projectRoot: string): string {
	const pathEnv = process.env.PATH || "";
	for (const dir of pathEnv.split(path.delimiter)) {
		if (!dir) continue;
		const candidate = path.join(dir, "flowforge");
		if (existsSync(candidate)) return candidate;
	}
	const local = path.join(projectRoot, "bin", "flowforge");
	if (existsSync(local)) return local;
	return "flowforge";
}

interface CliToolResult {
	content: { type: "text"; text: string }[];
	details: Record<string, unknown>;
}

function runCli(binary: string, args: string[]): Promise<CliToolResult> {
	return new Promise((resolve) => {
		execFile(binary, args, { timeout: 60000, maxBuffer: 32 * 1024 * 1024 }, (error, stdout, stderr) => {
			const parts: string[] = [];
			// Cobra's Print family defaults to stderr when no output writer is
			// set, so FlowForge CLI output can arrive on either stream; always
			// surface both.
			if (stdout) parts.push(String(stdout));
			if (stderr) parts.push(String(stderr));
			if (error) {
				const code = typeof (error as NodeJS.ErrnoException).code === "string" ? (error as NodeJS.ErrnoException).code : String((error as { code?: unknown }).code ?? "");
				parts.push(`flowforge ${args[0]} failed: ${error.message}${code ? ` (code ${code})` : ""}`);
			}
			resolve({
				content: [{ type: "text", text: parts.join("\n") || "(no output)" }],
				details: { binary, args },
			});
		});
	});
}

export default function flowforgeExtension(pi: ExtensionAPI): void {
	const projectRoot = resolveThroughLinks(process.cwd());
	const config = readFlowForgeConfig(path.join(projectRoot, ".flowforge", "config.yaml"));
	const globs = config.testFileGlobs.length > 0 ? config.testFileGlobs : DEFAULT_TEST_FILE_GLOBS;
	const globSource = config.testFileGlobs.length > 0 ? "agents.test_file_globs (.flowforge/config.yaml)" : "FlowForge defaults";

	// Licensed test-author entry.
	//
	// A grant is a short-lived, exact-file record the *main* session writes with
	// its normal file tools; this extension never issues one. Eligibility is the
	// native-runner child handshake captured here at factory time and never
	// re-read from the process environment, because background runners share the
	// process. Session identity (`sessionId`) and the per-start `instanceId` come
	// from `session_start`, so a resumed/reloaded run never inherits an old grant.
	//
	// Boundaries: this is not an OS sandbox. `bash` stays unintercepted, a child
	// can only ever consume the exact files its grant lists, all background
	// sessions are denied write/edit on the grant directory and its ancestors,
	// and every read/parse failure falls back to the default test-write denial.
	const eligible = process.env.PI_SUBAGENT_CHILD === "1" && process.env.PI_SUBAGENT_RUNNER_PID === String(process.pid);
	const grantDir = path.join(projectRoot, GRANT_DIR_REL);
	const session: { sessionId?: string; instanceId?: string } = {};

	const grantPathFor = (sessionId: string | undefined): string =>
		sessionId ? path.join(grantDir, `${sessionId}.json`) : "";

	const readGrantState = (): GrantState => {
		if (!eligible) {
			return { status: "unsupported", reason: "only a background native runner child session can consume test-author grants" };
		}
		const { sessionId, instanceId } = session;
		if (!sessionId || !instanceId) {
			return { status: "missing", reason: "this session has no live instance identity" };
		}
		let raw: string;
		try {
			raw = readFileSync(grantPathFor(sessionId), "utf8");
		} catch {
			return { status: "missing", reason: "no grant record for this session" };
		}
		return evaluateTestAuthorGrant(raw, { sessionId, instanceId, projectRoot });
	};

	pi.on("session_start", (_event, ctx) => {
		let sessionId: string | undefined;
		try {
			const live = ctx?.sessionManager?.getSessionId?.();
			sessionId = typeof live === "string" && SAFE_SESSION_ID.test(live) ? live : undefined;
		} catch {
			sessionId = undefined;
		}
		session.sessionId = sessionId;
		session.instanceId = randomUUID();
	});

	pi.on("session_shutdown", () => {
		delete session.sessionId;
		delete session.instanceId;
	});

	pi.registerTool({
		name: "flowforge_test_authority",
		label: "FlowForge Test Authority",
		description:
			"Read-only status of this session's licensed test-author entry: session id, per-start instance id, canonical project root, grant record path, whether this session may consume a grant at all, and the current grant state. It reports only and cannot issue or widen a grant.",
		promptSnippet: "Read this session's test-author grant status",
		parameters: Type.Object({}),
		execute: async () => {
			const sessionId = session.sessionId ?? "";
			const instanceId = session.instanceId ?? "";
			const grantPath = grantPathFor(session.sessionId);
			const state = readGrantState();
			const grant =
				state.status === "valid"
					? { status: state.status, role: state.grant.role, paths: state.grant.paths, expiresAt: state.grant.expiresAt, reason: "" }
					: { status: state.status, role: "", paths: [] as string[], expiresAt: "", reason: state.reason };
			const details = { sessionId, instanceId, projectRoot, grantPath, eligible, grant };
			return {
				content: [{ type: "text" as const, text: JSON.stringify(details, null, 2) }],
				details,
			};
		},
	});

	if (!config.disableTestGuard) {
		pi.on("tool_call", (event: { toolName?: string; input?: { path?: unknown } }) => {
			if (event.toolName !== "write" && event.toolName !== "edit") return;
			const target = typeof event.input?.path === "string" ? event.input.path : "";
			if (!target) return;
			let reason: string | undefined;
			try {
				reason = testWriteDenial(target, { projectRoot, grantDir, globs, globSource, eligible, readGrantState });
			} catch (error) {
				const message = error instanceof Error ? error.message : String(error);
				return { block: true, reason: `test-author authority check failed (${message}): refusing the write` };
			}
			if (!reason) return;
			return { block: true, reason };
		});
	}

	if (!config.disableRawScriptModelGuard) {
		pi.on("tool_call", (event: { toolName?: string; input?: Record<string, unknown> }) => {
			if (event.toolName !== "subagent") return;
			const input = event.input ?? {};
			const rawScript = typeof input.workflowScript === "string" || typeof input.workflowScriptPath === "string";
			if (!rawScript) return; // named agent / workflow / management calls are not intercepted
			const model = typeof input.model === "string" ? input.model.trim() : "";
			if (model) return; // explicit model override is honored
			return {
				block: true,
				reason:
					"raw-script delegation without an explicit model inherits the orchestrating session's model and bypasses the flowforge per-role model pins. Fix one of: (a) prefer a named agent dispatch (agent: 'flowforge-<role>') which carries its pinned model; (b) pass model explicitly on this call AND on each runs.run() inside the script. Disable with agents.disable_raw_script_model_guard: true in .flowforge/config.yaml.",
			};
		});
	}

	const binary = resolveFlowForgeBinary(projectRoot);
	const proposalsDir = path.join(config.docsDir, "proposals");

	pi.registerTool({
		name: "flowforge_frontier",
		label: "FlowForge Frontier",
		description: "Compute unblocked, ready-to-execute FlowForge tickets as JSON (wraps `flowforge frontier --json`).",
		promptSnippet: "List ready FlowForge tickets",
		parameters: Type.Object({
			feature: Type.Optional(Type.String({ description: "Proposal slug to scope output to (graph stays global; pass only when working one proposal)" })),
		}),
		execute: async (args) => {
			const cliArgs = ["frontier", "--json"];
			if (args?.feature) {
				cliArgs.push("--feature", args.feature);
			}
			return runCli(binary, cliArgs);
		},
	});

	pi.registerTool({
		name: "flowforge_check",
		label: "FlowForge Check",
		description: `Validate the FlowForge proposal graph in ${proposalsDir} (wraps \`flowforge check --dir ${proposalsDir}\`).`,
		promptSnippet: "Validate the FlowForge proposal graph",
		parameters: Type.Object({
			feature: Type.Optional(Type.String({ description: "Proposal slug to scope reported diagnostics to (validity stays global)" })),
		}),
		execute: async (args) => {
			const cliArgs = ["check", "--dir", proposalsDir];
			if (args?.feature) {
				cliArgs.push("--feature", args.feature);
			}
			return runCli(binary, cliArgs);
		},
	});
}
