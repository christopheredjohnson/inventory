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
	dungeon        *Dungeon
	currentFloor   *Floor
	showDebug      = true
	player         *Player
	enemies        []*Enemy
	enemyTemplates map[string]EnemyTemplate
	camera         rl.Camera2D

	tickTimer    float32 = 0
	tickInterval float32 = 0.6
	inv          *Inventory

	tilesetTexture rl.Texture2D
)

func main() {
	rl.InitWindow(ScreenWidth, ScreenHeight, "Real-Time Combat Game")
	defer rl.CloseWindow()
	rl.SetTargetFPS(60)

	tilesetTexture = rl.LoadTexture("assets/tileset.png")
	defer rl.UnloadTexture(tilesetTexture)

	inv = &Inventory{
		Cols:        5,
		Rows:        4,
		SlotSize:    32,
		SlotPadding: 4,
		Position:    rl.NewVector2(10, 10),
		// ItemTexture: playerTex,
		Slots: make([]ItemSlot, 5*4),
	}

	inv.AddItem(&Item{Name: "Coin", Stackable: true, IconRect: rl.NewRectangle(0, 0, 16, 16)}, 10)

	enemyTemplates = map[string]EnemyTemplate{
		"Orc": {
			Name:       "Orc",
			MaxHealth:  50,
			Texture:    tilesetTexture,
			Frame:      rl.NewRectangle(0, 0, TileSize, TileSize),
			FrameCount: 6,
			FrameSpeed: 0.2,
			AgroRadius: 10,
		},
		"Slime": {
			Name:       "Slime",
			MaxHealth:  25,
			Texture:    tilesetTexture,
			Frame:      rl.NewRectangle(0, 0, TileSize, TileSize),
			FrameCount: 4,
			FrameSpeed: 0.2,
			AgroRadius: 10,
		},
	}

	player = &Player{
		GridX:     2,
		GridY:     2,
		Pos:       rl.NewVector2(2*TileSize, 2*TileSize),
		Health:    100,
		MaxHealth: 100,
		Texture:   tilesetTexture,
		Animation: Animation{
			FrameCount: 8,
			FrameSpeed: 0.2,
		},
	}

	camera = rl.NewCamera2D(
		rl.NewVector2(0, 0),
		rl.NewVector2(float32(player.Pos.X+TileSize/2), float32(player.Pos.Y+TileSize/2)),
		0.0,
		3.0,
	)

	camera.Offset = rl.NewVector2(ScreenWidth/2, ScreenHeight/2)

	dungeon = &Dungeon{Floors: map[int]*Floor{}}
	dungeon.Floors[0] = GenerateRandomFloor(50, 40, 8)

	// fmt.Printf("%v", currentFloor)
	currentFloor = dungeon.Floors[0]

	player.GridX = currentFloor.SpawnPoint.X
	player.GridY = currentFloor.SpawnPoint.Y
	player.Pos = rl.NewVector2(float32(player.GridX*TileSize), float32(player.GridY*TileSize))

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

		if !isSolid(currentFloor, tx, ty, player, enemies) {
			player.Path = FindPath(currentFloor, TilePos{player.GridX, player.GridY}, TilePos{tx, ty}, player, currentFloor.Enemies)
		}
	}

	if rl.IsKeyPressed(rl.KeyF3) {
		showDebug = !showDebug
	}

	tickTimer += rl.GetFrameTime()
	if tickTimer >= tickInterval {
		tickTimer -= tickInterval
		player.PerformTick(currentFloor)
		for _, e := range currentFloor.Enemies {
			e.PerformTick(currentFloor, player)
		}
	}

	// animation frames
	player.Update()

	for _, e := range currentFloor.Enemies {
		e.Update()
	}

	live := currentFloor.Enemies[:0]
	for _, e := range currentFloor.Enemies {
		if e.Health > 0 {
			live = append(live, e)
		}
	}
	currentFloor.Enemies = live

	camera.Target = rl.NewVector2(player.Pos.X+TileSize/2, player.Pos.Y+TileSize/2)

	inv.Update()
}

func draw() {
	rl.BeginDrawing()
	rl.ClearBackground(rl.Magenta)

	rl.BeginMode2D(camera)

	for y := 0; y < currentFloor.Height; y++ {
		for x := 0; x < currentFloor.Width; x++ {
			tile := currentFloor.Tiles[y][x]
			// Optionally color solid tiles differently
			color := rl.DarkGray
			if tile.Solid {
				color = rl.Black
			}
			rl.DrawRectangle(int32(x*TileSize), int32(y*TileSize), TileSize, TileSize, color)
		}
	}

	player.Draw()

	for _, e := range currentFloor.Enemies {
		e.Draw()
	}

	rl.EndMode2D()
	// inv.Draw()
	rl.EndDrawing()
}
