package images

import (
	"net/http"
	"receipt_uploader/internal/models/image_meta"
)

type ServiceType interface {
	GenerateResizedImages(imageMeta *image_meta.ImageMeta, destDir string) error
	SaveUpload(bytes *[]byte, username, destDir string) (*image_meta.ImageMeta, error)
	ParseImage(r *http.Request) ([]byte, error)
	GetImage(imageMeta *image_meta.ImageMeta) ([]byte, string, error)
}
