package imagetools

import (
	"image"
	"image/color"
)

// Calculate the difference of two colors in RGBA
func ColorDiff(x, y color.Color) uint64 {
	d := uint64(0)
	x1 := Combine(x.RGBA())
	y1 := Combine(y.RGBA())
	for i, v := range x1 {
		u := uint64(v) - uint64(y1[i])
		d += u * u
	}
	return d
}

func PalleteDiff(x, y color.Palette) uint64 {
	d := uint64(0)
	if len(x) != len(y) {
		return d - 1
	}
	for i, v := range x {
		d += ColorDiff(v, y[i])
	}
	return d
}

// Palletize an image (m) into n colors, whose behavior may
// be undefined if m is larger than 1,073,774,592 pixels
func Palletize(m image.Image, n int) (res *image.Paletted) {
	b := m.Bounds()
	if b.Empty() {
		return new(image.Paletted)
	}
	res = &image.Paletted{
		Palette: make(color.Palette, n),
		Rect:    b,
		Stride:  b.Dx(),
		Pix:     make([]uint8, b.Dx()*b.Dy()),
	}
	for i := range n {
		res.Palette[i] = RandomRGBA64()
	}
	for t := 1; ; t++ {
		rgbas := make([][4]uint64, n)
		ns := make([]uint64, n)
		for p, c := range RangeImage(m) {
			i := res.Palette.Index(c)
			res.Pix[res.PixOffset(p.X, p.Y)] = uint8(i)
			cs := Combine(c.RGBA())
			for j := range rgbas[0] {
				rgbas[i][j] += uint64(cs[j])
			}
			ns[i]++
		}
		p1 := make(color.Palette, n)
		for i, v := range rgbas {
			if ni := ns[i]; ni == 0 {
				p1[i] = color.Transparent
			} else {
				p1[i] = color.RGBA64{
					R: uint16(v[0] / ni), G: uint16(v[1] / ni),
					B: uint16(v[2] / ni), A: uint16(v[3] / ni),
				}
			}
		}
		if SamePalettes(p1, res.Palette) {
			break
		}
		copy(res.Palette, p1)
	}
	return res
}

func Dimensions(m image.Image) (res map[string]*image.Gray) {
	res = make(map[string]*image.Gray, 4)
	defer func() { recover() }()
	switch m1 := CloneImage(m).(type) {
	case *image.RGBA:
		for i, k := range []string{"R", "G", "B", "A"} {
			res[k] = &image.Gray{
				Pix:    Step(m1.Pix[i:], 4),
				Stride: m1.Stride / 4,
				Rect:   m1.Rect,
			}
		}
	case *image.RGBA64:
		for i, k := range []string{"R", "G", "B", "A"} {
			res[k] = &image.Gray{
				Pix:    Step(m1.Pix[2*i:], 8),
				Stride: m1.Stride / 8,
				Rect:   m1.Rect,
			}
		}
	case *image.NRGBA:
		for i, k := range []string{"R", "G", "B", "A"} {
			res[k] = &image.Gray{
				Pix:    Step(m1.Pix[i:], 4),
				Stride: m1.Stride / 4,
				Rect:   m1.Rect,
			}
		}
	case *image.NRGBA64:
		for i, k := range []string{"R", "G", "B", "A"} {
			res[k] = &image.Gray{Pix: Step(m1.Pix[2*i:], 8),
				Stride: m1.Stride / 8,
				Rect:   m1.Rect,
			}
		}
	case *image.YCbCr:
		res["Y"] = &image.Gray{Pix: m1.Y, Stride: m1.YStride, Rect: m1.Rect}
		c := m1.Rect
		c.Max.X = m1.CStride
		c.Max.Y = len(m1.Cb) / m1.CStride
		res["Cb"] = &image.Gray{Pix: m1.Cb, Stride: m1.CStride, Rect: c}
		res["Cr"] = &image.Gray{Pix: m1.Cr, Stride: m1.CStride, Rect: c}
	case *image.NYCbCrA:
		y := &image.Gray{Pix: m1.Y, Stride: m1.YStride, Rect: m1.Rect}
		c := m1.Rect
		c.Max.X = m1.CStride
		c.Max.Y = len(m1.Cb) / m1.CStride
		res["Cb"] = &image.Gray{Pix: m1.Cb, Stride: m1.CStride, Rect: c}
		res["Cr"] = &image.Gray{Pix: m1.Cr, Stride: m1.CStride, Rect: c}
		a := &image.Gray{Pix: m1.A, Stride: m1.CStride, Rect: m1.Rect}
		for i, v := range a.Pix {
			y.Pix[i] = uint8(uint(y.Pix[i]) * uint(v) / 255)
		}
		res["Y"] = y
		res["A"] = a
	case *image.Gray:
		res["Gray"] = m1
	case *image.Gray16:
		res["Gray"] = &image.Gray{Pix: Step(m1.Pix, 2), Stride: m1.Stride / 2, Rect: m1.Rect}
	case *image.Alpha:
		res["Alpha"] = (*image.Gray)(m1)
	case *image.Alpha16:
		res["Alpha"] = &image.Gray{Pix: Step(m1.Pix, 2), Stride: m1.Stride / 2, Rect: m1.Rect}
	}
	return res
}
