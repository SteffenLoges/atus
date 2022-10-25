package helpers

import (
	"image"
	"image/color"
	"image/draw"
	"image/jpeg"
	"os"

	"golang.org/x/image/font"
	"golang.org/x/image/font/inconsolata"
	"golang.org/x/image/math/fixed"
)

type Image struct {
	*image.RGBA
}

func CreateSquareImage(size int) *Image {

	img := &Image{
		image.NewRGBA(image.Rect(0, 0, size, size)),
	}

	draw.Draw(img, img.Bounds(), image.White, image.Point{0, 0}, draw.Src)

	return img

}

func (img *Image) AddLabel(label string, c color.RGBA, offsetY int) {

	imageWidth := img.Bounds().Max.X
	imageHeight := img.Bounds().Max.Y

	// calc label position
	labelPosX := (imageWidth / 2) - (len(label) * (8 / 2))
	labelPosY := ((imageHeight / 2) - (16 / 2)) + offsetY

	point := fixed.Point26_6{
		X: fixed.I(labelPosX),
		Y: fixed.I(labelPosY),
	}

	d := &font.Drawer{
		Dst:  img,
		Src:  image.NewUniform(c),
		Face: inconsolata.Bold8x16,
		Dot:  point,
	}

	d.DrawString(label)

}

func (img *Image) Save(savePath string) error {

	f, err := os.Create(savePath)
	if err != nil {
		return err
	}
	defer f.Close()

	if err := jpeg.Encode(f, img, nil); err != nil {
		return err
	}

	return nil

}
