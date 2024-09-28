package futils

import (
	"fmt"
	"path/filepath"
	"strings"
)

// get file name without extension
func GetFileName(filePath string) string {
	base := filepath.Base(filePath)
	extension := filepath.Ext(filePath)
	fName := strings.TrimSuffix(base, extension)
	return fName
}

// GetOutputPath() generates a file path for the output file by taking an input file path,
// an output directory, and size is appended to the file name as suffix.
//
// Example:
// inputFilePath := "/path/to/input/file.txt"
// outputDir := "/path/to/output"
// size := "small"
// outputPath := "/path/to/output/file_small.txt"
func GetOutputPath(filePath string, outputDir string, size string) string {
	fName := GetFileName(filePath)
	extension := filepath.Ext(filePath)
	newFilename := fmt.Sprintf("%s_%s%s", fName, size, extension)
	return filepath.Join(outputDir, newFilename)
}
