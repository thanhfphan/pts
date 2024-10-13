package seamcarver

import (
	"image"
	"image/color"
	"math"
)

type SeamCarver struct {
	original image.Image

	pixels [][]*Pixel
}

func New(img image.Image) *SeamCarver {
	sc := &SeamCarver{
		original: img,
		pixels:   imageToArrayPixel(img),
	}

	sc.RecalculateEnergy()

	return sc
}

func imageToArrayPixel(img image.Image) [][]*Pixel {
	w, h := img.Bounds().Dx(), img.Bounds().Dy()
	arr := make([][]*Pixel, h)
	for y := 0; y < h; y++ {
		arr[y] = make([]*Pixel, w)
		for x := 0; x < w; x++ {
			arr[y][x] = &Pixel{
				C: img.At(x, y),
			}
		}
	}
	return arr
}

func arrayPixelToImage(arr [][]*Pixel) image.Image {
	h := len(arr)
	w := len(arr[0])

	img := image.NewRGBA(image.Rect(0, 0, w, h))

	for y := 0; y < h; y++ {
		for x := 0; x < w; x++ {
			img.Set(x, y, arr[y][x].C)
		}
	}

	return img
}

func (sc *SeamCarver) RecalculateEnergy() {
	h := len(sc.pixels)
	w := len(sc.pixels[0])

	for y := 0; y < h; y++ {
		for x := 0; x < w; x++ {
			// border
			if y == 0 || y == h-1 || x == 0 || x == w-1 {
				sc.pixels[y][x].E = 1000
				continue
			}

			rplus, gplus, bplus, _ := sc.pixels[y][x+1].C.RGBA()
			rminus, gminus, bminus, _ := sc.pixels[y][x-1].C.RGBA()
			rplus, gplus, bplus = rplus/257, gplus/257, bplus/257
			rminus, gminus, bminus = rminus/257, gminus/257, bminus/257

			deltaX := (rplus-rminus)*(rplus-rminus) + (gplus-gminus)*(gplus-gminus) + (bplus-bminus)*(bplus-bminus)

			// note result variable
			rplus, gplus, bplus, _ = sc.pixels[y+1][x].C.RGBA()
			rplus, gplus, bplus = rplus/257, gplus/257, bplus/257

			rminus, gminus, bminus, _ = sc.pixels[y-1][x].C.RGBA()
			rminus, gminus, bminus = rminus/257, gminus/257, bminus/257

			deltaY := (rplus-rminus)*(rplus-rminus) + (gplus-gminus)*(gplus-gminus) + (bplus-bminus)*(bplus-bminus)

			sc.pixels[y][x].E = math.Sqrt(float64(deltaX) + float64(deltaY))
		}
	}
}

// Picture represent the picture(current)
func (sc *SeamCarver) Picture() image.Image {
	return arrayPixelToImage(sc.pixels)
}

// Width of current picture
func (sc *SeamCarver) Width() int {
	return len(sc.pixels[0])
}

// Height of current picture
func (sc *SeamCarver) Height() int {
	return len(sc.pixels)
}

// Energy of pixel at column x and row y
func (sc *SeamCarver) Energy(x, y int) float64 {
	return sc.pixels[y][x].E
}

func (sc *SeamCarver) Color(x, y int) color.Color {
	return sc.pixels[y][x].C
}

// FindHorizontalSeam return sequence of indices for horizontal seam
func (sc *SeamCarver) FindHorizontalSeam() []int {
	panic("unimplement")
}

// FindVerticalSeam return sequence of indices for vertical seam
func (sc *SeamCarver) FindVerticalSeam() []int {
	panic("unimplement")
}

// RemoveHorizontalSeam from current picture
func (sc *SeamCarver) RemoveHorizontalSeam(seam []int) {
	panic("unimplement")
}

// RemoveVerticalSeam from current picture
func (sc *SeamCarver) RemoveVerticalSeam(seam []int) {
	panic("unimplement")
}
