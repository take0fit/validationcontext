package generator

import "testing"

func TestIsTargetGoFile(t *testing.T) {
	tests := []struct {
		path string
		want bool
	}{
		{path: "a.go", want: true},
		{path: "a_test.go", want: false},
		{path: "registry_init.go", want: false},
		{path: "a.txt", want: false},
	}

	for _, tt := range tests {
		if got := isTargetGoFile(tt.path); got != tt.want {
			t.Fatalf("isTargetGoFile(%q) = %v, want %v", tt.path, got, tt.want)
		}
	}
}

func TestIsIgnoredDir(t *testing.T) {
	if !isIgnoredDir("vendor") {
		t.Fatal("vendor should be ignored")
	}
	if !isIgnoredDir(".git") {
		t.Fatal(".git should be ignored")
	}
	if isIgnoredDir("internal") {
		t.Fatal("internal should not be ignored")
	}
}
