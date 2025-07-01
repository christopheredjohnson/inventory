package main

import rl "github.com/gen2brain/raylib-go/raylib"

type Player struct {
	GridX, GridY int
	Pos          rl.Vector2
	Health       int
	MaxHealth    int
	Path         []TilePos
	Texture      rl.Texture2D

	FrameIndex int
	FrameTimer float32
	FrameCount int
	FrameSpeed float32
}

func (p *Player) PerformTick() {
	for _, e := range enemies {
		if abs(e.GridX-p.GridX)+abs(e.GridY-p.GridY) == 1 {
			e.Health -= 20
			return
		}
	}

	if len(p.Path) > 0 {
		next := p.Path[0]
		if !isSolid(next.X, next.Y) {
			p.GridX = next.X
			p.GridY = next.Y
			p.Pos = rl.NewVector2(float32(p.GridX*TileSize), float32(p.GridY*TileSize))
			p.Path = p.Path[1:]
		}
		return
	}
}

func (p *Player) Update() {
	p.FrameTimer += rl.GetFrameTime()
	if p.FrameTimer >= p.FrameSpeed {
		p.FrameTimer = 0
		p.FrameIndex = (p.FrameIndex + 1) % p.FrameCount
	}
}

func (p *Player) Draw() {

	src := rl.NewRectangle(float32(p.FrameIndex*TileSize), 0, TileSize, TileSize)
	rl.DrawTextureRec(p.Texture, src, p.Pos, rl.White)

	rl.DrawRectangle(int32(p.Pos.X), int32(p.Pos.Y)-6, TileSize, 4, rl.DarkGray)
	hpPct := float32(p.Health) / float32(p.MaxHealth)
	rl.DrawRectangle(int32(p.Pos.X), int32(p.Pos.Y)-6, int32(hpPct*TileSize), 4, rl.Red)

	for _, pos := range player.Path {
		rl.DrawRectangleLines(int32(pos.X*TileSize), int32(pos.Y*TileSize), TileSize, TileSize, rl.Yellow)
	}
}
