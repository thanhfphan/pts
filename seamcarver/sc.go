package seamcarver

import (
	"fmt"
	"image"
	"image/color"
	"math"
	"strconv"
	"strings"
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
	h := len(sc.pixels)
	w := len(sc.pixels[0])

	for y := 0; y < h; y++ {
		for x := 0; x < w; x++ {
			// border
			if y == 0 || y == h-1 || x == 0 || x == w-1 {
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
			if seam[i] != j {
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
	// adj return the neighbors of the pixel at(x, y) to pixels (x − 1, y + 1), (x, y + 1), and (x + 1, y + 1),
	adj := func(y, x int) [][]int {
		r := [][]int{}
		if y+1 < height && x-1 >= 0 {
			r = append(r, []int{y + 1, x - 1})
		}
		if y+1 < height {
			r = append(r, []int{y + 1, x})
		}
		if y+1 < height && x+1 < width {
			r = append(r, []int{y + 1, x + 1})
		}

		return r
	}
	// travel graph
	distTo := map[string]float64{}
	pathTo := map[string]string{}
	var dfs func(int, int)
	dfs = func(y, x int) {
		k := encode(y, x)
		for _, arr := range adj(y, x) {
			visitkey := encode(arr[0], arr[1])
			pixel := pixels[arr[0]][arr[1]]
			val, has := distTo[visitkey]
			if !has || val > distTo[k]+pixel.E {
				distTo[visitkey] = distTo[k] + pixel.E
				pathTo[visitkey] = k
			}
			dfs(arr[0], arr[1])
		}
	}
	for x := 0; x < width; x++ {
		distTo[encode(0, x)] = 0.0
		dfs(0, x)
	}

	// get shortest path
	var smallest string
	for x := 0; x < width; x++ {
		k := encode(height-1, x)
		if smallest == "" || distTo[k] < distTo[smallest] {
			smallest = k
		}
	}
	path := []int{}
	current := smallest
	for current != "" {
		_, tmpx := decode(current)
		path = append(path, tmpx)
		current = pathTo[current]
	}
	// reverse
	for i, j := 0, len(path)-1; i < j; i, j = i+1, j-1 {
		path[i], path[j] = path[j], path[i]
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

func encode(y, x int) string {
	return fmt.Sprintf("%d-%d", y, x)
}

func decode(key string) (int, int) {
	items := strings.Split(key, "-")
	y, _ := strconv.Atoi(items[0])
	x, _ := strconv.Atoi(items[1])
	return y, x
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
