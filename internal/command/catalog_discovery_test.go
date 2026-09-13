package command

import (
	"os"
	"path/filepath"
	"testing"

	"flowforge/internal/config"
	"flowforge/internal/tracker"
)

func TestEvidenceOptionsFlowFromProjectConfig(t *testing.T) {
	projectRoot := t.TempDir()
	if err := initializeTestProject(projectRoot); err != nil {
		t.Fatal(err)
	}

	proposals := filepath.Join(projectRoot, "docs", "proposals", "legacy-proposal", "issues")
	if err := os.MkdirAll(proposals, 0755); err != nil {
		t.Fatal(err)
	}
	ticket := "---\nflowforge:\n  schema: 1\n  role: ticket\n---\n# 01: x\n\n**Status:** closed\n\n## Changes\n\n- [x] 1. Done without evidence\n\n## Constraints\n\n- Write set: app/\n"
	if err := os.WriteFile(filepath.Join(proposals, "01-x.md"), []byte(ticket), 0644); err != nil {
		t.Fatal(err)
	}

	scanDir := filepath.Join(projectRoot, "docs", "proposals")

	catalog, err := discoverProposalCatalog(scanDir)
	if err != nil {
		t.Fatal(err)
	}
	if !hasEvidenceCode(catalog.Diagnostics, "evidence-missing") {
		t.Fatal("expected evidence-missing without exemption config")
	}

	cfg, err := config.Load(projectRoot)
	if err != nil {
		t.Fatal(err)
	}
	cfg.Evidence.ExemptProposals = []string{"legacy-proposal"}
	if err := cfg.Save(projectRoot); err != nil {
		t.Fatal(err)
	}

	// discoverProposalCatalog resolves the project root from the working
	// directory, so run from inside the temp project.
	wd, err := os.Getwd()
	if err != nil {
		t.Fatal(err)
	}
	if err := os.Chdir(projectRoot); err != nil {
		t.Fatal(err)
	}
	defer func() { _ = os.Chdir(wd) }()

	catalog, err = discoverProposalCatalog(scanDir)
	if err != nil {
		t.Fatal(err)
	}
	if hasEvidenceCode(catalog.Diagnostics, "evidence-missing") {
		t.Fatal("exempt proposal must not report evidence diagnostics")
	}
}

func hasEvidenceCode(diagnostics []tracker.Diagnostic, code string) bool {
	for _, d := range diagnostics {
		if string(d.Code) == code {
			return true
		}
	}
	return false
}
