package ffimage

import (
	"context"
	"fmt"
	"math"
	"os"

	"github.com/gabriel-vasile/mimetype"
	ffmpeg "github.com/u2takey/ffmpeg-go"
	"gopkg.in/vansante/go-ffprobe.v2"
)

// ImageFormat
type ImageFormat string

const (
	ImageFormatUnknown = ""
	ImageFormatJPEG    = "jpg"
	ImageFormatJPEGXL  = "jxl"
	ImageFormatWEBP    = "webp"
	ImageFormatPNG     = "png"
	ImageFormatAVIF    = "avif"
	ImageFormatAPNG    = "apng"
	ImageFormatBMP     = "bmp"
	ImageFormatGIF     = "gif"
)

// ResizeType
type ResizeType int

const (
	ResizeTypeNone ResizeType = iota
	ResizeTypeUpscale
	ResizeTypeDownscale
)

// PositionType
type PositionType string

const (
	PositionTypeNone        PositionType = ""
	PositionTypeTopLeft     PositionType = "top_left"
	PositionTypeTop         PositionType = "top"
	PositionTypeTopRight    PositionType = "top_right"
	PositionTypeLeft        PositionType = "left"
	PositionTypeCenter      PositionType = "center"
	PositionTypeRight       PositionType = "right"
	PositionTypeBottomLeft  PositionType = "bottom_left"
	PositionTypeBottom      PositionType = "bottom"
	PositionTypeBottomRight PositionType = "bottom_right"
)

// Name   Worse  Best  Default  Usage
// JPEG    31     2      17   -qscale:v
// WEBP    0     100     75   -quality
// JPXL    0     100     90   -qscale:v
// AVIF    63     0      50   -crf
// PNG     -      -      -
// BMP     -      -      -
// GIF     -      -      -

// Image
type Image struct {
	Stream *ffprobe.Stream
	Width  int
	Height int
	Path   string
	Output *Output
	Silent bool
	isTemp bool
}

type Output struct {
	Path            string
	Args            []ffmpeg.KwArgs
	Filters         []*filter
	Quality         int
	Format          ImageFormat
	Loop            int
	IsPreserved     bool
	EXIF            string
	Codec           string
	BackgroundColor string
}

// NewImage
func NewImage(path string) (*Image, error) {
	image := &Image{
		Path: path,
		Output: &Output{
			BackgroundColor: "black",
			Args:            make([]ffmpeg.KwArgs, 0),
			Filters:         make([]*filter, 0),
		},
		Silent: true,
	}
	if err := image.loadImageSize(); err != nil {
		return nil, fmt.Errorf("load image size: %w", err)
	}
	image.addArg(ffmpeg.KwArgs{"map_metadata": "-1"})
	image.addFilter("format", ffmpeg.Args{"rgba"})
	return image, nil
}

func NewImageFromBytes(data []byte) (*Image, error) {
	mtype := mimetype.Detect(data)
	tmpFile, err := os.CreateTemp("", "ffimage.*."+mtype.Extension())
	if err != nil {
		return nil, fmt.Errorf("create temp: %w", err)
	}
	if _, err := tmpFile.Write(data); err != nil {
		return nil, fmt.Errorf("write: %w", err)
	}
	if err := tmpFile.Close(); err != nil {
		return nil, fmt.Errorf("close: %w", err)
	}
	img, err := NewImage(tmpFile.Name())
	if err != nil {
		return nil, fmt.Errorf("new image: %w", err)
	}
	img.isTemp = true

	return img, nil
}

// loadImageSize
func (i *Image) loadImageSize() error {
	data, err := ffprobe.ProbeURL(context.TODO(), i.Path)
	if err != nil {
		return fmt.Errorf("probe url: %w", err)
	}
	if len(data.Streams) == 0 || data.Streams[0].Width == 0 || data.Streams[0].Height == 0 || data.Streams[0] == nil {
		return fmt.Errorf("no valid stream found")
	}
	i.Stream = data.Streams[0]
	i.setWidthHeight(i.Stream.Width, i.Stream.Height)
	return nil
}

// filter
type filter struct {
	k    string
	args ffmpeg.Args
}

// addFilter
func (i *Image) addFilter(k string, args ffmpeg.Args) {
	i.Output.Filters = append(i.Output.Filters, &filter{k, args})
}

// addArg
func (i *Image) addArg(args ffmpeg.KwArgs) {
	i.Output.Args = append(i.Output.Args, args)
}

// setWidthHeight
func (i *Image) setWidthHeight(w, h int) {
	i.Width = w
	i.Height = h
}

// qualityFactor
func qualityFactor(min, max int, quality int, isLowerBetter bool) int {
	if isLowerBetter {
		return max + int((float64(quality)/100)*float64(min-max))
	}
	return min + int((float64(quality)/100)*float64(max-min))
}

// calcPosition
func (i *Image) calcPosition(origW, origH, w, h int, pos PositionType) (x int, y int) {
	switch pos {
	case PositionTypeTopLeft:
		return
	case PositionTypeTop:
		x = (origW / 2) - (w / 2)
	case PositionTypeTopRight:
		x = origW - w
	case PositionTypeLeft:
		y = (origH / 2) - (h / 2)
	case PositionTypeCenter:
		x = (origW / 2) - (w / 2)
		y = (origH / 2) - (h / 2)
	case PositionTypeRight:
		x = origW - w
		y = (origH / 2) - (h / 2)
	case PositionTypeBottomLeft:
		y = origH - h
	case PositionTypeBottom:
		x = (origW / 2) - (w / 2)
		y = origH - h
	case PositionTypeBottomRight:
		x = origW - w
		y = origH - h
	}
	x = int(math.Abs(float64(x)))
	y = int(math.Abs(float64(y)))
	return
}

// calcBestfit
func (i *Image) calcBestfit(origW, origH, w, h int, typ ResizeType) (newW, newH int) {
	ratio := float64(origW) / float64(origH)

	if origW > origH {
		newH = h
		newW = int((float64(newH) * ratio))

		if (typ == ResizeTypeDownscale && newW > w) || (typ == ResizeTypeUpscale && w > newW) {
			newW = w
			newH = int((float64(newW) / ratio))
		}
	} else {
		newW = w
		newH = int((float64(newW) / ratio))

		if (typ == ResizeTypeDownscale && newH > h) || (typ == ResizeTypeUpscale && h > newH) {
			newH = h
			newW = int((float64(newH) * ratio))
		}
	}
	return
}

// calcBestpad
func (i *Image) calcBestpad(origW, origH, w, h int) (newW, newH int) {
	var ratio float64

	if float64(origW)/float64(w) > float64(origH)/float64(h) {
		ratio = float64(w) / float64(origW)
	} else {
		ratio = float64(h) / float64(origH)
	}

	newW = int(float64(origW) * ratio)
	newH = int(float64(origH) * ratio)
	return
}
