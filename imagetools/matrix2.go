package imagetools

import (
	"image"
	"imagetools/types"
	"math"
	"math/cmplx"
)

func MinMax[T types.Real](m Matrix[T]) (a, b T) {
	if m = m.Clone(); len(m.val) <= 0 {
		return
	}
	a, b = m.val[0], m.val[0]
	for _, v := range m.val {
		a, b = min(a, v), max(b, v)
	}
	return
}

func Min[T types.Real](m Matrix[T]) T {
	if m = m.Clone(); len(m.val) <= 0 {
		return 0
	}
	a := m.val[0]
	for _, v := range m.val {
		a = min(a, v)
	}
	return a
}

func Max[T types.Real](m Matrix[T]) T {
	if m = m.Clone(); len(m.val) <= 0 {
		return 0
	}
	a := m.val[0]
	for _, v := range m.val {
		a = max(a, v)
	}
	return a
}

func ConvertMatrix[U, T types.Number](m Matrix[T]) Matrix[U] {
	return NewMatrix(m.x, m.y, ConvertSlice[U](m.val)...)
}
func MapMatrix[U, T types.Number](m Matrix[T], f func(T) U) Matrix[U] {
	return NewMatrix(m.x, m.y, ConvertSliceFunc(m.val, f)...)
}

func Normalize[T types.Real](m Matrix[T]) Matrix[uint8] {
	m1 := ConvertMatrix[float64](m)
	a, b := MinMax(m1)
	return ConvertMatrix[uint8](m1.Sub(a).Mul(255 / (b - a)))
}

func Absolutize[T types.Real](m Matrix[T]) Matrix[uint8] {
	m1 := MapMatrix(m, func(t T) float64 { return Abs(float64(t)) })
	_, b := MinMax(m1)
	return ConvertMatrix[uint8](m1.Mul(255 / b))
}

func LogAbsolutize[T types.Real](m Matrix[T]) Matrix[uint8] {
	m1 := MapMatrix(m, func(t T) float64 { return math.Log1p(Abs(float64(t))) })
	_, b := MinMax(m1)
	return ConvertMatrix[uint8](m1.Mul(255 / b))
}

func RGBA2Matrices(m image.Image) (r [4]Matrix[uint8]) {
	if m1 := RGBA(m); m1 != nil {
		for i := range r {
			r[i] = NewMatrix(m1.Rect.Dx(), m1.Rect.Dy(), Step(m1.Pix[i:], 4)...)
		}
	}
	return r
}

func Gray2Matrix(m *image.Gray) Matrix[uint8] {
	if m == nil {
		return Matrix[uint8]{}
	}
	m1 := Reduce(m.Pix, m.Stride, m.Rect, 1)
	return Matrix[uint8]{
		val: m1.Pix,
		x:   m1.Stride,
		y:   len(m1.Pix) / m1.Stride,
	}
}

func Matrix2Gray(m Matrix[uint8]) *image.Gray {
	m = m.Clone()
	return &image.Gray{
		Pix:    m.val,
		Stride: m.x,
		Rect:   image.Rect(0, 0, m.x, m.y),
	}
}

// Convert four matrices into RGBA image,
// undefind behavior if inconsistent dimensions
func Matrices2RGBA(r []Matrix[uint8]) (m *image.RGBA) {
	m = &image.RGBA{
		Pix:    make([]uint8, 4*len(r[0].val)),
		Stride: 4 * r[0].x,
		Rect:   image.Rect(0, 0, r[0].x, r[0].y),
	}
	for i, v := range r[:4] {
		for j, w := range v.val {
			if 4*j+i < len(m.Pix) {
				m.Pix[4*j+i] = w
			}
		}
	}
	return m
}

func Matrices2RGB(r []Matrix[uint8]) (m *image.RGBA) {
	m = &image.RGBA{
		Pix:    make([]uint8, 4*len(r[0].val)),
		Stride: 4 * r[0].x,
		Rect:   image.Rect(0, 0, r[0].x, r[0].y),
	}
	for i, v := range r[:3] {
		for j, w := range v.val {
			if 4*j+i < len(m.Pix) {
				m.Pix[4*j+i] = w
			}
		}
	}
	for i := 3; i < len(m.Pix); i += 4 {
		m.Pix[i] = 255
	}
	return m
}

func Shrink8[T types.Real](m Matrix[T]) Matrix[int8] {
	m = m.Clone()
	r := NewMatrix[int8](m.x, m.y)
	r.val = make([]int8, len(m.val))
	for i, v := range m.val {
		r.val[i] = int8(v)
		if v >= 127 {
			r.val[i] = 127
		} else if float64(v) <= -128 {
			r.val[i] = -128
		}
	}
	return r
}

func ShrinkU8[T types.Real](m Matrix[T]) Matrix[uint8] {
	m = m.Clone()
	r := NewMatrix[uint8](m.x, m.y)
	r.val = make([]uint8, len(m.val))
	for i, v := range m.val {
		r.val[i] = uint8(v)
		if v <= 0 {
			r.val[i] = 0
		} else if float64(v) >= 255 {
			r.val[i] = 255
		}
	}
	return r
}

func Convolve(i *image.RGBA, op Matrix[int], dx, dy int) *image.RGBA {
	ms := RGBA2Matrices(i)
	for i, v := range ms {
		c, _ := ConvertMatrix[int](v).Conv(op, dx, dy)
		ms[i] = ShrinkU8(c)
	}
	return Matrices2RGBA(ms[:])
}

// Fourier Matrix
func Fourier(n int) Matrix[complex128] {
	if n <= 0 {
		return Matrix[complex128]{}
	} else if n == 1 {
		return Matrix[complex128]{val: []complex128{1}, x: 1, y: 1}
	} else if n == 2 {
		return Matrix[complex128]{val: []complex128{1, 1, 1, -1}, x: 1, y: 1}
	}
	m := NewMatrix[complex128](n, n)
	j := 0
	for i := range m.val {
		t := float64(j) * 2 * math.Pi / float64(n)
		m.val[i] = complex(math.Cos(t), math.Sin(t))
		j = (j + i/n) % n
	}
	return m
}

func MakeComplex[T types.Real](m Matrix[T]) Matrix[complex128] {
	return MapMatrix(m, func(t T) complex128 { return complex(float64(t), 0) })
}
func MakeImag[T types.Real](m Matrix[T]) Matrix[complex128] {
	return MapMatrix(m, func(t T) complex128 { return complex(0, float64(t)) })
}

func GetReal(m Matrix[complex128]) Matrix[float64] {
	return MapMatrix(m, func(c complex128) float64 { return real(c) })
}
func GetImag(m Matrix[complex128]) Matrix[float64] {
	return MapMatrix(m, func(c complex128) float64 { return imag(c) })
}
func GetPhase(m Matrix[complex128]) Matrix[float64] {
	return MapMatrix(m, cmplx.Phase)
}
func GetAbs(m Matrix[complex128]) Matrix[float64] {
	return MapMatrix(m, cmplx.Abs)
}

func DFT[T types.Real](m Matrix[T]) Matrix[complex128] {
	m1, _ := Fourier(m.y).MulMat(MakeComplex(m))
	m1, _ = m1.MulMat(Fourier(m.x))
	return m1
}

func LR(m image.Image, w Matrix[uint], dx, dy int) (covs [4]Matrix[float64]) {
	if m.Bounds().Empty() {
		return
	}
	m2 := Split2(m, false)
	l, r := RGBA2Matrices(m2[0]), RGBA2Matrices(m2[1])
	for i := range 4 {
		covs[i] = NewMatrix[float64]((r[i].x-w.x)/dx+1, (r[i].y-w.y)/dy+1)
		for k, vr := range r[i].RangeSubMatrix(1, 1, w.x, w.y) {
			vl, _ := l[i].SubMatrix(k[0], k[1], k[0]+vr.x, k[1]+vr.y)
			covs[i].Assign(k[0], k[1], CoV(false, vl.val, vr.val))
		}
	}
	return covs
}

// Gausian funcion, e**(-x**2 / 2)
func Gauss(x float64) float64 {
	return math.Exp(-x / 2 * x)
}

func LaplaceGauss(n int, s float64) Matrix[float64] {
	m := NewMatrix[float64](n, n)
	dx := float64(n-1) / 2
	for i := range m.val {
		x, y := (float64(i%n)-dx)/s, (float64(i/n)-dx)/s
		m.val[i] = (x*x + y*y - 2) / s / s * Gauss(x) * Gauss(y)
	}
	return m
}

func HistogramizeMatrix(m Matrix[uint8]) Matrix[uint8] {
	h := [256]uint{}
	m = m.Clone()
	for _, v := range m.val {
		h[v]++
	}
	t := uint(0)
	for i, v := range h {
		h[i] = t*255 + v*uint(i)
		t += v
	}
	n := NewMatrix[uint8](m.x, m.y)
	for i, v := range m.val {
		k := h[v] / t
		if h[v]%t*2 >= t {
			k++
		}
		n.val[i] = uint8(k)
	}
	return n
}
