package uninstall

import (
	"fmt"
	"os"
	"path/filepath"

	"flowforge/internal/core"
)

type Result struct {
	Removed []string
	Errors  []error
}

func CleanBinary() (*Result, error) {
	result := &Result{}

	exe, err := os.Executable()
	if err != nil {
		return result, fmt.Errorf("finding executable: %w", err)
	}

	real, err := filepath.EvalSymlinks(exe)
	if err != nil {
		real = exe
	}

	if err := os.Remove(real); err != nil {
		if !os.IsNotExist(err) {
			result.Errors = append(result.Errors, fmt.Errorf("removing binary %s: %w", real, err))
		}
	} else {
		result.Removed = append(result.Removed, real)
	}

	oldPath := real + ".old"
	if err := os.Remove(oldPath); err != nil && !os.IsNotExist(err) {
		result.Errors = append(result.Errors, fmt.Errorf("removing backup %s: %w", oldPath, err))
	}

	return result, nil
}

func CleanConfig(homeDir string) (*Result, error) {
	result := &Result{}

	configDir := filepath.Join(homeDir, ".flowforge")
	if _, err := os.Stat(configDir); os.IsNotExist(err) {
		return result, nil
	}

	if err := os.RemoveAll(configDir); err != nil {
		result.Errors = append(result.Errors, fmt.Errorf("removing config dir: %w", err))
	} else {
		result.Removed = append(result.Removed, configDir)
	}

	return result, nil
}

func CleanProject(projectRoot string) (*Result, error) {
	result := &Result{}

	agentsBlockPath := filepath.Join(projectRoot, "AGENTS.md")
	if err := core.RemoveAgentsBlock(agentsBlockPath); err != nil {
		result.Errors = append(result.Errors, fmt.Errorf("removing agents block: %w", err))
	} else {
		result.Removed = append(result.Removed, "AGENTS.md FLOWFORGE block")
	}

	flowforgeDir := filepath.Join(projectRoot, ".flowforge")
	if _, err := os.Stat(flowforgeDir); err == nil {
		if err := os.RemoveAll(flowforgeDir); err != nil {
			result.Errors = append(result.Errors, fmt.Errorf("removing .flowforge dir: %w", err))
		} else {
			result.Removed = append(result.Removed, flowforgeDir)
		}
	}

	agentsSkills := filepath.Join(projectRoot, ".agents", "skills")
	if _, err := os.Stat(agentsSkills); err == nil {
		if err := os.RemoveAll(agentsSkills); err != nil {
			result.Errors = append(result.Errors, fmt.Errorf("removing .agents/skills: %w", err))
		} else {
			result.Removed = append(result.Removed, agentsSkills)
		}
	}

	agentsDir := filepath.Join(projectRoot, ".agents")
	if isEmptyDir(agentsDir) {
		os.Remove(agentsDir)
	}

	return result, nil
}

func isEmptyDir(path string) bool {
	entries, err := os.ReadDir(path)
	if err != nil {
		return false
	}
	return len(entries) == 0
}

func (r *Result) HasErrors() bool {
	return len(r.Errors) > 0
}

func (r *Result) ErrorSummary() string {
	if len(r.Errors) == 0 {
		return ""
	}
	var msg string
	for _, e := range r.Errors {
		msg += "  " + e.Error() + "\n"
	}
	return msg
}
