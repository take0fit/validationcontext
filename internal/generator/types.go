package generator

// Registration represents a single constructor registration
type Registration struct {
	ConstructorName string
	RegistrationKey string
	QualifiedName   string
	ParamType       string
}

// Import represents an import statement for the generated registry
type Import struct {
	Alias string
	Path  string
}

// RegistryData contains all data needed to generate a registry file
type RegistryData struct {
	PackageName   string
	Registrations []Registration
	Imports       []Import
}

// GenerateConfig contains configuration parsed from //go:generate comments
type GenerateConfig struct {
	OutputPath    string
	OutputPackage string
	Methods       []string
}

// PackageData represents aggregated data for a single package
type PackageData struct {
	PackageName   string
	PackagePath   string
	Registrations []Registration
	Config        *GenerateConfig
}
