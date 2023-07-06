package ffimage

import (
	"bytes"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strconv"

	ffmpeg "github.com/u2takey/ffmpeg-go"
)

// GetFrames returns the frame count of the image, returns 0 if it's not an animated image.
func (i *Image) GetFrames() int {
	frames, err := strconv.Atoi(i.Stream.NbFrames)
	if err != nil {
		frames = 0
	}
	return frames
}

// GetWidth returns the width of the image.
func (i *Image) GetWidth() int {
	return i.Stream.Width
}

// GetHeight returns the height of the image.
func (i *Image) GetHeight() int {
	return i.Stream.Height
}

// GetAspectRatio returns the aspect ratio of the image (w / h).
func (i *Image) GetAspectRatio() float64 {
	return float64(i.Width) / float64(i.Height)
}

// ResizeImage scales the size of an image to the given dimensions. The other parameter will be calculated if 0 is passed as either param. If ResizeType parameter is used both width and height must be given. The image will be stretched to exact size if ResizeType was left unspecified.
//
// - ResizeTypeUpscale: Given dimensions 400x400 an image of dimensions 300x225 would be scaled up to size 533x400.
//
// - ResizeTypeDownscale: Given dimensions 400x400 an image of dimensions 300x225 would be scaled up to size 400x300.
func (i *Image) ResizeImage(w, h int, typ ...ResizeType) *Image {
	ratio := float64(i.Width) / float64(i.Height)

	if w == 0 && h == 0 {
		return i
	}
	if len(typ) == 1 && typ[0] != ResizeTypeNone {
		// if w == 0 || h == 0 {
		// 	return i
		// }
		if w == 0 && h != 0 {
			w = h
		}
		if h == 0 && w != 0 {
			h = w
		}
		w, h = i.calcBestfit(i.Width, i.Height, w, h, typ[0])

	} else {
		if w == 0 {
			w = int(float64(h) * ratio)
		}
		if h == 0 {
			h = int(float64(w) / ratio)
		}
	}
	i.setWidthHeight(w, h)
	i.addFilter("scale", ffmpeg.Args{fmt.Sprintf("%d:%d", w, h)})
	return i
}

// ExtentImage comfortability method for setting image size. The method sets the image size and allows setting x,y coordinates where the new area begins.
//
// If "pos" is specified, "x" and "y" should be kept as 0.
func (i *Image) ExtentImage(w, h, x, y int, pos ...PositionType) *Image {
	if len(pos) == 1 && pos[0] != PositionTypeNone {
		x, y = i.calcPosition(i.Width, i.Height, w, h, pos[0])
	}
	i.setWidthHeight(w, h)
	i.addFilter("pad", ffmpeg.Args{fmt.Sprintf("%d:%d:%d:%d:%s", w, h, x, y, i.Output.BackgroundColor)})
	return i
}

// CropImage extracts a region of the image.
//
// If "pos" is specified, "x" and "y" should be kept as 0.
func (i *Image) CropImage(w, h, x, y int, pos ...PositionType) *Image {
	if len(pos) == 1 && pos[0] != PositionTypeNone {
		x, y = i.calcPosition(i.Width, i.Height, w, h, pos[0])
	}
	i.setWidthHeight(w, h)
	i.addFilter("crop", ffmpeg.Args{fmt.Sprintf("%d:%d:%d:%d", w, h, x, y)})
	return i
}

// CropThumbnailImage creates a fixed size thumbnail by first scaling the image up or down and cropping a specified area from the center.
func (i *Image) CropThumbnailImage(w, h int) *Image {
	i.ResizeImage(w, h, ResizeTypeUpscale).CropImage(w, h, 0, 0, PositionTypeCenter)
	i.setWidthHeight(w, h)
	return i
}

// ThumbnailImage creates a fixed size thumbnail and centered the image, the extented area will be filled with background color (black as default, can be set with SetBackgroundColor).
func (i *Image) ThumbnailImage(w, h int) *Image {
	imgW, imgH := i.calcBestpad(i.Width, i.Height, w, h)
	x, y := i.calcPosition(imgW, imgH, w, h, PositionTypeCenter)

	i.ResizeImage(imgW, imgH, ResizeTypeDownscale).ExtentImage(w, h, x, y)
	i.setWidthHeight(w, h)
	return i
}

// RotateImage rotates an image the specified number of degrees. Empty triangles left over from rotating the image are filled with the background color (black as default, can be set with SetBackgroundColor).
func (i *Image) RotateImage(degree int) *Image {
	i.addFilter("rotate", ffmpeg.Args{fmt.Sprintf("a=%d*PI/180:fillcolor=%s", degree, i.Output.BackgroundColor)})
	return i
}

// FlipImage creates a vertical mirror image.
func (i *Image) FlipImage() *Image {
	i.addFilter("vflip", ffmpeg.Args{})
	return i
}

// FlopImage creates a horizontal mirror image.
func (i *Image) FlopImage() *Image {
	i.addFilter("hflip", ffmpeg.Args{})
	return i
}

// SetBackgroundColor sets the object's default background color. Colors can be a name, hex, e.g. "black", "#FFFFFF", "0xFFFFFF" or rgba "#00000000".
func (i *Image) SetBackgroundColor(color string) *Image {
	i.Output.BackgroundColor = color
	return i
}

// SetLoop sets the repeat setting for animated image (e.g. GIF, WebP).
// - "-1" = no loop
// - "0" = infinite
// - "1" = loop once (play 2 times)
// - "2" = loop twice (play 3 times)
// - etc
func (i *Image) SetLoop(count int) *Image {
	i.Output.Loop = count
	return i
}

// SetQuality sets the quality for the image from 1 (low quality) to 100 (high quality). Native ffmpeg only works for: AVIF, JPEG, JPEGXL, WEBP.
//
// NOTE: for output format as PNG, the pngquant is required to be installed. The function does nothing for PNG if "pngquant" command was not found.
//
// NOTE: for output format as GIF, the gifsicle is required to be installed. The function does nothing for GIF if "gifsicle" command was not found.
func (i *Image) SetQuality(quality int) *Image {
	i.Output.Quality = quality
	return i

}

// SetImageFramerate changes the fps of the image, remains unchanged if the target fps is higher than current framerates. Filesize will be reduced if lower framerate was setted.
func (i *Image) SetImageFramerate(fps int) *Image {
	// TODO: https://github.com/eugeneware/ffprobe/issues/7
	// r_frame_rate=30/1
	// avg_frame_rate=438750/14777
	// i.addFilter("fps", ffmpeg.Args{fmt.Sprintf("%d", fps)})
	i.addArg(ffmpeg.KwArgs{"r": fmt.Sprintf("%d", fps)})
	return i
}

// DropFrames makes any animated images to static image by specifing `-vframes 1`.
func (i *Image) DropFrames() *Image {
	i.addArg(ffmpeg.KwArgs{"vframes": 1})
	return i
}

// SetImageFormat sets the output format and automatically decides the codec. Format will be detect automatically from the output filename if it wasn't been setted.
func (i *Image) SetImageFormat(format ImageFormat) *Image {
	switch format {
	case ImageFormatJPEG, ImageFormatPNG, ImageFormatBMP:
		i.DropFrames()
	}
	i.Output.Format = format
	return i
}

func (i *Image) suffixToFormat(ext string) ImageFormat {
	switch ext {
	case ".png":
		return ImageFormatPNG
	case ".apng":
		return ImageFormatAPNG
	case ".jpg", ".jpeg":
		return ImageFormatJPEG
	case ".gif":
		return ImageFormatGIF
	case ".webp":
		return ImageFormatWEBP
	case ".avif":
		return ImageFormatAVIF
	case ".bmp":
		return ImageFormatBMP
	case ".jxl":
		return ImageFormatJPEGXL
	}
	return ImageFormatUnknown
}

// WriteImage writes an image to the specified filename. If the filename parameter is empty string, the image is written to the source file.
func (i *Image) WriteImage(path string) error {
	i.Output.Path = path
	// Write to the input file if output path remains empty.
	if i.Output.Path == "" {
		i.Output.Path = i.Path
	}
	//
	isSameInputOutput := path == "" || path == i.Path
	var tmpFilename string

	if i.Output.Format == ImageFormatUnknown {
		i.Output.Format = i.suffixToFormat(filepath.Ext(path))
	}

	if i.Output.Format == ImageFormatUnknown {
		return errors.New("no image format")
	}

	i.buildQuality()
	i.buildLoop()
	i.buildBeforeEXIF()

	// Store the output to temp file if the output is the same as input,
	// because ffmpeg doesn't support the output to input.
	if isSameInputOutput {
		tmpFile, err := os.CreateTemp("", "*"+filepath.Ext(i.Path))
		if err != nil {
			return err
		}
		tmpFile.Close()
		tmpFilename = tmpFile.Name()
	}

	input := ffmpeg.Input(i.Path)

	for _, v := range i.Output.Filters {
		input = input.Filter(v.k, v.args)
	}

	// `palettegen` and `paletteuse` to keep the transparency for GIF.
	if i.Output.Format == ImageFormatGIF {
		split := input.Split()
		split1, split2 := split.Get("0"), split.Get("1")

		input = ffmpeg.Filter([]*ffmpeg.Stream{
			split1, split2.Filter("palettegen", ffmpeg.Args{})}, "paletteuse", ffmpeg.Args{})
	}

	buf := bytes.NewBuffer(nil)
	err := input.Output(path, i.Output.Args...).OverWriteOutput().Silent(i.Silent).WithErrorOutput(buf).Run()
	if err != nil {
		return errors.New(buf.String())
	}

	if isSameInputOutput {
		if err := os.Rename(tmpFilename, path); err != nil {
			return err
		}
	}

	i.buildAfterQuality()
	i.buildAfterEXIF()

	return nil
}

// PreserveEXIF preserves the metadata from the image, the metadata might contains GPS location or sensitive data, be caution.
//
// NOTE: Requires exiftool to be installed. The function does nothing if "exiftool" command was not found.
func (i *Image) PreserveEXIF() *Image {
	i.Output.IsPreserved = true
	return i
}
