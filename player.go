package main

import rl "github.com/gen2brain/raylib-go/raylib"

type Player struct {
	Pos       rl.Vector2
	GridX     int
	GridY     int
	TargetX   int
	TargetY   int
	Moving    bool
	MoveSpeed float32

	DestX   int
	DestY   int
	HasDest bool

	Health    int
	MaxHealth int
}

func NewPlayer() *Player {
	startX, startY := 2, 2
	return &Player{
		GridX:     startX,
		GridY:     startY,
		TargetX:   startX,
		TargetY:   startY,
		Pos:       rl.NewVector2(float32(startX*TileSize), float32(startY*TileSize)),
		MoveSpeed: 2.0,
		Health:    100,
		MaxHealth: 100,
	}
}

func (p *Player) Update() {
	// Smooth movement
	if p.Moving {
		targetPos := rl.NewVector2(float32(p.TargetX*TileSize), float32(p.TargetY*TileSize))
		direction := rl.Vector2Subtract(targetPos, p.Pos)
		if rl.Vector2Length(direction) < p.MoveSpeed {
			p.Pos = targetPos
			p.Moving = false
			p.GridX = p.TargetX
			p.GridY = p.TargetY
		} else {
			moveStep := rl.Vector2Scale(rl.Vector2Normalize(direction), p.MoveSpeed)
			p.Pos = rl.Vector2Add(p.Pos, moveStep)
		}
		return
	}

	// Handle mouse click to set destination
	if rl.IsMouseButtonPressed(rl.MouseLeftButton) {
		mouse := rl.GetMousePosition()
		worldPos := rl.GetScreenToWorld2D(mouse, camera)
		tx := int(worldPos.X) / TileSize
		ty := int(worldPos.Y) / TileSize

		if tx >= 0 && ty >= 0 && tx < MapWidth && ty < MapHeight && !world.IsSolid(tx, ty) {
			p.DestX = tx
			p.DestY = ty
			p.HasDest = true
		}
	}

	// Move one step toward destination
	if p.HasDest && (p.GridX != p.DestX || p.GridY != p.DestY) {
		nextX, nextY := p.GridX, p.GridY
		if p.DestX > p.GridX {
			nextX++
		} else if p.DestX < p.GridX {
			nextX--
		} else if p.DestY > p.GridY {
			nextY++
		} else if p.DestY < p.GridY {
			nextY--
		}

		if !world.IsSolid(nextX, nextY) {
			p.TargetX = nextX
			p.TargetY = nextY
			p.Moving = true
		} else {
			p.HasDest = false
		}
	} else {
		p.HasDest = false
	}
}

func (p *Player) Draw() {
	rl.DrawRectangleV(p.Pos, rl.NewVector2(TileSize, TileSize), rl.Red)

	// Optional: show destination box
	if p.HasDest {
		rl.DrawRectangleLines(int32(p.DestX*TileSize), int32(p.DestY*TileSize), TileSize, TileSize, rl.Green)
	}
}
