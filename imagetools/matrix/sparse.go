package imagetools

import types "imagetools/types"

type SparseMatrix[T types.Number] struct {
	x, y int          //row and column lengths
	val  map[Index2]T //values, y-by-x matrix
}
