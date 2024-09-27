package utils

import (
	"fmt"
	"image"
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"receipt_uploader/constants"
	"receipt_uploader/internal/models/configs"

	"github.com/google/uuid"
	"github.com/joho/godotenv"
)

// save the original upload into tmpDir
func SaveUploadedImage(r *http.Request, tmpDir string) (string, error) {
	log.Printf("SaveUploadedImage(tmpDir: %s)", tmpDir)

	parseErr := r.ParseMultipartForm(10 << 20) // Maximum 10 MB
	if parseErr != nil {
		return "", fmt.Errorf("r.ParseMultipartForm() failed, err: %s", parseErr.Error())
	}

	file, header, fromErr := r.FormFile("receipt")
	if fromErr != nil {
		return "", fmt.Errorf("r.FormFile() failed: %w", fromErr)
	}
	log.Printf("content-type: %s", header.Header.Get("Content-Type"))
	log.Printf("file size: %d", header.Size)

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
	log.Printf("written: %d", written)

	return tmpPath, nil
}

func SaveUpload(bytes []byte, tmpDir string) (string, error) {
	log.Printf("SaveUpload(tmpDir: %s)", tmpDir)

	fileName := uuid.New().String() + ".jpg"
	tmpPath := filepath.Join(tmpDir, fileName)

	tmpFile, createErr := os.Create(tmpPath)
	if createErr != nil {
		return "", fmt.Errorf("os.Create() failed: %w", createErr)
	}
	defer tmpFile.Close()

	_, copyErr := tmpFile.Write(bytes)
	if copyErr != nil {
		return "", fmt.Errorf("file.Write() failed: %w", copyErr)
	}

	return tmpPath, nil
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
		DIR_TMP:    os.Getenv("DIR_TMP"),
		DIR_IMAGES: os.Getenv("DIR_IMAGES"),
	}

	return config, nil
}

func InitServer(config *configs.Config) error {
	tmpErr := os.Mkdir(config.DIR_TMP, 0755)
	if tmpErr != nil {
		err := fmt.Errorf("os.Mkdir() failed, err: %s", tmpErr.Error())
		return err
	}
	log.Printf("folder %s has been created", config.DIR_TMP)

	imagesErr := os.Mkdir(config.DIR_IMAGES, 0755)
	if imagesErr != nil {
		err := fmt.Errorf("os.Mkdir() failed, err: %s", imagesErr.Error())
		return err
	}
	log.Printf("folder %s has been created", config.DIR_IMAGES)
	return nil
}
