package generator

import (
	"fmt"
	"io/fs"
	"path/filepath"
	"sort"
	"strings"
)

// Generator handles the code generation process
type Generator struct {
	verbose         bool
	explicitPkgPath string
	packageDataMap  map[string]*PackageData
}

// New creates a new Generator instance
func New(verbose bool, explicitPkgPath string) *Generator {
	return &Generator{
		verbose:         verbose,
		explicitPkgPath: explicitPkgPath,
		packageDataMap:  make(map[string]*PackageData),
	}
}

// ScanDirectory scans a directory for Go files and collects generation data.
func (g *Generator) ScanDirectory(targetDir string) (int, int, error) {
	files := make([]string, 0)
	err := filepath.WalkDir(targetDir, func(path string, d fs.DirEntry, walkErr error) error {
		if walkErr != nil {
			return walkErr
		}
		if d.IsDir() {
			if isIgnoredDir(d.Name()) {
				return filepath.SkipDir
			}
			return nil
		}
		if isTargetGoFile(path) {
			files = append(files, path)
		}
		return nil
	})
	if err != nil {
		return 0, 0, err
	}

	sort.Strings(files)

	fileCount := 0
	generateCommentCount := 0
	for _, path := range files {
		fileCount++
		if g.verbose {
			fmt.Printf("Processing file: %s\n", path)
		}
		commentCount := g.processFile(path)
		generateCommentCount += commentCount
	}

	return fileCount, generateCommentCount, nil
}

func isIgnoredDir(name string) bool {
	switch name {
	case ".git", "vendor":
		return true
	default:
		return false
	}
}

func isTargetGoFile(path string) bool {
	if !strings.HasSuffix(path, ".go") {
		return false
	}
	if strings.HasSuffix(path, "_test.go") {
		return false
	}
	if strings.HasSuffix(path, "registry_init.go") {
		return false
	}
	return true
}

// GenerateRegistries generates registry files for all collected packages.
func (g *Generator) GenerateRegistries() (int, error) {
	count := 0
	for packageDir, data := range g.packageDataMap {
		if len(data.Registrations) > 0 {
			if err := g.generateRegistryForPackage(data, packageDir); err != nil {
				return count, err
			}
			count++
		}
	}
	return count, nil
}
