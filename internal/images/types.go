package images

import "net/http"

type ServiceType interface {
	GenerateImages(srcPath, destDir string) error
	SaveUpload(bytes []byte, destDir string) (string, error)
	DecodeImage(r *http.Request) ([]byte, error)
	GetImage(receiptId, size, srcDir string) ([]byte, string, error)
}
