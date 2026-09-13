package tracker

import (
	"fmt"
	"os"
	"path/filepath"
	"regexp"
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
		diagnostics = append(diagnostics, evaluateEvidenceQuadruple(artifactBase, path, m[2], keys)...)
	}
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

func truncateChange(text string) string {
	text = strings.TrimSpace(text)
	runes := []rune(text)
	if len(runes) > 60 {
		return string(runes[:57]) + "..."
	}
	return text
}
