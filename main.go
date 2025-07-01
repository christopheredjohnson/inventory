// main.go
package main

import (
	rl "github.com/gen2brain/raylib-go/raylib"
)

const (
	ScreenWidth  = 800
	ScreenHeight = 600
	TileSize     = 16
	MapWidth     = 20
	MapHeight    = 15
)

type Tile struct {
	Solid bool
}

type TilePos struct {
	X, Y int
}

func (a TilePos) Equals(b TilePos) bool {
	return a.X == b.X && a.Y == b.Y
}

var (
	showDebug      = true
	worldTiles     [][]Tile
	player         *Player
	enemies        []*Enemy
	enemyTemplates map[string]EnemyTemplate
	camera         rl.Camera2D

	tickTimer    float32 = 0
	tickInterval float32 = 0.6
	inv          *Inventory
)

func main() {
	rl.InitWindow(ScreenWidth, ScreenHeight, "Real-Time Combat Game")
	defer rl.CloseWindow()
	rl.SetTargetFPS(60)

	playerTex := rl.LoadTexture("assets/Characters/Champions/Borg.png")
	defer rl.UnloadTexture(playerTex)

	orcTex := rl.LoadTexture("assets/Characters/Monsters/Orcs/Orc.png")
	defer rl.UnloadTexture(orcTex)

	iconTex := rl.LoadTexture("assets/User Interface/Icons-Essentials.png")
	defer rl.UnloadTexture(iconTex)
	inv = &Inventory{
		Cols:        5,
		Rows:        4,
		SlotSize:    32,
		SlotPadding: 4,
		Position:    rl.NewVector2(10, 10),
		ItemTexture: iconTex,
		Slots:       make([]ItemSlot, 5*4),
	}

	inv.AddItem(&Item{Name: "Coin", Stackable: true, IconRect: rl.NewRectangle(0, 0, 16, 16)}, 10)

	enemyTemplates = map[string]EnemyTemplate{
		"Orc": {
			Name:       "Orc",
			MaxHealth:  50,
			Texture:    orcTex,
			Frame:      rl.NewRectangle(0, 0, TileSize, TileSize),
			FrameCount: 4,
			FrameSpeed: 0.2,
			AgroRadius: 2,
		},
	}

	player = &Player{
		GridX:     2,
		GridY:     2,
		Pos:       rl.NewVector2(2*TileSize, 2*TileSize),
		Health:    100,
		MaxHealth: 100,
		Texture:   playerTex,
	}

	enemies = []*Enemy{
		NewEnemy(10, 10, enemyTemplates["Orc"]),
		NewEnemy(8, 8, enemyTemplates["Orc"]),
	}

	camera = rl.NewCamera2D(
		rl.NewVector2(0, 0),
		rl.NewVector2(float32(player.Pos.X+TileSize/2), float32(player.Pos.Y+TileSize/2)),
		0.0,
		3.0,
	)

	camera.Offset = rl.NewVector2(ScreenWidth/2, ScreenHeight/2)

	// Simple empty world
	worldTiles = make([][]Tile, MapHeight)
	for y := range worldTiles {
		worldTiles[y] = make([]Tile, MapWidth)
		for x := range worldTiles[y] {
			worldTiles[y][x] = Tile{Solid: false}
		}
	}

	for !rl.WindowShouldClose() {
		update()
		draw()
	}
}

func update() {

	mouseWorld := rl.GetScreenToWorld2D(rl.GetMousePosition(), camera)
	scroll := rl.GetMouseWheelMove()
	if scroll != 0 {
		// oldZoom := camera.Zoom
		camera.Zoom += scroll * 0.1

		// Clamp zoom
		if camera.Zoom < 0.5 {
			camera.Zoom = 0.5
		}
		if camera.Zoom > 3.0 {
			camera.Zoom = 3.0
		}

		// Adjust camera target to maintain mouse position
		diff := rl.Vector2Subtract(mouseWorld, rl.GetScreenToWorld2D(rl.GetMousePosition(), camera))
		camera.Target = rl.Vector2Add(camera.Target, diff)
	}

	if rl.IsMouseButtonPressed(rl.MouseLeftButton) {
		mouse := rl.GetMousePosition()
		worldPos := rl.GetScreenToWorld2D(mouse, camera)
		tx := int(worldPos.X) / TileSize
		ty := int(worldPos.Y) / TileSize

		if !isSolid(tx, ty) {
			player.Path = FindPath(TilePos{player.GridX, player.GridY}, TilePos{tx, ty})
		}
	}

	if rl.IsKeyPressed(rl.KeyF3) {
		showDebug = !showDebug
	}

	tickTimer += rl.GetFrameTime()
	if tickTimer >= tickInterval {
		tickTimer -= tickInterval
		player.PerformTick()
		for _, e := range enemies {
			e.PerformTick()
		}
	}

	live := enemies[:0]
	for _, e := range enemies {
		if e.Health > 0 {
			live = append(live, e)
		}
	}
	enemies = live

	camera.Target = rl.NewVector2(player.Pos.X+TileSize/2, player.Pos.Y+TileSize/2)

	inv.Update()
}

func draw() {
	rl.BeginDrawing()
	rl.ClearBackground(rl.Black)

	rl.BeginMode2D(camera)

	// Draw grid
	// for y := 0; y < MapHeight; y++ {
	// 	for x := 0; x < MapWidth; x++ {
	// 		rl.DrawRectangleLines(int32(x*TileSize), int32(y*TileSize), TileSize, TileSize, rl.DarkGray)
	// 	}
	// }
	player.Draw()

	for _, e := range enemies {
		e.Draw()
	}

	rl.EndMode2D()
	inv.Draw()
	rl.EndDrawing()
}
