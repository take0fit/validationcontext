package generator

type FunctionInfo struct {
	Name        string
	PackageName string
	Params      []ParamInfo
	Return      []ReturnInfo
	ImportPath  string
	FullKey     string
}

type ParamInfo struct {
	Name string
	Type string
}

type ReturnInfo struct {
	Name string
	Type string
}

type PackageInfo struct {
	Name      string
	Path      string
	Functions []*FunctionInfo
}

type RegistryInfo struct {
	PackageName string
	Functions   []*FunctionInfo
	ImportPath  string
}

// 元のコードで使用されている型定義
type Registration struct {
	ConstructorName string
	RegistrationKey string
	QualifiedName   string
	ParamType       string
}

type GenerateConfig struct {
	OutputFile     string
	OutputPath     string
	OutputPackage  string
	PackageName    string
	Methods        []string
	ExcludeMethods []string
	Verbose        bool
}

// 元のfile_generator.goで使用されているPackageData構造
type PackageData struct {
	PackageName   string
	PackagePath   string
	Registrations []Registration
	Config        *GenerateConfig
}

type Import struct {
	Path  string
	Alias string
}

type RegistryData struct {
	PackageName   string
	Registrations []Registration
	Imports       []Import
}
