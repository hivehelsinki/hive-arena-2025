package main

import (
	"fmt"
	"maps"
	"math/rand"
	"slices"
)

import . "hive-arena/common"

type Player struct {
	ID    int
	Name  string
	Token string
}

type GameSession struct {
	ID           string
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

func NewGameSession(id string, players int, mapdata MapData) *GameSession {

	tokens := generateTokens(players + 1)

	return &GameSession{
		ID:           id,
		AdminToken:   tokens[0],
		PlayerTokens: tokens[1:],
		State:        NewGameState(mapdata, players),
	}
}

func (game *GameSession) IsFull() bool {
	return len(game.Players) == game.State.NumPlayers
}
