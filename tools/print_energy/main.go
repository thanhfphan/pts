package main

import (
	"flag"
	"fmt"
	"image"
	"image/jpeg"
	"image/png"
	"os"
	"path/filepath"

	"thanhfphan.com/pts/seamcarver"
)

var path = flag.String("image", "demo.png", "the path to the image")

func main() {
	flag.Parse()

	pic, err := openImage(*path)
	if err != nil {
		panic(err)
	}
	fmt.Printf("Image is %d pixels wide by %d pixels high.\n", pic.Bounds().Dx(), pic.Bounds().Dy())

	sc := seamcarver.New(pic)
	fmt.Println("Printing energy calculated for each pixel.")
	for row := 0; row < sc.Height(); row++ {
		for col := 0; col < sc.Width(); col++ {
			fmt.Printf("%9.2f ", sc.Energy(col, row))
			// fmt.Printf("%2.0v ", sc.Color(col, row))
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
