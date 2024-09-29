package configs

import "receipt_uploader/constants"

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
		Height: 120,
		Name:   "small",
	}, {
		Width:  0,
		Height: 600,
		Name:   "medium",
	}, {
		Width:  0,
		Height: constants.IMAGE_SIZE_MIN_H,
		Name:   "large",
	},
}

type Config struct {
	ImagesDir  string // dir to store images
	UploadsDir string // dir to store uploads
	Port       string
	Dimensions Dimensions
}
