package common

import (
	"fmt"
	"strconv"
	"strings"
)

func (c Coords) String() string {
	return fmt.Sprintf("%d,%d", c.Row, c.Col)
}

func FromString(s string) (Coords, error) {
	parts := strings.Split(s, ",")
	if len(parts) != 2 {
		return Coords{}, fmt.Errorf("invalid coordinate string format: %s", s)
	}
	row, err := strconv.Atoi(parts[0])
	if err != nil {
		return Coords{}, fmt.Errorf("invalid row value: %w", err)
	}
	col, err := strconv.Atoi(parts[1])
	if err != nil {
		return Coords{}, fmt.Errorf("invalid col value: %w", err)
	}
	return Coords{Row: row, Col: col}, nil
}

func (t Terrain) String() string {
	return []string{"invalid", "empty", "rock", "field"}[t]
}

func (d Direction) String() string {
	return []string{"E", "SE", "SW", "W", "NW", "NE"}[d]
}

func (k SpawnKind) String() string {
	return []string{"hive", "bee"}[k]
}
