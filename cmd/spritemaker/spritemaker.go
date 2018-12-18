package main

import (
	"flag"
	"fmt"
	"github.com/GerryLon/go-toolkit/common"
	"io/ioutil"
	"os"
	"path/filepath"
)

// make css sprite image
// spritemaker -s sourceDir -o /path/to/sprite.png
// TODO: -l layout
// TODO: -q quality
func main() {
	s := flag.String("s", "", "source dir")
	o := flag.String("o", "", "output image name(with path)")
	flag.Parse()

	*s = "E:/desktop/wallpaper"

	images, err := getImagesFromDir(*s)
	if err != nil {
		fmt.Println(err)
		return
	}

	writeSpriteImage(*o, images)
}

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
			images = append(images, filepath.Join(src, file.Name()))
		}
	}

	if len(images) == 0 {
		return images, fmt.Errorf("no image found in source dir:%s", src)
	}

	return images, nil
}

func writeSpriteImage(dst string, images []string) error {

	out := dst

	if out == "" {
		out = common.MustGetCWD()
	}

	_, err := os.Create(dst)
	if err != nil {
		if os.IsNotExist(err) {
			return fmt.Errorf("dst dir: %s is not exist", dst)
		} else {
			return fmt.Errorf("write image to dst dir: %s, err: %v", dst, err)
		}
	}

	return nil
}
