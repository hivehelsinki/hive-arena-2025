package common

import (
	"fmt"
	"strconv"
	"strings"
)

func (c Coords) String() string {
	return fmt.Sprintf("%d,%d", c.Row, c.Col)
}

func (c *Coords) FromString(s string) error {
	parts := strings.Split(s, ",")
	if len(parts) != 2 {
		return fmt.Errorf("invalid coordinate string format: %s", s)
	}
	row, err := strconv.Atoi(parts[0])
	if err != nil {
		return fmt.Errorf("invalid row value: %w", err)
	}
	col, err := strconv.Atoi(parts[1])
	if err != nil {
		return fmt.Errorf("invalid col value: %w", err)
	}
	c.Row = row
	c.Col = col
	return nil
}

func (c Coords) MarshalText() ([]byte, error) {
	return []byte(c.String()), nil
}

func (c *Coords) UnmarshalText(b []byte) error {
	return c.FromString(string(b))
}
