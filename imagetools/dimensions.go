package imagetools

import (
	"image"
)

func Palletize(m image.Image, n int) (res *image.Paletted) {
	return nil
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
