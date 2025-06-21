package main

import (
	"flag"
	"fmt"
	"log"

	"github.com/take0fit/validationcontext/internal/generator"
)

func main() {
	var (
		targetDir = flag.String("dir", ".", "directory to scan for go files")
		verbose   = flag.Bool("v", false, "verbose output")
	)
	flag.Parse()

	fmt.Printf("Scanning directory: %s\n", *targetDir)

	gen := generator.New(*verbose)

	fileCount, generateCommentCount, err := gen.ScanDirectory(*targetDir)
	if err != nil {
		log.Fatalf("failed to scan directory: %v", err)
	}

	fmt.Printf("Processed %d Go files\n", fileCount)
	fmt.Printf("Found %d //go:generate comments\n", generateCommentCount)

	registryCount, err := gen.GenerateRegistries()
	if err != nil {
		log.Fatalf("failed to generate registries: %v", err)
	}

	if fileCount == 0 {
		fmt.Println("No Go files found in the target directory")
	}
	if generateCommentCount == 0 {
		fmt.Println("No //go:generate voauto-gen or //go:generate validationcontext comments found")
		fmt.Println("\nTo use this tool, add a comment like this to your Go file:")
		fmt.Println("//go:generate voauto-gen")
		fmt.Println("or")
		fmt.Println("//go:generate voauto-gen -output=custom_registry.go -methods=NewUser,NewProduct")
	}

	fmt.Printf("Generated %d registry files\n", registryCount)
}
