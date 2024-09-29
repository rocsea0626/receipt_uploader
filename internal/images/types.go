package images

import "net/http"

type ServiceType interface {
	GenerateImages(bytes *[]byte, srcPath, destDir string) error
	SaveUpload(bytes *[]byte, destDir string) (string, error)
	ParseImage(r *http.Request) ([]byte, error)
	GetImage(receiptId, size, srcDir string) ([]byte, string, error)
}
