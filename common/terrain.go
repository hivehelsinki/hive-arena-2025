package common

import (
	"os"
	"strconv"
	"strings"
)

//go:generate stringer -type Terrain
type Terrain int

const (
	INVALID Terrain = iota
	EMPTY
	ROCK
	FIELD
)

func (t Terrain) IsWalkable() bool {
	return t == EMPTY || t == FIELD
}

// Doubled coordinates system (https://www.redblobgames.com/grids/hexagons/)
// Pointy tops (horizontal rows)
// Top-left corner is 0,0
// Rows increase by 1, (vertical) columns increase by 2

type Coords struct {
	Row int
	Col int
}

func (c Coords) Neighbour(dir Direction) Coords {
	offset := directionToOffset[dir]
	return Coords{Row: c.Row + offset.Row, Col: c.Col + offset.Col}
}

func (c Coords) Neighbours() []Coords {
	neighbors := make([]Coords, 0, 6)
	for _, offset := range directionToOffset {
		neighbors = append(neighbors, Coords{Row: c.Row + offset.Row, Col: c.Col + offset.Col})
	}
	return neighbors
}

func abs(x int) int {
	return max(x, -x)
}

func (c Coords) Distance(b Coords) int {
	dcol := abs(c.Col - b.Col)
	drow := abs(c.Row - b.Row)
	return drow + max(0, (dcol-drow)/2)
}

//go:generate stringer -type Direction
type Direction int

const (
	E Direction = iota
	SE
	SW
	W
	NW
	NE
)

var directionToOffset = map[Direction]Coords{
	E:  {0, 2},
	NE: {-1, 1},
	NW: {-1, -1},
	W:  {0, -2},
	SW: {1, -1},
	SE: {1, 1},
}

type Spawn struct {
	Kind   EntityType
	Player int
	Coords Coords
}

type MapData struct {
	Map    map[Coords]Terrain
	Spawns []Spawn
}

var charToTerrain = map[rune]Terrain{
	'.': EMPTY,
	'F': FIELD,
	'R': ROCK,
}

var charToSpawn = map[rune]EntityType{
	'H': HIVE,
	'B': BEE,
}

func LoadMap(path string) (MapData, error) {
	content, err := os.ReadFile(path)
	if err != nil {
		return MapData{}, err
	}

	lines := strings.Split(string(content), "\n")
	gameMap := make(map[Coords]Terrain)
	spawns := []Spawn{}

	for trow, line := range lines {
		for tcol, char := range line {
			if terrain, ok := charToTerrain[char]; ok {
				coords := Coords{trow, tcol / 2}
				gameMap[coords] = terrain
			} else if kind, ok := charToSpawn[char]; ok {
				coords := Coords{trow, tcol / 2}
				playerStr := string(line[tcol+1])
				player, _ := strconv.Atoi(playerStr)

				spawns = append(spawns, Spawn{
					Kind:   kind,
					Player: player,
					Coords: coords,
				})
				gameMap[coords] = EMPTY
			}
		}
	}

	return MapData{
		Map:    gameMap,
		Spawns: spawns,
	}, nil
}
