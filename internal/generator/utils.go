package generator

import (
	"go/build"
	"os"
	"path/filepath"
	"strings"
)

// findModuleRoot finds the root directory of the Go module
func findModuleRoot(dir string) string {
	for {
		if _, err := os.Stat(filepath.Join(dir, "go.mod")); err == nil {
			return dir
		}
		parent := filepath.Dir(dir)
		if parent == dir {
			// Reached filesystem root without finding go.mod
			return ""
		}
		dir = parent
	}
}

// getPackagePath determines the package path relative to the module root
func getPackagePath(moduleRoot, packageDir string) string {
	if moduleRoot == "" {
		// Fallback to GOPATH-style package detection
		gopath := os.Getenv("GOPATH")
		if gopath != "" {
			for _, path := range filepath.SplitList(gopath) {
				srcDir := filepath.Join(path, "src")
				if rel, err := filepath.Rel(srcDir, packageDir); err == nil && !strings.HasPrefix(rel, "..") {
					return filepath.ToSlash(rel)
				}
			}
		}

		// Try to use go/build to get package info
		if pkg, err := build.ImportDir(packageDir, build.FindOnly); err == nil {
			return pkg.ImportPath
		}

		return ""
	}

	// Get module name from go.mod
	moduleName := getModuleName(moduleRoot)
	if moduleName == "" {
		return ""
	}

	// Get relative path from module root
	relPath, err := filepath.Rel(moduleRoot, packageDir)
	if err != nil {
		return ""
	}

	if relPath == "." {
		return moduleName
	}

	return moduleName + "/" + filepath.ToSlash(relPath)
}

// getModuleName extracts the module name from go.mod file
func getModuleName(moduleRoot string) string {
	goModPath := filepath.Join(moduleRoot, "go.mod")
	content, err := os.ReadFile(goModPath)
	if err != nil {
		return ""
	}

	lines := strings.Split(string(content), "\n")
	for _, line := range lines {
		line = strings.TrimSpace(line)
		if strings.HasPrefix(line, "module ") {
			parts := strings.Fields(line)
			if len(parts) >= 2 {
				return parts[1]
			}
		}
	}

	return ""
}

// getCurrentPackagePath tries to determine the current package path
// from the current working directory
func getCurrentPackagePath() string {
	pwd, err := os.Getwd()
	if err != nil {
		return ""
	}

	moduleRoot := findModuleRoot(pwd)
	return getPackagePath(moduleRoot, pwd)
}
