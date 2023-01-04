package imageprocessor

import (
	"bytes"
	"fmt"
	"image"
	"image/color"

	"gopkg.in/fogleman/gg.v1"
)

type ImageProcessor struct {
}

func (i ImageProcessor) ProcessImageEntry(author, title string, data []byte) ([]byte, error) {
	source, _, err := image.Decode(bytes.NewReader(data))
	if err != nil {
		return []byte{}, err
	}

	bounds := source.Bounds()
	percent_of_the_bounds := float64(bounds.Dy()) * (float64(3) / float64(100))
	rgba_canvas := image.NewRGBA(image.Rect(0, 0, bounds.Dx(), bounds.Dy()+int(percent_of_the_bounds)))

	gg_context := gg.NewContextForRGBA(rgba_canvas)
	gg_context.DrawImage(source, 0, 0)
	gg_context.DrawRectangle(float64(source.Bounds().Min.X), float64(source.Bounds().Max.Y), float64(bounds.Dx()), percent_of_the_bounds)
	gg_context.SetColor(color.Black)
	gg_context.Fill()

	if err != nil {
		return []byte{}, err
	}
	gg_context.LoadFontFace("./cmd/Roboto-Light.ttf", 36)
	gg_context.SetColor(color.White)
	gg_context.DrawString(fmt.Sprintf("%s - %s", author, title), 0, float64(bounds.Dy())+(percent_of_the_bounds/2)+gg_context.FontHeight()/2)

	buffer := new(bytes.Buffer)
	err = gg_context.EncodePNG(buffer)
	return buffer.Bytes(), err
}

func GetNewImageProcesor() ImageProcessor {
	return ImageProcessor{}
}
