package common

import (
	"fmt"
	"slices"
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

func (c *Coords) UnmarshalText(b []byte) error {
	str := string(b)
	coords, err := FromString(str)
	if err != nil {
		return err
	}

	*c = coords
	return nil
}

func (t Terrain) String() string {
	return []string{"INVALID", "EMPTY", "ROCK", "FIELD"}[t]
}

func (t Terrain) MarshalText() (text []byte, err error) {
	return []byte(t.String()), nil
}

var directionStrings = []string{"E", "SE", "SW", "W", "NW", "NE"}

func (d Direction) String() string {
	return directionStrings[d]
}

func (d Direction) MarshalText() (text []byte, err error) {
	return []byte(directionStrings[d]), nil
}

func (d *Direction) UnmarshalText(b []byte) error {
	str := string(b)
	index := slices.Index(directionStrings, str)

	if index >= 0 {
		*d = Direction(index)
		return nil
	}

	return fmt.Errorf("Could not unmarshal Direction: %s", str)
}

func (t EntityType) String() string {
	return []string{"WALL", "HIVE", "BEE"}[t]
}

func (t EntityType) MarshalText() (text []byte, err error) {
	return []byte(t.String()), nil
}

var orderTypeStrings = []string{
	"MOVE",
	"ATTACK",
	"BUILD_WALL",
	"BUILD_HIVE",
	"FORAGE",
	"SPAWN",
}

func (t OrderType) String() string {
	return orderTypeStrings[t]
}

func (t OrderType) MarshalText() (text []byte, err error) {
	return []byte(t.String()), nil
}

func (t *OrderType) UnmarshalText(b []byte) error {
	str := string(b)
	index := slices.Index(orderTypeStrings, str)

	if index >= 0 {
		*t = OrderType(index)
		return nil
	}

	return fmt.Errorf("Could not unmarshal OrderType: %s", str)
}
