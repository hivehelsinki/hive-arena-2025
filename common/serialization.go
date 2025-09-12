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

func (c Coords) MarshalText() (text []byte, err error) {
	return []byte(c.String()), nil
}

func (t Terrain) String() string {
	return []string{"INVALID", "EMPTY", "ROCK", "FIELD"}[t]
}

func (t Terrain) MarshalText() (text []byte, err error) {
	return []byte(t.String()), nil
}

func (d Direction) String() string {
	return []string{"E", "SE", "SW", "W", "NW", "NE"}[d]
}

func (d Direction) MarshalText() (text []byte, err error) {
	return []byte(d.String()), nil
}

func (t EntityType) String() string {
	return []string{"WALL", "HIVE", "BEE"}[t]
}

func (t EntityType) MarshalText() (text []byte, err error) {
	return []byte(t.String()), nil
}
