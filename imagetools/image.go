package imagetools

import (
	"image"
	"image/color"
)

// Generic base for most of Go standard image types
type BasicImage struct {
	Pix []uint8
	// Stride is the Pix stride (in bytes) between vertically adjacent pixels.
	Stride int
	// Rect is the image's bounds.
	Rect image.Rectangle
}

// Reduce a generic image into its minimal clone
// with each pixel being n bytes long
func Reduce(p []uint8, s int, r image.Rectangle, n int) *BasicImage {
	if r.Empty() {
		return &BasicImage{Pix: []uint8{}}
	}
	x, y := r.Dx()*n, r.Dy()
	pix := make([]uint8, x*y)
	for i := range y {
		copy(pix[i*x:(i+1)*x], p[i*s:])
	}
	return &BasicImage{
		Pix:    pix,
		Stride: x,
		Rect:   r.Sub(r.Min),
	}
}

// Clone an image, keeping minimum slice length but the underlying type.
func CloneImage(m image.Image) image.Image {
	defer func() { recover() }()
	switch m1 := m.(type) {
	case nil:
		return nil
	case *image.RGBA:
		return (*image.RGBA)(Reduce(m1.Pix, m1.Stride, m1.Rect, 4))
	case *image.RGBA64:
		return (*image.RGBA64)(Reduce(m1.Pix, m1.Stride, m1.Rect, 8))
	case *image.NRGBA:
		return (*image.NRGBA)(Reduce(m1.Pix, m1.Stride, m1.Rect, 4))
	case *image.NRGBA64:
		return (*image.NRGBA64)(Reduce(m1.Pix, m1.Stride, m1.Rect, 8))
	case *image.Gray:
		return (*image.Gray)(Reduce(m1.Pix, m1.Stride, m1.Rect, 1))
	case *image.Gray16:
		return (*image.Gray16)(Reduce(m1.Pix, m1.Stride, m1.Rect, 2))
	case *image.Alpha:
		return (*image.Alpha)(Reduce(m1.Pix, m1.Stride, m1.Rect, 1))
	case *image.Alpha16:
		return (*image.Alpha16)(Reduce(m1.Pix, m1.Stride, m1.Rect, 2))
	case *image.CMYK:
		return (*image.CMYK)(Reduce(m1.Pix, m1.Stride, m1.Rect, 4))
	case *image.YCbCr:
		return CloneYCC(m1)
	case *image.NYCbCrA:
		ycc := CloneYCC(&m1.YCbCr)
		b := &image.NYCbCrA{
			YCbCr:   *ycc,
			AStride: ycc.YStride,
			A:       Reduce(m1.A, m1.AStride, m1.Rect, 1).Pix,
		}
		return b
	case image.Rectangle:
		return m1.Canon()
	case *image.Rectangle:
		return m1.Canon()
	case *image.Uniform:
		return &image.Uniform{C: m1.C}
	case *image.Paletted:
		b := Reduce(m1.Pix, m1.Stride, m1.Rect, 1)
		p := m1.Palette
		if len(m1.Palette) > 256 {
			p = p[0:256:256]
		}
		return &image.Paletted{
			Palette: p,
			Pix:     b.Pix,
			Rect:    b.Rect,
			Stride:  b.Stride,
		}
	default: //Hard to clone
		return RGBA(m)
	}
}

func CloneYCC(m *image.YCbCr) *image.YCbCr {
	b := &image.YCbCr{
		Rect:           m.Rect.Sub(m.Rect.Min),
		SubsampleRatio: m.SubsampleRatio,
	}
	c, r := Reduce(m.Y, m.YStride, b.Rect, 1), b.Rect
	b.Y, b.YStride = c.Pix, c.Stride
	x, y := r.Max.X, r.Max.Y
	switch m.SubsampleRatio {
	case image.YCbCrSubsampleRatio440, 440, 0440,
		image.YCbCrSubsampleRatio420, 420, 0420,
		image.YCbCrSubsampleRatio410, 410, 0410:
		y /= 2
	}
	switch m.SubsampleRatio {
	case image.YCbCrSubsampleRatio422, 422, 0422,
		image.YCbCrSubsampleRatio420, 420, 0420:
		x /= 2
	case image.YCbCrSubsampleRatio411, 411, 0411,
		image.YCbCrSubsampleRatio410, 410, 0410:
		x /= 4
	}
	b.Cb = make([]uint8, x*y)
	b.Cr = make([]uint8, x*y)
	b.CStride = x
	for i := range y {
		copy(b.Cb[i*x:(i+1)*x], m.Cb[i*m.CStride:])
		copy(b.Cr[i*x:(i+1)*x], m.Cr[i*m.CStride:])
	}
	return b
}

// Regularize an image into RGBA, with Rect starting at (0,0)
func RGBA(m image.Image) *image.RGBA {
	switch m1 := m.(type) {
	case nil:
		return nil
	case *image.RGBA:
		return (*image.RGBA)(Reduce(m1.Pix, m1.Stride, m1.Rect, 4))
	case *image.RGBA64:
		return (*image.RGBA)(Reduce(Step(m1.Pix, 2), m1.Stride/2, m1.Rect, 4))
	case *image.NRGBA:
		ret := (*image.RGBA)(Reduce(m1.Pix, m1.Stride, m1.Rect, 4))
		for i := 0; i < len(ret.Pix); i += 4 {
			p1 := ConvertSlice[uint16](m1.Pix[i : i+4])
			for j := range 3 {
				ret.Pix[i+j] = uint8(p1[j] * p1[3] / 255)
			}
		}
		return ret
	case *image.Gray:
		pix := make([]byte, 0, 4*len(m1.Pix))
		for _, v := range m1.Pix {
			pix = append(pix, v, v, v, 255)
		}
		return &image.RGBA{
			Pix:    pix,
			Stride: m1.Stride * 4,
			Rect:   m1.Rect.Sub(m1.Rect.Min),
		}
	case *image.Alpha:
		pix := make([]byte, 0, 4*len(m1.Pix))
		for _, v := range m1.Pix {
			pix = append(pix, v, v, v, v)
		}
		return &image.RGBA{
			Pix:    pix,
			Stride: m1.Stride * 4,
			Rect:   m1.Rect.Sub(m1.Rect.Min),
		}
	}
	rect := m.Bounds().Sub(m.Bounds().Min)
	pix := make([]byte, 4*rect.Max.X*rect.Max.Y)
	i := 0
	for _, c := range RangeImage(m) {
		r, g, b, a := c.RGBA()
		pix[i] = uint8(r >> 8)
		pix[i+1] = uint8(g >> 8)
		pix[i+2] = uint8(b >> 8)
		pix[i+3] = uint8(a >> 8)
		i += 4
		if i >= len(pix) {
			break
		}
	}
	return &image.RGBA{Pix: pix, Stride: 4 * rect.Max.X, Rect: rect}
}

// Noop series: Used to avoid nil pointer dereference
func Noop0(func() bool)               {}
func Noop1[T any](func(T) bool)       {}
func Noop2[T, U any](func(T, U) bool) {}

func RangeImage(i image.Image) func(func(image.Point, color.Color) bool) {
	if i == nil {
		return Noop2[image.Point, color.Color]
	} else if b := i.Bounds(); b.Empty() {
		return Noop2[image.Point, color.Color]
	} else {
		return func(yield func(image.Point, color.Color) bool) {
			for y := b.Min.Y; y < b.Max.Y; y++ {
				for x := b.Min.X; x < b.Max.X; x++ {
					if !yield(image.Point{x, y}, i.At(x, y)) {
						return
					}
				}
			}
		}
	}
}

func RangeRGBA64(i image.Image) func(func(image.Point, color.RGBA64) bool) {
	if i1, _ := i.(image.RGBA64Image); i1 == nil {
		return Noop2[image.Point, color.RGBA64]
	} else if b := i1.Bounds(); b.Empty() {
		return Noop2[image.Point, color.RGBA64]
	} else {
		return func(yield func(image.Point, color.RGBA64) bool) {
			for y := b.Min.Y; y < b.Max.Y; y++ {
				for x := b.Min.X; x < b.Max.X; x++ {
					if !yield(image.Point{x, y}, i1.RGBA64At(x, y)) {
						return
					}
				}
			}
		}
	}
}

func Pixel(i image.Image, x, y int) color.RGBA64 {
	switch j := i.(type) {
	case nil:
		return color.RGBA64{}
	case image.RGBA64Image:
		return j.RGBA64At(x, y)
	default:
		r, g, b, a := j.At(x, y).RGBA()
		return color.RGBA64{uint16(r), uint16(g), uint16(b), uint16(a)}
	}
}

func Uint16BE(h, l uint8) uint {
	return uint(h)<<8 | uint(l)
}

func Histogramize(m *image.Gray) *image.Gray {
	h := [256]uint{}
	m = (*image.Gray)(Reduce(m.Pix, m.Stride, m.Rect, 1))
	for _, v := range m.Pix {
		h[v]++
	}
	t := uint(0)
	for i, v := range h {
		h[i] = t*255 + v*uint(i)
		t += v
	}
	for i, v := range m.Pix {
		k := h[v] / t
		if h[v]%t*2 >= t {
			k++
		}
		m.Pix[i] = uint8(k)
	}
	return m
}
