package main

import rl "github.com/gen2brain/raylib-go/raylib"

type Player struct {
	GridX, GridY int
	Pos          rl.Vector2
	Health       int
	MaxHealth    int
	Path         []TilePos
	Texture      rl.Texture2D
	Animation    Animation
}

func (p *Player) PerformTick(floor *Floor) {
	// Attack adjacent enemy
	for _, e := range floor.Enemies {
		if abs(e.GridX-p.GridX)+abs(e.GridY-p.GridY) == 1 {
			e.Health -= 20
			return
		}
	}

	// Move to next tile if path is valid
	if len(p.Path) > 0 {
		next := p.Path[0]
		if !isSolid(floor, next.X, next.Y, p, floor.Enemies) {
			p.GridX = next.X
			p.GridY = next.Y
			p.Pos = rl.NewVector2(float32(p.GridX*TileSize), float32(p.GridY*TileSize))
			p.Path = p.Path[1:]
		}
	}
}

func (p *Player) Update() {
	p.Animation.Update()
}

func (p *Player) Draw() {

	frame := p.Animation.FrameRectFromGrid(8, 4, 16, 32) // FrameX shifts column

	src := frame
	dest := rl.NewRectangle(
		p.Pos.X,    // X stays the same
		p.Pos.Y-16, // shift Y up by 16px to align bottom to grid
		16,         // draw width
		32,         // draw height
	)

	origin := rl.NewVector2(0, 0)
	rotation := float32(0)

	rl.DrawTexturePro(p.Texture, src, dest, origin, rotation, rl.White)

	rl.DrawRectangle(int32(p.Pos.X), int32(p.Pos.Y)-6, TileSize, 4, rl.DarkGray)
	hpPct := float32(p.Health) / float32(p.MaxHealth)
	rl.DrawRectangle(int32(p.Pos.X), int32(p.Pos.Y)-6, int32(hpPct*TileSize), 4, rl.Red)

	for _, pos := range player.Path {
		rl.DrawRectangleLines(int32(pos.X*TileSize), int32(pos.Y*TileSize), TileSize, TileSize, rl.Yellow)
	}
}
