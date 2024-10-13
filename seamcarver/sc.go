package seamcarver

import (
	"image"
)

type SeamCarver struct {
	original image.Image

	currentImage [][]Pixel
}

func New(img image.Image) *SeamCarver {
	sc := &SeamCarver{
		original:     img,
		currentImage: imageToArrayPixel(img),
	}

	sc.RecalculateEnergy()

	return sc
}

func imageToArrayPixel(img image.Image) [][]Pixel {
	w, h := img.Bounds().Dx(), img.Bounds().Dy()
	arr := make([][]Pixel, h)
	for y := 0; y < h; y++ {
		arr[y] = make([]Pixel, w)
		for x := 0; x < w; x++ {
			arr[y][x] = Pixel{
				C: img.At(x, y),
			}
		}
	}
	return arr
}

func arrayPixelToImage(arr [][]Pixel) image.Image {
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
	// TODO: implement
}

// Picture represent the picture(current)
func (sc *SeamCarver) Picture() image.Image {
	return arrayPixelToImage(sc.currentImage)
}

// Width of current picture
func (sc *SeamCarver) Width() int {
	return len(sc.currentImage[0])
}

// Height of current picture
func (sc *SeamCarver) Height() int {
	return len(sc.currentImage)
}

// Energy of pixel at column x and row y
func (sc *SeamCarver) Energy(x, y int) float64 {
	return sc.currentImage[y][x].E
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
