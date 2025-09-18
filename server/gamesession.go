package main

import (
	"fmt"
	"maps"
	"math/rand"
	"slices"
	"sync"
	"time"
)

import . "hive-arena/common"

type Player struct {
	ID    int
	Name  string
	Token string
}

type GameSession struct {
	mutex sync.Mutex

	ID           string
	Map          string
	CreatedDate  time.Time
	AdminToken   string
	PlayerTokens []string
	Players      []Player
	State        *GameState
}

func generateTokens(count int) []string {
	tokens := make(map[string]bool)

	for len(tokens) < count {
		tokens[fmt.Sprintf("%x", rand.Uint64())] = true
	}

	return slices.Collect(maps.Keys(tokens))
}

func NewGameSession(id string, players int, mapname string, mapdata MapData) *GameSession {

	tokens := generateTokens(players + 1)

	return &GameSession{
		ID:           id,
		Map:          mapname,
		CreatedDate:  time.Now(),
		AdminToken:   tokens[0],
		PlayerTokens: tokens[1:],
		State:        NewGameState(mapdata, players),
	}
}

func (game *GameSession) IsFull() bool {
	return len(game.Players) == game.State.NumPlayers
}
