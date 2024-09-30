package image_worker

import (
	"fmt"
	"os"
	"path/filepath"
	"receipt_uploader/internal/images"
	"receipt_uploader/internal/logging"
	"receipt_uploader/internal/models/image_meta"
	"time"
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
	logging.Debugf("ResizeImages(srcDir: %s, destDir: %s)", srcDir, destDir)

	fName, fileErr := getFile(srcDir)
	if fileErr != nil {
		return fmt.Errorf("getFile() failed, err: %w", fileErr)
	}
	if fName != "" {
		startTime := time.Now()
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
		elapsedTime := time.Since(startTime)
		logging.Infof("ResizeImages() completes with %d ms", elapsedTime.Milliseconds())
	}

	return nil
}

// getFile scans the specified directory for files and returns the name of the first file found.
//
// Parameters:
//   - dir: A string representing the path to the directory to be scanned.
//
// Returns:
//   - A string containing the name of the first file found in the directory.
//   - An error indicating any problems encountered while reading the directory.
//
// Example:
//
//	fileName, err := getFile("/path/to/directory")
//	if err != nil {
//	    log.Fatalf("Error: %v", err)
//	}
//	fmt.Printf("Found file: %s\n", fileName)
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
