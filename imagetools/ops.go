package imagetools

var (
	Roberts = []Matrix[int]{
		NewMatrix(2, 2, 1, 0, 0, -1),
		NewMatrix(2, 2, 0, 1, -1, 0),
	}
	Sobel = []Matrix[int]{
		NewMatrix(3, 3, -1, -2, -1, 0, 0, 0, 1, 2, 1),
		NewMatrix(3, 3, -2, -1, 0, -1, 0, 1, 0, 1, 2),
		NewMatrix(3, 3, -1, 0, 1, -2, 0, 2, -1, 0, 1),
		NewMatrix(3, 3, 0, 1, 2, -1, 0, 1, -2, -1, 0),
	}
	Prewitt = []Matrix[int]{
		NewMatrix(3, 3, -1, -1, -1, 0, 0, 0, 1, 1, 1),
		NewMatrix(3, 3, -1, -1, 0, -1, 0, 1, 0, 1, 1),
		NewMatrix(3, 3, -1, 0, 1, -1, 0, 1, -1, 0, 1),
		NewMatrix(3, 3, 0, 1, 1, -1, 0, 1, -1, -1, 0),
	}
	Laplace   = NewMatrix(3, 3, -0, -1, -0, -1, 4, -1, -0, -1, -0)
	Laplace8  = NewMatrix(3, 3, -1, -1, -1, -1, 8, -1, -1, -1, -1)
	Laplace12 = NewMatrix(3, 3, -1, -2, -1, -2, 12, -2, -1, -2, -1)
)
