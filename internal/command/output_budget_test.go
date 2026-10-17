package command

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

// budgetFixture builds an S1-shaped notebook (197+ tickets across two
// features): mostly ready work with complete execution contracts, a small
// blocked chain, and closed tickets carrying completion evidence. Nothing in
// it generates gap noise -- the assertion is purely about channel volume.
func budgetFixture(t *testing.T) string {
	t.Helper()
	root := t.TempDir()
	const fm = "---\nflowforge:\n  schema: 1\n  role: ticket\n---\n"
	execContract := "\n\n## Execution detail\n\n### Verified contracts\n\n- output_budget.go owns the contract.\n\n### Execution scenarios\n\n- Success: channel within budget.\n- Failure: test fails with byte count.\n\n### Expected tests\n\n- go test ./internal/command/...\n\n### Generated artifacts\n\n- Not applicable.\n\n### Conventions\n\n- Files own content.\n"
	evidence := "\n## Completion evidence\n\n- Delivered earlier; evidence on file.\n"

	writeTicket := func(feature, id string, blockedBy, status string) {
		issues := filepath.Join(root, feature, "issues")
		if err := os.MkdirAll(issues, 0o755); err != nil {
			t.Fatal(err)
		}
		body := fm + "# " + id + ": budget fixture ticket\n\n**Status:** " + status + "\n"
		if blockedBy != "" {
			body += "**Blocked by:** " + blockedBy + "\n"
		}
		body += execContract
		if status == "closed" {
			body += evidence
		}
		if err := os.WriteFile(filepath.Join(issues, id+".md"), []byte(body), 0o644); err != nil {
			t.Fatal(err)
		}
	}

	// 120 ready tickets (60 per feature), independent.
	for i := 0; i < 60; i++ {
		writeTicket("feat-a", padID(i), "", "open")
		writeTicket("feat-b", padID(i), "", "open")
	}
	// 15 blocked tickets in feat-a chained on existing ready ids.
	for i := 0; i < 15; i++ {
		writeTicket("feat-a", "b"+padID(i), padID(i), "open")
	}
	// 65 closed tickets with completion evidence.
	for i := 0; i < 65; i++ {
		writeTicket("feat-b", "c"+padID(i), "", "closed")
	}
	return root
}

func padID(i int) string {
	return strings.Repeat("0", 2-len(fmt2(i))) + fmt2(i)
}

func fmt2(i int) string {
	if i < 10 {
		return string(rune('0' + i))
	}
	tens, ones := i/10, i%10
	return string(rune('0'+tens)) + string(rune('0'+ones))
}

// captureAgentFacingChannel runs one command with the given flags and
// returns the combined stdout+stderr byte count: the number an agent host
// actually receives. Real stdio is piped (not cmd.SetOut) because the text
// sections print via fmt.Printf, and package flag vars are assigned AFTER the
// command is created -- flag registration resets them to defaults.
func captureAgentFacingChannel(t *testing.T, root string, kind string, jsonMode bool) int {
	t.Helper()
	var run func() error
	if kind == "frontier" {
		cmd := newFrontierCmd()
		frontierDir, frontierJSON, frontierFeature, frontierFull = root, jsonMode, "", false
		run = func() error { return cmd.RunE(cmd, nil) }
	} else {
		cmd := newCheckCmd()
		checkDir, checkJSON, checkFeature, checkFull = root, jsonMode, "", false
		run = func() error { return cmd.RunE(cmd, nil) }
	}

	reader, writer, err := os.Pipe()
	if err != nil {
		t.Fatal(err)
	}
	savedStdout, savedStderr := os.Stdout, os.Stderr
	os.Stdout, os.Stderr = writer, writer
	runErr := run()
	writer.Close()
	os.Stdout, os.Stderr = savedStdout, savedStderr
	data := make([]byte, 0, 2048)
	buf := make([]byte, 4096)
	for {
		n, readErr := reader.Read(buf)
		data = append(data, buf[:n]...)
		if readErr != nil {
			break
		}
	}
	if runErr != nil {
		t.Fatalf("%s failed: %v\noutput so far: %s", kind, runErr, data)
	}
	return len(data)
}

func TestMaxAgentOutputBytesLocked(t *testing.T) {
	if MaxAgentOutputBytes != 16*1024 {
		t.Fatalf("MaxAgentOutputBytes drifted to %d; changing it requires a design revision", MaxAgentOutputBytes)
	}
}

func TestAgentFacingChannelsWithinBudget(t *testing.T) {
	root := budgetFixture(t)
	// Agent-facing entries carry the --dir prefix; chdir into the fixture and
	// scan "." so entry paths stay representative instead of absolute-length.
	previousWD, err := os.Getwd()
	if err != nil {
		t.Fatal(err)
	}
	if err := os.Chdir(root); err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() { _ = os.Chdir(previousWD) })

	channels := []struct {
		kind     string
		jsonMode bool
	}{
		{"frontier", false},
		{"frontier", true},
		{"check", false},
		{"check", true},
	}
	for _, channel := range channels {
		got := captureAgentFacingChannel(t, ".", channel.kind, channel.jsonMode)
		if got > MaxAgentOutputBytes {
			t.Fatalf("%s (json=%v) emitted %d bytes > MaxAgentOutputBytes %d: regression in output convergence", channel.kind, channel.jsonMode, got, MaxAgentOutputBytes)
		}
		t.Logf("%s json=%v: %d bytes", channel.kind, channel.jsonMode, got)
	}
}
