package command

import (
	"fmt"
	"io"
	"io/fs"
	"os"
	"path/filepath"
	"runtime"
)

func deployManagedAssets(targetDir string) error {
	assetsDir, cleanup, err := locateAssetsDir()
	if err != nil {
		return err
	}
	defer cleanup()

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

// locateAssetsDir returns the path to the assets directory and a cleanup function.
// It first tries to extract assets from the embedded filesystem (for standalone binaries).
// If that fails, it falls back to filesystem-based lookup (for development).
func locateAssetsDir() (string, func(), error) {
	noop := func() {}

	// Try embedded assets first (standalone binary)
	if dir, err := extractEmbeddedAssets(); err == nil {
		return dir, func() { os.RemoveAll(dir) }, nil
	}

	// Fallback: filesystem-based lookup (development)
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
			return candidate, noop, nil
		}
	}

	return "", noop, fmt.Errorf("flowforge assets not found; expected assets next to the executable or in the source checkout")
}

// extractEmbeddedAssets extracts the embedded assets filesystem to a temporary directory.
func extractEmbeddedAssets() (string, error) {
	tmpDir, err := os.MkdirTemp("", "flowforge-assets-")
	if err != nil {
		return "", fmt.Errorf("creating temp dir for embedded assets: %w", err)
	}

	if err := fs.WalkDir(embeddedAssets, "assets", func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}

		rel, err := filepath.Rel("assets", path)
		if err != nil {
			return err
		}

		targetPath := filepath.Join(tmpDir, rel)

		if d.IsDir() {
			return os.MkdirAll(targetPath, 0755)
		}

		data, err := fs.ReadFile(embeddedAssets, path)
		if err != nil {
			return fmt.Errorf("reading embedded file %s: %w", path, err)
		}

		if err := os.MkdirAll(filepath.Dir(targetPath), 0755); err != nil {
			return fmt.Errorf("creating directory for %s: %w", targetPath, err)
		}

		if err := os.WriteFile(targetPath, data, 0644); err != nil {
			return fmt.Errorf("writing %s: %w", targetPath, err)
		}

		return nil
	}); err != nil {
		os.RemoveAll(tmpDir)
		return "", fmt.Errorf("extracting embedded assets: %w", err)
	}

	return tmpDir, nil
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
