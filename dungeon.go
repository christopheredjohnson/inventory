package main

import "math/rand"

type Dungeon struct {
	Floors       map[int]*Floor
	CurrentFloor int
}

type Floor struct {
	Width, Height int
	Tiles         [][]Tile
	Enemies       []*Enemy
	Items         []*Item
	SpawnPoint    TilePos
}

func NewEmptyFloor(width, height int) *Floor {
	tiles := make([][]Tile, height)
	for y := range tiles {
		tiles[y] = make([]Tile, width)
		for x := range tiles[y] {
			tiles[y][x] = Tile{Solid: true}
		}
	}
	return &Floor{
		Width:  width,
		Height: height,
		Tiles:  tiles,
	}
}

type Room struct {
	X, Y, W, H int
}

func (r Room) Center() TilePos {
	return TilePos{r.X + r.W/2, r.Y + r.H/2}
}

func (r Room) Intersects(other Room) bool {
	return r.X <= other.X+other.W && r.X+r.W >= other.X &&
		r.Y <= other.Y+other.H && r.Y+r.H >= other.Y
}

func GenerateRandomFloor(width, height, maxRooms int) *Floor {
	floor := NewEmptyFloor(width, height)

	rooms := []Room{}

	for i := 0; i < maxRooms; i++ {
		w := rand.Intn(6) + 4 // room width: 4–9
		h := rand.Intn(6) + 4 // room height: 4–9
		x := rand.Intn(width - w - 1)
		y := rand.Intn(height - h - 1)

		newRoom := Room{X: x, Y: y, W: w, H: h}

		overlap := false
		for _, other := range rooms {
			if newRoom.Intersects(other) {
				overlap = true
				break
			}
		}
		if overlap {
			continue
		}

		createRoom(floor, newRoom)

		if len(rooms) > 0 {
			prevCenter := rooms[len(rooms)-1].Center()
			newCenter := newRoom.Center()

			// Randomly connect horizontal then vertical, or vice versa
			if rand.Intn(2) == 0 {
				createHTunnel(floor, prevCenter.X, newCenter.X, prevCenter.Y)
				createVTunnel(floor, prevCenter.Y, newCenter.Y, newCenter.X)
			} else {
				createVTunnel(floor, prevCenter.Y, newCenter.Y, prevCenter.X)
				createHTunnel(floor, prevCenter.X, newCenter.X, newCenter.Y)
			}
		} else {
			// First room = spawn point
			floor.SpawnPoint = newRoom.Center()
		}

		rooms = append(rooms, newRoom)
	}

	// Optional: Place enemies in rooms
	for _, room := range rooms[1:] {
		c := room.Center()
		floor.Enemies = append(floor.Enemies, NewEnemy(c.X, c.Y, enemyTemplates["Orc"]))
	}

	return floor
}

func createRoom(floor *Floor, r Room) {
	for y := r.Y; y < r.Y+r.H; y++ {
		for x := r.X; x < r.X+r.W; x++ {
			floor.Tiles[y][x].Solid = false
		}
	}
}

func createHTunnel(floor *Floor, x1, x2, y int) {
	if x1 > x2 {
		x1, x2 = x2, x1
	}
	for x := x1; x <= x2; x++ {
		floor.Tiles[y][x].Solid = false
	}
}

func createVTunnel(floor *Floor, y1, y2, x int) {
	if y1 > y2 {
		y1, y2 = y2, y1
	}
	for y := y1; y <= y2; y++ {
		floor.Tiles[y][x].Solid = false
	}
}
