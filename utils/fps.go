package utils

type Fps struct {
	lastTime  float64
	countTime float64
	count     uint
}

// NewFps 构建一个Fps，传入当前GLFW时间
func NewFps(lastTime float64) *Fps {
	return &Fps{
		lastTime:  lastTime,
		countTime: 0.0,
		count:     0,
	}
}

//func (f *Fps) Get(currentTime float64) (ok bool, fps uint) {
//	//currentTime := glfw.GetTime()
//	deltaTime := currentTime - f.lastTime
//	f.countTime += deltaTime
//	f.count++
//	//println(int(f.countTime))
//	if f.countTime >= 1.0 {
//		f.countTime = 0.0
//		out := f.count
//		f.count = 0
//		f.lastTime = currentTime
//		return true, out
//	} else {
//		return false, 0
//	}
//
//}

// Get FPS计数器
func (f *Fps) Get(currentTime float64) (ok bool, fps int) {
	//currentTime = glfw.GetTime()
	deltaTime := currentTime - f.lastTime
	f.countTime += deltaTime
	if f.countTime >= 1.0 {
		// 一秒内的计数次数除时间
		fps = int(float64(f.count) / (currentTime - f.lastTime))
		// 归零
		f.lastTime = currentTime
		f.countTime = 0.0
		f.count = 0
		return true, fps
	} else {
		f.count++
		return false, 0
	}
}
