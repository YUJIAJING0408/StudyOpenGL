package utils

import (
	"fmt"
	"math"
	"testing"
)

func TestMC(t *testing.T) {
	// mc测试
	var mc = &MCF{
		Function: func(x float64) float64 {
			return math.Sqrt(1-math.Pow(x, 2)) * 4
		},
		AreaLeft:    0.0,
		AreaRight:   1.0,
		Integration: 0,
		Time:        20000,
	}
	mc.DeIntCalc()
	fmt.Printf("%0.5f\n", mc.Integration)
}
