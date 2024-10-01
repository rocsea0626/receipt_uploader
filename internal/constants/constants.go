package constants

import "time"

const (
	PORT                 = ":8080"
	ROOT_DIR_IMAGES      = "receipts"              // root dir to store all uplaoded and converted photos
	MAX_UPLOAD_SIZE      = int64(10 * 1024 * 1024) // Maximum 10 MB
	HTTP_ERR_MSG_500     = "internal server error"
	HTTP_ERR_MSG_400     = "invalid image"
	HTTP_ERR_MSG_403     = "access forbidden"
	HTTP_ERR_MSG_404     = "image not found"
	HTTP_ERR_MSG_405     = "method not allowed"
	IMAGE_SIZE_MIN_W     = 600
	IMAGE_SIZE_MIN_H     = 800
	IMAGE_WORKER_TIMEOUT = 2 * time.Second
	QUEUE_CAPACITY       = 100
)
