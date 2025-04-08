// Matrix package
// May have undefined behavior if row or column count exceeds 65535 (32b device)
// or 2147483647 (64b device)
package imagetools

import (
	"fmt"
	"imagetools/types"
)

type Index2 [2]int

func (p Index2) String() string {
	return "(" + types.FormatNumber(p[0], 10, 0) + "," + types.FormatNumber(p[1], 10, 0) + ")"
}

type Matrix[T types.Number] struct {
	val []T //values, y-by-x matrix
	x   int //column count, row length
	y   int //row count, column length
}

// Identity matrix of type T
func IdentityMatrix[T types.Number](n int) Matrix[T] {
	v := make([]T, n*n)
	for i := range n {
		v[i*(n+1)] = 1
	}
	return Matrix[T]{val: v, x: n, y: n}
}

// Create a new r-by-c matrix, zero matrix by default
func NewMatrix[T types.Number](x, y int, t ...T) Matrix[T] {
	if x <= 0 || y <= 0 {
		return Matrix[T]{}
	}
	v := make([]T, x*y)
	copy(v, t)
	return Matrix[T]{val: v, x: x, y: y}
}

// Clone a matrix
func (m Matrix[T]) Clone() Matrix[T] {
	if m.x <= 0 || m.y <= 0 {
		return Matrix[T]{}
	}
	v := make([]T, m.x*m.y)
	copy(v, m.val)
	return Matrix[T]{val: v, x: m.x, y: m.y}
}

// Clone a matrix
func (m Matrix[T]) Empty() bool {
	return m.x <= 0 || m.y <= 0
}

// Expand (or corp) a matrix with zeroes filled
//
//	 dir:
//		0 1 2
//		3 4 5
//		6 7 8
func (m Matrix[T]) Expand(x, y int, dir int) (m1 Matrix[T]) {
	if dir >= 9 || dir < 0 {
		dir = 0
	}
	m1 = NewMatrix[T](x, y)
	dx := [3]int{0, (x - m.x) / 2, x - m.x}[dir/3]
	dy := [3]int{0, (y - m.y) / 2, y - m.y}[dir%3]
	for i := range m.y {
		if r := i + dx; r >= 0 && r < y {
			for j := range m.x {
				if c := j + dy; c >= 0 && c < y {
					m1.val[r*x+c] = m.val[i*m.x+j]
				}
			}
		}
	}
	return m1
}

// Get a particular row
func (m Matrix[T]) Row(y int) []T {
	return Step(m.val[y*m.x:], 1)
}

// Get a particular column
func (m Matrix[T]) Column(x int) []T {
	return Step(m.val[x:], m.x)
}

// Get the dimension lengths of the matrix
func (m Matrix[T]) Dims() Index2 {
	return Index2{m.x, m.y}
}

// Get the element at (x,y)
func (m Matrix[T]) At(x, y int) (T, error) {
	if x < 0 || x >= m.x || y < 0 || y >= m.y {
		return 0, ErrOutOfBounds
	} else if a := m.x*y + x; a > len(m.val) {
		return 0, nil
	} else {
		return m.val[a], nil
	}
}

// Assign the element to t at (x,y)
func (m Matrix[T]) Assign(x, y int, t T) error {
	if x < 0 || x >= m.x || y < 0 || y >= m.y {
		return DimensionError{
			Op:   "matrix.Assign",
			Dims: []Index2{{x, y}},
			Why:  ErrOutOfBounds,
		}
	}
	m.val[m.x*y+x] = t
	return nil
}

// Iterate the matrix
func (m Matrix[T]) Range(dx, dy int) func(func(Index2, T) bool) {
	m = m.Clone()
	return func(yield func(Index2, T) bool) {
		for y := 0; y < m.y; y += dy {
			for x := 0; x < m.x; x += dx {
				if !yield(Index2{x, y}, m.val[y*m.x+x]) {
					return
				}
			}
		}
	}
}

// Step the matrix
func (m Matrix[T]) Step(dx, dy int) Matrix[T] {
	m = m.Clone()
	if m.x <= 0 || m.y <= 0 {
		return Matrix[T]{}
	} else if m.x == 1 && m.y == 1 {
		return NewMatrix(1, 1, m.val[0])
	}
	m1, i := NewMatrix[T]((m.x-1)/dx+1, (m.y-1)/dy+1), 0
	for y := 0; y < m.y; y += dy {
		for x := 0; x < m.x; x += dx {
			m1.val[i] = m.val[y*m.x+x]
			i++
		}
	}
	return m1
}

// Sub-matrix of index [x0,x1) * [y0,y1)
func (m Matrix[T]) SubMatrix(x0, y0, x1, y1 int) (Matrix[T], error) {
	if x0 < 0 || x0 >= x1 || x1 > m.x || y0 < 0 || y0 >= y1 || y1 > m.x {
		return Matrix[T]{}, &DimensionError{
			Dims: []Index2{{x0, y0}, {x1, y1}},
			Op:   "SubMatrix",
			Why:  ErrOutOfBounds,
		}
	}
	m1 := NewMatrix[T](x1-x0, y1-y0)
	for i := range m1.val {
		x2, y2 := i%m1.x, i/m1.x
		m1.val[i] = m.val[(y2+y0)*m.x+(x2+x0)]
	}
	return m1, nil
}

// Iterate the matrix in the form of sub-matrices sized (x1,y1)
func (m Matrix[T]) RangeSubMatrix(dx, dy, x1, y1 int) func(func(Index2, Matrix[T]) bool) {
	m = m.Clone()
	return func(yield func(Index2, Matrix[T]) bool) {
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

// Add number t to each element of m
func (m Matrix[T]) Add(t T) Matrix[T] {
	m = m.Clone()
	for i := range m.val {
		m.val[i] += t
	}
	return m
}

// Add matrix n to m,
func (m Matrix[T]) AddElem(n Matrix[T]) (Matrix[T], error) {
	if m.Empty() {
		return n.Clone(), nil
	} else if n.Empty() {
		return m.Clone(), nil
	} else if m.x != n.x || m.y != n.y {
		return Matrix[T]{}, &DimensionError{
			Op:   "Add",
			Dims: []Index2{{m.x, m.y}, {n.x, n.y}},
			Why:  ErrDimensions,
		}
	}
	m = m.Clone()
	for i, v := range n.val {
		if i >= len(m.val) {
			break
		}
		m.val[i] += v
	}
	return m, nil
}

// Subtract number t from each element of m
func (m Matrix[T]) Sub(t T) Matrix[T] {
	m = m.Clone()
	for i := range m.val {
		m.val[i] -= t
	}
	return m
}

// Subtract matrix n from m
func (m Matrix[T]) SubElem(n Matrix[T]) (Matrix[T], error) {
	m = m.Clone()
	if n.Empty() {
		return m, nil
	} else if m.x != n.x || m.y != n.y {
		return Matrix[T]{}, &DimensionError{
			Op:   "Sub",
			Dims: []Index2{{m.x, m.y}, {n.x, n.y}},
			Why:  ErrDimensions,
		}
	}
	for i, v := range n.val {
		if i >= len(m.val) {
			break
		}
		m.val[i] -= v
	}
	return m, nil
}

// Multiply matrix m by t
func (m Matrix[T]) Mul(t T) Matrix[T] {
	m = m.Clone()
	for i := range m.val {
		m.val[i] *= t
	}
	return m
}

func (m Matrix[T]) MulElem(n Matrix[T]) (Matrix[T], error) {
	if m.Empty() {
		return n.Clone(), nil
	} else if n.Empty() {
		return m.Clone(), nil
	} else if m.x != n.x || m.y != n.y {
		return Matrix[T]{}, &DimensionError{
			Op:   "MulElem",
			Dims: []Index2{{m.x, m.y}, {n.x, n.y}},
			Why:  ErrDimensions,
		}
	}
	m = m.Clone()
	for i, v := range n.val {
		if i >= len(m.val) {
			break
		}
		m.val[i] *= v
	}
	return m, nil
}

// Multiply matrix m by another matrix n
func (m Matrix[T]) MulMat(n Matrix[T]) (Matrix[T], error) {
	m = m.Clone()
	if m.x != n.y {
		return Matrix[T]{}, &DimensionError{
			Op:   "MulMat",
			Dims: []Index2{{m.x, m.y}, {n.x, n.y}},
			Why:  ErrDimensions,
		}
	}
	r := NewMatrix[T](n.x, m.y)
	for i := range m.y {
		for j := range m.x {
			for k := range n.x {
				r.val[i*r.x+k] += m.val[i*m.x+j] * n.val[j*n.x+k]
			}
		}
	}
	return r, nil
}

// Divide matrix m by t
func (m Matrix[T]) Div(t T) Matrix[T] {
	m = m.Clone()
	for i := range m.val {
		m.val[i] /= t
	}
	return m
}
func (m Matrix[T]) DivElem(n Matrix[T]) (Matrix[T], error) {
	m = m.Clone()
	if n.Empty() {
		return m, nil
	}
	if m.x != n.x || m.y != n.y {
		return Matrix[T]{}, &DimensionError{
			Op:   "DivElem",
			Dims: []Index2{{m.x, m.y}, {n.x, n.y}},
			Why:  ErrDimensions,
		}
	}
	var err error
	for i, v := range n.val {
		if v == 0 && err == nil {
			err = &DimensionError{
				Op:   "DivElem",
				Dims: []Index2{{i % len(n.val), i / len(n.val)}},
				Why:  ErrDivideBy0,
			}
		}
		if i >= len(m.val) {
			break
		}
		m.val[i] /= v
	}
	return m, err
}

// Check if two matrices are equal
func (m Matrix[T]) Equal(n Matrix[T]) bool {
	defer func() { recover() }()
	m = m.Clone()
	if m.x != n.x || m.y != n.y {
		return false
	}
	for i, v := range m.val {
		if n.val[i] != v {
			return false
		}
	}
	return true
}

// Convolution
func (m Matrix[T]) Conv(n Matrix[T], dx, dy int) (Matrix[T], error) {
	if m.Empty() || n.Empty() {
		return m, nil
	}
	if m.x < n.x || m.y < n.y || dx <= 0 || dy <= 0 {
		return Matrix[T]{}, &DimensionError{
			Op:   "Conv",
			Dims: []Index2{},
			Why:  ErrDimensions,
		}
	}
	m, n = m.Clone(), n.Clone()
	r := NewMatrix[T]((m.x-n.x)/dx+1, (m.y-n.y)/dy+1)
	for p, v := range r.Range(dx, dy) {
		v = 0
		for p1, v1 := range n.Range(1, 1) {
			a := (p[1]+p1[1])*m.x + (p[0] + p1[0])
			if a < len(m.val) {
				v += m.val[a] * v1
			}
		}
		r.val[p[0]+r.x*p[1]] = v
	}
	return r, nil
}

// Elementary Transformation of the first kind:
// Swap two rows or columns, (a) and (b)
func (m Matrix[T]) Elem1(col bool, a, b int) (Matrix[T], error) {
	m = m.Clone()
	if col {
		if a >= m.x || b >= m.x {
			return m, ErrOutOfBounds
		} else if a != b {
			for i := range m.y {
				c := i * m.x
				m.val[c+a], m.val[c+b] = m.val[c+b], m.val[c+a]
			}
		}
	} else {
		if a >= m.y || b >= m.y {
			return m, ErrOutOfBounds
		} else if a != b {
			c, d := m.val[a*m.x:], m.val[b*m.x:]
			for i := range m.x {
				c[i], d[i] = d[i], c[i]
			}
		}
	}
	return m, nil
}

// Elementary Transformation of the second kind:
// Multiply row/column (a) by non-zero number
func (m Matrix[T]) Elem2(col bool, a int, t T) (Matrix[T], error) {
	m = m.Clone()
	if t == 0 {
		return m, ErrOutOfBounds
	}
	if col {
		if a >= m.x {
			return m, ErrOutOfBounds
		}
		for i := range m.y {
			m.val[i*m.x+a] *= t
		}
	} else {
		if a >= m.y {
			return m, ErrOutOfBounds
		}
		for i := range m.x {
			m.val[a*m.x+i] *= t
		}
	}
	return m, nil
}

// Elementary Transformation of the third kind:
// Add t times one row/column (b) to (a)
func (m Matrix[T]) Elem3(col bool, a, b int, t T) (Matrix[T], error) {
	if col {
		if a >= m.x || b >= m.x {
			return m, ErrOutOfBounds
		} else if t == 0 {
			return m, nil //no-op
		}
		for i := range m.x {
			m.val[i*m.y+a] += t * m.val[i*m.y+b]
		}
	} else {
		if a >= m.x || b >= m.x {
			return m, ErrOutOfBounds
		} else if t == 0 {
			return m, nil //no-op
		}
		for i := range m.y {
			m.val[a*m.y+i] += t * m.val[b*m.y+i]
		}
	}
	return m, nil
}

// Transverse of a matrix
func (m Matrix[T]) Trans() Matrix[T] {
	t := NewMatrix[T](m.y, m.x)
	for i, v := range m.val {
		t.val[i%m.x*m.y+i/m.x] = v
	}
	return t
}

// Inverse of a matrix
func (m Matrix[T]) Inv() (Matrix[T], error) {
	m = m.Clone()
	if m.x != m.y {
		return Matrix[T]{}, &DimensionError{
			Dims: []Index2{{m.x, m.y}},
			Op:   "Inverse",
			Why:  ErrNotSquare,
		}
	}
	switch m.x {
	case 0:
		return m, nil
	case 1:
		return Matrix[T]{val: []T{1 / m.val[0]}, x: 1, y: 1}, nil
	case 2:
		return Matrix[T]{val: []T{m.val[3], -m.val[1], -m.val[2], m.val[0]}, x: 2, y: 2}, nil
	}
	return m, nil
}

// Determinant of a matrix
func (m Matrix[T]) Det() (T, error) {
	m = m.Clone()
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
		if m.val[k] == 0 {
			for j := i + 1; j < m.y; j++ {
				if m.val[j*m.x+i] != 0 {
					m.Elem1(false, j, i)
					s = -s
				}
			}
		}
		if m.val[k] == 0 {
			return 0, nil
		}
		for j := i + 1; j < m.y; j++ {
			x := m.val[j*m.x+i]
			if x != 0 {
				m.Elem2(false, j, m.val[k])
				m.Elem3(false, j, i, x)
			}
		}
	}
	for i := range m.y {
		s *= m.val[i*m.x+i]
	}
	return s, nil
}

/*
Format the matrix

Flags:
  - '#': print the type and dimension information
  - '+': print '+' if the number is not negative
  - '-':
*/
func (m Matrix[T]) Format(f fmt.State, r rune) {
	defer func() { recover() }()
	m = m.Clone()
	if f.Flag('#') {
		fmt.Fprintf(f, "%T(%d,%d)", m, m.x, m.y)
	}
	s := fmt.Sprintf("%"+string(r), T(0))
	if s[0] == '%' && s[1] == '!' {
		f.Write([]byte("[" + s + "]"))
		return
	}
	if m.x <= 0 || m.y <= 0 {
		f.Write([]byte("[]"))
		return
	}
	prec, ok := f.Precision()
	if !ok {
		prec = 6
	}
	wid, _ := f.Width()
	wid = min(wid, m.y)
	if wid > 0 { // multi-line output
		outs := make([][]string, m.y)
		ls := make([]int, wid)
		for i := range m.y {
			outs[i] = make([]string, m.x)
			for j := range m.x {
				s = FormatNumber(m.val[i*m.x+j], byte(r), prec)
				outs[i][j] = s
				ls[i%wid] = max(ls[i%wid], len(s))
			}
		}
		for i, row := range outs {
			if i == 0 {
				f.Write([]byte{'['})
			} else if wid < m.y {
				f.Write([]byte{'\n'})
			}
			for j, v := range row {
				if j > 0 {
					f.Write([]byte{','})
				}
				if j%wid == 0 {
					f.Write([]byte{'\n', '\t'})
				}
				f.Write(Fill[byte](' ', ls[j%wid]-len(v)))
				f.Write([]byte(v))
			}
		}
	} else { //one-line output
		for i, v := range m.val {
			if i == 0 {
				f.Write([]byte{'['})
			} else if i%m.x == 0 {
				f.Write([]byte{';'})
			} else {
				f.Write([]byte{','})
			}
			f.Write([]byte(FormatNumber(v, byte(r), prec)))
		}
	}
	f.Write([]byte{']'})
}
