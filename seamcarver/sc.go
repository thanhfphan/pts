package seamcarver

import (
	"image"
	"image/color"
	"math"
)

type SeamCarver struct {
	picture image.Image

	energy [][]float64
}

func New(img image.Image) *SeamCarver {
	sc := &SeamCarver{
		picture: img,
	}

	sc.recalculateEnergy()

	return sc
}

// Picture represent the picture(current)
func (sc *SeamCarver) Picture() image.Image {
	return sc.picture
}

// Width of current picture
func (sc *SeamCarver) Width() int {
	w, _ := sc.picture.Bounds().Dx(), sc.picture.Bounds().Dy()
	return w
}

// Height of current picture
func (sc *SeamCarver) Height() int {
	_, h := sc.picture.Bounds().Dx(), sc.picture.Bounds().Dy()
	return h
}

// Energy of pixel at column x and row y
func (sc *SeamCarver) Energy(x, y int) float64 {
	return sc.energy[y][x]
}

func (sc *SeamCarver) Color(x, y int) color.Color {
	return sc.picture.At(x, y)
}

// FindHorizontalSeam return sequence of indices for horizontal seam
func (sc *SeamCarver) FindHorizontalSeam() []int {
	// - transpose
	tp := transpose(sc.energy)
	// - find vertical seam
	return retrieveSeamPath(tp)
}

// FindVerticalSeam return sequence of indices for vertical seam
func (sc *SeamCarver) FindVerticalSeam() []int {
	return retrieveSeamPath(sc.energy)
}

// RemoveHorizontalSeam from current picture
func (sc *SeamCarver) RemoveHorizontalSeam(seam []int) {
	newImg := transposeImage(sc.picture)
	removeVerticalSeam(newImg, seam)
	sc.picture = transposeImage(newImg) // transpose back
	sc.recalculateEnergy()
}

// RemoveVerticalSeam from current picture
func (sc *SeamCarver) RemoveVerticalSeam(seam []int) {
	sc.picture = removeVerticalSeam(sc.picture, seam)
	sc.recalculateEnergy()
}

func (sc *SeamCarver) recalculateEnergy() {
	width, height := sc.picture.Bounds().Dx(), sc.picture.Bounds().Dy()
	sc.energy = make([][]float64, height)
	for i, _ := range sc.energy {
		sc.energy[i] = make([]float64, width)
	}

	for y := 0; y < height; y++ {
		for x := 0; x < width; x++ {
			// border
			if y == 0 || y == height-1 || x == 0 || x == width-1 {
				sc.energy[y][x] = 1000
				continue
			}

			deltaX := delta(sc.picture.At(x+1, y), sc.picture.At(x-1, y))
			deltaY := delta(sc.picture.At(x, y+1), sc.picture.At(x, y-1))

			sc.energy[y][x] = math.Sqrt(float64(deltaX) + float64(deltaY))
		}
	}
}

func removeVerticalSeam(img image.Image, seam []int) image.Image {
	w, h := img.Bounds().Dx(), img.Bounds().Dy()
	newImg := image.NewRGBA(image.Rect(0, 0, w-1, h))

	for i := 0; i < h; i++ {
		var c int
		for j := 0; j < w; j++ {
			if i < len(seam) && seam[i] != j {
				newImg.Set(c, i, img.At(j, i))
				c++
			}
		}
	}

	return newImg
}

// retrieveSeamPath travel from top to bottom
func retrieveSeamPath(energy [][]float64) []int {
	height, width := len(energy), len(energy[0])
	cost := make([][]float64, height)
	for i := range cost {
		cost[i] = make([]float64, width)
	}

	for i := 0; i < width; i++ {
		cost[0][i] = energy[0][i]
	}

	for row := 1; row < height; row++ {
		for col := 0; col < width; col++ {
			cost[row][col] = cost[row-1][col]
			if col-1 > 0 {
				cost[row][col] = math.Min(cost[row][col], cost[row-1][col-1])
			}
			if col+1 < width {
				cost[row][col] = math.Min(cost[row][col], cost[row-1][col+1])
			}

			cost[row][col] = cost[row][col] + energy[row][col]
		}
	}

	path := make([]int, height)
	mincost := math.MaxFloat64
	for col := 0; col < width; col++ {
		if cost[height-1][col] < mincost {
			mincost = cost[height-1][col]
			path[height-1] = col
		}
	}

	for row := height - 2; row >= 0; row-- {
		lc := path[row+1] // last column
		minCol := lc
		if lc-1 > 0 && cost[row][lc-1] < cost[row][minCol] {
			minCol = lc - 1
		}
		if lc+1 < width && cost[row][lc+1] < cost[row][minCol] {
			minCol = lc + 1
		}

		path[row] = minCol
	}

	return path
}

// transpose return new matrix after transpose
func transpose(energy [][]float64) [][]float64 {
	rows, cols := len(energy), len(energy[0])
	newpixels := make([][]float64, cols)
	for i := range newpixels {
		newpixels[i] = make([]float64, rows)
	}
	for i := 0; i < rows; i++ {
		for j := 0; j < cols; j++ {
			newpixels[j][i] = energy[i][j]
		}
	}

	return newpixels
}

func transposeImage(img image.Image) image.Image {
	w, h := img.Bounds().Dx(), img.Bounds().Dy()
	newimg := image.NewRGBA(image.Rect(0, 0, h, w))
	for y := 0; y < h; y++ {
		for x := 0; x < w; x++ {
			newimg.Set(y, x, img.At(x, y))
		}
	}

	return newimg
}

func delta(c1, c2 color.Color) uint32 {
	r, g, b, _ := c1.RGBA()
	r, g, b = r/257, g/257, b/257

	r2, g2, b2, _ := c2.RGBA()
	r2, g2, b2 = r2/257, g2/257, b2/257

	d := (r-r2)*(r-r2) + (g-g2)*(g-g2) + (b-b2)*(b-b2)

	return d
}
