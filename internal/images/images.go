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
	"receipt_uploader/internal/models/configs"
	"receipt_uploader/internal/utils"
	"strings"

	"github.com/google/uuid"
	"github.com/nfnt/resize"
)

type Service struct {
	Dimensions *configs.Dimensions
}

func NewService(d *configs.Dimensions) ServiceType {
	return &Service{
		Dimensions: d,
	}
}

func (s *Service) GenerateResizedImages(srcPath, destDir string) error {
	logging.Debugf("GenerateResizedImages(srcPath: %s, destDir: %s)", srcPath, destDir)

	fileBytes, readErr := os.ReadFile(srcPath)
	if readErr != nil {
		return fmt.Errorf("os.ReadFile() failed: %v", readErr)
	}

	mkErr := os.MkdirAll(destDir, 0755)
	if mkErr != nil {
		err := fmt.Errorf("os.Mkdir() failed, err: %s", mkErr.Error())
		return err
	}

	copyDestPath := utils.GenerateDestPath(srcPath, destDir, "")
	copyErr := saveImage(&fileBytes, copyDestPath)
	if copyErr != nil {
		return fmt.Errorf("saveImage(copyDestPath: %s) failed, err: %s", copyDestPath, copyErr.Error())
	}

	img, _, decodeErr := image.Decode(bytes.NewReader(fileBytes))
	if decodeErr != nil {
		return fmt.Errorf("image.Decode() failed, err: %s", decodeErr.Error())
	}

	for _, d := range *s.Dimensions {
		resizedImg, resizeErr := resizeImage(&img, d.Width, d.Height)
		if resizeErr != nil {
			return fmt.Errorf(
				"resizeImage(srcPath: %s, width: %d, height: %d) failed, err: %s",
				srcPath, d.Width, d.Height, resizeErr.Error(),
			)
		}

		destPath := utils.GenerateDestPath(srcPath, destDir, d.Name)
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
	fName := receiptId + ".jpg"
	if size != "" {
		fName = receiptId + "_" + size + ".jpg"
	}

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
