package common

import (
	"image"
	"image/gif"
	"image/jpeg"
	"image/png"
	"io"
	"os"
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
