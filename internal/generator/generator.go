package generator

import (
	"fmt"
	"os"
	"path/filepath"
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

// ScanDirectory scans a directory for Go files and collects generation data
func (g *Generator) ScanDirectory(targetDir string) (int, int, error) {
	fileCount := 0
	generateCommentCount := 0

	err := filepath.Walk(targetDir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if strings.HasSuffix(path, ".go") &&
			!strings.HasSuffix(path, "_test.go") &&
			!strings.HasSuffix(path, "registry_init.go") {
			fileCount++
			if g.verbose {
				fmt.Printf("Processing file: %s\n", path)
			}

			commentCount := g.processFile(path)
			generateCommentCount += commentCount
		}
		return nil
	})

	return fileCount, generateCommentCount, err
}

// GenerateRegistries generates registry files for all collected packages
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
