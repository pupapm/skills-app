package main

import "math"

func cap01(x float64) float64 {
	if x < 0 {
		return 0
	}
	if x > 1 {
		return 1
	}
	return x
}

func capRange(x, lo, hi float64) float64 {
	if x < lo {
		return lo
	}
	if x > hi {
		return hi
	}
	return x
}

func safeDiv(a, b float64) float64 {
	if b <= 0 {
		return 0
	}
	return a / b
}

func maxInt(a, b int) int {
	if a > b {
		return a
	}
	return b
}

func round2(x float64) float64 {
	return math.Round(x*100) / 100
}
