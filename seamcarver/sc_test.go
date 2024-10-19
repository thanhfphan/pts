package seamcarver

import (
	"image"
	"image/color"
	"testing"
)

func TestRetrieveSeamPath(t *testing.T) {
	tests := []struct {
		name   string
		energy [][]float64
		want   []int
	}{
		{
			name: "simple case",
			energy: [][]float64{
				{1, 2, 3},
				{4, 5, 6},
				{7, 8, 9},
			},
			want: []int{0, 0, 0},
		},
		{
			name: "another simple case",
			energy: [][]float64{
				{10, 11, 12},
				{10, 1, 10},
				{11, 10, 12},
			},
			want: []int{0, 1, 1},
		},
		{
			name: "complex case",
			energy: [][]float64{
				{1, 2, 3, 4},
				{4, 3, 2, 1},
				{1, 2, 3, 4},
				{4, 3, 2, 1},
			},
			want: []int{0, 1, 0, 1},
		},
		{
			name: "single column",
			energy: [][]float64{
				{1},
				{2},
				{3},
			},
			want: []int{0, 0, 0},
		},
		{
			name: "single row",
			energy: [][]float64{
				{1, 2, 3},
			},
			want: []int{0},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := retrieveSeamPath(tt.energy)
			if !equal(got, tt.want) {
				t.Errorf("retrieveSeamPath() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestTransposeImage(t *testing.T) {
	tests := []struct {
		name string
		img  image.Image
		want image.Image
	}{
		{
			name: "3x2 image",
			img: func() image.Image {
				img := image.NewRGBA(image.Rect(0, 0, 3, 2))
				img.Set(0, 0, color.RGBA{255, 0, 0, 255})
				img.Set(1, 0, color.RGBA{0, 255, 0, 255})
				img.Set(2, 0, color.RGBA{0, 0, 255, 255})
				img.Set(0, 1, color.RGBA{255, 255, 0, 255})
				img.Set(1, 1, color.RGBA{0, 255, 255, 255})
				img.Set(2, 1, color.RGBA{255, 0, 255, 255})
				return img
			}(),
			want: func() image.Image {
				img := image.NewRGBA(image.Rect(0, 0, 2, 3))
				img.Set(0, 0, color.RGBA{255, 0, 0, 255})
				img.Set(0, 1, color.RGBA{0, 255, 0, 255})
				img.Set(0, 2, color.RGBA{0, 0, 255, 255})
				img.Set(1, 0, color.RGBA{255, 255, 0, 255})
				img.Set(1, 1, color.RGBA{0, 255, 255, 255})
				img.Set(1, 2, color.RGBA{255, 0, 255, 255})
				return img
			}(),
		},
		{
			name: "1x1 image",
			img: func() image.Image {
				img := image.NewRGBA(image.Rect(0, 0, 1, 1))
				img.Set(0, 0, color.RGBA{123, 123, 123, 255})
				return img
			}(),
			want: func() image.Image {
				img := image.NewRGBA(image.Rect(0, 0, 1, 1))
				img.Set(0, 0, color.RGBA{123, 123, 123, 255})
				return img
			}(),
		},
		{
			name: "2x3 image",
			img: func() image.Image {
				img := image.NewRGBA(image.Rect(0, 0, 2, 3))
				img.Set(0, 0, color.RGBA{255, 0, 0, 255})
				img.Set(1, 0, color.RGBA{0, 255, 0, 255})
				img.Set(0, 1, color.RGBA{0, 0, 255, 255})
				img.Set(1, 1, color.RGBA{255, 255, 0, 255})
				img.Set(0, 2, color.RGBA{0, 255, 255, 255})
				img.Set(1, 2, color.RGBA{255, 0, 255, 255})
				return img
			}(),
			want: func() image.Image {
				img := image.NewRGBA(image.Rect(0, 0, 3, 2))
				img.Set(0, 0, color.RGBA{255, 0, 0, 255})
				img.Set(1, 0, color.RGBA{0, 0, 255, 255})
				img.Set(2, 0, color.RGBA{0, 255, 255, 255})
				img.Set(0, 1, color.RGBA{0, 255, 0, 255})
				img.Set(1, 1, color.RGBA{255, 255, 0, 255})
				img.Set(2, 1, color.RGBA{255, 0, 255, 255})
				return img
			}(),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := transposeImage(tt.img)
			if !imagesEqual(got, tt.want) {
				t.Errorf("transposeImage() = %v, want %v", got, tt.want)
			}
			back := transposeImage(got)
			if !imagesEqual(back, tt.img) {
				t.Errorf("transpose back should equal the original image back = %v, want %v", back, tt.img)
			}
		})
	}
}

func imagesEqual(a, b image.Image) bool {
	if !a.Bounds().Eq(b.Bounds()) {
		return false
	}
	for y := 0; y < a.Bounds().Dy(); y++ {
		for x := 0; x < a.Bounds().Dx(); x++ {
			if a.At(x, y) != b.At(x, y) {
				return false
			}
		}
	}
	return true
}

func equal(a, b []int) bool {
	if len(a) != len(b) {
		return false
	}
	for i := range a {
		if a[i] != b[i] {
			return false
		}
	}
	return true
}
