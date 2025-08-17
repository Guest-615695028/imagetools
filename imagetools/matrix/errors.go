package imagetools

type BasicError string

func (s BasicError) Error() string {
	return string(s)
}

const (
	ErrDimensions  BasicError = "inconsistent dimensions"
	ErrEmptyMatrix BasicError = "empty matrix"
	ErrNotSquare   BasicError = "not square matrix"
	ErrOutOfBounds BasicError = "out of bounds"
	ErrDivideBy0   BasicError = "division by zero"
	ErrLargeKernel BasicError = "the kernel is too large"
	ErrInvalidStep BasicError = "the step is not positive interger"
)

type DimensionError struct {
	Dims []Index2
	Op   string
	Why  error
}

// Matrix dimensional error
func (e DimensionError) Error() string {
	s := "matrix dimension error: " + e.Op + "("
	for i, v := range e.Dims {
		if i > 0 {
			s += ","
		}
		s += v.String()
	}
	return s + ")" + e.Why.Error()
}
func (e DimensionError) Unwrap() error {
	return e.Why
}
