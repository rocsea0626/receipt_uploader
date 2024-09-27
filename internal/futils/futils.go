package futils

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
)

func CopyFile(srcPath, dstPath string) error {
	srcFile, err := os.Open(srcPath)
	if err != nil {
		return fmt.Errorf("os.Open() failed: %w", err)
	}
	defer srcFile.Close()

	destFile, err := os.Create(dstPath)
	if err != nil {
		return fmt.Errorf("os.Create() failed: %w", err)
	}
	defer destFile.Close()

	_, err = io.Copy(destFile, srcFile)
	if err != nil {
		return fmt.Errorf("io.Copy() failed: %w", err)
	}

	return nil
}

// get file name without extension
func GetFileName(filePath string) string {
	base := filepath.Base(filePath)
	extension := filepath.Ext(filePath)
	fName := strings.TrimSuffix(base, extension)
	return fName
}

// getOutputPath() generates a file path for the output file by taking an input file path,
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
