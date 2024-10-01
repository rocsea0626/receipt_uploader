package tasks

import "receipt_uploader/internal/models/image_meta"

type ResizeTask struct {
	ImageMeta image_meta.ImageMeta
	DestDir   string
}
