package validationcontext

import (
	"fmt"
	"os"
	"path/filepath"
)

// ValidateFilePath checks if the value is a valid file path.
func (vc *ValidationContext) ValidateFilePath(value, field, errMsg string) {
	if _, err := os.Stat(value); err != nil {
		if os.IsNotExist(err) {
			if errMsg != "" {
				vc.AddError(field, errMsg)
				return
			}
			vc.AddError(field, fmt.Sprintf("%sには、有効なファイルパスを指定してください。", field))
		}
	}
}

// ValidateFileExtension checks if the file has a valid extension.
func (vc *ValidationContext) ValidateFileExtension(file *os.File, field string, validExtensions []string, errMsg string) {
	ext := filepath.Ext(file.Name())
	for _, validExt := range validExtensions {
		if ext == validExt {
			return
		}
	}
	if errMsg != "" {
		vc.AddError(field, errMsg)
		return
	}
	vc.AddError(field, fmt.Sprintf("%sには、有効な拡張子（%v）を持つファイルを指定してください。", field, validExtensions))
}
