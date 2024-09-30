package image_worker

import (
	"fmt"
	"os"
	"path/filepath"
	"receipt_uploader/internal/images"
	"receipt_uploader/internal/models/image_meta"
)

type Service struct {
	ImageService images.ServiceType
}

func NewService(service images.ServiceType) ServiceType {
	return &Service{
		ImageService: service,
	}
}

func (s *Service) ResizeImages(srcDir, destDir string) error {
	fName, fileErr := getFile(srcDir)
	if fileErr != nil {
		return fmt.Errorf("getFile() failed, err: %w", fileErr)
	}
	if fName != "" {
		imageFile, imageErr := image_meta.FromUploadDir(filepath.Join(srcDir, fName))
		if imageErr != nil {
			return fmt.Errorf("metainfo.FromPath() failed, err: %w", imageErr)
		}

		genErr := s.ImageService.GenerateResizedImages(imageFile, destDir)
		if genErr != nil {
			return fmt.Errorf("s.ImageService.GenerateResizedImages() failed, err: %w", genErr)
		}
		removeErr := os.Remove(imageFile.Path)
		if removeErr != nil {
			return fmt.Errorf("os.Remove() failed, err: %w", removeErr)
		}
	}

	return nil
}

func getFile(dir string) (string, error) {
	entries, err := os.ReadDir(dir)
	if err != nil {
		return "", fmt.Errorf("os.ReadDir() failed, err: %w", err)
	}

	for _, entry := range entries {
		if !entry.IsDir() {
			return entry.Name(), nil
		}
	}

	return "", nil
}
