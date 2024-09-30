package image_meta

import (
	"errors"
	"fmt"
	"path/filepath"
	"receipt_uploader/internal/logging"
	"strings"

	"github.com/google/uuid"
)

// ImageMeta represents the metadata associated with an image file.
type ImageMeta struct {
	Path      string `json:"path"`      // Full path to the image file
	Dir       string `json:"dir"`       // Directory containing the image file
	FileName  string `json:"fileName"`  // The name of the image file including extension
	Extension string `json:"extension"` // File extension (e.g., ".jpg")
	Username  string `json:"username"`  // Username of the uploader
	ReceiptID string `json:"receiptId"` // Unique identifier for the receipt
}

// FromUploadDir creates an ImageMeta object from the full file path of an uploaded image in config.DIR_UPLOADS folder
func FromUploadDir(fullPath string) (*ImageMeta, error) {
	fileName := filepath.Base(fullPath)
	extension := filepath.Ext(fullPath)
	usernameAndUuid := strings.TrimSuffix(fileName, extension)
	tokens := strings.Split(usernameAndUuid, "#")
	if len(tokens) != 2 {
		return nil, errors.New("invalid path, missing username")
	}

	return &ImageMeta{
		Path:      fullPath,
		Dir:       filepath.Dir(fullPath),
		FileName:  fileName,
		Extension: extension,
		Username:  tokens[0],
		ReceiptID: tokens[1],
	}, nil
}

// FromFormData creates an ImageMeta object from upload request, based on the provided username, extension,
// and config.DIR_UPLOADS directory
func FromFormData(username, extension, uploadDir string) *ImageMeta {
	receiptId := strings.ReplaceAll(uuid.New().String(), "-", "")
	fileName := username + "#" + receiptId + "." + extension
	path := filepath.Join(uploadDir, fileName)

	imgFile, _ := FromUploadDir(path) // Ignoring error here as it's a new path
	return imgFile
}

// FromGetRequset constructs an ImageMeta object from the provided receiptID, size,
// username in GET request and source directory.
func FromGetRequset(receiptID, size, username, srcDir string) *ImageMeta {
	logging.Debugf("FromGetRequset(receiptID: %s, size: %s, username: %s, srcDir: %s)", receiptID, size, username, srcDir)

	extension := ".jpg" // Defaulting to JPEG format
	fName := receiptID + extension
	if size != "" {
		fName = receiptID + "_" + size + extension
	}

	return &ImageMeta{
		Path:      filepath.Join(srcDir, username, fName),
		Dir:       filepath.Join(srcDir, username),
		Extension: extension,
		ReceiptID: receiptID,
		Username:  username,
		FileName:  fName,
	}
}

// GetResizedPath constructs a file path for a resized image based on its metadata.
// If a size is specified, it appends the size to the file name, otherwise it uses the original receiptID.
func GetResizedPath(imgFile *ImageMeta, destDir, size string) string {
	logging.Debugf("GetResizedPath(destDir: %s, ext: %s)", destDir, imgFile.Extension)

	newFilename := fmt.Sprintf("%s%s", imgFile.ReceiptID, imgFile.Extension)
	if size != "" {
		newFilename = fmt.Sprintf("%s_%s%s", imgFile.ReceiptID, size, imgFile.Extension)
	}

	fPath := filepath.Join(destDir, newFilename)
	return fPath
}
