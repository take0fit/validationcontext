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

var generateCommentRegexp = regexp.MustCompile(`//go:generate\s+(?:voauto-gen|validationcontext)(?:\s+(.+))?`)

// processFile processes a single Go file and collects registration data.
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

	for _, commentGroup := range file.Comments {
		for _, comment := range commentGroup.List {
			if g.verbose {
				fmt.Printf("  Comment: %s\n", comment.Text)
			}

			parsedConfig, ok := parseGenerateComment(comment.Text)
			if !ok {
				continue
			}

			generateCommentCount++
			if g.verbose {
				fmt.Printf("Found generate comment in %s: %s\n", filename, comment.Text)
			}
			config = parsedConfig
			break
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

// ensurePackageData ensures package data exists for the given directory.
func (g *Generator) ensurePackageData(packageDir string, file *ast.File, config *GenerateConfig) {
	if g.packageDataMap[packageDir] == nil {
		var packagePath string

		if g.explicitPkgPath != "" {
			if !strings.Contains(g.explicitPkgPath, "/") {
				currentPkgPath := getCurrentPackagePath()
				if currentPkgPath != "" {
					packagePath = currentPkgPath
					if g.verbose {
						fmt.Printf("  Using current directory package path: %s\n", packagePath)
					}
				} else {
					packagePath = g.explicitPkgPath
					if g.verbose {
						fmt.Printf("  Using GOPACKAGE: %s\n", packagePath)
					}
				}
			} else {
				packagePath = g.explicitPkgPath
				if g.verbose {
					fmt.Printf("  Using explicit package path: %s\n", packagePath)
				}
			}
		} else {
			moduleDir := findModuleRoot(packageDir)
			packagePath = getPackagePath(moduleDir, packageDir)
			if g.verbose {
				fmt.Printf("  Auto-detected package path: %s\n", packagePath)
			}
		}

		if g.verbose {
			fmt.Printf("  Package directory: %s\n", packageDir)
			fmt.Printf("  Final package path: %s\n", packagePath)
		}

		g.packageDataMap[packageDir] = &PackageData{
			PackageName:   file.Name.Name,
			PackagePath:   packagePath,
			Registrations: []Registration{},
			Config:        config,
		}
	}
}

// parseGenerateComment parses //go:generate comment and extracts configuration.
func parseGenerateComment(comment string) (*GenerateConfig, bool) {
	matches := generateCommentRegexp.FindStringSubmatch(comment)
	if len(matches) == 0 {
		return nil, false
	}

	args := ""
	if len(matches) > 1 {
		args = matches[1]
	}

	config := &GenerateConfig{
		OutputPath:    "registry_init.go",
		OutputPackage: "",
	}

	parts := strings.Fields(args)
	for i := 0; i < len(parts); i++ {
		part := parts[i]
		switch {
		case strings.HasPrefix(part, "-output="):
			value := strings.TrimSpace(strings.TrimPrefix(part, "-output="))
			if value != "" {
				config.OutputPath = value
			}
		case strings.HasPrefix(part, "-package="):
			config.OutputPackage = strings.TrimSpace(strings.TrimPrefix(part, "-package="))
		case strings.HasPrefix(part, "-methods="):
			methodsStr := strings.TrimSpace(strings.TrimPrefix(part, "-methods="))
			for i+1 < len(parts) && !strings.HasPrefix(parts[i+1], "-") {
				i++
				methodsStr += strings.TrimSpace(parts[i])
			}
			if methodsStr == "" {
				continue
			}
			methods := strings.Split(methodsStr, ",")
			for _, method := range methods {
				trimmed := strings.TrimSpace(method)
				if trimmed != "" {
					config.Methods = append(config.Methods, trimmed)
				}
			}
		}
	}

	return config, true
}
