package main

import (
	rl "github.com/gen2brain/raylib-go/raylib"
)

const (
	ScreenWidth  = 800
	ScreenHeight = 600
	TileSize     = 32
	MapWidth     = 20
	MapHeight    = 15
)

var (
	world  *World
	player *Player
	camera rl.Camera2D
	ui     *UI
)

func main() {
	rl.InitWindow(ScreenWidth, ScreenHeight, "Tile-Based Game")
	defer rl.CloseWindow()
	rl.SetTargetFPS(60)

	tilesheet := rl.LoadTexture("assets/tiles.png")
	defer rl.UnloadTexture(tilesheet)

	world = NewWorld(tilesheet)
	player = NewPlayer()
	ui = NewUI(10, 10, 200, 60)
	camera = rl.NewCamera2D(
		rl.NewVector2(0, 0), // offset (filled in below)
		rl.NewVector2(0, 0), // target (we’ll update this each frame)
		0.0,                 // rotation
		1.0,                 // zoom
	)
	// Optional: center player on screen initially
	camera.Offset = rl.NewVector2(ScreenWidth/2, ScreenHeight/2)

	for !rl.WindowShouldClose() {
		update()
		draw()
	}
}

func update() {
	player.Update()

	// Center camera on the middle of the player
	camera.Target = rl.NewVector2(
		player.Pos.X+TileSize/2,
		player.Pos.Y+TileSize/2,
	)
}

func draw() {
	rl.BeginDrawing()
	rl.ClearBackground(rl.Black)
	rl.BeginMode2D(camera)

	world.Draw()
	player.Draw()

	rl.EndMode2D()

	ui.Draw()

	rl.EndDrawing()
}
