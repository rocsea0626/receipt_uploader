package utils

import (
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"receipt_uploader/constants"
	"receipt_uploader/internal/logging"
	"receipt_uploader/internal/models/configs"
	"strings"

	"github.com/joho/godotenv"
)

// get file name without extension
func GetFileName(filePath string) string {
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
	logging.Debugf("size: %s", size)
	fName := GetFileName(srcPath)
	extension := filepath.Ext(srcPath)
	newFilename := fmt.Sprintf("%s_%s%s", fName, size, extension)
	return filepath.Join(destDir, newFilename)
}

func ValidateGetImageRequest(r *http.Request) (string, string, error) {
	receiptID := strings.TrimPrefix(r.URL.Path, "/receipts/")

	size := r.URL.Query().Get("size")
	if size != constants.IMAGE_SIZE_SMALL && size != constants.IMAGE_SIZE_MEDIUM && size != constants.IMAGE_SIZE_LARGE {
		return "", "", fmt.Errorf("invalid size parameter, size: %s", size)
	}

	return receiptID, size, nil
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
		Port:      os.Getenv("PORT"),
		ImagesDir: filepath.Join(constants.ROOT_DIR_IMAGES, os.Getenv("DIR_IMAGES")),
	}

	return config, nil
}
