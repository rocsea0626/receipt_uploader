package image_worker

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"receipt_uploader/internal/constants"
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
	Timeout      time.Duration
}

func NewService(config *configs.Config, service images.ServiceType) ServiceType {
	return &Service{
		Interval:     config.Interval,
		SrcDir:       config.UploadsDir,
		DestDir:      config.ResizedDir,
		ImageService: service,
		Timeout:      constants.IMAGE_WORKER_TIMEOUT,
	}
}

func (s *Service) Start(stopChan <-chan struct{}) {
	fmt.Println("starting worker...")
	var wg sync.WaitGroup
	resizeChan := make(chan struct{}, 1)

	for {
		select {
		case <-stopChan:
			fmt.Println("Stopping image worker...")
			wg.Wait()
			return
		case <-time.After(s.Interval):
			logging.Infof("Wake up after sleeping, starting processing...")

			select {
			case resizeChan <- struct{}{}: // allows only 1 instance of resizeImages() running
				wg.Add(1)
				go func() {
					defer func() {
						<-resizeChan
						wg.Done()
					}()

					resizeErr := resizeImages(s.SrcDir, s.DestDir, s.Timeout, s.ImageService)
					if resizeErr != nil {
						logging.Errorf("resizeImages() failed, err: %s", resizeErr.Error())
					}
					// time.Sleep(500 * time.Millisecond)
				}()
			default:
				logging.Infof("resizeImages() is already in progress, skipping this cycle...")
			}
		}
	}
}

func resizeImages(srcDir, destDir string, timeout time.Duration, imagesService images.ServiceType) error {
	logging.Debugf("resizeImages(srcDir: %s, destDir: %s, timeout: %v)", srcDir, destDir, timeout)

	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	errChan := make(chan error, 1)

	go func() {
		defer func() {
			close(errChan)
		}()

		fName, fileErr := getFirstFile(srcDir)
		if fileErr != nil {
			errChan <- fmt.Errorf("getFirstFile() failed, err: %w", fileErr)
		}
		if fName != "" {
			startTime := time.Now()
			imageMeta, imageErr := image_meta.FromUploadDir(filepath.Join(srcDir, fName))
			if imageErr != nil {
				errChan <- fmt.Errorf("metainfo.FromPath() failed, err: %w", imageErr)
				return
			}

			genErr := imagesService.GenerateResizedImages(imageMeta, destDir)
			if genErr != nil {
				errChan <- fmt.Errorf("s.ImageService.GenerateResizedImages() failed, err: %w", genErr)
				return
			}
			removeErr := os.Remove(imageMeta.Path)
			if removeErr != nil {
				errChan <- fmt.Errorf("os.Remove() failed, err: %w", removeErr)
				return
			}

			elapsedTime := time.Since(startTime)
			logging.Infof("resizeImages() completes with %d ms", elapsedTime.Milliseconds())
		}

	}()

	select {
	case genErr := <-errChan:
		return genErr
	case <-ctx.Done():
		return fmt.Errorf("resizeImages() timed out")
	}
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
