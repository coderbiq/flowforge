package tracker

import (
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"sort"
	"strings"
)

// Options carries catalog discovery configuration from project config.
type Options struct {
	// ExemptProposals lists proposal directory names whose tickets skip
	// evidence diagnostics entirely (legacy adoption path).
	ExemptProposals []string
}

var (
	changesSectionRe   = regexp.MustCompile(`(?mi)^##\s+Changes\s*$`)
	constraintsSection = regexp.MustCompile(`(?mi)^##\s+Constraints\s*$`)
	nextSectionRe      = regexp.MustCompile(`(?mi)^##\s+`)
	checkedItemRe      = regexp.MustCompile(`^(\s*)-\s+\[x\]\s+(.*)$`)
	uncheckedItemRe    = regexp.MustCompile(`^\s*-\s+\[\s\]\s+`)
	evidenceKeyRe      = regexp.MustCompile(`^\s*-\s+(cmd|exit|output|artifact):\s*(.*)$`)
	writeSetRe         = regexp.MustCompile(`(?mi)^.*Write set:.*$`)
)

// resolveArtifactBase walks upward from the discovery root to find the
// repository root (a directory containing .git or .flowforge). Evidence
// artifact paths are repository-relative; when no repository marker is
// found the discovery root itself is the base.
func resolveArtifactBase(root string) string {
	dir := root
	for {
		for _, marker := range []string{".git", ".flowforge"} {
			if info, err := os.Stat(filepath.Join(dir, marker)); err == nil && (info.IsDir() || marker == ".git") {
				return dir
			}
		}
		parent := filepath.Dir(dir)
		if parent == dir {
			return root
		}
		dir = parent
	}
}

// discoverEvidenceDiagnostics validates that every checked Change in the
// ticket's Changes section carries a command-anchored evidence quadruple.
// Pure documentation tickets (no Write set in Constraints) are skipped, and
// quadruples on unchecked items are legal intermediate state.
func discoverEvidenceDiagnostics(artifactBase, path, body string) []Diagnostic {
	changes := sectionBetween(body, changesSectionRe)
	if changes == "" {
		return nil
	}
	if !hasWriteSet(body) {
		return nil
	}

	var diagnostics []Diagnostic
	nonZeroByCmd := map[string]int{}
	lines := strings.Split(changes, "\n")
	for i := 0; i < len(lines); i++ {
		m := checkedItemRe.FindStringSubmatch(lines[i])
		if m == nil {
			continue
		}
		keys := map[string]string{}
		for j := i + 1; j < len(lines); j++ {
			line := lines[j]
			if checkedItemRe.MatchString(line) || uncheckedItemRe.MatchString(line) || strings.TrimSpace(line) == "" {
				if strings.TrimSpace(line) == "" {
					continue
				}
				break
			}
			if km := evidenceKeyRe.FindStringSubmatch(line); km != nil {
				keys[km[1]] = cleanEvidenceValue(km[2])
			}
		}
		if keys["cmd"] != "" && keys["exit"] != "" && keys["exit"] != "0" {
			nonZeroByCmd[normalizeEvidenceCmd(keys["cmd"])]++
		}
		diagnostics = append(diagnostics, evaluateEvidenceQuadruple(artifactBase, path, m[2], keys)...)
	}
	diagnostics = append(diagnostics, repeatFailureDiagnostics(path, nonZeroByCmd)...)
	return diagnostics
}

func evaluateEvidenceQuadruple(artifactBase, path, changeText string, keys map[string]string) []Diagnostic {
	if len(keys) == 0 {
		return []Diagnostic{warning(DiagnosticEvidenceMissing, path, "Checked Change has no evidence quadruple: "+truncateChange(changeText))}
	}
	var missing []string
	for _, key := range []string{"cmd", "exit", "output", "artifact"} {
		if keys[key] == "" {
			missing = append(missing, key)
		}
	}
	if len(missing) > 0 {
		return []Diagnostic{warning(DiagnosticEvidenceIncomplete, path, fmt.Sprintf("Evidence quadruple missing keys %v: %s", missing, truncateChange(changeText)))}
	}
	if keys["exit"] != "0" {
		return []Diagnostic{warning(DiagnosticEvidenceExitNonzero, path, "Evidence exit code is not 0, Change must stay unchecked: "+truncateChange(changeText))}
	}
	artifactPath := filepath.Join(artifactBase, keys["artifact"])
	if _, err := os.Stat(artifactPath); err != nil {
		return []Diagnostic{warning(DiagnosticEvidenceArtifactMissing, path, "Evidence artifact does not exist in repository: "+keys["artifact"])}
	}
	return nil
}

func sectionBetween(body string, heading *regexp.Regexp) string {
	loc := heading.FindStringIndex(body)
	if loc == nil {
		return ""
	}
	rest := body[loc[1]:]
	if next := nextSectionRe.FindStringIndex(rest); next != nil {
		rest = rest[:next[0]]
	}
	return rest
}

func hasWriteSet(body string) bool {
	constraints := sectionBetween(body, constraintsSection)
	return constraints != "" && writeSetRe.MatchString(constraints)
}

func cleanEvidenceValue(raw string) string {
	value := strings.TrimSpace(raw)
	value = strings.Trim(value, "`")
	value = strings.Trim(value, "\"")
	return strings.TrimSpace(value)
}

// evidenceRepeatFailureThreshold is the number of non-zero exits for one
// normalized command within a single ticket that triggers
// DiagnosticEvidenceRepeatFailure (research reference: opencode
// DOOM_LOOP_THRESHOLD=3, Gemini CLI aborts after 5 consecutive repeats).
const evidenceRepeatFailureThreshold = 3

// normalizeEvidenceCmd trims leading/trailing whitespace and collapses
// consecutive whitespace runs so trivially different spellings of one
// command count as the same command.
func normalizeEvidenceCmd(cmd string) string {
	return strings.Join(strings.Fields(cmd), " ")
}

// repeatFailureDiagnostics reports each command whose non-zero exits across
// one ticket's checked Changes reach the repeat-failure threshold. It
// complements DiagnosticEvidenceExitNonzero (single failed Change) with the
// repeated-failure pattern; both diagnostics can coexist.
func repeatFailureDiagnostics(path string, nonZeroByCmd map[string]int) []Diagnostic {
	var repeated []string
	for cmd, count := range nonZeroByCmd {
		if count >= evidenceRepeatFailureThreshold {
			repeated = append(repeated, cmd)
		}
	}
	sort.Strings(repeated)
	var diagnostics []Diagnostic
	for _, cmd := range repeated {
		diagnostics = append(diagnostics, warning(
			DiagnosticEvidenceRepeatFailure,
			path,
			fmt.Sprintf("Command reported non-zero exit %d times across checked Changes (threshold %d): %s", nonZeroByCmd[cmd], evidenceRepeatFailureThreshold, cmd),
		))
	}
	return diagnostics
}

func truncateChange(text string) string {
	text = strings.TrimSpace(text)
	runes := []rune(text)
	if len(runes) > 60 {
		return string(runes[:57]) + "..."
	}
	return text
}
