package test_utils

import (
	"image"
	"image/color"
	"image/png"
	"os"
)

func CreateTestImage(filePath string) error {
	img := image.NewRGBA(image.Rect(0, 0, 800, 600))
	for x := 0; x < img.Bounds().Dx(); x++ {
		for y := 0; y < img.Bounds().Dy(); y++ {
			img.Set(x, y, color.RGBA{uint8(x % 256), uint8(y % 256), 0, 255})
		}
	}
	out, err := os.Create(filePath)
	if err != nil {
		return err
	}
	defer out.Close()
	return png.Encode(out, img)
}
