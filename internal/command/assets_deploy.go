package command

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"runtime"
)

func deployManagedAssets(targetDir string) error {
	assetsDir, err := locateAssetsDir()
	if err != nil {
		return err
	}

	if err := copyDir(filepath.Join(assetsDir, "skills"), filepath.Join(targetDir, ".agents", "skills"), true); err != nil {
		return fmt.Errorf("deploying skills: %w", err)
	}

	if err := copyDir(filepath.Join(assetsDir, "templates"), filepath.Join(targetDir, ".flowforge", "templates"), true); err != nil {
		return fmt.Errorf("deploying templates: %w", err)
	}

	agentRules := filepath.Join(assetsDir, "AGENTS.md")
	if _, err := os.Stat(agentRules); err == nil {
		targetPath := filepath.Join(targetDir, "AGENTS.md")
		if _, err := os.Stat(targetPath); os.IsNotExist(err) {
			if err := copyFile(agentRules, targetPath, false); err != nil {
				return fmt.Errorf("deploying AGENTS.md: %w", err)
			}
		} else if err != nil {
			return fmt.Errorf("checking target AGENTS.md: %w", err)
		}
	} else if !os.IsNotExist(err) {
		return fmt.Errorf("checking asset AGENTS.md: %w", err)
	}

	return nil
}

func locateAssetsDir() (string, error) {
	var candidates []string

	if executable, err := os.Executable(); err == nil {
		exeDir := filepath.Dir(executable)
		candidates = append(candidates,
			filepath.Join(exeDir, "assets"),
			filepath.Join(exeDir, "..", "assets"),
		)
	}

	if _, file, _, ok := runtime.Caller(0); ok {
		repoRoot := filepath.Clean(filepath.Join(filepath.Dir(file), "..", ".."))
		candidates = append(candidates, filepath.Join(repoRoot, "assets"))
	}

	if wd, err := os.Getwd(); err == nil {
		candidates = append(candidates, filepath.Join(wd, "assets"))
	}

	for _, candidate := range candidates {
		if isAssetsDir(candidate) {
			return candidate, nil
		}
	}

	return "", fmt.Errorf("flowforge assets not found; expected assets next to the executable or in the source checkout")
}

func isAssetsDir(path string) bool {
	info, err := os.Stat(filepath.Join(path, "skills"))
	if err != nil || !info.IsDir() {
		return false
	}
	return true
}

func copyDir(srcDir, dstDir string, overwrite bool) error {
	info, err := os.Stat(srcDir)
	if os.IsNotExist(err) {
		return nil
	}
	if err != nil {
		return fmt.Errorf("checking source directory %s: %w", srcDir, err)
	}
	if !info.IsDir() {
		return fmt.Errorf("source path is not a directory: %s", srcDir)
	}

	return filepath.WalkDir(srcDir, func(path string, entry os.DirEntry, walkErr error) error {
		if walkErr != nil {
			return walkErr
		}
		if entry.Name() == ".gitkeep" {
			return nil
		}

		rel, err := filepath.Rel(srcDir, path)
		if err != nil {
			return err
		}
		targetPath := filepath.Join(dstDir, rel)

		if entry.IsDir() {
			return os.MkdirAll(targetPath, 0755)
		}

		return copyFile(path, targetPath, overwrite)
	})
}

func copyFile(srcPath, dstPath string, overwrite bool) error {
	if !overwrite {
		if _, err := os.Stat(dstPath); err == nil {
			return nil
		} else if !os.IsNotExist(err) {
			return fmt.Errorf("checking target file %s: %w", dstPath, err)
		}
	}

	src, err := os.Open(srcPath)
	if err != nil {
		return fmt.Errorf("opening source file %s: %w", srcPath, err)
	}
	defer func() {
		if closeErr := src.Close(); closeErr != nil {
			fmt.Fprintf(os.Stderr, "warning: closing source file %s: %v\n", srcPath, closeErr)
		}
	}()

	if err := os.MkdirAll(filepath.Dir(dstPath), 0755); err != nil {
		return fmt.Errorf("creating target directory: %w", err)
	}

	dst, err := os.OpenFile(dstPath, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0644)
	if err != nil {
		return fmt.Errorf("opening target file %s: %w", dstPath, err)
	}

	if _, err := io.Copy(dst, src); err != nil {
		if closeErr := dst.Close(); closeErr != nil {
			return fmt.Errorf("copying file %s: %w (closing target: %v)", dstPath, err, closeErr)
		}
		return fmt.Errorf("copying file %s: %w", dstPath, err)
	}

	if err := dst.Close(); err != nil {
		return fmt.Errorf("closing target file %s: %w", dstPath, err)
	}

	return nil
}
