package imagetools

import (
	"image"
	"image/color"
	"imagetools/types"
)

type Cropper interface {
	image.Image
	SubImage(r image.Rectangle) image.Image
}
type Setter interface {
	image.Image
	Set(x, y int, c color.Color)
}
type RGBA64Setter interface {
	image.Image
	SetRGBA64(x, y int, c color.RGBA64)
}
type RGBA64CapableImage interface {
	image.RGBA64Image
	RGBA64Setter
}

type CroppedImage struct {
	image.Image
	image.Rectangle
}

func (m CroppedImage) Bounds() image.Rectangle {
	return m.Rectangle
}

func (m CroppedImage) At(x, y int) color.Color {
	if (image.Point{x, y}.In(m.Rectangle)) {
		return m.Image.At(x, y)
	} else {
		return color.Transparent
	}
}

func (m CroppedImage) ColorModel() color.Model {
	return m.Image.ColorModel()
}

func (m CroppedImage) RGBA64At(x, y int) color.RGBA64 {
	if (image.Point{x, y}).In(m.Rectangle) {
		return color.RGBA64{0, 0, 0, 0}
	} else if m1, ok := m.Image.(image.RGBA64Image); ok {
		return m1.RGBA64At(x, y)
	} else {
		r, g, b, a := m.Image.At(x, y).RGBA()
		return color.RGBA64{uint16(r), uint16(g), uint16(b), uint16(a)}
	}
}

func Crop(m image.Image, r image.Rectangle) image.Image {
	switch m1 := m.(type) {
	case Cropper:
		return CloneImage(m1.SubImage(r))
	case image.Rectangle:
		return m1.Intersect(r)
	}
	return CroppedImage{Image: m, Rectangle: r}
}

// Range [a, b] for n steps, avoiding overflow
func RangeQR(a, b int, n int) func(func(int, int) bool) {
	if n <= 0 {
		return nil
	}
	d, e := (b-a)/n, (b-a)%n
	return func(yield func(k, v int) bool) {
		for q, r := a, 0; q < b; a = q {
			if !yield(q, r) {
				return
			}
			r += e
			q += d + r/n
			r %= n
		}
	}
}

// Split the image into y evens down and x evens across
func SplitN(im image.Image, x, y int) [][]image.Image {
	if x <= 0 || y <= 0 || im == nil {
		return [][]image.Image{{}}
	} else if x == 1 && y == 1 {
		return [][]image.Image{{CloneImage(im)}}
	}
	b := im.Bounds()
	if b.Empty() {
		return [][]image.Image{{}}
	}
	dx, dy := b.Dx(), b.Dy()
	if dx%x != 0 || dy%y != 0 {
		return [][]image.Image{{}}
	}
	dx /= x
	dy /= y
	images := make([][]image.Image, y)
	for m := range y {
		images[m] = make([]image.Image, x)
		y0 := b.Min.Y + m*dy
		for n := range x {
			x0 := b.Min.X + n*dx
			images[m][n] = CloneImage(Crop(im, image.Rect(x0, y0, x0+dx, y0+dy)))
		}
	}
	return images
}

// Split an image into two halves
func Split2(im image.Image, vert bool) [2]image.Image {
	b := im.Bounds()
	a := b
	if vert {
		a.Max.Y = (b.Min.Y + b.Max.Y) / 2
		b.Min.Y = a.Max.Y
	} else {
		a.Max.X = (b.Min.X + b.Max.X) / 2
		b.Min.X = a.Max.X
	}
	return [2]image.Image{Crop(im, a), Crop(im, b)}
}

func Gcd[I types.Integer](x, y I) I {
	for x != 0 && y != 0 {
		x, y = y, x%y
	}
	return types.Abs(x + y)
}
func Lcm[I types.Integer](x, y I) I {
	if x == 0 || y == 0 {
		return 0
	} else {
		return x / Gcd(x, y) * y //avoiding overflow of x*y
	}
}
