package images

import (
	"bytes"
	"fmt"
	"image"
	"image/jpeg"
	"net/http"
	"os"
	"path/filepath"
	"receipt_uploader/internal/constants"
	"receipt_uploader/internal/logging"
	"receipt_uploader/internal/models/configs"
	"receipt_uploader/internal/models/http_requests"
	"receipt_uploader/internal/models/image_meta"

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

// GenerateResizedImages generates resized versions of an image based on
// specified dimensions and saves them to a given destination directory.
//
// This method reads the original image file specified by imageMeta.Path,
// creates a directory structure for the specified username, and then
// generates resized images according to predefined dimensions.
// The resized images are saved in the destination directory.
//
// Parameters:
//   - imageMeta: A pointer to an image_meta.ImageMeta struct that contains
//     metadata about the image, including its path and associated username.
//   - destDir: A string representing the base directory where the resized
//     images will be saved, organized under the user's name.
//
// Returns:
//   - An error if any of the following operations fail:
//
// Example:
//
//	func main() {
//	    service := &Service{}
//	    imageMeta := &image_meta.ImageMeta{
//	        Path: "path/to/original/image.jpg",
//	        Username: "john_doe",
//	    }
//	    err := service.GenerateResizedImages(imageMeta, "path/to/destination")
//	    if err != nil {
//	        log.Fatalf("Error generating resized images: %v", err)
//	    }
//	    log.Println("Resized images generated successfully.")
//	}
func (s *Service) GenerateResizedImages(imageMeta *image_meta.ImageMeta, destDir string) error {
	logging.Infof("GenerateResizedImages(srcPath: %s, destDir: %s)", imageMeta.Path, destDir)

	fileBytes, readErr := os.ReadFile(imageMeta.Path)
	if readErr != nil {
		return fmt.Errorf("os.ReadFile() failed: %v", readErr)
	}

	destDir = filepath.Join(destDir, imageMeta.Username)
	mkErr := os.MkdirAll(destDir, 0755)
	if mkErr != nil {
		err := fmt.Errorf("os.Mkdir() failed, err: %s", mkErr.Error())
		return err
	}

	copyDestPath := image_meta.GetResizedPath(imageMeta, destDir, "")
	logging.Debugf("copyDestPath: %s)", copyDestPath)

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
				imageMeta.Path, d.Width, d.Height, resizeErr.Error(),
			)
		}

		destPath := image_meta.GetResizedPath(imageMeta, destDir, d.Name)
		logging.Debugf("destPath: %s", destPath)
		saveErr := saveImage(&resizedImg, destPath)
		if saveErr != nil {
			return fmt.Errorf("saveImage(destPath: %s) failed, err: %s", destPath, saveErr.Error())
		}
	}

	return nil
}

// ParseImage processes an HTTP request to extract and validate an uploaded image file.
//
// This method parses a multipart form from the incoming HTTP request and attempts to
// retrieve the file associated with the form field named "receipt". It validates the
// content type, image format, and dimensions of the image based on predefined constants.
//
// The function performs the following tasks:
// - Parses the incoming multipart form data with a specified maximum upload size.
// - Retrieves the uploaded file from the form.
// - Reads the file's content and decodes it to check if it is a valid image.
// - Validates the image format to ensure it is a JPEG image.
// - Validates the dimensions of the image against specified minimum width and height.
//
// Parameters:
//   - r: A pointer to an http.Request that contains the uploaded image
//     in a multipart form data format.
//
// Returns:
//   - A byte slice containing the raw image data if successful.
//   - An error if any of the following fail:
//
// Example:
//
//	func handler(w http.ResponseWriter, r *http.Request) {
//	    service := &Service{}
//	    imgData, err := service.ParseImage(r)
//	    if err != nil {
//	        http.Error(w, err.Error(), http.StatusBadRequest)
//	        return
//	    }
//	    // Process imgData as needed
//	}
func (s *Service) ParseImage(r *http.Request) ([]byte, error) {
	parseErr := r.ParseMultipartForm(constants.MAX_UPLOAD_SIZE)
	if parseErr != nil {
		return nil, fmt.Errorf("r.ParseMultipartForm() failed, err: %s", parseErr.Error())
	}
	uploadRequest, reqErr := http_requests.ParseUploadRequest(r)
	if reqErr != nil {
		return nil, fmt.Errorf("http_requests.FromRequest() failed: %w", reqErr)
	}

	payloadSize := int64(len(uploadRequest.Payload))
	if payloadSize > constants.MAX_UPLOAD_SIZE {
		return nil, fmt.Errorf("image size is too big, payloadSize=%d", payloadSize)
	}

	img, format, decodeErr := image.Decode(bytes.NewReader(uploadRequest.Payload))
	if decodeErr != nil {
		return nil, fmt.Errorf("image.Decode() failed, err: %s", decodeErr.Error())
	}
	if img.Bounds().Dx() < constants.IMAGE_SIZE_MIN_W || img.Bounds().Dy() < constants.IMAGE_SIZE_MIN_H {
		return nil, fmt.Errorf("invalid image size, minHeight=%d, minWidth=%d", constants.IMAGE_SIZE_MIN_H, constants.IMAGE_SIZE_MIN_W)
	}

	if format != "jpeg" {
		return nil, fmt.Errorf("invalid image format, format=%s, only jpeg format is allowed", format)
	}

	return uploadRequest.Payload, nil
}

func (s *Service) SaveUpload(bytes *[]byte, username, uploadDir string) (*image_meta.ImageMeta, error) {
	logging.Debugf("SaveUpload(len(bytes): %d, uploadDir: %s)", len(*bytes), uploadDir)

	mkErr := os.MkdirAll(uploadDir, 0755)
	if mkErr != nil {
		err := fmt.Errorf("os.Mkdir() failed, err: %s", mkErr.Error())
		return nil, err
	}

	extension := "jpg"
	imageMeta := image_meta.FromFormData(username, extension, uploadDir)
	saveImage(bytes, imageMeta.Path)

	return imageMeta, nil
}

func (s *Service) GetImage(imageMeta *image_meta.ImageMeta) ([]byte, string, error) {
	logging.Debugf("GetImage(imageMeta.Path: %s, filaName: %s)", imageMeta.Path, imageMeta.FileName)

	fileBytes, readErr := os.ReadFile(imageMeta.Path)
	if readErr != nil {
		if os.IsNotExist(readErr) {
			return nil, "", readErr
		}
		return nil, "", fmt.Errorf("os.ReadFile() failed: %v", readErr)
	}
	return fileBytes, imageMeta.FileName, nil
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
