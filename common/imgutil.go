package common

import (
	"image"
	"image/gif"
	"image/jpeg"
	"image/png"
	"io"
	"os"
	"path/filepath"
)

func GetImageSize(imagePath string) (Size, error) {
	file, err := os.Open(imagePath)
	if err != nil {
		return Size{}, err
	}
	defer file.Close()

	image, _, err := image.DecodeConfig(file)
	if err != nil {
		return Size{}, err
	}

	return Size{image.Width, image.Height}, nil
}

func MustGetImageSize(imagePath string) Size {
	size, err := GetImageSize(imagePath)
	if err != nil {
		panic(err)
	}
	return size
}

func ImageDecode(suffix string) func(r io.Reader) (image.Image, error) {
	tmp := suffix
	if tmp[0] == '.' {
		tmp = tmp[1:]
	}

	switch tmp {
	case "jpeg", "jpg":
		return jpeg.Decode
	case "png":
		return png.Decode
	case "gif":
		return gif.Decode

	default:
		panic("invalid suffix: " + suffix)
	}

}

func ImageEncode(suffix string) func(w io.Writer, m image.Image, options interface{}) error {
	tmp := suffix
	if tmp[0] == '.' {
		tmp = tmp[1:]
	}

	switch tmp {
	case "jpeg", "jpg":
		return func(w io.Writer, m image.Image, options interface{}) error {
			return jpeg.Encode(w, m, options.(*jpeg.Options))
		}

	// png ignore quality
	case "png":
		return func(w io.Writer, m image.Image, options interface{}) error {
			return png.Encode(w, m)
		}

	case "gif":
		return func(w io.Writer, m image.Image, options interface{}) error {
			return gif.Encode(w, m, options.(*gif.Options))
		}
	default:
		panic("invalid suffix: " + suffix)
	}
}

func IsJPEG(imagePath string) bool {
	ext := filepath.Ext(imagePath)
	if ext == ".jpg" || ext == ".jpeg" {
		return true
	}
	return false
}
