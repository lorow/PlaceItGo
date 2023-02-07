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
	percentOfTheBounds := float64(bounds.Dy()) * (float64(3) / float64(100))
	rgbaCanvas := image.NewRGBA(image.Rect(0, 0, bounds.Dx(), bounds.Dy()+int(percentOfTheBounds)))

	ggContext := gg.NewContextForRGBA(rgbaCanvas)
	ggContext.DrawImage(source, 0, 0)
	ggContext.DrawRectangle(float64(source.Bounds().Min.X), float64(source.Bounds().Max.Y), float64(bounds.Dx()), percentOfTheBounds)
	ggContext.SetColor(color.Black)
	ggContext.Fill()

	if err != nil {
		return []byte{}, err
	}
	ggContext.LoadFontFace("./cmd/Roboto-Light.ttf", 36)
	ggContext.SetColor(color.White)
	ggContext.DrawString(fmt.Sprintf("%s - %s", author, title), 0, float64(bounds.Dy())+(percentOfTheBounds/2)+ggContext.FontHeight()/2)

	buffer := new(bytes.Buffer)
	err = ggContext.EncodePNG(buffer)
	return buffer.Bytes(), err
}

func GetNewImageProcessor() ImageProcessor {
	return ImageProcessor{}
}
