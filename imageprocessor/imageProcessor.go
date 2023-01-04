package imageprocessor

import (
	"bytes"
	"fmt"
	"image"
	"image/draw"
	"image/png"

	"golang.org/x/image/font"
	"golang.org/x/image/font/basicfont"
	"golang.org/x/image/math/fixed"
)

type ImageProcessor struct {
}

func (i ImageProcessor) ProcessImageEntry(author, title string, data []byte) ([]byte, error) {
	source, _, err := image.Decode(bytes.NewReader(data))
	if err != nil {
		return []byte{}, err
	}

	// first, we create a mutable canvas containing the base image and is also a bit larger
	// so that we can fit in the text
	bounds := source.Bounds()
	percent_of_the_bounds := float32(bounds.Dy()) * (float32(3) / float32(100))
	rgba_canvas := image.NewRGBA(image.Rect(0, 0, bounds.Dx(), bounds.Dy()+int(percent_of_the_bounds)))
	// then fill it with black color and copy the downloaded image onto it
	draw.Draw(rgba_canvas, rgba_canvas.Bounds(), image.NewUniform(image.Black), image.Point{}, draw.Src)
	draw.Draw(rgba_canvas, source.Bounds(), source, bounds.Min, draw.Src)
	// and then draw the text
	point := fixed.Point26_6{X: fixed.I(0), Y: fixed.I(bounds.Dy() + int(percent_of_the_bounds/2))}

	drawer := &font.Drawer{
		Dst:  rgba_canvas,
		Src:  image.NewUniform(image.White),
		Face: basicfont.Face7x13,
		Dot:  point,
	}

	drawer.DrawString(fmt.Sprintf("%s - %s", author, title))

	buffer := new(bytes.Buffer)
	err = png.Encode(buffer, rgba_canvas)
	return buffer.Bytes(), err
}

func GetNewImageProcesor() ImageProcessor {
	return ImageProcessor{}
}
