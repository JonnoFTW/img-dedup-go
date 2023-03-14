package dct

import (
	"log"
	"math"
)

var (
	cosCache = make([]float64, 0)
)

func getCacheIdx(n int, N int, k int) int {
	return n*N + k
}
func getCachedCos(n int, N int, k int) float64 {
	// Cache is only valid for a single value of N
	// Since this program uses a single hash type per run, there should only be 1 value of N
	// so the cache should only be generated once
	Np := N + 1
	cacheSize := Np * Np
	if len(cosCache) != cacheSize {
		cosCache = make([]float64, cacheSize)
		for _n := 0; _n < N; _n++ {
			for _k := 0; _k < N; _k++ {
				cosCache[getCacheIdx(_n, N, _k)] = math.Cos(((math.Pi * float64(_k)) * float64(2*_n+1)) / float64(2*N))
			}
		}
	}
	return cosCache[getCacheIdx(n, N, k)]
}
func dct2(val float64, n int, N int, k int) float64 {
	return val * getCachedCos(n, N, k)
}

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
					sum += dct2(vals[n][colIdx], n, N, k)
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
					sum += dct2(x[n], n, N, k)
				}
				out[rowIdx][k] = 2. * sum
			}
		}
	} else {
		log.Fatalf("Invalid axis specified: %d. Must be 0 or 1", axis)
	}

	return out
}
