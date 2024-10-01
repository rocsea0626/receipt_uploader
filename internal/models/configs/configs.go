package configs

import (
	"receipt_uploader/internal/constants"
	"time"
)

// defines resized image's size and name of the size
type Dimension struct {
	Width  int
	Height int
	Name   string // small, medium, large
}

type Dimensions []Dimension

var AllowedDimensions = Dimensions{
	{
		Width:  0, // set to 0 to keep original ratio of image
		Height: constants.IMAGE_SIZE_MIN_H - 680,
		Name:   "small",
	}, {
		Width:  0,
		Height: constants.IMAGE_SIZE_MIN_H - 200,
		Name:   "medium",
	}, {
		Width:  0,
		Height: constants.IMAGE_SIZE_MIN_H,
		Name:   "large",
	},
}

type Config struct {
	ResizedDir    string // dir to store resize images
	UploadsDir    string // dir to store uploads
	Port          string
	Dimensions    Dimensions    // allowed resizing options
	Interval      time.Duration // Interval for worker in seconds
	Mode          string        // dev, qa, release
	QueueCapacity int
}
