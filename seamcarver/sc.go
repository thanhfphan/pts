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
	return sc.picture.Bounds().Dx()
}

// Height of current picture
func (sc *SeamCarver) Height() int {
	return sc.picture.Bounds().Dy()
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
	newImg = removeVerticalSeam(newImg, seam)
	sc.picture = transposeImage(newImg) // transpose back
	sc.recalculateEnergy()
}

// RemoveVerticalSeam from current picture
func (sc *SeamCarver) RemoveVerticalSeam(seam []int) {
	sc.picture = removeVerticalSeam(sc.picture, seam)
	sc.recalculateEnergy()
}

func (sc *SeamCarver) InsertVerticalSearm(n int) {
	copyimg := copyImage(sc.picture)
	deleteSeam := [][]int{}
	for i := 0; i < n; i++ {
		seam := sc.FindVerticalSeam()
		deleteSeam = append(deleteSeam, seam)
		sc.picture = removeVerticalSeam(sc.picture, seam)
		sc.recalculateEnergy()
	}

	sc.picture = copyimg
	sc.recalculateEnergy()
	for _, seam := range deleteSeam {
		sc.picture = insertVerticalSeam(sc.picture, seam)
		sc.recalculateEnergy()
	}
}

func (sc *SeamCarver) InsertHorizontalSeam(n int) {
	sc.picture = transposeImage(sc.picture)
	sc.InsertVerticalSearm(n)
	sc.picture = transposeImage(sc.picture)
}

func insertVerticalSeam(img image.Image, seam []int) image.Image {
	width, height := img.Bounds().Dx(), img.Bounds().Dy()
	newImg := image.NewRGBA(image.Rect(0, 0, width+1, height))

	for y := 0; y < height; y++ {
		var c int
		for x := 0; x < width; x++ {
			newImg.Set(c, y, img.At(x, y))
			c++
			if seam[y] == x {
				// Duplicate the pixel at the seam
				left := img.At(max(x-1, 0), y)
				right := img.At(min(x+1, width-1), y)
				newImg.Set(c, y, averageColor(left, right))
				c++
			}
		}
	}

	return newImg
}

func (sc *SeamCarver) recalculateEnergy() {
	width, height := sc.picture.Bounds().Dx(), sc.picture.Bounds().Dy()
	sc.energy = make([][]float64, height)
	for i := range sc.energy {
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
		copy(cost[i], energy[i])
	}

	for row := 1; row < height; row++ {
		for col := 0; col < width; col++ {
			minCost := cost[row-1][col]
			if col > 0 {
				minCost = math.Min(minCost, cost[row-1][col-1])
			}
			if col < width-1 {
				minCost = math.Min(minCost, cost[row-1][col+1])
			}
			cost[row][col] = minCost + energy[row][col]
		}
	}

	path := make([]int, height)
	minCost := math.MaxFloat64
	for col := 0; col < width; col++ {
		if cost[height-1][col] < minCost {
			minCost = cost[height-1][col]
			path[height-1] = col
		}
	}

	for row := height - 2; row >= 0; row-- {
		lastCol := path[row+1]
		minCol := lastCol
		if lastCol > 0 && cost[row][lastCol-1] < cost[row][minCol] {
			minCol = lastCol - 1
		}
		if lastCol < width-1 && cost[row][lastCol+1] < cost[row][minCol] {
			minCol = lastCol + 1
		}
		path[row] = minCol
	}

	return path
}

// transpose return new matrix after transpose
func transpose(energy [][]float64) [][]float64 {
	rows, cols := len(energy), len(energy[0])
	transposed := make([][]float64, cols)
	for col := range transposed {
		transposed[col] = make([]float64, rows)
		for row := range energy {
			copy(transposed[col], energy[row])
		}
	}

	return transposed
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

func copyImage(img image.Image) image.Image {
	bounds := img.Bounds()
	newImg := image.NewRGBA(image.Rect(0, 0, bounds.Dx(), bounds.Dy()))
	for y := 0; y < bounds.Dy(); y++ {
		for x := 0; x < bounds.Dx(); x++ {
			newImg.Set(x, y, img.At(x, y))
		}
	}

	return newImg
}

func delta(c1, c2 color.Color) uint32 {
	r1, g1, b1, _ := c1.RGBA()
	r1, g1, b1 = r1/257, g1/257, b1/257

	r2, g2, b2, _ := c2.RGBA()
	r2, g2, b2 = r2/257, g2/257, b2/257

	delta := (r1-r2)*(r1-r2) + (g1-g2)*(g1-g2) + (b1-b2)*(b1-b2)

	return delta
}

func averageColor(left, right color.Color) color.Color {
	leftR, leftG, leftB, leftA := left.RGBA()
	rightR, rightG, rightB, rightA := right.RGBA()

	avgR := uint8((leftR + rightR) / 257 / 2)
	avgG := uint8((leftG + rightG) / 257 / 2)
	avgB := uint8((leftB + rightB) / 257 / 2)
	avgA := uint8((leftA + rightA) / 257 / 2)

	return color.RGBA{R: avgR, G: avgG, B: avgB, A: avgA}
}
