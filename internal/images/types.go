package images

import "net/http"

type ServiceType interface {
	GenerateImages(srcPath string, destDir string) error
	SaveUpload(bytes []byte, tmpDir string) (string, error)
	DecodeImage(r *http.Request) ([]byte, error)
}
