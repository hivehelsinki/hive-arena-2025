package common

import (
	"fmt"
	"strconv"
	"strings"
)

func (c Coords) String() string {
	return fmt.Sprintf("%d,%d", c.Row, c.Col)
}

func CoordsFromString(s string) (Coords, error) {
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
	coords, err := CoordsFromString(str)
	if err != nil {
		return err
	}

	*c = coords
	return nil
}
