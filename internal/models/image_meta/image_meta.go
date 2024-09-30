package image_meta

import (
	"errors"
	"fmt"
	"path/filepath"
	"receipt_uploader/internal/logging"
	"strings"

	"github.com/google/uuid"
)

type ImageMeta struct {
	Path      string `json:"path"`
	Dir       string `json:"dir"`
	FileName  string `json:"fileName"`
	Extension string `json:"extension"` // i.e. ".jpg"
	Username  string `json:"username"`
	ReceiptID string `json:"receiptId"`
}

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

func FromFormData(username, extension, uploadDir string) *ImageMeta {

	receiptId := strings.ReplaceAll(uuid.New().String(), "-", "")
	fileName := username + "#" + receiptId + "." + extension
	path := filepath.Join(uploadDir, fileName)

	imgFile, _ := FromUploadDir(path)
	return imgFile
}

func FromGetRequset(receiptID, size, username, srcDir string) *ImageMeta {
	logging.Debugf("FromGetRequset(receiptID: %s, size: %s, username: %s, srcDir: %s)", receiptID, size, username, srcDir)

	extension := ".jpg"
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

func GetResizedPath(imgFile *ImageMeta, destDir, size string) string {
	logging.Debugf("GetResizedPath(destDir: %s, ext: %s)", destDir, imgFile.Extension)

	newFilename := fmt.Sprintf("%s%s", imgFile.ReceiptID, imgFile.Extension)
	if size != "" {
		newFilename = fmt.Sprintf("%s_%s%s", imgFile.ReceiptID, size, imgFile.Extension)
	}

	fPath := filepath.Join(destDir, newFilename)
	return fPath
}
