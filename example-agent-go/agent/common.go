package agent

import (
	"fmt"
	"strings"
	"strconv"
	"errors"
)

type Coords struct {
	Row int
	Col int
}

type Entity struct {
	Type   string
	Hp     uint
	Player uint
}

type Hex struct {
	Terrain   string
	Resources uint
	Influence uint
	Entity    *Entity
}

type GameState struct {
	NumPlayers          uint
	Turn                uint
	Hexes               map[Coords]Hex
	PlayerResources     []uint
	lastInfluenceChange uint
	Winners             map[uint]bool
	GameOver            bool
}

type Order struct {
	Type string `json:"type"`
	Coords Coords `json:"coords"`
	Direction string `json:"direction"`
}

func (c Coords) MarshalText() (text []byte, err error) {
	str := fmt.Sprintf("%d,%d", c.Row, c.Col)
	return ([]byte)(str), nil
}

func (c *Coords) UnmarshalText(text []byte) error {
	parts := strings.Split(string(text), ",")

	if len(parts) != 2 {
		return errors.New("Bad coords")
	}

	var err1, err2 error
	c.Row, err1 = strconv.Atoi(parts[0])
	c.Col, err2 = strconv.Atoi(parts[1])

	if err1 != nil { return err1 }
	if err2 != nil { return err2 }

	return nil
}
