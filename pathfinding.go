package main

type Node struct {
	Pos  TilePos
	G, H int
	From *Node
}

func FindPath(start, goal TilePos) []TilePos {

	open := []*Node{}
	closed := map[TilePos]bool{}
	startNode := &Node{Pos: start}
	open = append(open, startNode)

	var getLowestF = func(nodes []*Node) *Node {
		lowest := nodes[0]
		for _, n := range nodes {
			if (n.G + n.H) < (lowest.G + lowest.H) {
				lowest = n
			}
		}
		return lowest
	}

	var neighbors = func(pos TilePos) []struct {
		Pos  TilePos
		Cost int
	} {
		dirs := []struct {
			X, Y int
			Cost int
		}{
			{1, 0, 10}, {-1, 0, 10}, {0, 1, 10}, {0, -1, 10},
			{1, 1, 14}, {1, -1, 14}, {-1, 1, 14}, {-1, -1, 14},
		}
		var out []struct {
			Pos  TilePos
			Cost int
		}
		for _, d := range dirs {
			nx, ny := pos.X+d.X, pos.Y+d.Y
			if !isSolid(nx, ny) {
				out = append(out, struct {
					Pos  TilePos
					Cost int
				}{TilePos{nx, ny}, d.Cost})
			}
		}
		return out
	}

	// Chebyshev distance (since we allow diagonal moves)
	heuristic := func(a, b TilePos) int {
		dx := abs(a.X - b.X)
		dy := abs(a.Y - b.Y)
		return 10 * max(dx, dy)
	}

	for len(open) > 0 {
		current := getLowestF(open)
		if current.Pos == goal {
			path := []TilePos{}
			for n := current; n != nil; n = n.From {
				path = append([]TilePos{n.Pos}, path...)
			}
			return path[1:] // skip starting tile
		}

		open = removeNode(open, current)
		closed[current.Pos] = true

		for _, n := range neighbors(current.Pos) {
			if closed[n.Pos] {
				continue
			}

			cost := current.G + n.Cost
			found := false
			for _, node := range open {
				if node.Pos == n.Pos {
					if cost < node.G {
						node.G = cost
						node.From = current
					}
					found = true
					break
				}
			}
			if !found {
				open = append(open, &Node{
					Pos:  n.Pos,
					G:    cost,
					H:    heuristic(n.Pos, goal),
					From: current,
				})
			}
		}
	}
	return nil
}

func removeNode(slice []*Node, target *Node) []*Node {
	out := []*Node{}
	for _, n := range slice {
		if n != target {
			out = append(out, n)
		}
	}
	return out
}

func isSolid(x, y int) bool {
	if x < 0 || y < 0 || x >= MapWidth || y >= MapHeight {
		return true
	}
	if worldTiles[y][x].Solid {
		return true
	}
	for _, e := range enemies {
		if e.GridX == x && e.GridY == y && e.Health > 0 {
			return true
		}
	}
	// BLOCK PLAYER TILE
	if player.GridX == x && player.GridY == y {
		return true
	}
	return false
}
