package generator

import (
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"log"
	"path/filepath"
	"regexp"
	"strings"
)

// processFile processes a single Go file and collects registration data
func (g *Generator) processFile(filename string) int {
	fset := token.NewFileSet()
	file, err := parser.ParseFile(fset, filename, nil, parser.ParseComments)
	if err != nil {
		log.Printf("failed to parse file %s: %v", filename, err)
		return 0
	}

	packageDir := filepath.Dir(filename)
	generateCommentCount := 0
	var config *GenerateConfig

	// Look for //go:generate comments
	for _, commentGroup := range file.Comments {
		for _, comment := range commentGroup.List {
			if g.verbose {
				fmt.Printf("  Comment: %s\n", comment.Text)
			}

			if strings.Contains(comment.Text, "//go:generate voauto-gen") ||
				strings.Contains(comment.Text, "//go:generate validationcontext") {
				generateCommentCount++
				fmt.Printf("Found generate comment in %s: %s\n", filename, comment.Text)

				config = parseGenerateComment(comment.Text)
				break
			}
		}
		if config != nil {
			break
		}
	}

	if config != nil {
		g.ensurePackageData(packageDir, file, config)
		packageData := g.packageDataMap[packageDir]
		registrations := collectRegistrationsFromFile(file, config, packageData.PackagePath, g.verbose)
		packageData.Registrations = append(packageData.Registrations, registrations...)

		if g.verbose {
			fmt.Printf("  Added %d registrations from %s\n", len(registrations), filename)
		}
	}

	return generateCommentCount
}

// ensurePackageData ensures package data exists for the given directory
func (g *Generator) ensurePackageData(packageDir string, file *ast.File, config *GenerateConfig) {
	if g.packageDataMap[packageDir] == nil {
		moduleDir := findModuleRoot(packageDir)
		packagePath := getPackagePath(moduleDir, packageDir)

		g.packageDataMap[packageDir] = &PackageData{
			PackageName:   file.Name.Name,
			PackagePath:   packagePath,
			Registrations: []Registration{},
			Config:        config,
		}
	}
}

// parseGenerateComment parses //go:generate comment and extracts configuration
func parseGenerateComment(comment string) *GenerateConfig {
	re := regexp.MustCompile(`//go:generate\s+(?:voauto-gen|validationcontext)(?:\s+(.+))?`)
	matches := re.FindStringSubmatch(comment)
	if len(matches) < 2 {
		return &GenerateConfig{
			OutputPath:    "registry_init.go",
			OutputPackage: "",
		}
	}

	args := matches[1]
	config := &GenerateConfig{
		OutputPath:    "registry_init.go",
		OutputPackage: "",
	}

	parts := strings.Fields(args)
	for _, part := range parts {
		if strings.HasPrefix(part, "-output=") {
			config.OutputPath = strings.TrimPrefix(part, "-output=")
		} else if strings.HasPrefix(part, "-package=") {
			config.OutputPackage = strings.TrimPrefix(part, "-package=")
		} else if strings.HasPrefix(part, "-methods=") {
			methodsStr := strings.TrimPrefix(part, "-methods=")
			config.Methods = strings.Split(methodsStr, ",")
			for i, method := range config.Methods {
				config.Methods[i] = strings.TrimSpace(method)
			}
		}
	}

	return config
}
