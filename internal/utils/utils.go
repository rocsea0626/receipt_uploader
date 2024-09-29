package utils

import (
	"fmt"
	"image"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"receipt_uploader/constants"
	"receipt_uploader/internal/logging"
	"receipt_uploader/internal/models/configs"
	"strings"

	"github.com/google/uuid"
	"github.com/joho/godotenv"
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

// save the original upload into tmpDir
func SaveUploadedImage(r *http.Request, tmpDir string) (string, error) {
	logging.Debugf("SaveUploadedImage(tmpDir: %s)", tmpDir)

	parseErr := r.ParseMultipartForm(10 << 20) // Maximum 10 MB
	if parseErr != nil {
		return "", fmt.Errorf("r.ParseMultipartForm() failed, err: %s", parseErr.Error())
	}

	file, header, fromErr := r.FormFile("receipt")
	if fromErr != nil {
		return "", fmt.Errorf("r.FormFile() failed: %w", fromErr)
	}
	logging.Debugf("content-type: %s", header.Header.Get("Content-Type"))
	logging.Debugf("file size: %d", header.Size)

	defer file.Close()

	img, _, decodeErr := image.Decode(file)
	if decodeErr != nil {
		return "", fmt.Errorf("image.Decode() failed, err: %s", decodeErr.Error())
	}
	if img.Bounds().Dx() < constants.IMAGE_SIZE_MIN_W || img.Bounds().Dy() < constants.IMAGE_SIZE_MIN_H {
		return "", fmt.Errorf("invalid image size, minHeight=%d, minWidth=%d", constants.IMAGE_SIZE_MIN_H, constants.IMAGE_SIZE_MIN_W)
	}

	fileName := uuid.New().String() + ".jpg"
	tmpPath := filepath.Join(tmpDir, fileName)

	tmpFile, createErr := os.Create(tmpPath)
	if createErr != nil {
		return "", fmt.Errorf("os.Create() failed: %w", createErr)
	}
	defer tmpFile.Close()

	_, seekErr := file.Seek(0, io.SeekStart)
	if seekErr != nil {
		return "", fmt.Errorf("file.Seek() failed: %w", seekErr)
	}

	written, copyErr := io.Copy(tmpFile, file)
	if copyErr != nil {
		return "", fmt.Errorf("io.Copy() failed: %w", copyErr)
	}
	logging.Debugf("written: %d", written)

	return tmpPath, nil
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
		Port:         os.Getenv("PORT"),
		UploadedDir:  filepath.Join(constants.ROOT_DIR_IMAGES, os.Getenv("DIR_UPLOAD")),
		GeneratedDir: filepath.Join(constants.ROOT_DIR_IMAGES, os.Getenv("DIR_GENERATED")),
	}

	return config, nil
}
