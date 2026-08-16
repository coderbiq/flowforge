package command

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/spf13/cobra"

	"flowforge/internal/core"
)

// HostEvidence is the read-only observation used by subagent status. It is
// deliberately separate from host intent: disk evidence never authorizes a
// write, enable, or deletion.
type HostEvidence struct {
	Host       string   `json:"host"`
	Detected   bool     `json:"detected"`
	Sources    []string `json:"sources,omitempty"`
	Intent     string   `json:"intent"`
	Registered int      `json:"registered"`
}

func newSubagentCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "subagent",
		Short: "Manage explicit subagent host authorization",
		Long:  "Manage explicit subagent host authorization. Host presence is observed separately from user authorization.",
	}
	cmd.AddCommand(newSubagentEnableCmd(), newSubagentDisableCmd(), newSubagentStatusCmd())
	return cmd
}

func newSubagentEnableCmd() *cobra.Command {
	var hosts []string
	var dryRun bool
	cmd := &cobra.Command{
		Use:   "enable",
		Short: "Explicitly authorize one or more subagent hosts",
		Long:  "Enable requires one or more explicit --host values (opencode or codex). It never infers hosts from disk. Host intent, rendered files, dynamic manifest entries, and the AGENTS orchestration block are committed only after every requested host reconciles successfully.",
		Args:  cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			if err := validateExplicitHosts(hosts); err != nil {
				return err
			}
			svc, err := currentConfigService()
			if err != nil {
				return fmt.Errorf("finding project root: %w", err)
			}
			defer svc.Close()
			return syncProject(cmd, svc.ProjectRoot(), syncOptions{forced: hosts, dryRun: dryRun, explicitEnable: true})
		},
	}
	cmd.Flags().StringSliceVar(&hosts, "host", nil, "Required host authorization (opencode,codex); repeatable")
	cmd.Flags().BoolVar(&dryRun, "dry-run", false, "Preview rendering and reconciliation without changing files or the manifest")
	return cmd
}

func newSubagentDisableCmd() *cobra.Command {
	var hosts []string
	var dryRun bool
	cmd := &cobra.Command{
		Use:   "disable",
		Short: "Explicitly remove authorization for registered host facilities",
		Long:  "Disable is explicit authorization to delete only manifest-registered host facilities; it never deletes unregistered files. Use --dry-run to preview the later backup and deletion plan.",
		Args:  cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			if len(hosts) > 0 {
				if err := validateHosts(hosts); err != nil {
					return err
				}
			}
			svc, err := currentConfigService()
			if err != nil {
				return fmt.Errorf("finding project root: %w", err)
			}
			defer svc.Close()
			plan, err := disableProject(svc.ProjectRoot(), hosts, dryRun)
			if err != nil {
				return err
			}
			if dryRun {
				fmt.Fprintf(cmd.OutOrStdout(), "Disable dry-run (backup %s):\n", plan.BackupDir)
				for _, status := range plan.Statuses {
					fmt.Fprintf(cmd.OutOrStdout(), "- %s\n", status)
				}
				fmt.Fprintln(cmd.OutOrStdout(), "No files or manifest entries were changed.")
				return nil
			}
			fmt.Fprintf(cmd.OutOrStdout(), "Disabled hosts: %s\n", strings.Join(plan.Hosts, ", "))
			for _, status := range plan.Statuses {
				fmt.Fprintf(cmd.OutOrStdout(), "- %s\n", status)
			}
			fmt.Fprintf(cmd.OutOrStdout(), "Backup: %s\n", plan.BackupDir)
			return nil
		},
	}
	cmd.Flags().StringSliceVar(&hosts, "host", nil, "Required host authorization to remove (opencode,codex); repeatable")
	cmd.Flags().BoolVar(&dryRun, "dry-run", false, "Preview without changing files or the manifest")
	return cmd
}

func newSubagentStatusCmd() *cobra.Command {
	var requested []string
	cmd := &cobra.Command{
		Use:   "status",
		Short: "Show host evidence and authorization without changing files",
		Long:  "Status is read-only: it reads the v2 manifest and host evidence from disk. It does not migrate, save intent, render host files, or edit AGENTS.md.",
		Args:  cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			if len(requested) > 0 {
				if err := validateHosts(requested); err != nil {
					return err
				}
			}
			svc, err := currentConfigService()
			if err != nil {
				return fmt.Errorf("finding project root: %w", err)
			}
			defer svc.Close()
			manifest, err := core.LoadProjectManifest(svc.ProjectRoot())
			if err != nil {
				return fmt.Errorf("loading project manifest read-only: %w", err)
			}
			evidence, err := DetectHostEvidence(svc.ProjectRoot(), manifest)
			if err != nil {
				return err
			}
			if len(requested) > 0 {
				wanted := make(map[string]bool, len(requested))
				for _, host := range requested {
					wanted[host] = true
				}
				filtered := evidence[:0]
				for _, item := range evidence {
					if wanted[item.Host] {
						filtered = append(filtered, item)
					}
				}
				evidence = filtered
			}
			if isJSONOutput(cmd) {
				data, err := json.Marshal(evidence)
				if err != nil {
					return fmt.Errorf("encoding subagent status: %w", err)
				}
				_, err = fmt.Fprintln(cmd.OutOrStdout(), string(data))
				return err
			}
			fmt.Fprintf(cmd.OutOrStdout(), "Subagent status (read-only, mode=%s)\n", manifest.Mode)
			for _, item := range evidence {
				fmt.Fprintf(cmd.OutOrStdout(), "- %s: intent=%s detected=%t registered=%d", item.Host, item.Intent, item.Detected, item.Registered)
				if len(item.Sources) > 0 {
					fmt.Fprintf(cmd.OutOrStdout(), " evidence=%s", strings.Join(item.Sources, ","))
				}
				fmt.Fprintln(cmd.OutOrStdout())
			}
			return nil
		},
	}
	cmd.Flags().StringSliceVar(&requested, "host", nil, "Optional host filter (opencode,codex)")
	return cmd
}

func validateExplicitHosts(hosts []string) error {
	if len(hosts) == 0 {
		return fmt.Errorf("at least one explicit --host is required (opencode or codex)")
	}
	return validateHosts(hosts)
}

func validateHosts(hosts []string) error {
	seen := make(map[string]bool, len(hosts))
	for _, host := range hosts {
		if !validHost(host) {
			return fmt.Errorf("unsupported host %q; choose opencode or codex", host)
		}
		if seen[host] {
			return fmt.Errorf("duplicate host %q", host)
		}
		seen[host] = true
	}
	return nil
}

// DetectHostEvidence only reads manifest data and known host paths. Callers
// must not use it as an implicit enable signal.
func DetectHostEvidence(root string, manifest *core.ProjectManifest) ([]HostEvidence, error) {
	if manifest == nil {
		return nil, fmt.Errorf("manifest is required")
	}
	registered := map[string]int{core.HostOpenCode: 0, core.HostCodex: 0}
	for _, entry := range manifest.Files {
		if entry.Type == "opencode_agent" {
			registered[core.HostOpenCode]++
		}
		if entry.Type == "codex_agent" {
			registered[core.HostCodex]++
		}
	}
	paths := map[string][]string{
		core.HostOpenCode: {".opencode", "opencode.json", "opencode.jsonc"},
		core.HostCodex:    {".codex", ".codex/config.toml"},
	}
	result := make([]HostEvidence, 0, 2)
	for _, host := range []string{core.HostOpenCode, core.HostCodex} {
		item := HostEvidence{Host: host, Intent: hostIntent(manifest, host), Registered: registered[host]}
		for _, path := range paths[host] {
			info, err := os.Stat(filepath.Join(root, path))
			if os.IsNotExist(err) {
				continue
			}
			if err != nil {
				return nil, fmt.Errorf("checking %s host evidence %s: %w", host, path, err)
			}
			if path == ".opencode" || path == ".codex" {
				if !info.IsDir() {
					return nil, fmt.Errorf("invalid %s host evidence %s", host, path)
				}
			}
			item.Detected = true
			item.Sources = append(item.Sources, path)
		}
		result = append(result, item)
	}
	return result, nil
}

func hostIntent(manifest *core.ProjectManifest, host string) string {
	if host == core.HostOpenCode {
		return manifest.HostIntent.OpenCode
	}
	return manifest.HostIntent.Codex
}
