package images

import (
	"bytes"
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
	"strings"

	"github.com/google/uuid"
	"github.com/nfnt/resize"
)

type Service struct {
}

func NewService() ServiceType {
	return &Service{}
}

func (s *Service) GenerateImages(srcPath, destDir string) error {
	log.Printf("GenerateImages(srcPath: %s, destDir: %s)", srcPath, destDir)

	mkErr := os.MkdirAll(destDir, 0755)
	if mkErr != nil {
		err := fmt.Errorf("os.Mkdir() failed, err: %s", mkErr.Error())
		return err
	}

	smallImg, resizeSmallErr := resizeImage(srcPath, constants.IMAGE_SIZE_W_S, constants.IMAGE_SIZE_H_S)
	if resizeSmallErr != nil {
		return fmt.Errorf(
			"ResizeImage(srcPath: %s, width: %d, height: %d) failed, err: %s",
			srcPath, constants.IMAGE_SIZE_H_S, constants.IMAGE_SIZE_H_S, resizeSmallErr.Error(),
		)
	}
	saveSmallErr := saveImage(&smallImg, futils.GetOutputPath(srcPath, destDir, "small"))
	if saveSmallErr != nil {
		return fmt.Errorf("saveImage() failed, err: %s", saveSmallErr.Error())
	}

	mediumImg, resizeErr2 := resizeImage(srcPath, constants.IMAGE_SIZE_W_M, constants.IMAGE_SIZE_H_M)
	if resizeErr2 != nil {
		return fmt.Errorf(
			"ResizeImage(srcPath: %s, width: %d, height: %d) failed, err: %s",
			srcPath, constants.IMAGE_SIZE_W_M, constants.IMAGE_SIZE_H_M, resizeErr2.Error(),
		)
	}
	saveMediumErr := saveImage(&mediumImg, futils.GetOutputPath(srcPath, destDir, "medium"))
	if saveMediumErr != nil {
		return fmt.Errorf("saveImage() failed, err: %s", saveMediumErr.Error())
	}

	fileBytes, readErr := os.ReadFile(srcPath)
	if readErr != nil {
		return fmt.Errorf("saveImage() failed, err: %s", readErr.Error())
	}
	saveLargeErr := saveImage(&fileBytes, futils.GetOutputPath(srcPath, destDir, "large"))
	if saveLargeErr != nil {
		return fmt.Errorf("saveImage() failed, err: %s", saveLargeErr.Error())
	}

	return nil
}

func (s *Service) DecodeImage(r *http.Request) ([]byte, error) {
	parseErr := r.ParseMultipartForm(constants.MAX_UPLOAD_SIZE)
	if parseErr != nil {
		return nil, fmt.Errorf("r.ParseMultipartForm() failed, err: %s", parseErr.Error())
	}

	file, header, fromErr := r.FormFile("receipt")
	if fromErr != nil {
		return nil, fmt.Errorf("r.FormFile() failed: %w", fromErr)
	}
	log.Printf("content-type: %s", header.Header.Get("Content-Type"))
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

func (s *Service) SaveUpload(bytes []byte, destDir string) (string, error) {
	log.Printf("SaveUpload(len(bytes): %d, destDir: %s)", len(bytes), destDir)

	fileName := strings.ReplaceAll(uuid.New().String()+".jpg", "-", "")
	destPath := filepath.Join(destDir, fileName)

	saveImage(&bytes, destPath)

	return destPath, nil
}

func (s *Service) GetImage(receiptId, size, srcDir string) ([]byte, string, error) {
	log.Printf("GetImage(receiptId: %s, size: %s, srcDir: %s)", receiptId, size, srcDir)

	fName := fmt.Sprintf("%s_%s.jpg", receiptId, size)
	fPath := fmt.Sprintf("%s/%s", srcDir, fName)
	log.Printf("fPath: %s", fPath)

	fileBytes, readErr := os.ReadFile(fPath)
	if readErr != nil {
		if os.IsNotExist(readErr) {
			return nil, "", readErr
		}
		return nil, "", fmt.Errorf("os.ReadFile() failed: %v", readErr)
	}
	return fileBytes, fName, nil
}

func resizeImage(srcPath string, width, height int) ([]byte, error) {
	log.Printf("resizeImage(srcPath: %s, width: %d, height: %d)", srcPath, width, height)

	file, err := os.Open(srcPath)
	if err != nil {
		return nil, fmt.Errorf("os.Open() failed: %v", err)
	}
	defer file.Close()

	img, _, err := image.Decode(file)
	if err != nil {
		return nil, fmt.Errorf("image.Decode() failed: %v", err)
	}

	var buf bytes.Buffer
	resizedImg := resize.Resize(uint(width), uint(height), img, resize.Lanczos3)
	encodeErr := jpeg.Encode(&buf, resizedImg, nil)
	if encodeErr != nil {
		return nil, fmt.Errorf("jpeg.Encode() failed, err: %s", encodeErr.Error())
	}

	return buf.Bytes(), nil
}

func saveImage(bytes *[]byte, destPath string) error {
	log.Printf("saveImage(len(bytes): %d, destPath: %s)", len(*bytes), destPath)

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
