package command

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"time"

	"github.com/spf13/cobra"

	"flowforge/internal/config"
	"flowforge/internal/tracker"
)

var (
	proposalDir string
)

// newProposalCmd builds the `flowforge proposal` command group with the
// close/reopen lifecycle subcommands. Design authority:
// ff-wiki/proposals/roadmap-management/design.md (D4).
func newProposalCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "proposal",
		Short: "Manage proposal lifecycle (close/reopen with library-wide link rewriting)",
	}
	cmd.AddCommand(newProposalCloseCmd(), newProposalReopenCmd())
	return cmd
}

func newProposalCloseCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "close <slug>",
		Short: "Archive a finished proposal into proposals/closed/<year>/ and rewrite library links",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			return runProposalClose(cmd, args[0])
		},
	}
	cmd.Flags().StringVarP(&proposalDir, "dir", "d", "", "Proposals directory to operate on (default: <docs_dir>/proposals)")
	return cmd
}

func newProposalReopenCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "reopen <slug>",
		Short: "Move an archived proposal back to proposals/ and rewrite library links",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			return runProposalReopen(cmd, args[0])
		},
	}
	cmd.Flags().StringVarP(&proposalDir, "dir", "d", "", "Proposals directory to operate on (default: <docs_dir>/proposals)")
	return cmd
}

// proposalsDirOrDefault resolves the proposals root from --dir or the config
// discovery used by status/frontier (internal/command/status.go:20-26).
func proposalsDirOrDefault() (string, error) {
	if proposalDir != "" {
		return proposalDir, nil
	}
	return config.ResolveProposalsDir(".")
}

// runProposalClose implements `flowforge proposal close <slug>`:
// locate the live directory, verify all tickets are terminal, move it under
// proposals/closed/<year>/, then rewrite library-wide references.
func runProposalClose(cmd *cobra.Command, slug string) error {
	proposalsRoot, err := proposalsDirOrDefault()
	if err != nil {
		return err
	}
	slug = filepath.Base(filepath.Clean(slug))

	liveDir := filepath.Join(proposalsRoot, slug)
	if info, err := os.Stat(liveDir); err != nil || !info.IsDir() {
		return fmt.Errorf("live proposal not found: %s", liveDir)
	}
	archiveRoot := filepath.Join(proposalsRoot, "closed")
	if info, err := os.Stat(archiveRoot); err == nil && !info.IsDir() {
		return fmt.Errorf("closed path is not a directory: %s", archiveRoot)
	}

	if err := requireAllTicketsTerminal(liveDir, slug); err != nil {
		return err
	}

	year := time.Now().UTC().Year()
	yearDir := filepath.Join(archiveRoot, fmt.Sprintf("%d", year))
	target := filepath.Join(yearDir, slug)
	if _, err := os.Stat(target); err == nil {
		return fmt.Errorf("archive target already exists: %s", target)
	}
	if err := os.MkdirAll(yearDir, 0o755); err != nil {
		return fmt.Errorf("creating archive directory %s: %w", yearDir, err)
	}
	if err := os.Rename(liveDir, target); err != nil {
		return fmt.Errorf("moving %s to %s: %w", liveDir, target, err)
	}

	docsDir := filepath.Dir(proposalsRoot)
	rewritten, err := RewriteReferences(docsDir, "proposals/"+slug, "proposals/closed/"+fmt.Sprintf("%d", year)+"/"+slug)
	if err != nil {
		return err
	}

	cmd.Printf("Closed %s → %s (tickets: all terminal; links rewritten in %d file(s))\n", liveDir, target, rewritten)
	return nil
}

// runProposalReopen implements `flowforge proposal reopen <slug>`: find the
// unique archived copy under closed/*/, refuse a live-name collision, move it
// back to the proposals root, then rewrite library-wide references.
func runProposalReopen(cmd *cobra.Command, slug string) error {
	proposalsRoot, err := proposalsDirOrDefault()
	if err != nil {
		return err
	}
	slug = filepath.Base(filepath.Clean(slug))

	liveDir := filepath.Join(proposalsRoot, slug)
	if info, err := os.Stat(liveDir); err == nil {
		if !info.IsDir() {
			return fmt.Errorf("live path is not a directory: %s", liveDir)
		}
		return fmt.Errorf("live proposal already exists: %s", liveDir)
	}

	archiveRoot := filepath.Join(proposalsRoot, "closed")
	matches, err := findArchivedProposals(archiveRoot, slug)
	if err != nil {
		return err
	}
	if len(matches) == 0 {
		return fmt.Errorf("no archived proposal %q under %s", slug, archiveRoot)
	}
	if len(matches) > 1 {
		return fmt.Errorf("multiple archived proposals named %q: %s", slug, strings.Join(matches, ", "))
	}
	source := matches[0]

	docsDir := filepath.Dir(proposalsRoot)
	relSource, err := filepath.Rel(docsDir, source)
	if err != nil {
		return fmt.Errorf("resolving %s: %w", source, err)
	}
	if err := os.Rename(source, liveDir); err != nil {
		return fmt.Errorf("moving %s to %s: %w", source, liveDir, err)
	}

	rewritten, err := RewriteReferences(docsDir, filepath.ToSlash(relSource), "proposals/"+slug)
	if err != nil {
		return err
	}

	cmd.Printf("Reopened %s → %s (links rewritten in %d file(s))\n", source, liveDir, rewritten)
	return nil
}

// requireAllTicketsTerminal fails when any ticket discovered under dir is not
// in a terminal status, listing every offender.
func requireAllTicketsTerminal(dir, slug string) error {
	issues, err := tracker.DiscoverIssues(dir)
	if err != nil {
		return fmt.Errorf("discovering issues in %s: %w", dir, err)
	}
	var open []string
	for _, issue := range issues {
		if issue.Status.IsTerminal() {
			continue
		}
		rel, err := filepath.Rel(dir, issue.FilePath)
		if err != nil {
			rel = issue.FilePath
		}
		open = append(open, fmt.Sprintf("%s (%s)", rel, issue.Status))
	}
	if len(open) == 0 {
		return nil
	}
	sort.Strings(open)
	return fmt.Errorf("proposal %q has %d non-terminal ticket(s): %s", slug, len(open), strings.Join(open, ", "))
}

// findArchivedProposals returns every directory named slug directly under a
// year directory of <archiveRoot> (closed/<year>/<slug>).
func findArchivedProposals(archiveRoot, slug string) ([]string, error) {
	entries, err := os.ReadDir(archiveRoot)
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			return nil, nil
		}
		return nil, fmt.Errorf("reading %s: %w", archiveRoot, err)
	}
	var matches []string
	for _, entry := range entries {
		if !entry.IsDir() {
			continue
		}
		candidate := filepath.Join(archiveRoot, entry.Name(), slug)
		if info, err := os.Stat(candidate); err == nil && info.IsDir() {
			matches = append(matches, candidate)
		}
	}
	sort.Strings(matches)
	return matches, nil
}
