package main

import rl "github.com/gen2brain/raylib-go/raylib"

type Animation struct {
	FrameX, FrameY int // current frame position in tile grid
	FrameIndex     int
	FrameTimer     float32
	FrameCount     int
	FrameSpeed     float32
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

func (a *Animation) FrameRectFromGrid(tileX, tileY, tileW, tileH int) rl.Rectangle {
	return rl.NewRectangle(
		float32((tileX+a.FrameIndex)*tileW),
		float32(tileY*tileH),
		float32(tileW),
		float32(tileH),
	)
}
