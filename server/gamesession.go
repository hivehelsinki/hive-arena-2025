package main

import . "hive-arena/common"

type GameSession struct {

}

func NewGameSession(players int, mapdata MapData) *GameSession {

	return &GameSession{}
}

func (game *GameSession) IsFull() bool {
	return false
}
