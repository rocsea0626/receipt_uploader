package images

import (
	"fmt"
	"image"
	"image/jpeg"
	"log"
	"os"
	"receipt_uploader/constants"
	"receipt_uploader/internal/futils"

	"github.com/nfnt/resize"
)

func GenerateImages(srcPath string, destDir string) error {
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

func resizeImage(srcPath string, width, height int) (image.Image, error) {
	file, err := os.Open(srcPath)
	if err != nil {
		return nil, fmt.Errorf("os.Open() failed: %v", err)
	}
	defer file.Close()

	img, _, err := image.Decode(file)
	if err != nil {
		return nil, fmt.Errorf("image.Decode() failed: %v", err)
	}

	resizedImg := resize.Resize(uint(width), uint(height), img, resize.Lanczos3)
	return resizedImg, nil
}

func saveImage(img *image.Image, destPath string) error {
	outFile, createErr := os.Create(destPath)
	if createErr != nil {
		return fmt.Errorf("os.Create() failed, err: %s", createErr.Error())
	}
	defer outFile.Close()

	EncodeErr := jpeg.Encode(outFile, *img, nil)
	if EncodeErr != nil {
		return fmt.Errorf("jpeg.Encode() failed, err: %s", EncodeErr.Error())
	}

	return nil
}
