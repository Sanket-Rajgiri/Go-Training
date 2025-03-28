package main

import (
	"fmt"
)

func Sqrt(x int) float64 {
	res, low := 1, 1
	hi := x / 2
	for low <= hi {
		mid := (low + (hi-low)/2)
		if mid*mid <= x {
			res = mid
			low = mid + 1
		} else {
			hi = mid - 1
		}
	}
	return float64(res)
}

func main() {
	fmt.Println(Sqrt(49))
}
