package utils

import (
	"encoding/json"
	"fmt"
	"image"
	"image/jpeg"
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"receipt_uploader/constants"
	"receipt_uploader/internal/futils"
	"receipt_uploader/internal/models/configs"
	"receipt_uploader/internal/models/http_responses"
	"strings"

	"github.com/google/uuid"
	"github.com/joho/godotenv"
	"github.com/nfnt/resize"
)

func saveImage(img image.Image, outputFilePath string) error {
	outFile, createErr := os.Create(outputFilePath)
	if createErr != nil {
		return fmt.Errorf("os.Create() failed, err: %s", createErr.Error())
	}
	defer outFile.Close()

	EncodeErr := jpeg.Encode(outFile, img, nil)
	if EncodeErr != nil {
		return fmt.Errorf("jpeg.Encode() failed, err: %s", EncodeErr.Error())
	}

	return nil
}

func resizeImage(filePath string, width, height uint) (image.Image, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return nil, fmt.Errorf("os.Open() failed: %v", err)
	}
	defer file.Close()

	img, _, err := image.Decode(file)
	if err != nil {
		return nil, fmt.Errorf("image.Decode() failed: %v", err)
	}

	resizedImg := resize.Resize(width, height, img, resize.Lanczos3)
	return resizedImg, nil
}

// getOutputPath() generates a file path for the output file by taking an input file path,
// an output directory, and size is appended to the file name as suffix.
//
// Example:
// inputFilePath := "/path/to/input/file.txt"
// outputDir := "/path/to/output"
// size := "small"
// outputPath := "/path/to/output/file_small.txt"
func getOutputPath(filePath string, outputDir string, size string) string {
	fName := GetFileName(filePath)
	extension := filepath.Ext(filePath)
	newFilename := fmt.Sprintf("%s_%s%s", fName, size, extension)
	return filepath.Join(outputDir, newFilename)
}

// save the original upload into tmpDir
func SaveUploadedImage(r *http.Request, tmpDir string) (string, error) {
	file, _, fromErr := r.FormFile("receipt")
	if fromErr != nil {
		return "", fromErr
	}
	defer file.Close()

	fileName := uuid.New().String() + ".jpg"
	tmpPath := filepath.Join(tmpDir, fileName)

	tmpFile, createErr := os.Create(tmpPath)
	if createErr != nil {
		err := fmt.Errorf("os.Create() failed: %w", createErr)
		return "", err
	}
	defer tmpFile.Close()

	_, copyErr := io.Copy(tmpFile, file)
	if copyErr != nil {
		err := fmt.Errorf("io.Copy() failed: %w", createErr)
		return "", err
	}

	return tmpPath, nil
}

func GenerateImages(inputPath string, outputDir string) error {
	log.Printf("GenerateImages(inputPath: %s)", inputPath)

	smallImg, resizeSmallErr := resizeImage(inputPath, constants.IMAGE_SIZE_W_S, constants.IMAGE_SIZE_H_S)
	if resizeSmallErr != nil {
		err := fmt.Errorf(
			"ResizeImage(inputPath: %s, width: %d, height: %d) failed, err: %s",
			inputPath, constants.IMAGE_SIZE_H_S, constants.IMAGE_SIZE_H_S, resizeSmallErr.Error(),
		)
		return err
	}
	saveSmallErr := saveImage(smallImg, getOutputPath(inputPath, outputDir, "small"))
	if saveSmallErr != nil {
		err := fmt.Errorf("saveImage() failed, err: %s", saveSmallErr.Error())
		return err
	}

	mediumImg, resizeErr2 := resizeImage(inputPath, constants.IMAGE_SIZE_W_S, constants.IMAGE_SIZE_H_S)
	if resizeErr2 != nil {
		err := fmt.Errorf(
			"ResizeImage(inputPath: %s, width: %d, height: %d) failed, err: %s",
			inputPath, constants.IMAGE_SIZE_H_S, constants.IMAGE_SIZE_H_S, resizeErr2.Error(),
		)
		return err
	}
	saveMediumErr := saveImage(mediumImg, getOutputPath(inputPath, outputDir, "medium"))
	if saveMediumErr != nil {
		err := fmt.Errorf("saveImage() failed, err: %s", saveMediumErr.Error())
		return err
	}

	saveLargeErr := futils.CopyFile(inputPath, getOutputPath(inputPath, outputDir, "large"))
	if saveLargeErr != nil {
		err := fmt.Errorf("futils.CopyFile() failed, err: %s", saveLargeErr.Error())
		return err
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

func SendErrorResponse(w http.ResponseWriter, resp *http_responses.ErrorResponse, status int) {
	errMap := map[string]string{
		"error": resp.Error,
	}
	sendResponse(w, &errMap, status)
}

func SendUploadResponse(w http.ResponseWriter, resp *http_responses.UploadResponse) {
	respMap := map[string]string{
		"receiptId": resp.ReceiptID,
	}
	sendResponse(w, &respMap, http.StatusCreated)
}

func sendResponse(w http.ResponseWriter, response *map[string]string, status int) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(*response)
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

func InitServer() error {
	tmpErr := os.Mkdir("tmp", 0755)
	if tmpErr != nil {
		err := fmt.Errorf("os.Mkdir() failed, err: %s", tmpErr.Error())
		return err
	}

	imagesErr := os.Mkdir("images", 0755)
	if imagesErr != nil {
		err := fmt.Errorf("os.Mkdir() failed, err: %s", imagesErr.Error())
		return err
	}
	return nil
}
