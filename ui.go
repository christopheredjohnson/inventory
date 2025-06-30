package main

import (
	"fmt"

	rl "github.com/gen2brain/raylib-go/raylib"
)

type UI struct {
	X, Y   int
	Width  int
	Height int
}

func NewUI(x, y, width, height int) *UI {
	return &UI{
		X:      x,
		Y:      y,
		Width:  width,
		Height: height,
	}
}

func (ui *UI) Draw() {
	// Background box
	rl.DrawRectangle(int32(ui.X), int32(ui.Y), int32(ui.Width), int32(ui.Height), rl.Fade(rl.Black, 0.6))
	rl.DrawRectangleLines(int32(ui.X), int32(ui.Y), int32(ui.Width), int32(ui.Height), rl.White)

	DrawHealthBar(ui.X+8, ui.Y+40, 180, 16, player.Health, player.MaxHealth)
}

func DrawHealthBar(x, y, width, height int, current, max int) {
	// Border
	rl.DrawRectangleLines(int32(x), int32(y), int32(width), int32(height), rl.White)

	// Background
	rl.DrawRectangle(int32(x), int32(y), int32(width), int32(height), rl.DarkGray)

	// Filled bar
	healthPercent := float32(current) / float32(max)
	barWidth := int(float32(width) * healthPercent)
	rl.DrawRectangle(int32(x), int32(y), int32(barWidth), int32(height), rl.Red)

	// Text (e.g. "75 / 100")
	text := fmt.Sprintf("%d / %d", current, max)
	textX := int32(x) + (int32(width)-rl.MeasureText(text, 14))/2
	textY := y + (height-14)/2
	rl.DrawText(text, int32(textX), int32(textY), 14, rl.White)
}
