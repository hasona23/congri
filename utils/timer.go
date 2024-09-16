package utils

import "github.com/hajimehoshi/ebiten/v2"

type Timer struct {
	Time         float32
	current_time float32
}

// time in seconds
func NewTimer(time float32) Timer {
	return Timer{Time: time, current_time: 0}
}

func (t Timer) GetCurrentTime() float32 {
	return t.current_time
}
func (timer *Timer) Ticked() bool {
	if timer.current_time >= timer.Time {
		timer.Reset()
		return true
	}
	return false
}
func (timer *Timer) Reset() {
	timer.current_time = 0
}
func (timer *Timer) UpdateTimer() {

	timer.current_time += 1 / float32(ebiten.TPS())

}
