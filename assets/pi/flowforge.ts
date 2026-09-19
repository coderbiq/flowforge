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
 * 2. Native LLM tools `flowforge_frontier` / `flowforge_check` wrapping the
 *    FlowForge CLI (PATH first, then <project>/bin/flowforge).
 */

import { execFile } from "node:child_process";
import { existsSync, readFileSync } from "node:fs";
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

interface FlowForgeProjectConfig {
	docsDir: string;
	testFileGlobs: string[];
	disableTestGuard: boolean;
}

/**
 * Minimal line-level reader for the flat `.flowforge/config.yaml` subset
 * FlowForge needs: top-level `docs_dir`, and inside the `agents:` block
 * `test_file_globs` (block or inline list) and `disable_test_guard`.
 * Missing file, missing `agents:` key, or missing keys fall back to
 * defaults without error.
 */
export function readFlowForgeConfig(configPath: string): FlowForgeProjectConfig {
	const result: FlowForgeProjectConfig = { docsDir: "docs", testFileGlobs: [], disableTestGuard: false };
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
	const projectRoot = process.cwd();
	const config = readFlowForgeConfig(path.join(projectRoot, ".flowforge", "config.yaml"));
	const globs = config.testFileGlobs.length > 0 ? config.testFileGlobs : DEFAULT_TEST_FILE_GLOBS;
	const globSource = config.testFileGlobs.length > 0 ? "agents.test_file_globs (.flowforge/config.yaml)" : "FlowForge defaults";

	if (!config.disableTestGuard) {
		pi.on("tool_call", (event: { toolName?: string; input?: { path?: unknown } }) => {
			if (event.toolName !== "write" && event.toolName !== "edit") return;
			const target = typeof event.input?.path === "string" ? event.input.path : "";
			if (!target) return;
			const matched = matchesAnyGlob(target, globs, projectRoot);
			if (!matched) return;
			return {
				block: true,
				reason: `protected test file (matched ${matched}, source: ${globSource}): preset acceptance tests are authored by Plan; if the task requires editing it, report STATUS: BLOCKED instead`,
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
		parameters: Type.Object({}),
		execute: async () => runCli(binary, ["frontier", "--json"]),
	});

	pi.registerTool({
		name: "flowforge_check",
		label: "FlowForge Check",
		description: `Validate the FlowForge proposal graph in ${proposalsDir} (wraps \`flowforge check --dir ${proposalsDir}\`).`,
		promptSnippet: "Validate the FlowForge proposal graph",
		parameters: Type.Object({}),
		execute: async () => runCli(binary, ["check", "--dir", proposalsDir]),
	});
}
