package main

import rl "github.com/gen2brain/raylib-go/raylib"

type EnemyTemplate struct {
	Name       string
	MaxHealth  int
	Texture    rl.Texture2D
	Frame      rl.Rectangle
	FrameCount int
	FrameSpeed float32
	AgroRadius int // number of tiles
}

type Enemy struct {
	GridX, GridY int
	Pos          rl.Vector2
	Health       int
	Template     EnemyTemplate
	Path         []TilePos
	Animation    Animation
}

func NewEnemy(x, y int, tmpl EnemyTemplate) *Enemy {
	return &Enemy{
		GridX:    x,
		GridY:    y,
		Pos:      rl.NewVector2(float32(x*TileSize), float32(y*TileSize)),
		Health:   tmpl.MaxHealth,
		Template: tmpl,
		Animation: Animation{
			FrameCount: tmpl.FrameCount,
			FrameSpeed: tmpl.FrameSpeed,
		},
	}
}

func (e *Enemy) PerformTick(floor *Floor, player *Player) {
	if e.Health <= 0 {
		return
	}

	dx := player.GridX - e.GridX
	dy := player.GridY - e.GridY
	dist := max(abs(dx), abs(dy))

	if dist <= e.Template.AgroRadius {
		// Attack if adjacent
		if abs(dx) <= 1 && abs(dy) <= 1 && !(dx == 0 && dy == 0) {
			player.Health -= 10
			e.Path = nil
			return
		}

		playerPos := TilePos{player.GridX, player.GridY}
		enemyPos := TilePos{e.GridX, e.GridY}

		// Recalculate path if it's empty or outdated
		if len(e.Path) == 0 || !e.Path[len(e.Path)-1].Equals(playerPos) {
			// Get walkable tiles around player
			dirs := []TilePos{{1, 0}, {-1, 0}, {0, 1}, {0, -1}, {1, 1}, {-1, 1}, {1, -1}, {-1, -1}}
			var targets []TilePos
			for _, d := range dirs {
				tx, ty := playerPos.X+d.X, playerPos.Y+d.Y
				if !isSolid(floor, tx, ty, player, floor.Enemies) {
					targets = append(targets, TilePos{tx, ty})
				}
			}

			// Try path to closest available spot
			for _, target := range targets {
				if path := FindPath(floor, enemyPos, target, player, floor.Enemies); path != nil {
					e.Path = path
					break
				}
			}
		}

		// Follow path
		if len(e.Path) > 0 {
			next := e.Path[0]
			if !isSolid(floor, next.X, next.Y, player, floor.Enemies) {
				e.GridX = next.X
				e.GridY = next.Y
				e.Pos = rl.NewVector2(float32(e.GridX*TileSize), float32(e.GridY*TileSize))
				e.Path = e.Path[1:]
			}
		}
	}
}

func (e *Enemy) Update() {
	e.Animation.Update()
}

func (e *Enemy) Draw() {
	if e.Health <= 0 {
		return
	}

	frame := e.Animation.FrameRect(rl.NewRectangle(0, 0, 16, 32)) // source sprite

	src := frame
	dest := rl.NewRectangle(
		e.Pos.X,    // X stays the same
		e.Pos.Y-16, // shift Y up by 16px to align bottom to grid
		16,         // draw width
		32,         // draw height
	)

	origin := rl.NewVector2(0, 0)
	rotation := float32(0)

	rl.DrawTexturePro(e.Template.Texture, src, dest, origin, rotation, rl.White)

	// frame := e.Animation.FrameRect(e.Template.Frame)
	// rl.DrawTextureRec(e.Template.Texture, frame, e.Pos, rl.White)

	barW := TileSize
	barH := 4
	rl.DrawRectangle(int32(e.Pos.X), int32(e.Pos.Y)-6, int32(barW), int32(barH), rl.DarkGray)
	hpPct := float32(e.Health) / float32(e.Template.MaxHealth)
	rl.DrawRectangle(int32(e.Pos.X), int32(e.Pos.Y)-6, int32(hpPct*float32(barW)), int32(barH), rl.Green)

	if showDebug {
		tileRange := e.Template.AgroRadius
		tileX := e.GridX
		tileY := e.GridY

		left := (tileX - tileRange) * TileSize
		top := (tileY - tileRange) * TileSize
		size := (tileRange*2 + 1) * TileSize

		rl.DrawRectangleLines(int32(left), int32(top), int32(size), int32(size), rl.ColorAlpha(rl.Red, 0.6))
	}
}
