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

	"thanhfphan.com/pts/seamcarver"
)

var path = flag.String("image", "demo.png", "the path to the image")

func main() {
	flag.Parse()

	picture, err := openImage(*path)
	if err != nil {
		panic(err)
	}

	fmt.Printf("%s (%d-by-%d image)\n", *path, picture.Bounds().Dx(), picture.Bounds().Dy())
	fmt.Println()

	carver := seamcarver.New(picture)

	verticalSeam := carver.FindVerticalSeam()
	fmt.Println("Vertical seam:", verticalSeam)
	fmt.Println("Printing energy calculated for each pixel.")
	for row := 0; row < carver.Height(); row++ {
		for col := 0; col < carver.Width(); col++ {
			var c string
			if verticalSeam[row] == col {
				c = "*" // marked
			}
			fmt.Printf("%10.2f%s", carver.Energy(col, row), c)
		}
		fmt.Println()
	}

	fmt.Println()

	horizontalSeam := carver.FindHorizontalSeam()
	fmt.Println("Horizontal seam:", horizontalSeam)
	fmt.Println("Printing energy calculated for each pixel.")
	for row := 0; row < carver.Height(); row++ {
		for col := 0; col < carver.Width(); col++ {
			var c string
			if horizontalSeam[col] == row {
				c = "*" // marked
			}
			fmt.Printf("%10.2f%s", carver.Energy(col, row), c)
		}
		fmt.Println()
	}
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
