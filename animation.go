package main

import rl "github.com/gen2brain/raylib-go/raylib"

type Animation struct {
	FrameIndex int
	FrameTimer float32
	FrameCount int
	FrameSpeed float32
}

func (a *Animation) Update() {
	a.FrameTimer += rl.GetFrameTime()

	if a.FrameTimer >= a.FrameSpeed {
		a.FrameTimer = 0
		a.FrameIndex = (a.FrameIndex + 1) % a.FrameCount
	}
}


func (a *Animation) FrameRect(base rl.Rectangle) rl.Rectangle {
	base.X = float32(a.FrameIndex) * base.Width
	return base
}
