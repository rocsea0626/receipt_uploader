package utils

import (
	"fmt"
	"image"
	"image/jpeg"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"github.com/google/uuid"
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

func ResizeImage(filePath string, width, height uint) (image.Image, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return nil, fmt.Errorf("error opening image file: %v", err)
	}
	defer file.Close()

	img, _, err := image.Decode(file)
	if err != nil {
		return nil, fmt.Errorf("error decoding image: %v", err)
	}

	resizedImg := resize.Resize(width, height, img, resize.Lanczos3)
	return resizedImg, nil
}

func SaveUploadImage(r *http.Request) (string, error) {
	file, _, fromErr := r.FormFile("receipt")
	if fromErr != nil {
		return "", fromErr
	}
	defer file.Close()

	uploadPath := "./uploads/"
	fileName := uuid.New().String() + ".png"

	destFile, createErr := os.Create(filepath.Join(uploadPath, fileName))
	if createErr != nil {
		return "", createErr
	}
	defer destFile.Close()

	_, copyErr := io.Copy(destFile, file)
	if copyErr != nil {
		return "", copyErr
	}

	return fileName, nil
}

// func GenerateImages(inputFilePath string, outputDir string) error {
// 	inputFilePath := filepath.Join(uploadDir, fileName)

// 	smallImg, resizeSmallErr := ResizeImage(inputFilePath, constants.IMAGE_SIZE_W_S, constants.IMAGE_SIZE_H_S)
// 	if resizeSmallErr != nil {
// 		err := fmt.Errorf(
// 			"ResizeImage(inputFilePath: %s, width: %d, height: %d) failed, err: %s",
// 			inputFilePath, constants.IMAGE_SIZE_H_S, constants.IMAGE_SIZE_H_S, resizeSmallErr.Error(),
// 		)
// 		return err
// 	}
// 	smallOutputFilePath := filepath.Join(outputDir, fileName)

// 	saveImage(smallImg, outputFilePath)

// 	mediumImg, resizeErr2 := ResizeImage(inputFilePath, constants.IMAGE_SIZE_W_S, constants.IMAGE_SIZE_H_S)
// 	if resizeErr2 != nil {
// 		err := fmt.Errorf(
// 			"ResizeImage(inputFilePath: %s, width: %d, height: %d) failed, err: %s",
// 			inputFilePath, constants.IMAGE_SIZE_H_S, constants.IMAGE_SIZE_H_S, resizeErr2.Error(),
// 		)
// 		return err
// 	}
// 	saveImage(smallImg, outputFilePath)

// }

func AppendSuffix(filePath string, suffix string) string {
	extension := filepath.Ext(filePath)
	base := strings.TrimSuffix(filePath, extension)
	newFilename := fmt.Sprintf("%s_%s%s", base, suffix, extension)
	return newFilename
}
