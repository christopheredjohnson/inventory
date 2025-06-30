package main

import rl "github.com/gen2brain/raylib-go/raylib"

type Tile struct {
	Frame rl.Rectangle // Position in the tilesheet
	Solid bool         // Whether the tile is walkable
}

type World struct {
	Tiles     [][]Tile
	Tilesheet rl.Texture2D
}

func NewWorld(tilesheet rl.Texture2D) *World {
	world := &World{
		Tiles:     make([][]Tile, MapHeight),
		Tilesheet: tilesheet,
	}

	for y := 0; y < MapHeight; y++ {
		world.Tiles[y] = make([]Tile, MapWidth)
		for x := 0; x < MapWidth; x++ {
			world.Tiles[y][x] = getTileByType("grass")
		}
	}

	return world
}

func (w *World) Draw() {
	for y := 0; y < MapHeight; y++ {
		for x := 0; x < MapWidth; x++ {
			tile := w.Tiles[y][x]
			dst := rl.NewVector2(float32(x*TileSize), float32(y*TileSize))
			rl.DrawTextureRec(w.Tilesheet, tile.Frame, dst, rl.White)

			if tile.Solid {
				rl.DrawRectangleLines(int32(dst.X), int32(dst.Y), TileSize, TileSize, rl.DarkGray)
			}
		}
	}
}

func (w *World) IsSolid(x, y int) bool {
	if x < 0 || y < 0 || x >= MapWidth || y >= MapHeight {
		return true
	}
	if w.Tiles[y][x].Solid {
		return true
	}
	for _, e := range enemies {
		if e.GridX == x && e.GridY == y {
			return true
		}
	}
	return false
}

func getTileByType(tileType string) Tile {
	switch tileType {
	default:
		return Tile{
			Frame: rl.NewRectangle(32, 0, TileSize, TileSize),
			Solid: false,
		}
	}
}
