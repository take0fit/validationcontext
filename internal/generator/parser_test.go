package generator

import (
	"reflect"
	"testing"
)

func TestParseGenerateComment(t *testing.T) {
	tests := []struct {
		name       string
		comment    string
		wantOK     bool
		wantOutput string
		wantPkg    string
		wantMethod []string
	}{
		{
			name:       "Simple",
			comment:    "//go:generate voauto-gen",
			wantOK:     true,
			wantOutput: "registry_init.go",
		},
		{
			name:       "WithOptions",
			comment:    "//go:generate voauto-gen -output=custom.go -package=foo -methods=NewUser, NewBook",
			wantOK:     true,
			wantOutput: "custom.go",
			wantPkg:    "foo",
			wantMethod: []string{"NewUser", "NewBook"},
		},
		{
			name:       "LegacyCommandAlias",
			comment:    "//go:generate validationcontext -methods=NewUser",
			wantOK:     true,
			wantOutput: "registry_init.go",
			wantMethod: []string{"NewUser"},
		},
		{
			name:    "NotGenerateComment",
			comment: "// go:generate voauto-gen",
			wantOK:  false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, ok := parseGenerateComment(tt.comment)
			if ok != tt.wantOK {
				t.Fatalf("ok = %v, want %v", ok, tt.wantOK)
			}
			if !tt.wantOK {
				return
			}
			if got.OutputPath != tt.wantOutput {
				t.Fatalf("OutputPath = %q, want %q", got.OutputPath, tt.wantOutput)
			}
			if got.OutputPackage != tt.wantPkg {
				t.Fatalf("OutputPackage = %q, want %q", got.OutputPackage, tt.wantPkg)
			}
			if !reflect.DeepEqual(got.Methods, tt.wantMethod) {
				t.Fatalf("Methods = %#v, want %#v", got.Methods, tt.wantMethod)
			}
		})
	}
}
