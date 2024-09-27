package constants

const (
	PORT             = ":8080"
	OUTPUT_DIR       = "images"
	HTTP_ERR_MSG_500 = "internal server error"
	HTTP_ERR_MSG_400 = "invalid image"
	HTTP_ERR_MSG_404 = "image not found"
	HTTP_ERR_MSG_405 = "method not allowed"
	IMAGE_SIZE_MIN_W = 600
	IMAGE_SIZE_MIN_H = 800
	IMAGE_SIZE_W_S   = 0 // set to 0 to keep original ratio of image
	IMAGE_SIZE_H_S   = 120
	IMAGE_SIZE_W_M   = 0 // set to 0 to keep original ratio of image
	IMAGE_SIZE_H_M   = 600
)
