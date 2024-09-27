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

	"github.com/google/uuid"
	"github.com/nfnt/resize"
)

type Service struct {
}

func NewService() ServiceType {
	return &Service{}
}

func (s *Service) GenerateImages(srcPath string, destDir string) error {
	log.Printf("GenerateImages(srcPath: %s, destDir: %s)", srcPath, destDir)

	smallImg, resizeSmallErr := resizeImage(srcPath, constants.IMAGE_SIZE_W_S, constants.IMAGE_SIZE_H_S)
	if resizeSmallErr != nil {
		err := fmt.Errorf(
			"ResizeImage(srcPath: %s, width: %d, height: %d) failed, err: %s",
			srcPath, constants.IMAGE_SIZE_H_S, constants.IMAGE_SIZE_H_S, resizeSmallErr.Error(),
		)
		return err
	}
	saveSmallErr := saveImage(&smallImg, futils.GetOutputPath(srcPath, destDir, "small"))
	if saveSmallErr != nil {
		err := fmt.Errorf("saveImage() failed, err: %s", saveSmallErr.Error())
		return err
	}

	mediumImg, resizeErr2 := resizeImage(srcPath, constants.IMAGE_SIZE_W_M, constants.IMAGE_SIZE_H_M)
	if resizeErr2 != nil {
		err := fmt.Errorf(
			"ResizeImage(srcPath: %s, width: %d, height: %d) failed, err: %s",
			srcPath, constants.IMAGE_SIZE_H_S, constants.IMAGE_SIZE_H_S, resizeErr2.Error(),
		)
		return err
	}
	saveMediumErr := saveImage(&mediumImg, futils.GetOutputPath(srcPath, destDir, "medium"))
	if saveMediumErr != nil {
		err := fmt.Errorf("saveImage() failed, err: %s", saveMediumErr.Error())
		return err
	}

	saveLargeErr := futils.CopyFile(srcPath, futils.GetOutputPath(srcPath, destDir, "large"))
	if saveLargeErr != nil {
		err := fmt.Errorf("futils.CopyFile() failed, err: %s", saveLargeErr.Error())
		return err
	}

	return nil
}

func (s *Service) DecodeImage(r *http.Request) ([]byte, error) {

	parseErr := r.ParseMultipartForm(10 << 20) // Maximum 10 MB
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

func (s *Service) SaveUpload(bytes []byte, tmpDir string) (string, error) {
	log.Printf("SaveUpload(tmpDir: %s)", tmpDir)

	fileName := uuid.New().String() + ".jpg"
	destPath := filepath.Join(tmpDir, fileName)

	saveImage(&bytes, destPath)

	return destPath, nil
}

func resizeImage(srcPath string, width, height int) ([]byte, error) {
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

// func saveImage(img *image.Image, destPath string) error {
// 	outFile, createErr := os.Create(destPath)
// 	if createErr != nil {
// 		return fmt.Errorf("os.Create() failed, err: %s", createErr.Error())
// 	}
// 	defer outFile.Close()

// 	EncodeErr := jpeg.Encode(outFile, *img, nil)
// 	if EncodeErr != nil {
// 		return fmt.Errorf("jpeg.Encode() failed, err: %s", EncodeErr.Error())
// 	}

// 	return nil
// }

func saveImage(bytes *[]byte, destPath string) error {
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
