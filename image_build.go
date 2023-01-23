package ffimage

import (
	"encoding/json"
	"fmt"
	"os"
	"os/exec"

	ffmpeg "github.com/u2takey/ffmpeg-go"
)

// buildQuality
func (i *Image) buildQuality() *Image {
	if i.Output.Quality == 0 {
		return i
	}
	switch i.Output.Format {
	case ImageFormatAVIF:
		q := qualityFactor(0, 63, i.Output.Quality, true)
		i.addArg(ffmpeg.KwArgs{"crf": q})

	case ImageFormatJPEG:
		q := qualityFactor(2, 31, i.Output.Quality, true)
		i.addArg(ffmpeg.KwArgs{"qscale:v": q})

	case ImageFormatJPEGXL:
		q := qualityFactor(0, 100, i.Output.Quality, false)
		i.addArg(ffmpeg.KwArgs{"qscale:v": q})

	case ImageFormatWEBP:
		q := qualityFactor(0, 100, i.Output.Quality, false)
		i.addArg(ffmpeg.KwArgs{"quality": q})
	}
	return i
}

// buildLoop
func (i *Image) buildLoop() *Image {
	if i.Output.Format == ImageFormatAPNG {
		i.addArg(ffmpeg.KwArgs{"plays": i.Output.Loop})
	} else {
		i.addArg(ffmpeg.KwArgs{"loop": i.Output.Loop})
	}
	return i
}

// buildAfterQuality
func (i *Image) buildAfterQuality() *Image {
	if i.Output.Quality == 0 {
		return i
	}
	switch i.Output.Format {
	case ImageFormatPNG:
		q := qualityFactor(0, 100, i.Output.Quality, false)
		exec.Command("pngquant", "--quality", fmt.Sprintf("0-%d", q), "-f", i.Output.Path, "-o", i.Output.Path).Run()

	case ImageFormatGIF:
		q := qualityFactor(0, 100, i.Output.Quality, false)
		exec.Command("gifsicle", "-O3", fmt.Sprintf("--lossy=%d", q), i.Output.Path, "-o", i.Output.Path).Run()
	}
	return i
}

// buildBeforeEXIF
func (i *Image) buildBeforeEXIF() *Image {
	if !i.Output.IsPreserved {
		return i
	}
	tmpFile, err := os.CreateTemp("", "")
	if err != nil {
		return i
	}
	b, err := exec.Command("exiftool", "-json", i.Path).Output()
	if err != nil {
		return i
	}
	j := make([]map[string]interface{}, 0)
	if err := json.Unmarshal(b, &j); err != nil {
		return i
	}
	if len(j) == 0 {
		return i
	}
	// Set SourceFile as * so the data extracted from exiftool can import to any file.
	j[0]["SourceFile"] = "*"
	//
	b, err = json.Marshal(j)
	if err != nil {
		return i
	}
	if _, err := tmpFile.Write(b); err != nil {
		return i
	}
	if err := tmpFile.Close(); err != nil {
		return i
	}
	i.Output.EXIF = tmpFile.Name()
	return i
}

// buildAfterEXIF
func (i *Image) buildAfterEXIF() *Image {
	if !i.Output.IsPreserved {
		return i
	}
	if err := exec.Command("exiftool", "-overwrite_original", fmt.Sprintf("-json=%s", i.Output.EXIF), i.Output.Path).Run(); err != nil {
		return i
	}
	return i
}
