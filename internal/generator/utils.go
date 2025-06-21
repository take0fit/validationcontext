package generator

import (
	"os"
	"path/filepath"
	"strings"
)

// findModuleRoot finds the root directory containing go.mod
func findModuleRoot(dir string) string {
	for {
		if _, err := os.Stat(filepath.Join(dir, "go.mod")); err == nil {
			return dir
		}
		parent := filepath.Dir(dir)
		if parent == dir {
			break
		}
		dir = parent
	}
	return ""
}

// getPackagePath constructs the full package path from module root and package directory
func getPackagePath(moduleRoot, packageDir string) string {
	if moduleRoot == "" {
		return ""
	}

	rel, err := filepath.Rel(moduleRoot, packageDir)
	if err != nil {
		return ""
	}

	moduleName := getModuleName(moduleRoot)
	if moduleName == "" {
		return ""
	}

	if rel == "." {
		return moduleName
	}

	return moduleName + "/" + filepath.ToSlash(rel)
}

// getModuleName extracts module name from go.mod file
func getModuleName(moduleRoot string) string {
	modPath := filepath.Join(moduleRoot, "go.mod")
	content, err := os.ReadFile(modPath)
	if err != nil {
		return ""
	}

	lines := strings.Split(string(content), "\n")
	for _, line := range lines {
		line = strings.TrimSpace(line)
		if strings.HasPrefix(line, "module ") {
			return strings.TrimSpace(strings.TrimPrefix(line, "module"))
		}
	}
	return ""
}
