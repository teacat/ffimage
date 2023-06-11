package ffimage

import (
	"fmt"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func newImage(a *assert.Assertions, name string) *Image {
	img, err := NewImage(fmt.Sprintf("./test/%s", name))
	a.NoError(err)

	return img
}

func newOutput(name string) string {
	return fmt.Sprintf("./test/output/%s", name)
}

func TestMain(test *testing.T) {
	a := assert.New(test)

	a.NoError(os.MkdirAll("./test/output", 0644))
}

func TestWidthHeight(test *testing.T) {
	a := assert.New(test)
	img := newImage(a, "source.png")

	a.Equal(431, img.GetWidth())
	a.Equal(324, img.GetHeight())
}

func TestResize(test *testing.T) {
	a := assert.New(test)
	img, output := newImage(a, "source.png"), newOutput("resize-300x300.png")

	// 300x300
	err := img.ResizeImage(300, 300).WriteImage(output)
	a.NoError(err)

	img, err = NewImage(output)
	a.NoError(err)

	a.Equal(300, img.GetWidth())
	a.Equal(300, img.GetHeight())
}

func TestResizeWidth(test *testing.T) {
	a := assert.New(test)
	img, output := newImage(a, "source.png"), newOutput("resize-300x.png")

	// 300x
	err := img.ResizeImage(300, 0).WriteImage(output)
	a.NoError(err)

	img, err = NewImage(output)
	a.NoError(err)

	a.Equal(300, img.GetWidth())
	a.Equal(225, img.GetHeight())
}

func TestResizeHeight(test *testing.T) {
	a := assert.New(test)
	img, output := newImage(a, "source.png"), newOutput("resize-x300.png")

	// x300
	err := img.ResizeImage(0, 300).WriteImage(output)
	a.NoError(err)

	img, err = NewImage(output)
	a.NoError(err)

	a.Equal(399, img.GetWidth())
	a.Equal(300, img.GetHeight())
}

func TestResizeBestfitDownscale(test *testing.T) {
	a := assert.New(test)
	img, output := newImage(a, "source.png"), newOutput("resize-300x300-downscale.png")

	// 300x300
	err := img.ResizeImage(300, 300, ResizeTypeDownscale).WriteImage(output)
	a.NoError(err)

	img, err = NewImage(output)
	a.NoError(err)

	a.Equal(300, img.GetWidth())
	a.Equal(225, img.GetHeight())
}

func TestResizeBestfitUpscale(test *testing.T) {
	a := assert.New(test)
	img, output := newImage(a, "source.png"), newOutput("resize-300x300-upscale.png")

	// 300x300
	err := img.ResizeImage(300, 300, ResizeTypeUpscale).WriteImage(output)
	a.NoError(err)

	img, err = NewImage(output)
	a.NoError(err)

	a.Equal(399, img.GetWidth())
	a.Equal(300, img.GetHeight())
}

func TestExtentImage(test *testing.T) {
	a := assert.New(test)
	img, output := newImage(a, "source.png"), newOutput("extent-500x500.png")
	err := img.ExtentImage(500, 500, 0, 0).WriteImage(output)
	a.NoError(err)

	img, output = newImage(a, "source.png"), newOutput("extent-top_left-500x500.png")
	err = img.ExtentImage(500, 500, 0, 0, PositionTypeTopLeft).WriteImage(output)
	a.NoError(err)

	img, output = newImage(a, "source.png"), newOutput("extent-top-500x500.png")
	err = img.ExtentImage(500, 500, 0, 0, PositionTypeTop).WriteImage(output)
	a.NoError(err)

	img, output = newImage(a, "source.png"), newOutput("extent-top_right-500x500.png")
	err = img.ExtentImage(500, 500, 0, 0, PositionTypeTopRight).WriteImage(output)
	a.NoError(err)

	img, output = newImage(a, "source.png"), newOutput("extent-left-500x500.png")
	err = img.ExtentImage(500, 500, 0, 0, PositionTypeLeft).WriteImage(output)
	a.NoError(err)

	img, output = newImage(a, "source.png"), newOutput("extent-center-500x500.png")
	err = img.ExtentImage(500, 500, 0, 0, PositionTypeCenter).WriteImage(output)
	a.NoError(err)

	img, output = newImage(a, "source.png"), newOutput("extent-right-500x500.png")
	err = img.ExtentImage(500, 500, 0, 0, PositionTypeRight).WriteImage(output)
	a.NoError(err)

	img, output = newImage(a, "source.png"), newOutput("extent-bottom_left-500x500.png")
	err = img.ExtentImage(500, 500, 0, 0, PositionTypeBottomLeft).WriteImage(output)
	a.NoError(err)

	img, output = newImage(a, "source.png"), newOutput("extent-bottom-500x500.png")
	err = img.ExtentImage(500, 500, 0, 0, PositionTypeBottom).WriteImage(output)
	a.NoError(err)

	img, output = newImage(a, "source.png"), newOutput("extent-bottom_right-500x500.png")
	err = img.ExtentImage(500, 500, 0, 0, PositionTypeBottomRight).WriteImage(output)
	a.NoError(err)

	img, err = NewImage(output)
	a.NoError(err)

	a.Equal(500, img.GetWidth())
	a.Equal(500, img.GetHeight())
}

func TestExtentImageSmaller(test *testing.T) {
	a := assert.New(test)
	img, output := newImage(a, "source.png"), newOutput("extent-100x100.png")

	err := img.ExtentImage(100, 100, 0, 0).WriteImage(output)
	a.Error(err)

	img, err = NewImage(output)
	a.Error(err)
}

func TestRotateImage(test *testing.T) {
	a := assert.New(test)
	img, output := newImage(a, "source.png"), newOutput("rotate-90deg.png")
	err := img.RotateImage(90).WriteImage(output)
	a.NoError(err)

	img, output = newImage(a, "source.png"), newOutput("rotate-30deg.png")
	err = img.RotateImage(30).WriteImage(output)
	a.NoError(err)

	img, output = newImage(a, "source.png"), newOutput("rotate-180deg.png")
	err = img.RotateImage(180).WriteImage(output)
	a.NoError(err)

	img, output = newImage(a, "source.png"), newOutput("rotate-360deg.png")
	err = img.RotateImage(360).WriteImage(output)
	a.NoError(err)

	img, output = newImage(a, "source.png"), newOutput("rotate-720deg.png")
	err = img.RotateImage(720).WriteImage(output)
	a.NoError(err)

	img, err = NewImage(output)
	a.NoError(err)

	a.Equal(431, img.GetWidth())
	a.Equal(324, img.GetHeight())
}

func TestFlopFlipImage(test *testing.T) {
	a := assert.New(test)
	img, output := newImage(a, "source.png"), newOutput("flop.png")
	err := img.FlopImage().WriteImage(output)
	a.NoError(err)

	img, output = newImage(a, "source.png"), newOutput("flip.png")
	err = img.FlipImage().WriteImage(output)
	a.NoError(err)

	img, output = newImage(a, "source.png"), newOutput("flop_flip.png")
	err = img.FlopImage().FlipImage().WriteImage(output)
	a.NoError(err)

	img, err = NewImage(output)
	a.NoError(err)

	a.Equal(431, img.GetWidth())
	a.Equal(324, img.GetHeight())
}

func TestCropImage(test *testing.T) {
	a := assert.New(test)

	img, output := newImage(a, "source.png"), newOutput("crop-200x200.png")
	err := img.CropImage(200, 200, 0, 0).WriteImage(output)
	a.NoError(err)

	img, output = newImage(a, "source.png"), newOutput("crop-top_left-200x200.png")
	err = img.CropImage(200, 200, 0, 0, PositionTypeTopLeft).WriteImage(output)
	a.NoError(err)

	img, output = newImage(a, "source.png"), newOutput("crop-top-200x200.png")
	err = img.CropImage(200, 200, 0, 0, PositionTypeTop).WriteImage(output)
	a.NoError(err)

	img, output = newImage(a, "source.png"), newOutput("crop-top_right-200x200.png")
	err = img.CropImage(200, 200, 0, 0, PositionTypeTopRight).WriteImage(output)
	a.NoError(err)

	img, output = newImage(a, "source.png"), newOutput("crop-left-200x200.png")
	err = img.CropImage(200, 200, 0, 0, PositionTypeLeft).WriteImage(output)
	a.NoError(err)

	img, output = newImage(a, "source.png"), newOutput("crop-center-200x200.png")
	err = img.CropImage(200, 200, 0, 0, PositionTypeCenter).WriteImage(output)
	a.NoError(err)

	img, output = newImage(a, "source.png"), newOutput("crop-right-200x200.png")
	err = img.CropImage(200, 200, 0, 0, PositionTypeRight).WriteImage(output)
	a.NoError(err)

	img, output = newImage(a, "source.png"), newOutput("crop-bottom_left-200x200.png")
	err = img.CropImage(200, 200, 0, 0, PositionTypeBottomLeft).WriteImage(output)
	a.NoError(err)

	img, output = newImage(a, "source.png"), newOutput("crop-bottom-200x200.png")
	err = img.CropImage(200, 200, 0, 0, PositionTypeBottom).WriteImage(output)
	a.NoError(err)

	img, output = newImage(a, "source.png"), newOutput("crop-bottom_right-200x200.png")
	err = img.CropImage(200, 200, 0, 0, PositionTypeBottomRight).WriteImage(output)
	a.NoError(err)

	img, err = NewImage(output)
	a.NoError(err)

	a.Equal(200, img.GetWidth())
	a.Equal(200, img.GetHeight())
}

func TestCropImageLarger(test *testing.T) {
	a := assert.New(test)
	img, output := newImage(a, "source.png"), newOutput("crop-768x768x0x0.png")

	err := img.CropImage(768, 768, 0, 0, PositionTypeCenter).WriteImage(output)
	a.Error(err)

	img, err = NewImage(output)
	a.Error(err)
}

func TestCropThumbnailImage(test *testing.T) {
	a := assert.New(test)
	img, output := newImage(a, "source.png"), newOutput("crop_thumbnail-300x300.png")

	err := img.CropThumbnailImage(300, 300).WriteImage(output)
	a.NoError(err)

	img, err = NewImage(output)
	a.NoError(err)

	a.Equal(300, img.GetWidth())
	a.Equal(300, img.GetHeight())
}

func TestThumbnailImage(test *testing.T) {
	a := assert.New(test)
	img, output := newImage(a, "source.png"), newOutput("thumbnail-300x300.png")

	err := img.ThumbnailImage(300, 300).WriteImage(output)
	a.NoError(err)

	img, err = NewImage(output)
	a.NoError(err)

	a.Equal(300, img.GetWidth())
	a.Equal(300, img.GetHeight())
}

func TestThumbnailImageBackgroundColor(test *testing.T) {
	a := assert.New(test)
	img, output := newImage(a, "source.png"), newOutput("thumbnail-transparent-300x300.png")
	err := img.SetBackgroundColor("#00000000").ThumbnailImage(300, 300).WriteImage(output)
	a.NoError(err)

	img, output = newImage(a, "source.png"), newOutput("thumbnail-blue-300x300.png")
	err = img.SetBackgroundColor("blue").ThumbnailImage(300, 300).WriteImage(output)
	a.NoError(err)

	img, err = NewImage(output)
	a.NoError(err)

	a.Equal(300, img.GetWidth())
	a.Equal(300, img.GetHeight())
}

func TestSetQualityPNG(test *testing.T) {
	a := assert.New(test)
	img, output := newImage(a, "source.png"), newOutput("quality-png-10.png")
	err := img.SetQuality(10).CropThumbnailImage(300, 300).WriteImage(output)
	a.NoError(err)

	img, output = newImage(a, "source.png"), newOutput("quality-png-30.png")
	err = img.SetQuality(30).CropThumbnailImage(300, 300).WriteImage(output)
	a.NoError(err)

	img, output = newImage(a, "source.png"), newOutput("quality-png-50.png")
	err = img.SetQuality(50).CropThumbnailImage(300, 300).WriteImage(output)
	a.NoError(err)

	img, output = newImage(a, "source.png"), newOutput("quality-png-80.png")
	err = img.SetQuality(80).CropThumbnailImage(300, 300).WriteImage(output)
	a.NoError(err)

	img, output = newImage(a, "source.png"), newOutput("quality-png-100.png")
	err = img.SetQuality(100).CropThumbnailImage(300, 300).WriteImage(output)
	a.NoError(err)

	img, err = NewImage(output)
	a.NoError(err)

	a.Equal(300, img.GetWidth())
	a.Equal(300, img.GetHeight())
}

func TestSetQualityJPG(test *testing.T) {
	a := assert.New(test)
	img, output := newImage(a, "source.png"), newOutput("quality-jpg-10.jpg")
	err := img.SetQuality(10).CropThumbnailImage(300, 300).WriteImage(output)
	a.NoError(err)

	img, output = newImage(a, "source.png"), newOutput("quality-jpg-30.jpg")
	err = img.SetQuality(30).CropThumbnailImage(300, 300).WriteImage(output)
	a.NoError(err)

	img, output = newImage(a, "source.png"), newOutput("quality-jpg-50.jpg")
	err = img.SetQuality(50).CropThumbnailImage(300, 300).WriteImage(output)
	a.NoError(err)

	img, output = newImage(a, "source.png"), newOutput("quality-jpg-80.jpg")
	err = img.SetQuality(80).CropThumbnailImage(300, 300).WriteImage(output)
	a.NoError(err)

	img, output = newImage(a, "source.png"), newOutput("quality-jpg-100.jpg")
	err = img.SetQuality(100).CropThumbnailImage(300, 300).WriteImage(output)
	a.NoError(err)

	img, err = NewImage(output)
	a.NoError(err)

	a.Equal(300, img.GetWidth())
	a.Equal(300, img.GetHeight())
}

func TestPreserveEXIF(test *testing.T) {
	a := assert.New(test)
	img, output := newImage(a, "source.jpg"), newOutput("exif.jpg")

	err := img.PreserveEXIF().CropThumbnailImage(300, 300).WriteImage(output)
	a.NoError(err)

	img, err = NewImage(output)
	a.NoError(err)

	a.Equal(300, img.GetWidth())
	a.Equal(300, img.GetHeight())
}

func TestSetLoop(test *testing.T) {
	a := assert.New(test)
	img, output := newImage(a, "source.gif"), newOutput("loop-1.gif")

	err := img.SetLoop(1).WriteImage(output)
	a.NoError(err)

	img, err = NewImage(output)
	a.NoError(err)

	a.Equal(96, img.GetWidth())
	a.Equal(96, img.GetHeight())
}

func TestSetImageFramerate(test *testing.T) {
	a := assert.New(test)
	img, output := newImage(a, "source.gif"), newOutput("fps-1.gif")

	err := img.SetImageFramerate(1).WriteImage(output)
	a.NoError(err)

	img, err = NewImage(output)
	a.NoError(err)

	a.Equal(96, img.GetWidth())
	a.Equal(96, img.GetHeight())
}

func TestGetFrames(test *testing.T) {
	a := assert.New(test)
	img := newImage(a, "source.gif")
	a.Equal(60, img.GetFrames())

	img = newImage(a, "source.jpg")
	a.Equal(0, img.GetFrames())
}

func TestDropFrames(test *testing.T) {
	a := assert.New(test)
	img, output := newImage(a, "source.gif"), newOutput("drop-frame.gif")

	err := img.DropFrames().WriteImage(output)
	a.NoError(err)

	img, err = NewImage(output)
	a.NoError(err)

	a.Equal(1, img.GetFrames())
}

func TestPNGConvertion(test *testing.T) {
	a := assert.New(test)
	img, output := newImage(a, "source.png"), newOutput("png-to-jpg.jpg")
	err := img.SetImageFormat(ImageFormatJPEG).WriteImage(output)
	a.NoError(err)

	img, output = newImage(a, "source.png"), newOutput("png-to-gif.gif")
	err = img.SetImageFormat(ImageFormatGIF).WriteImage(output)
	a.NoError(err)

	img, output = newImage(a, "source.png"), newOutput("png-to-bmp.bmp")
	err = img.SetImageFormat(ImageFormatBMP).WriteImage(output)
	a.NoError(err)

	img, output = newImage(a, "source.png"), newOutput("png-to-avif.avif")
	err = img.SetImageFormat(ImageFormatAVIF).WriteImage(output)
	a.NoError(err)

	img, output = newImage(a, "source.png"), newOutput("png-to-webp.webp")
	err = img.SetImageFormat(ImageFormatWEBP).WriteImage(output)
	a.NoError(err)

	img, output = newImage(a, "source.png"), newOutput("png-to-apng.apng")
	err = img.SetImageFormat(ImageFormatAPNG).WriteImage(output)
	a.NoError(err)

	img, err = NewImage(output)
	a.NoError(err)

	a.Equal(300, img.GetWidth())
	a.Equal(300, img.GetHeight())
}

func TestGIFConvertion(test *testing.T) {
	a := assert.New(test)
	img, output := newImage(a, "source.gif"), newOutput("gif-to-jpg.jpg")
	err := img.SetImageFormat(ImageFormatJPEG).WriteImage(output)
	a.NoError(err)

	img, output = newImage(a, "source.gif"), newOutput("gif-to-png.png")
	err = img.SetImageFormat(ImageFormatPNG).WriteImage(output)
	a.NoError(err)

	img, output = newImage(a, "source.gif"), newOutput("gif-to-bmp.bmp")
	err = img.SetImageFormat(ImageFormatBMP).WriteImage(output)
	a.NoError(err)

	img, output = newImage(a, "source.gif"), newOutput("gif-to-avif.avif")
	err = img.SetImageFormat(ImageFormatAVIF).WriteImage(output)
	a.NoError(err)

	img, output = newImage(a, "source.gif"), newOutput("gif-to-webp.webp")
	err = img.SetImageFormat(ImageFormatWEBP).WriteImage(output)
	a.NoError(err)

	img, output = newImage(a, "source.gif"), newOutput("gif-to-apng.apng")
	err = img.SetImageFormat(ImageFormatAPNG).WriteImage(output)
	a.NoError(err)

	img, err = NewImage(output)
	a.NoError(err)

	a.Equal(300, img.GetWidth())
	a.Equal(300, img.GetHeight())
}
