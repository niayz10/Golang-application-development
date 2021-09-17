package main

import (
	"fmt"
)

func Sqrt(x float64) float64 {
	z :=1.0
	for i:=0 ; i<10 ; i++ {
		z -= (z*z - x) / (2*z)
		fmt.Println(z)
	}
	return z
}

func main() {
	for x := 1.0 ; x < 4 ; x++ {
		fmt.Println("The answers for value of", x)
		fmt.Println(Sqrt(x))
		}
}
