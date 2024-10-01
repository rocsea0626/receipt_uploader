package image_worker

import (
	"fmt"
	"os"
	"path/filepath"
	"receipt_uploader/internal/images"
	"receipt_uploader/internal/logging"
	"receipt_uploader/internal/models/configs"
	"receipt_uploader/internal/models/image_meta"
	"sync"
	"time"
)

type Service struct {
	Interval     time.Duration
	SrcDir       string
	DestDir      string
	ImageService images.ServiceType
}

func NewService(config *configs.Config, service images.ServiceType) ServiceType {
	return &Service{
		Interval:     config.Interval,
		SrcDir:       config.UploadsDir,
		DestDir:      config.ResizedDir,
		ImageService: service,
	}
}

func (s *Service) Start(stopChan <-chan struct{}) {
	fmt.Println("starting worker...")
	var wg sync.WaitGroup

	for {
		select {
		case <-stopChan:
			fmt.Println("Stopping image worker...")
			wg.Wait()
			return
		case <-time.After(s.Interval):
			logging.Debugf("Wake up after sleeping, starting processing...")

			wg.Add(1)
			defer wg.Done()
			resizeErr := resizeImages(s.SrcDir, s.DestDir, s.ImageService)
			if resizeErr != nil {
				logging.Errorf("resizeImages() failed, err: %s", resizeErr.Error())
			}
		}
	}
}

func resizeImages(srcDir, destDir string, imagesService images.ServiceType) error {
	logging.Debugf("resizeImages(srcDir: %s, destDir: %s)", srcDir, destDir)

	fName, fileErr := getFirstFile(srcDir)
	if fileErr != nil {
		return fmt.Errorf("getFirstFile() failed, err: %w", fileErr)
	}
	if fName != "" {
		startTime := time.Now()
		imageMeta, imageErr := image_meta.FromUploadDir(filepath.Join(srcDir, fName))
		if imageErr != nil {
			return fmt.Errorf("metainfo.FromPath() failed, err: %w", imageErr)
		}

		genErr := imagesService.GenerateResizedImages(imageMeta, destDir)
		if genErr != nil {
			return fmt.Errorf("s.ImageService.GenerateResizedImages() failed, err: %w", genErr)
		}
		removeErr := os.Remove(imageMeta.Path)
		if removeErr != nil {
			return fmt.Errorf("os.Remove() failed, err: %w", removeErr)
		}
		elapsedTime := time.Since(startTime)
		logging.Infof("resizeImages() completes with %d ms", elapsedTime.Milliseconds())
	}

	return nil
}

// getFirstFile scans the specified directory for files and returns the name of the first file found.
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
func getFirstFile(dir string) (string, error) {
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
