package images

import "net/http"

type ServiceType interface {
	GenerateImages(srcPath string) error
	SaveUpload(bytes []byte) (string, error)
	DecodeImage(r *http.Request) ([]byte, error)
}
