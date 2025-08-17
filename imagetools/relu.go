package imagetools

import (
	"math"
)

func Relu(a float64) func(x float64) float64 {
	return func(x float64) float64 {
		if x >= 0 {
			return x
		} else {
			return a * x
		}
	}
}

func Sigmoid(x float64) float64 {
	return 1 / (1 + math.Exp(-x))
}

func Softplus(x float64) float64 {
	return math.Log1p(math.Exp(x))
}

func Mish(x float64) float64 {
	y := math.Exp(x) + 1
	return x * (1 - 2/(y*y+1))
}

func Gradient(x0 float64, f func(float64) float64) (x, y float64) {
	x = x0
	y = f(x)
	return
}
