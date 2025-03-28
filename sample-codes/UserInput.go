package main

import (
	"fmt"
	"math"
)

func main() {
	var a, b float64

	fmt.Scanf("%f %f", &a, &b)
	c := math.Abs((a * b) - (a + b))
	fmt.Print(c)
}
