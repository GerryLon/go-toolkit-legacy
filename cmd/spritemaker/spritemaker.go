package main

import (
	"flag"
	"fmt"
	"github.com/GerryLon/go-toolkit/common"
	"image"
	"image/draw"
	"image/jpeg"
	"io/ioutil"
	"os"
	"path/filepath"
)

const (
	defaultSpriteName = "sprite.png"
	defaultPadding    = 2
	defaultQuality    = 75
)

// make css sprite image
// example: ./spritemaker.exe -s E:/test/sprite/7185 -o E:/test/sprite/sprite.png
// spritemaker OPTIONS
// OPTIONS:
// -s sourceDir
// -o /path/to/sprite.png
// -p padding(default 2)
// -q quality(default 75, only for jpeg format, valid value:[1, 100])
// TODO: -l layout
func main() {
	s := flag.String("s", "", "source dir")
	o := flag.String("o", "", "output image name(with path, default is ./sprite.png)")
	p := flag.Int("p", 2, "padding")
	q := flag.Int("q", 75, "quality")
	flag.Parse()

	// *s = "E:/test/sprite/7185"
	// *o = "E:/test/sprite/sprite.png"
	images, err := getImagesFromDir(*s)
	if err != nil {
		fmt.Println(err)
		return
	}

	err = writeSpriteImage(*o, images, *p, *q)
	if err != nil {
		fmt.Println(err)
		return
	}
}

// 从给定的源文件获取图片(只获取一级,不递归)
func getImagesFromDir(src string) ([]string, error) {
	var images []string

	if src == "" {
		return images, fmt.Errorf("source src is empty")
	}

	fileInfo, err := os.Stat(src)
	if err != nil {
		if os.IsNotExist(err) {
			return images, fmt.Errorf("source dir: %s is not exist", src)
		} else {
			return images, fmt.Errorf("get images from source dir: %s, err: %v", src, err)
		}
	}

	if !fileInfo.IsDir() {
		return images, fmt.Errorf("source %s is not a directory", src)
	}

	// // this is recursively
	// filepath.Walk(src, func(path string, info os.FileInfo, err error) error {
	// 	if err != nil {
	// 		return err
	// 	}
	//
	// 	if info.IsDir() {
	// 		return nil
	// 	}
	//
	// 	fmt.Println(path)
	//
	// 	return nil
	// })
	files, err := ioutil.ReadDir(src)
	if err != nil {
		return images, fmt.Errorf("get images from source dir: %s, err: %v", src, err)
	}

	imageSuffixes := []string{".jpeg", ".jpg", ".png", "gif"}
	for _, file := range files {
		if file.IsDir() {
			continue
		}

		if common.HasSuffix(file.Name(), imageSuffixes) {

			// multi images with absolute path
			images = append(images, filepath.Join(src, file.Name()))
		}
	}

	if len(images) == 0 {
		return images, fmt.Errorf("no image found in source dir:%s", src)
	}

	return images, nil
}

// 将指定的一些小图拼成大图
func writeSpriteImage(dst string, images []string, padding int, quality int) error {
	out := dst

	// param check
	if out == "" {
		out = filepath.Join(common.MustGetCWD(), defaultSpriteName)
	}
	if padding < 0 {
		padding = defaultPadding
	}
	if quality < 1 || quality > 100 {
		quality = defaultQuality
	}

	// create sprite image
	sprite, err := os.Create(out)
	if err != nil {
		return fmt.Errorf("write image to: %s, err: %v", dst, err)
	}
	defer sprite.Close()

	// 计算最终拼出的图的大小
	dstImgSize := calcDstImgSize(images, padding)
	spriteImage := image.NewNRGBA(image.Rect(0, 0, dstImgSize.Width, dstImgSize.Height))
	tmpRect := image.Rectangle{Min: image.Point{X: 0, Y: 0}, Max: image.Point{X: 0, Y: 0}} // 临时变量

	for index, img := range images {
		decode := common.ImageDecode(filepath.Ext(img))
		file, err := os.Open(img)
		if err != nil {
			return err
		}

		i, err := decode(file)
		file.Close()
		if err != nil {
			return err
		}
		tmpRect.Max.X += i.Bounds().Max.X
		tmpRect.Min.X = tmpRect.Max.X - i.Bounds().Max.X
		tmpRect.Min.Y = 0
		if i.Bounds().Max.Y > tmpRect.Max.Y {
			tmpRect.Max.Y = i.Bounds().Max.Y
		}

		// 将padding算上
		if index > 0 {
			tmpRect.Min.X += padding
			tmpRect.Max.X += padding
		}

		// fmt.Printf("%s bounds: %v\n", img, tmpRect)
		draw.Draw(spriteImage, tmpRect, i, image.Point{0, 0}, draw.Over)
	}

	encodeFn := common.ImageEncode(filepath.Ext(out))

	// jpeg with quality
	var encodeOptions interface{}
	if common.IsJPEG(out) {
		encodeOptions = &jpeg.Options{Quality: quality}
	}
	return encodeFn(sprite, spriteImage, encodeOptions)
}

func calcDstImgSize(images []string, padding int) common.Size {
	size := common.Size{}

	// 当前策略, 一字排开, 高度取最高的图片的高
	for _, img := range images {
		tmpSize := common.MustGetImageSize(img)
		size.Width += tmpSize.Width

		if tmpSize.Height > size.Height {
			size.Height = tmpSize.Height
		}
	}

	// 加上 n - 1个padding
	// + len(images)是为了避免图片尺寸是小数的问题, 最后拼出来的图可能不全
	size.Width += (len(images)-1)*padding + len(images)

	return size
}
