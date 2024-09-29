package images

import (
	"bytes"
	"fmt"
	"image"
	"image/jpeg"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"receipt_uploader/constants"
	"receipt_uploader/internal/logging"
	"receipt_uploader/internal/utils"
	"strings"

	"github.com/google/uuid"
	"github.com/nfnt/resize"
)

type Service struct {
}

func NewService() ServiceType {
	return &Service{}
}

func (s *Service) GenerateImages(fileBytes *[]byte, srcPath, destDir string) error {
	logging.Debugf("GenerateImages(len(fileBytes): %d, destDir: %s)", len(*fileBytes), destDir)

	mkErr := os.MkdirAll(destDir, 0755)
	if mkErr != nil {
		err := fmt.Errorf("os.Mkdir() failed, err: %s", mkErr.Error())
		return err
	}

	// // create image with original size
	saveLargeErr := saveImage(fileBytes, utils.GenerateDestPath(srcPath, destDir, "large"))
	if saveLargeErr != nil {
		return fmt.Errorf("saveImage() failed, err: %s", saveLargeErr.Error())
	}

	sizes := []string{"small", "medium"}
	img, _, decodeErr := image.Decode(bytes.NewReader(*fileBytes))
	if decodeErr != nil {
		return fmt.Errorf("image.Decode() failed, err: %s", decodeErr.Error())
	}

	width := 0
	height := 0
	for _, size := range sizes {

		if size == "small" {
			width = constants.IMAGE_SIZE_W_S
			height = constants.IMAGE_SIZE_H_S
		}

		if size == "meidum" {
			width = constants.IMAGE_SIZE_W_M
			height = constants.IMAGE_SIZE_H_M
		}

		if size == "large" {
			width = constants.IMAGE_SIZE_W_L
			height = constants.IMAGE_SIZE_H_L
		}

		resizedImg, resizeErr := resizeImage(&img, width, height)
		if resizeErr != nil {
			return fmt.Errorf(
				"resizeImage(srcPath: %s, width: %d, height: %d) failed, err: %s",
				srcPath, constants.IMAGE_SIZE_H_S, constants.IMAGE_SIZE_H_S, resizeErr.Error(),
			)
		}

		destPath := utils.GenerateDestPath(srcPath, destDir, size)
		logging.Debugf("destPath: %s", destPath)
		saveErr := saveImage(&resizedImg, destPath)
		if saveErr != nil {
			return fmt.Errorf("saveImage(destPath: %s) failed, err: %s", destPath, saveErr.Error())
		}
	}

	return nil
}

func (s *Service) ParseImage(r *http.Request) ([]byte, error) {
	parseErr := r.ParseMultipartForm(constants.MAX_UPLOAD_SIZE)
	if parseErr != nil {
		return nil, fmt.Errorf("r.ParseMultipartForm() failed, err: %s", parseErr.Error())
	}

	file, header, fromErr := r.FormFile("receipt")
	if fromErr != nil {
		return nil, fmt.Errorf("r.FormFile() failed: %w", fromErr)
	}
	logging.Debugf("content-type: %s", header.Header.Get("Content-Type"))
	defer file.Close()

	fileBytes, err := io.ReadAll(file)
	if err != nil {
		return nil, fmt.Errorf("failed to read file: %w", err)
	}

	img, format, decodeErr := image.Decode(bytes.NewReader(fileBytes))
	if decodeErr != nil {
		return nil, fmt.Errorf("image.Decode() failed, err: %s", decodeErr.Error())
	}
	if img.Bounds().Dx() < constants.IMAGE_SIZE_MIN_W || img.Bounds().Dy() < constants.IMAGE_SIZE_MIN_H {
		return nil, fmt.Errorf("invalid image size, minHeight=%d, minWidth=%d", constants.IMAGE_SIZE_MIN_H, constants.IMAGE_SIZE_MIN_W)
	}

	if format != "jpeg" {
		return nil, fmt.Errorf("invalid image format, format=%s, only jpeg format is allowed", format)
	}

	return fileBytes, nil
}

func (s *Service) SaveUpload(bytes *[]byte, destDir string) (string, error) {
	logging.Debugf("SaveUpload(len(bytes): %d, destDir: %s)", len(*bytes), destDir)

	mkErr := os.MkdirAll(destDir, 0755)
	if mkErr != nil {
		err := fmt.Errorf("os.Mkdir() failed, err: %s", mkErr.Error())
		return "", err
	}

	fileName := strings.ReplaceAll(uuid.New().String()+".jpg", "-", "")
	destPath := filepath.Join(destDir, fileName)

	saveImage(bytes, destPath)

	return destPath, nil
}

func (s *Service) GetImage(receiptId, size, srcDir string) ([]byte, string, error) {
	logging.Debugf("GetImage(receiptId: %s, size: %s, srcDir: %s)", receiptId, size, srcDir)

	fName := fmt.Sprintf("%s_%s.jpg", receiptId, size)
	fPath := fmt.Sprintf("%s/%s", srcDir, fName)
	logging.Debugf("fPath: %s", fPath)

	fileBytes, readErr := os.ReadFile(fPath)
	if readErr != nil {
		if os.IsNotExist(readErr) {
			return nil, "", readErr
		}
		return nil, "", fmt.Errorf("os.ReadFile() failed: %v", readErr)
	}
	return fileBytes, fName, nil
}

func resizeImage(img *image.Image, width, height int) ([]byte, error) {
	logging.Debugf("resizeImage(width: %d, height: %d)", width, height)

	var buf bytes.Buffer
	resizedImg := resize.Resize(uint(width), uint(height), *img, resize.Lanczos3)
	encodeErr := jpeg.Encode(&buf, resizedImg, nil)
	if encodeErr != nil {
		return nil, fmt.Errorf("jpeg.Encode() failed, err: %s", encodeErr.Error())
	}

	return buf.Bytes(), nil
}

func saveImage(bytes *[]byte, destPath string) error {
	logging.Debugf("saveImage(len(bytes): %d, destPath: %s)", len(*bytes), destPath)

	destFile, createErr := os.Create(destPath)
	if createErr != nil {
		return fmt.Errorf("os.Create() failed, err: %s", createErr.Error())
	}
	defer destFile.Close()

	_, writeErr := destFile.Write(*bytes)
	if writeErr != nil {
		return fmt.Errorf("file.Write() failed: %w", writeErr)
	}

	return nil
}
