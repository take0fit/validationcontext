package generator

import (
	"bytes"
	"fmt"
	"go/format"
	"os"
	"path/filepath"
	"strings"
	"text/template"
)

// generateRegistryForPackage generates a registry file for a specific package
func (g *Generator) generateRegistryForPackage(data *PackageData, packageDir string) error {
	if g.verbose {
		fmt.Printf("Generating registry for package %s with %d registrations\n", data.PackageName, len(data.Registrations))
	}

	packageName := data.Config.OutputPackage
	if packageName == "" {
		packageName = data.PackageName
	}

	imports := make(map[string]Import)

	// Add import if outputting to different package
	if packageName != data.PackageName && data.PackagePath != "" {
		alias := data.PackageName
		imports[data.PackagePath] = Import{
			Alias: alias,
			Path:  data.PackagePath,
		}

		// Update QualifiedName
		for i := range data.Registrations {
			data.Registrations[i].QualifiedName = alias + "." + data.Registrations[i].ConstructorName
		}
	}

	// Create import slice
	importSlice := make([]Import, 0, len(imports))
	for _, imp := range imports {
		importSlice = append(importSlice, imp)
	}

	// Generate code from template
	registryData := RegistryData{
		PackageName:   packageName,
		Registrations: data.Registrations,
		Imports:       importSlice,
	}

	tmpl, err := template.New("registry").Parse(registryTemplate)
	if err != nil {
		return fmt.Errorf("failed to parse template: %w", err)
	}

	var buf bytes.Buffer
	if err := tmpl.Execute(&buf, registryData); err != nil {
		return fmt.Errorf("failed to execute template: %w", err)
	}

	// Format code
	formatted, err := format.Source(buf.Bytes())
	if err != nil {
		return fmt.Errorf("failed to format source: %w", err)
	}

	// Determine output path
	outputPath := data.Config.OutputPath
	if !filepath.IsAbs(outputPath) {
		outputPath = filepath.Join(packageDir, outputPath)
	}

	// Create directory if it doesn't exist
	outputDir := filepath.Dir(outputPath)
	if err := os.MkdirAll(outputDir, 0755); err != nil {
		return fmt.Errorf("failed to create output directory: %w", err)
	}

	// Write to file
	if err := os.WriteFile(outputPath, formatted, 0644); err != nil {
		return fmt.Errorf("failed to write file: %w", err)
	}

	fmt.Printf("Generated %s with %d registrations\n", outputPath, len(data.Registrations))
	return nil
}

// generateRegistrationKey generates a unique registration key
func generateRegistrationKey(packagePath, constructorName string) string {
	if packagePath == "" {
		return constructorName
	}

	parts := strings.Split(packagePath, "/")
	if len(parts) == 0 {
		return constructorName
	}

	var keyParts []string
	if len(parts) >= 3 {
		keyParts = parts[len(parts)-3:]
	} else if len(parts) >= 2 {
		keyParts = parts[len(parts)-2:]
	} else {
		keyParts = []string{parts[len(parts)-1]}
	}

	packageKey := strings.Join(keyParts, "_")
	return packageKey + "_" + constructorName
}
