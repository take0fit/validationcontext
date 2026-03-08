package validationcontext

import (
	"fmt"
	"os"
	"path/filepath"
)

// ValidateFilePath checks if the value is a valid file path.
func (vc *ValidationContext) ValidateFilePath(value, field, errMsg string) {
	if _, statErr := os.Stat(value); statErr != nil {
		vc.addErrorWithMessage(field, errMsg, fmt.Sprintf("%sには、有効なファイルパスを指定してください。", field))
	}
}

// ValidateFileExtension checks if the file has a valid extension.
func (vc *ValidationContext) ValidateFileExtension(file *os.File, field string, validExtensions []string, errMsg string) {
	if file == nil {
		vc.addErrorWithMessage(field, errMsg, fmt.Sprintf("%sには、有効なファイルを指定してください。", field))
		return
	}

	ext := filepath.Ext(file.Name())
	for _, validExt := range validExtensions {
		if ext == validExt {
			return
		}
	}
	vc.addErrorWithMessage(field, errMsg, fmt.Sprintf("%sには、有効な拡張子（%v）を持つファイルを指定してください。", field, validExtensions))
}

// ValidateFileSize checks if the file size is within the specified limit.
func (vc *ValidationContext) ValidateFileSize(file *os.File, field string, maxSize int64, errMsg string) {
	if file == nil {
		vc.addErrorWithMessage(field, errMsg, fmt.Sprintf("%sには、有効なファイルを指定してください。", field))
		return
	}

	fileInfo, statErr := file.Stat()
	if statErr != nil {
		vc.addErrorWithMessage(field, errMsg, fmt.Sprintf("%sのファイル情報の取得に失敗しました: %v", field, statErr))
		return
	}

	if fileInfo.Size() > maxSize {
		vc.addErrorWithMessage(field, errMsg, fmt.Sprintf("%sのファイルサイズは%dMB以下でなければなりません", field, maxSize/(1024*1024)))
	}
}
