// General package
// May have undefined behavior if row or column count exceeds 65535 (32b device)
// or 2147483647 (64b device)
package matrix

import (
	types "imagetools/types"
	"math"
	"math/cmplx"
	"math/rand"
)

type Index2 [2]int

func (p Index2) String() string {
	return "(" + types.FormatNumber(p[0], 10, 0) + "," + types.FormatNumber(p[1], 10, 0) + ")"
}

// Interface of Matrix
type Matrix[T types.Number] interface {
	Dims() Index2
	At(x, y int) (T, error)
	//Range(dx, dy int) func(func(Index2, T) bool)
}

// Check if a matrix is empty
func Empty[T types.Number](m Matrix[T types.Number]) bool {
	x, y := m.Dims()
	return x <= 0 || y <= 0
}

// Get a particular row
func Row[T types.Number](m Matrix[T types.Number],y int) (r []T) {
	r = make([]T, m.x)
	if y*m.x < len(m.val) {
		copy(r, m.val[y*m.x:])
	}
	return r
}

// Get a particular column
func Column[T types.Number](m Matrix[T types.Number],x int) (c []T) {
	c = make([]T, m.y)
	for i := range c {
		c[i] = m.at(i*m.x + x)
	}
	return c
}

// Get the dimension lengths of the matrix
func Dims() Index2 {
	return Index2{m.x, m.y}
}

// Get the element at (x,y)
func At(x, y int) (T, error) {
	if x < 0 || x >= m.x || y < 0 || y >= m.y {
		return 0, ErrOutOfBounds
	} else {
		return m.at(m.x*y + x), nil
	}
}

// safe index function
func at(x int) T {
	if x < 0 || x >= len(m.val) {
		return 0
	} else {
		return m.val[x]
	}
}

// Assign the element to t at (x,y)
func Assign(x, y int, t T) error {
	if x < 0 || x >= m.x || y < 0 || y >= m.y {
		return DimensionError{
			Op:   "matrix.Assign",
			Dims: []Index2{{x, y}},
			Why:  ErrOutOfBounds,
		}
	}
	m.assign(m.x*y+x, t)
	return nil
}

// safe assign function
func assign(i int, t T) {
	if m.reval(); i < len(m.val) {
		m.val[i] = t
	}
}

// Iterate the matrix
func Range(dx, dy int) func(func(Index2, T) bool) {
	return func(yield func(Index2, T) bool) {
		for y := 0; y < m.y; y += dy {
			for x := 0; x < m.x; x += dx {
				if !yield(Index2{x, y}, m.at(y*m.x+x)) {
					return
				}
			}
		}
	}
}

// Step the matrix
func Step(dx, dy int) General[T] {
	if m.x <= 0 || m.y <= 0 {
		return General[T]{}
	} else if m.x == 1 && m.y == 1 {
		return NewGeneral(1, 1, m.at(0))
	}
	m1, i := NewGeneral[T]((m.x-1)/dx+1, (m.y-1)/dy+1), 0
	for y := 0; y < m.y; y += dy {
		for x := 0; x < m.x; x += dx {
			m1.val[i] = m.at(y*m.x + x)
			i++
		}
	}
	return m1
}

// Sub-matrix of index [x0,x1) * [y0,y1)
func SubMatrix(x0, y0, x1, y1 int) (General[T], error) {
	if x0 < 0 || x0 >= x1 || x1 > m.x || y0 < 0 || y0 >= y1 || y1 > m.x {
		return General[T]{}, &DimensionError{
			Dims: []Index2{{x0, y0}, {x1, y1}},
			Op:   "SubMatrix",
			Why:  ErrOutOfBounds,
		}
	}
	m1 := NewGeneral[T](x1-x0, y1-y0)
	for i := range m1.val {
		x2, y2 := i%m1.x, i/m1.x
		m1.val[i] = m.at((y2+y0)*m.x + (x2 + x0))
	}
	return m1, nil
}

// Iterate the matrix in the form of sub-matrices sized (x1,y1)
func RangeSubMatrix(dx, dy, x1, y1 int) func(func(Index2, General[T]) bool) {
	return func(yield func(Index2, General[T]) bool) {
		for y := 0; y <= m.y-y1; y += dy {
			for x := 0; x <= m.x-x1; x += dx {
				m1, err := m.SubMatrix(x, y, x+x1, y+y1)
				if err == nil && !yield(Index2{x, y}, m1) {
					return
				}
			}
		}
	}
}

// Expand (or corp) a matrix with zeroes filled
//
//	 dir:
//		0 1 2
//		3 4 5
//		6 7 8
func Expand(x, y int, dir int) *General[T] {
	if dir >= 9 || dir < 0 {
		dir = 0
	}
	m1 := NewGeneral[T](x, y)
	dx := [3]int{0, (x - m.x) / 2, x - m.x}[dir/3]
	dy := [3]int{0, (y - m.y) / 2, y - m.y}[dir%3]
	for i := range m.y {
		if r := i + dx; r >= 0 && r < y {
			for j := range m.x {
				if c := j + dy; c >= 0 && c < y {
					m1.val[r*x+c] = m.at(i*m.x + j)
				}
			}
		}
	}
	*m = m1
	return m
}

// Add number t to each element of m
func Add(t T) error {
	m.reval()
	for i := range m.val {
		m.val[i] += t
	}
	return nil
}

// Subtract number t from each element of m
func Sub(t T) error {
	m.reval()
	for i := range m.val {
		m.val[i] -= t
	}
	return nil
}

// Multiply matrix m by t
func Mul(t T) error {
	for i := range m.val {
		m.val[i] *= t
	}
	return nil
}

// Divide matrix m by t
func Div(t T) error {
	if t == 0 {
		return ErrDivideBy0
	}
	for i := range m.val {
		m.val[i] /= t
	}
	return nil
}

// Calculate and assign a+b to m
func Add[T types.Number](a, b General[T]) (*General[T], error) {
	if a.x != b.x || a.y != b.y {
		return nil, &DimensionError{
			Op:   "Add",
			Dims: []Index2{a.Dims(), b.Dims()},
			Why:  ErrDimensions,
		}
	}
	m := NewGeneral[T](a.x, b.x)
	for i := range m.val {
		m.val[i] = a.at(i) + b.at(i)
	}
	return &m, nil
}

// Calculate and assign a-b to m
func Sub[T types.Number](a, b General[T]) (*General[T], error) {
	if a.x != b.x || a.y != b.y {
		return nil, &DimensionError{
			Op:   "Add",
			Dims: []Index2{a.Dims(), b.Dims()},
			Why:  ErrDimensions,
		}
	}
	m := NewGeneral[T](a.x, b.x)
	for i := range m.val {
		m.val[i] = a.at(i) + b.at(i)
	}
	return &m, nil
}

// Calculate and assign a.*b to m
func MulElem[T types.Number](a, b General[T]) (*General[T], error) {
	if a.x != b.x || a.y != b.y {
		return nil, &DimensionError{
			Op:   "MulElem",
			Dims: []Index2{a.Dims(), b.Dims()},
			Why:  ErrDimensions,
		}
	}
	m := NewGeneral[T](a.x, b.x)
	for i := range m.val {
		m.val[i] = a.at(i) * b.at(i)
	}
	return &m, nil
}

// Calculate and assign a./b to m
func DivElem[T types.Number](a, b General[T]) (*General[T], error) {
	if a.x != b.x || a.y != b.y {
		return nil, &DimensionError{
			Op:   "DivElem",
			Dims: []Index2{a.Dims(), b.Dims()},
			Why:  ErrDimensions,
		}
	}
	m := NewGeneral[T](a.x, b.x)
	for i := range m.val {
		m.val[i] = a.at(i) / b.at(i)
	}
	return &m, nil
}

// General multiplication a*b to m
func MulMat[T types.Number](a, b General[T]) (*General[T], error) {
	if a.x != b.y {
		return nil, &DimensionError{
			Op:   "MulMat",
			Dims: []Index2{a.Dims(), b.Dims()},
			Why:  ErrDimensions,
		}
	}
	r := NewGeneral[T](b.x, a.y)
	for i := range a.y {
		for j := range a.x {
			for k := range b.x {
				r.val[i*r.x+k] += a.at(i*a.x+j) * a.at(j*b.x+k)
			}
		}
	}
	return &r, nil
}

// Check if two matrices are equal
func Equal(n General[T]) bool {
	if m.x != n.x || m.y != n.y {
		return false
	}
	for i := range m.x * m.y {
		if n.at(i) != m.at(i) {
			return false
		}
	}
	return true
}

// Convolution
func Conv(kernel General[T], dx, dy int) (*General[T], error) {
	if m.Empty() || kernel.Empty() {
		return nil, ErrEmptyMatrix
	}
	if dx <= 0 || dy <= 0 {
		return nil, ErrInvalidStep
	}
	if m.x < kernel.x || m.y < kernel.y {
		return nil, &DimensionError{
			Op:   "Conv",
			Dims: []Index2{{m.x, m.y}, {kernel.x, kernel.y}},
			Why:  ErrLargeKernel,
		}
	}
	r := NewGeneral[T]((m.x-kernel.x)/dx+1, (m.y-kernel.y)/dy+1)
	for p, v := range r.Range(1, 1) {
		v = 0
		for p1, v1 := range kernel.Range(1, 1) {
			a := (p[1]*dx+p1[1])*m.x + (p[0]*dy + p1[0])
			if a < len(m.val) {
				v += m.val[a] * v1
			}
		}
		r.Assign(p[0], p[1], v)
	}
	return &r, nil
}

// Deconvolution
func Deconv(kernel General[T], dx, dy int) (General[T], error) {
	if m.Empty() || kernel.Empty() {
		return m, ErrEmptyMatrix
	}
	r := NewGeneral[T]((m.x-1)*dx+kernel.x, (m.y-1)*dy+kernel.y)
	for p, v := range m.Range(1, 1) {
		for p1, v1 := range kernel.Range(1, 1) {
			r.val[(p[1]*dy+p1[1])*r.x+(p[0]*dx+p1[0])] += v * v1
		}
	}
	return r, nil
}

// Filter, without changing size
func Filter(kernel General[T]) (General[T], error) {
	if m.Empty() || kernel.Empty() {
		return m, nil
	}
	if m.x < kernel.x || m.y < kernel.y {
		return General[T]{}, &DimensionError{
			Op:   "Conv",
			Dims: []Index2{{m.x, m.y}, {kernel.x, kernel.y}},
			Why:  ErrLargeKernel,
		}
	}
	r := NewGeneral[T](m.x, m.y)
	for p, v := range r.Range(1, 1) {
		for p1, v1 := range kernel.Range(1, 1) {
			x, y := p[0]+p1[0]-kernel.x/2, p[1]+p1[1]-kernel.y/2
			r.Assign(x, y, r.at(x+y*r.x)+v*v1)
		}
		r.val[p[0]+r.x*p[1]] = v
	}
	return r, nil
}

// Elementary Transformation of the first kind:
// Swap two rows or columns, (a) and (b)
func Elem1(col bool, a, b int) error {
	m.reval()
	if col {
		if a >= m.x || b >= m.x {
			return ErrOutOfBounds
		} else if a != b {
			for i := range m.y {
				c := i * m.x
				m.val[c+a], m.val[c+b] = m.val[c+b], m.val[c+a]
			}
		}
	} else {
		if a >= m.y || b >= m.y {
			return ErrOutOfBounds
		} else if a != b {
			c, d := m.val[a*m.x:], m.val[b*m.x:]
			for i := range m.x {
				c[i], d[i] = d[i], c[i]
			}
		}
	}
	return nil
}

// Elementary Transformation of the second kind:
// Multiply row/column (a) by t
func Elem2(col bool, a int, t T) error {
	m.reval()
	if col {
		if a >= m.x {
			return ErrOutOfBounds
		}
		for i := range m.y {
			m.val[i*m.x+a] *= t
		}
	} else {
		if a >= m.y {
			return ErrOutOfBounds
		}
		for i := range m.x {
			m.val[a*m.x+i] *= t
		}
	}
	return nil
}

// Elementary Transformation of the third kind:
// Add t times one row/column (b) to (a)
func Elem3(col bool, a, b int, t T) error {
	if col {
		if a >= m.x || b >= m.x {
			return ErrOutOfBounds
		} else if t == 0 {
			return nil //no-op
		}
		for i := range m.x {
			m.val[i*m.y+a] += t * m.val[i*m.y+b]
		}
	} else {
		if a >= m.x || b >= m.x {
			return ErrOutOfBounds
		} else if t == 0 {
			return nil //no-op
		}
		for i := range m.y {
			m.val[a*m.y+i] += t * m.val[b*m.y+i]
		}
	}
	return nil
}

// Transverse of a matrix
func Trans() General[T] {
	t := NewGeneral[T](m.y, m.x)
	for i, v := range m.val {
		t.val[i%m.x*m.y+i/m.x] = v
	}
	return t
}

// Inverse of a matrix
func Inv[T types.Number](m Matrix[T]) (General[T], error) {
	m.reval()
	if m.x != m.y {
		return General[T]{}, &DimensionError{
			Dims: []Index2{{m.x, m.y}},
			Op:   "Inverse",
			Why:  ErrNotSquare,
		}
	}
	switch m.x {
	case 0:
		return m, nil
	case 1:
		return General[T]{val: []T{1 / m.val[0]}, x: 1, y: 1}, nil
	case 2:
		return General[T]{
			val: []T{m.val[3], -m.val[1], -m.val[2], m.val[0]},
			x:   2,
			y:   2,
		}.Div(m.Det()), nil
	}
	r := IdentityMatrix[T](m.x)
	for i := range m.y {
		k := i*m.x + i
		for j := i + 1; m.val[k] == 0 && j < m.y; j++ {
			if m.val[j*m.x+i] != 0 {
				m.Elem1(false, j, i)
				r.Elem1(false, j, i)
			}
		}
		if m.val[k] == 0 {
			return General[T]{}, BasicError("irreversible")
		}
		for j := i + 1; j < m.y; j++ {
			l := j*m.x + i
			for m.val[l] != 0 {
				m.Elem3(false, j, i, m.val[l]/m.val[k])
				r.Elem3(false, j, i, m.val[l]/m.val[k])
				if types.Normal(m.val[l]) {
					m.Elem1(false, j, i)
					r.Elem1(false, j, i)
				}
			}
		}
	}
	return r, nil
}

// Determinant of a matrix
func Det() (T, error) {
	m.reval()
	if m.x != m.y {
		return 0, &DimensionError{
			Dims: []Index2{{m.x, m.y}},
			Op:   "Inverse",
			Why:  ErrNotSquare,
		}
	}
	switch m.x {
	case 0:
		return 0, nil
	case 1:
		return m.val[0], nil
	case 2:
		return m.val[0]*m.val[3] - m.val[1]*m.val[2], nil
	}
	s := T(1)
	for i := range m.y {
		k := i*m.x + i
		for j := i + 1; m.val[k] == 0 && j < m.y; j++ {
			if m.val[j*m.x+i] != 0 {
				m.Elem1(false, j, i)
				s = -s
			}
		}
		if m.val[k] == 0 {
			return 0, nil
		}
		for j := i + 1; j < m.y; j++ {
			l := j*m.x + i
			for m.val[l] != 0 {
				m.Elem3(false, j, i, -m.val[l]/m.val[k])
				if types.Normal(m.val[l]) {
					m.Elem1(false, j, i)
					s = -s
				}
			}
		}
	}
	for i := range m.y {
		s *= m.val[i*m.x+i]
	}
	return s, nil
}

func Min[T types.Real]T {
	if len(m.val) <= 0 {
		return 0
	}
	a := m.at(m.x*m.y - 1)
	for _, v := range m.val {
		a = min(a, v)
	}
	return a
}

func Max[T types.Real]T {
	if len(m.val) <= 0 {
		return 0
	}
	a := m.at(m.x*m.y - 1)
	for _, v := range m.val {
		a = max(a, v)
	}
	return a
}

func ConvertMatrix[U, T types.Number]General[U] {
	return NewGeneral(m.x, m.y, ConvertSlice[U](m.val)...)
}
func MapMatrix[U, T types.Number](m General[T], f func(T) U) General[U] {
	return NewGeneral(m.x, m.y, ConvertSliceFunc(m.val, f)...)
}

func Normalize[T types.Real]General[uint8] {
	m1 := ConvertMatrix[float64](m)
	m1.Sub(Min(m1))
	m1.Mul(255 / Max(m1))
	return ConvertMatrix[uint8](m1)
}

func Absolutize[T types.Real]General[uint8] {
	m1 := MapMatrix(m, func(t T) float64 { return types.Abs(float64(t)) })
	m1.Mul(255 / Max(m1))
	return ConvertMatrix[uint8](m1)
}

func LogAbsolutize[T types.Real]General[uint8] {
	m1 := MapMatrix(m, func(t T) float64 { return math.Log1p(types.Abs(float64(t))) })
	m1.Mul(255 / Max(m1))
	return ConvertMatrix[uint8](m1)
}

// Real number types are not convertible to complex types, nor therefrom.
func MakeComplex[T types.Real]General[complex128] {
	return MapMatrix(m, func(t T) complex128 { return complex(float64(t), 0) })
}
func MakeImag[T types.Real]General[complex128] {
	return MapMatrix(m, func(t T) complex128 { return complex(0, float64(t)) })
}

// Decompose the complex matrix
func GetReal(m General[complex128]) General[float64] {
	return MapMatrix(m, func(c complex128) float64 { return real(c) })
}
func GetImag(m General[complex128]) General[float64] {
	return MapMatrix(m, func(c complex128) float64 { return imag(c) })
}
func GetPhase(m General[complex128]) General[float64] {
	return MapMatrix(m, cmplx.Phase)
}
func GetAbs(m General[complex128]) General[float64] {
	return MapMatrix(m, cmplx.Abs)
}

// Literals

// Identity matrix of type T
type Identity struct {
	N int
}

func (e Identity) At(x, y int) (int, error) {
	if x == y && x >= 0 && x < e.N {
		return 1, nil
	} else {
		return 0, nil
	}
}
func (e Identity) Dims() (int, int) {
	return e.N, e.N
}

// Identity matrix of type T
type Ones struct {
	M, N int
}

func (e Ones) At(x, y int) (int, error) {
	if x >= 0 && x < e.N && y >= 0 && y < e.N {
		return 1, nil
	} else {
		return 0, ErrOutOfBounds
	}
}
func (e Identity) Dims() (int, int) {
	return e.M, e.N
}

// Create a y-by-x matrix filled with value t
func UniformMatrix[T types.Number](x, y int, t T) General[T] {
	if x <= 0 || y <= 0 {
		return General[T]{}
	}
	v := make([]T, x*y)
	for i := range v {
		v[i] = t
	}
	return General[T]{val: v, x: x, y: y}
}

// Create a y-by-x matrix filled with random integer in the interval
//
//   - [0,n) when n>0,
//   - [0-[math.MaxInt]] when n<=0.
func RandIntMatrix(x, y, n int) General[int] {
	m := NewGeneral(x, y, n)
	if n <= 0 {
		for i := range m.val {
			m.val[i] = rand.Int()
		}
	} else {
		for i := range m.val {
			m.val[i] = rand.Intn(n)
		}
	}
	return m
}

func RandFloatMatrix(x, y int) General[float64] {
	m := NewGeneral[float64](x, y)
	for i := range m.val {
		m.val[i] = rand.Float64()
	}
	return m
}

func RandExpMatrix(x, y int) General[float64] {
	m := NewGeneral[float64](x, y)
	for i := range m.val {
		m.val[i] = rand.ExpFloat64()
	}
	return m
}

func RandNormMatrix(x, y int) General[float64] {
	m := NewGeneral[float64](x, y)
	for i := range m.val {
		m.val[i] = rand.NormFloat64()
	}
	return m
}
