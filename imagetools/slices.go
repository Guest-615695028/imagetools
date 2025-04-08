package imagetools

import (
	"math"
	"reflect"
	"slices"
	"imagetools/types"
)

// Generic operator < with NaN always less than other values
func Compare[T types.Ordered](x, y T) int {
	xNaN, yNaN := x != x, y != y
	switch {
	case xNaN && yNaN || x == y:
		return 0
	case xNaN || x < y:
		return -1
	case yNaN || x > y:
		return +1
	default:
		return 0
	}
}

// Sort and crop s so that
//
//	a <= s[0] && s[len(s)-1] <= b
//
// May share storage with s
func SortAndCrop[S ~[]E, E Ordered](s S, a, b E) S {
	return SortAndCropFunc(s, a, b, Compare)
}

func SortAndCropFunc[S ~[]E, E any](s S, a, b E, f func(E, E) int) S {
	if len(s) <= 0 || f(a, a) != 0 || f(b, b) != 0 || f(*new(E), *new(E)) != 0 {
		return s
	} else if f(a, b) > 0 {
		s = SortAndCropFunc(s, b, a, f)
		slices.Reverse(s)
		return s
	} else if f(a, b) == 0 {
		// Lower time complexity to O(n)
		u := make([]E, 0, len(s))
		for _, v := range s {
			if f(v, a) == 0 && f(v, b) == 0 {
				u = append(u, v)
			}
		}
		return u
	}
	slices.SortFunc(s, f)
	i, j := 0, len(s)-1
	if f(s[j], a) < 0 || f(s[i], b) > 0 {
		return []E{}
	}
	for f(s[i], a) < 0 {
		i++
	}
	for f(s[j], b) > 0 {
		j--
	}
	if i > j {
		return []E{}
	} else {
		return s[i : j+1 : j+2-i]
	}
}

// Fill generate a slice of n-times x.
func Fill[E any](x E, n int) (s []E) {
	if n <= 0 {
		return []E{}
	}
	s = make([]E, n)
	for n--; n >= 0; n-- {
		s[n] = x
	}
	return
}

// FillPointer generate a slice of n-times pointer to values equal to x,
// isolated to each other.
func FillPointer[E any](x E, n int) (s []*E) {
	s = make([]*E, n)
	for n--; n >= 0; n-- {
		s[n] = new(E)
		*s[n] = x
	}
	return
}

// Step a slice by n into independent storage
func Step[S ~[]E, E any](s S, n int) S {
	if n <= 0 || len(s) <= 0 {
		return S{}
	} else if len(s) <= n {
		return S{s[0]}
	} else if n == 1 {
		return append(S{}, s...)
	}
	r := make(S, (len(s)-1)/n+1)
	for i := range r {
		r[i] = s[i*n]
	}
	return r
}

// Make a byte-reversed copy
func ReverseBytes[I types.Integer, S ~[]I](b S) S {
	c := slices.Clone(b)
	for i, v := range b {
		c[i] = ^v
	}
	return c
}

// Make a reversed copy
func Reverse[S ~[]E, E any](s S) S {
	c := slices.Clone(s)
	slices.Reverse(c)
	return c
}

func SortAndCropUnique[S ~[]E, E Ordered](s S, a, b E) S {
	s = SortAndCrop(s, a, b)
	r := make(S, 0, len(s))
	for i, v := range s {
		if i == 0 || v > s[i-1] {
			r = append(r, v)
		}
	}
	return r
}

// Convert a slice of real numbers to another type of slice
func ConvertSlice[R, E any, S ~[]E](s S) []R {
	r, tr := make([]R, len(s)), reflect.TypeFor[R]()
	for i, v := range s {
		if vv := reflect.ValueOf(v); vv.Type().ConvertibleTo(tr) {
			r[i] = vv.Convert(tr).Interface().(R)
		}
	}
	return r
} // Convert a slice of real numbers to another type of slice
func ConvertSliceFunc[R, E any, S ~[]E](s S, f func(E) R) []R {
	r := make([]R, len(s))
	for i, v := range s {
		r[i] = f(v)
	}
	return r
}

func Connect[S ~[]E, E any](ss ...S) (s S) {
	for _, v := range ss {
		s = append(s, v...)
	}
	return s
}

// Calculate the p-power average value of x
func Mean[E types.Real](p float64, x ...E) E {
	n, m := len(x), 0.0
	if n <= 0 {
		return 0
	}
	if math.IsInf(p, 1) {
		return slices.Max(x)
	}
	if math.IsInf(p, -1) {
		return slices.Min(x)
	}
	if p == 0 {
		m = 1.0
		for _, v := range x {
			m *= math.Pow(float64(v), 1/float64(n))
		}
	} else {
		for _, v := range x {
			m += math.Pow(float64(v), p) / float64(n)
		}
	}
	return E(m)
}

// Calculate the medain value of x
func Median[E types.Real](x ...E) E {
	if len(x) <= 0 {
		return 0
	}
	slices.Sort(x)
	return x[len(x)/2]
}

// Calculate the most frequent value of x.
// If there are more than one, it will return the smallest
func Mode[E types.Real](x ...E) E {
	if len(x) <= 0 {
		return 0
	}
	slices.Sort(x)
	var m, p E
	n, k := 0, 1
	for _, v := range x {
		if p == v {
			k++
		} else {
			if k > n {
				n, m = k, p
			}
			p, k = v, 1
		}
	}
	return m
}

// Calculate the variance of x
func Variance[E types.Real](s bool, x ...E) float64 {
	n := len(x)
	if n <= 0 {
		return 0
	}
	var s1, s2 float64
	for _, v := range x {
		s2 += float64(v) / float64(n) * float64(v)
		s1 += float64(v) / float64(n)
	}
	if s2 -= s1 * s1; s {
		s2 /= 1 - 1/float64(n)
	}
	return s2
}

// calculate the covariance of x and y.
func CoV[X ~[]E1, Y ~[]E2, E1, E2 types.Real](s bool, x X, y Y) float64 {
	n := min(len(x), len(y))
	sx, sy, s2 := 0.0, 0.0, 0.0
	for i := range n {
		s2 += float64(x[i]) / float64(n) * float64(y[i])
		sx += float64(x[i]) / float64(n)
		sy += float64(y[i]) / float64(n)
	}
	if s2 -= sx * sy; s {
		s2 /= 1 - 1/float64(n)
	}
	return s2
}
