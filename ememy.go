package main

import (
	rl "github.com/gen2brain/raylib-go/raylib"
)

type EnemyTemplate struct {
	Name       string
	MaxHealth  int
	Texture    rl.Texture2D
	Frame      rl.Rectangle // Initial frame (x, y, width, height)
	FrameCount int          // Total frames in animation (horizontally)
	FrameSpeed float32      // Seconds per frame
}

type Enemy struct {
	Name         string
	GridX, GridY int
	Pos          rl.Vector2
	Health       int
	MaxHealth    int

	Texture      rl.Texture2D
	Frame        rl.Rectangle
	FrameCount   int
	FrameSpeed   float32
	AnimTimer    float32
	CurrentFrame int
}

func NewEnemy(x, y int, template EnemyTemplate) *Enemy {
	return &Enemy{
		Name:         template.Name,
		GridX:        x,
		GridY:        y,
		Pos:          rl.NewVector2(float32(x*TileSize), float32(y*TileSize)),
		Health:       template.MaxHealth,
		MaxHealth:    template.MaxHealth,
		Texture:      template.Texture,
		Frame:        template.Frame,
		FrameCount:   template.FrameCount,
		FrameSpeed:   template.FrameSpeed,
		AnimTimer:    0,
		CurrentFrame: 0,
	}
}

func (e *Enemy) Draw() {

	// Health bar
	barWidth := TileSize
	barHeight := 6
	healthPct := float32(e.Health) / float32(e.MaxHealth)
	rl.DrawRectangle(int32(e.Pos.X), int32(e.Pos.Y)-10, int32(barWidth), int32(barHeight), rl.DarkGray)
	rl.DrawRectangle(int32(e.Pos.X), int32(e.Pos.Y)-10, int32(float32(barWidth)*healthPct), int32(barHeight), rl.Green)
	rl.DrawTextureRec(e.Texture, e.Frame, e.Pos, rl.White)
}

func (e *Enemy) Update() {
	e.AnimTimer += rl.GetFrameTime()
	if e.AnimTimer >= e.FrameSpeed {
		e.AnimTimer = 0
		e.CurrentFrame = (e.CurrentFrame + 1) % e.FrameCount
		e.Frame.X = float32(e.CurrentFrame) * e.Frame.Width
	}
}
