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

	sc.recalculateEnergy()

	return sc
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
	// - transpose
	tp := transpose(sc.pixels)
	// - find vertical seam
	return findShortestPath(tp)
}

// FindVerticalSeam return sequence of indices for vertical seam
func (sc *SeamCarver) FindVerticalSeam() []int {
	return findShortestPath(sc.pixels)
}

// RemoveHorizontalSeam from current picture
func (sc *SeamCarver) RemoveHorizontalSeam(seam []int) {
	newpixels := transpose(sc.pixels)
	newpixels = removeVerticalSeam(newpixels, seam)
	sc.pixels = transpose(newpixels)
	sc.recalculateEnergy()
}

// RemoveVerticalSeam from current picture
func (sc *SeamCarver) RemoveVerticalSeam(seam []int) {
	sc.pixels = removeVerticalSeam(sc.pixels, seam)
	sc.recalculateEnergy()
}

func (sc *SeamCarver) recalculateEnergy() {
	height := len(sc.pixels)
	if height == 0 {
		return
	}
	width := len(sc.pixels[0])

	for y := 0; y < height; y++ {
		for x := 0; x < width; x++ {
			// border
			if y == 0 || y == height-1 || x == 0 || x == width-1 {
				sc.pixels[y][x].E = 1000
				continue
			}

			deltaX := calcDeltaSquare(sc.pixels[y][x+1].C, sc.pixels[y][x-1].C)
			deltaY := calcDeltaSquare(sc.pixels[y+1][x].C, sc.pixels[y-1][x].C)

			sc.pixels[y][x].E = math.Sqrt(float64(deltaX) + float64(deltaY))
		}
	}
}

func removeVerticalSeam(input [][]*Pixel, seam []int) [][]*Pixel {
	h := len(input)
	if h == 0 {
		return nil
	}
	w := len(input[0])
	newpixels := make([][]*Pixel, h)
	for i := range newpixels {
		newpixels[i] = make([]*Pixel, w-1) // minus 1
	}

	for i := 0; i < h; i++ {
		var c int
		for j := 0; j < w; j++ {
			if i < len(seam) && seam[i] != j {
				newpixels[i][c] = input[i][j]
				c++
			}
		}
	}

	return newpixels
}

// findShortestPath travel from top to bottom
func findShortestPath(pixels [][]*Pixel) []int {
	height, width := len(pixels), len(pixels[0])
	cost := make([][]float64, height)
	for i := range cost {
		cost[i] = make([]float64, width)
	}

	for i := 0; i < width; i++ {
		cost[0][i] = pixels[0][i].E
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

			cost[row][col] = cost[row][col] + pixels[row][col].E
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
func transpose(pixels [][]*Pixel) [][]*Pixel {
	rows, cols := len(pixels), len(pixels[0])
	newpixels := make([][]*Pixel, cols)
	for i := range newpixels {
		newpixels[i] = make([]*Pixel, rows)
	}
	for i := 0; i < rows; i++ {
		for j := 0; j < cols; j++ {
			newpixels[j][i] = pixels[i][j]
		}
	}

	return newpixels
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

func calcDeltaSquare(c1, c2 color.Color) uint32 {
	r, g, b, _ := c1.RGBA()
	r, g, b = r/257, g/257, b/257

	r2, g2, b2, _ := c2.RGBA()
	r2, g2, b2 = r2/257, g2/257, b2/257

	d := (r-r2)*(r-r2) + (g-g2)*(g-g2) + (b-b2)*(b-b2)

	return d
}
