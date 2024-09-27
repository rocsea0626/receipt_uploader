package futils

import (
	"fmt"
	"io"
	"os"
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
