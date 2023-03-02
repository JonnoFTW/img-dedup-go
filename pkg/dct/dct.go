package dct

import (
	"log"
	"math"
)

func Dct2(vals [][]float64, axis int) [][]float64 {
	// https://en.wikipedia.org/wiki/Discrete_cosine_transform#DCT-II
	// https://docs.scipy.org/doc/scipy/reference/generated/scipy.fftpack.dct.html
	out := make([][]float64, len(vals))
	for i := 0; i < len(out); i++ {
		out[i] = make([]float64, len(vals[i]))
	}
	if axis == 0 {
		// col-wise
		N := len(vals) // no of rows
		for colIdx := 0; colIdx < len(vals[0]); colIdx++ {
			for k := 0; k < N; k++ {
				sum := 0.0
				for n := 0; n < N; n++ { // go down the row
					sum += vals[n][colIdx] * math.Cos(((math.Pi*float64(k))*float64(2*n+1))/float64(2*N))
				}
				out[k][colIdx] = 2. * sum
			}
		}
	} else if axis == 1 {
		// row-wise
		for rowIdx, x := range vals {
			N := len(x) // no of cols
			for k := 0; k < N; k++ {
				sum := 0.0
				for n := 0; n < N; n++ {
					sum += x[n] * math.Cos(((math.Pi*float64(k))*float64(2*n+1))/float64(2*N))
				}
				out[rowIdx][k] = 2. * sum
			}
		}
	} else {
		log.Fatalf("Invalid axis specified: %d. Must be 0 or 1", axis)
	}

	return out
}
