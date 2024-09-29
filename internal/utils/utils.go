package utils

import (
	"fmt"
	"os"
	"path/filepath"
	"receipt_uploader/constants"
	"receipt_uploader/internal/models/configs"
	"strings"

	"github.com/joho/godotenv"
)

// extract file name without extension
func ExtractFileName(filePath string) string {
	base := filepath.Base(filePath)
	extension := filepath.Ext(filePath)
	fName := strings.TrimSuffix(base, extension)
	return fName
}

// GenerateDestPath() generates a file path for the output file from srcPath
//
// Example:
// srcPath := "/path/to/input/file.jpg"
// destDir := "/path/to/output"
// size := "small"
// output := "/path/to/output/file_small.jpg"
func GenerateDestPath(srcPath string, destDir string, size string) string {
	fName := ExtractFileName(srcPath)
	extension := filepath.Ext(srcPath)
	newFilename := fmt.Sprintf("%s%s", fName, extension)
	if size != "" {
		newFilename = fmt.Sprintf("%s_%s%s", fName, size, extension)
	}
	return filepath.Join(destDir, newFilename)
}

func LoadConfig() (*configs.Config, error) {
	env := os.Getenv("env")
	envFile := ".env"
	if env == "dev" {
		envFile = ".env.development"
	}
	loadErr := godotenv.Load(envFile)
	if loadErr != nil {
		return nil, loadErr
	}

	config := &configs.Config{
		Port:       os.Getenv("PORT"),
		ImagesDir:  filepath.Join(constants.ROOT_DIR_IMAGES, os.Getenv("DIR_IMAGES")),
		UploadsDir: filepath.Join(constants.ROOT_DIR_IMAGES, os.Getenv("DIR_UPLOADS")),
		Dimensions: configs.AllowedDimensions,
	}

	return config, nil
}
