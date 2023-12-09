package utils

import (
	"math/rand"
)

type MCF struct {
	Function    func(float64) float64 //被积函数
	AreaLeft    float64               //左边界
	AreaRight   float64               //右边界
	Integration float64               //积分结果
	Time        int                   //随机次数，次数越大精度越高，性能开销越大
}

// 利用大数统计的思想近似求解函数最大值
func (mc *MCF) getMax() (maxNum float64) {
	var time = 10000 //默认次数10000次
	maxNum = mc.Function(mc.AreaLeft)
	d := mc.AreaRight - mc.AreaLeft
	if mc.Time > 0 {
		time = mc.Time
	}
	for i := 0; i < time; i++ {
		x := mc.AreaLeft + d*rand.Float64()
		if mc.Function(x) > maxNum {
			maxNum = mc.Function(x)
		}
	}
	return maxNum
}

// 利用大数统计的思想近似求解函数最小值
func (mc *MCF) getMin() (minNum float64) {
	var time = 10000 //默认次数10000次
	minNum = mc.Function(mc.AreaLeft)
	d := mc.AreaRight - mc.AreaLeft
	if mc.Time > 0 {
		time = mc.Time
	}
	for i := 0; i < time; i++ {
		x := mc.AreaLeft + d*rand.Float64()
		if mc.Function(x) < minNum {
			minNum = mc.Function(x)
		}
	}
	return minNum
}

// DeIntCalc 利用蒙特卡罗方法，通过面积比例关系近似求解函数定积分
func (mc *MCF) DeIntCalc() float64 {
	var time = 10000
	var num float64 = 0
	d1 := mc.AreaRight - mc.AreaLeft
	d2 := mc.getMax() - mc.getMin()
	if mc.Time > 0 {
		time = mc.Time
	}
	for i := 0; i < time; i++ {
		x := mc.AreaLeft + d1*rand.Float64()
		y := mc.getMin() + d2*rand.Float64()
		if mc.Function(x) > 0 {
			if y < mc.Function(x) && y > 0 {
				num++
			}
		} else {
			if y > mc.Function(x) && y < 0 {
				num--
			}
		}
	}
	s := d1 * d2
	result := s * (num / float64(time))
	mc.Integration = result
	return result
}

//	func main(){
//		fmt.Println(DeIntCalc(function,0,1))//0.332323442141
//	}
//
// 代替Rand的拟蒙特卡洛，Van der Corput
func VanDerCorput(x int, base int) (r float64) {
	bk := 1 / base
	for x > 0 {
		r += float64((x % base) * bk)
		x /= base
		bk /= base
	}
	return r
}
