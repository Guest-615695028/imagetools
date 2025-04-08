package imagetools

import (
	"math/bits"
	"reflect"
	"slices"
	"strconv"
)

func Cond[C comparable, T any](c C, t, f T) T {
	if c != *new(C) {
		return t
	} else {
		return f
	}
}

// Ignore the returned error
func NoError[T any](t T, err error) T { return t }

func Bool(a any) bool {
	switch b := a.(type) {
	case nil:
		return false
	case bool:
		return b
	case int, int8, int16, int32, int64,
		uint, uint8, uint16, uint32, uint64, uintptr,
		float32, float64, complex64, complex128:
		return b != 0
	case string:
		return b != ""
	default:
		defer func() { recover() }()
		v := reflect.ValueOf(b)
		return v.IsValid() && !v.IsZero()
	}
}

func Abs[T Real](t T) T {
	if t < 0 {
		return -t
	} else {
		return t + 0 //avoid -0
	}
}
func Sign[T Real](t T) int {
	if t < 0 {
		return -1
	} else if t > 0 {
		return 1
	} else {
		return 0
	}
}
func IsNan[T Real](t T) bool {
	return t != t
}
func IsInf[T Real](t T) bool {
	return t != 0 && t == t/2
}

// Safe algorithm for x*y.n
func MulDiv(x, y, n uint) (uint, uint) {
	h, l := bits.Mul(x, y)
	return bits.Div(h, l, n)
}

// Format number into string, parameters:
//   - t: the number
//   - f: the formatter character as for fmt.Printf[]
//   - prec: the precision, or base
//   - pos: show plus sign or space
func FormatNumber[T Number](t T, f byte, prec int) (s string) {
	defer func() { recover() }()
	switch v := reflect.ValueOf(t); v.Kind() {
	case reflect.Int, reflect.Int8, reflect.Int16,
		reflect.Int32, reflect.Int64:
		s = FormatInt(v.Int(), f, prec)
	case reflect.Uint, reflect.Uint8, reflect.Uint16,
		reflect.Uint32, reflect.Uint64, reflect.Uintptr:
		s = FormatInt(v.Uint(), f, prec)
	default:
		switch f {
		case 'b', 'e', 'E', 'f', 'g', 'G', 'x', 'X':
			//no_op
		case 'F', 'v':
			f = 'f'
		case 'B':
			f = 'b'
		default:
			return "<invalid floating-point format>"
		}
		switch v.Kind() {
		case reflect.Float32:
			s = strconv.FormatFloat(v.Float(), f, prec, 32)
		case reflect.Float64:
			s = strconv.FormatFloat(v.Float(), f, prec, 64)
		case reflect.Complex64:
			s = strconv.FormatComplex(v.Complex(), f, prec, 64)
		case reflect.Complex128:
			s = strconv.FormatComplex(v.Complex(), f, prec, 128)
		default:
			s = v.String()
		}
	}
	if f >= 'A' && f <= 'Z' {
		s = ToUpper(s)
	}
	return s
}

func FormatInt[I Integer](i I, f byte, l int) string {
	if f <= 1 || f == 'd' || f == 'v' {
		f = 10
	}
	if l <= 1 && f <= 10 && i >= 0 && i < I(f) {
		return string('0' + i)
	}
	return formatInt(uint64(Abs(i)), l, i < 0, f)
}

func formatInt(a uint64, l int, s bool, f byte) string {
	ds, n := "0123456789abcdefghijklmnopqrstuvwxyz", uint64(f)
	switch f {
	case 'B', 'b':
		n = 2
	case 'O', 'o':
		n = 8
	case 'x':
		n = 16
	case 'X':
		ds, n = "0123456789ABCDEF", 16
	}
	if n < 2 || n > 36 {
		return "<illegal base>"
	}
	res := make([]byte, 0, 24)
	for a > 0 || len(res) < l {
		res = append(res, ds[a%n])
		a /= n
	}
	switch f {
	case 'b', 'B', 'o', 'O', 'x', 'X':
		res = append(res, f, '0')
	}
	if s {
		res = append(res, '-')
	}
	slices.Reverse(res)
	return string(res)
}

// String tools

func Repeat(r rune, n int) string {
	if n <= 0 {
		return ""
	} else if n == 1 {
		return string(r)
	}
	s := make([]rune, n)
	for i := range n {
		s[i] = r
	}
	return string(s)
}

func ToLower(s string) string {
	b := []byte(s)
	for i, v := range b {
		if v >= 'A' && v <= 'Z' {
			b[i] += 32
		}
	}
	return string(b)
}

func ToUpper(s string) string {
	b := []byte(s)
	for i, v := range b {
		if v >= 'a' && v <= 'z' {
			b[i] -= 32
		}
	}
	return string(b)
}

func Compare[T Ordered](x, y T) int {
	xn, yn := x != x, y != y
	switch {
	case xn && yn || x == y:
		return 0
	case xn || x < y:
		return -1
	case yn || x > y:
		return 1
	}
	return 0
}
func CompareNumber[X, Y Real](x X, y Y) int {
	switch {
	case x == 0 && y == 0:
		return 0
	case x <= 0 && y >= 0:
		return -1
	case x >= 0 && y <= 0:
		return 1
	case x > 0: //y>0
		fx, fy := float64(x), float64(y)
		return Compare(fx, fy)
		//return uint64(x) == uint64(y) || float64(x) == float64(y)
	default: // x<0 && y<0
		return Compare(float64(x), float64(y))
	}
	return 0
}
