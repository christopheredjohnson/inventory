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
			if x == 0 || x == MapWidth-1 || y == 0 || y == MapHeight-1 {
				world.Tiles[y][x] = getTileByType("wall")
			} else {
				world.Tiles[y][x] = getTileByType("floor")
			}
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
	return w.Tiles[y][x].Solid
}

func getTileByType(tileType string) Tile {
	switch tileType {
	case "wall":
		return Tile{
			Frame: rl.NewRectangle(32, 160, TileSize, TileSize), // wall tile in tilesheet
			Solid: true,
		}
	case "floor":
		fallthrough
	default:
		return Tile{
			Frame: rl.NewRectangle(0, 128, TileSize, TileSize), // grass/floor tile
			Solid: false,
		}
	}
}
