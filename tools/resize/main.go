package main

import (
	"flag"
	"fmt"
	"image"
	_ "image/gif"
	"image/jpeg"
	"image/png"
	"os"
	"path/filepath"
	"time"

	"thanhfphan.com/pts/seamcarver"
)

var (
	path = flag.String("input", "mountain.jpg", "the path to the image")
	row  = flag.Int("row", 0, "the number row to remove")
	col  = flag.Int("col", 10, "the number column to remove")
)

func main() {
	flag.Parse()

	picture, err := openImage(*path)
	if err != nil {
		panic(err)
	}

	carver := seamcarver.New(picture)
	fmt.Printf("The image '%s' has %d rows and %d columns\n", *path, carver.Height(), carver.Width())

	start := time.Now()
	for i := 0; i < *col; i++ {
		seam := carver.FindVerticalSeam()
		carver.RemoveVerticalSeam(seam)
		fmt.Printf("Removed %d colums\n", i+1)
	}

	for i := 0; i < *row; i++ {
		seam := carver.FindHorizontalSeam()
		carver.RemoveHorizontalSeam(seam)
		fmt.Printf("Removed %d rows \n", i+1)
	}

	output := "out.jpg"
	file, err := os.Create(output)
	if err != nil {
		panic(err)
	}
	defer file.Close()

	options := jpeg.Options{Quality: 80}
	err = jpeg.Encode(file, carver.Picture(), &options)
	if err != nil {
		panic(err)
	}

	fmt.Printf("Output image '%s' has %d rows and %d columns\n", output, carver.Height(), carver.Width())
	duration := time.Since(start)
	fmt.Printf("Resize took %v seconds\n", duration.Seconds())
}

func openImage(path string) (image.Image, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	ext := filepath.Ext(path)
	switch ext {
	case ".png":
		img, err := png.Decode(file)
		if err != nil {
			return nil, fmt.Errorf("png.Decode '%s' error: %w", path, err)
		}

		return img, nil
	case ".jpg", ".jpeg":
		img, err := jpeg.Decode(file)
		if err != nil {
			return nil, fmt.Errorf("jpeg.Decode '%s' error: %w", path, err)
		}

		return img, nil
	default:
		// Try to decode any format supported by image.Decode
		img, _, err := image.Decode(file)
		if err != nil {
			return nil, fmt.Errorf("try to decode image '%s' error: %w", path, err)
		}

		return img, nil
	}

}
