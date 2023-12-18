package utils

import (
	"fmt"
	"github.com/chewxy/math32"
	"testing"
)

func TestMC(t *testing.T) {
	// mc测试
	var mc = &MCF{
		Function: func(x float32) float32 {
			return math32.Sqrt(1-math32.Pow(x, 2)) * 4
		},
		AreaLeft:    0.0,
		AreaRight:   1.0,
		Integration: 0,
		Time:        10000,
	}
	mc.DeIntCalc()
	fmt.Printf("%0.5f\n", mc.Integration)
}
