package main

import (
	"fmt"
	"os"
)

import . "hive-arena/agent"

func think(state *GameState, player uint) []Order {

	var orders []Order

	//fmt.Printf("%+v\n", state)

	for coords, hex := range state.Hexes {
		unit := hex.Entity

		if unit != nil && unit.Type == "BEE" && unit.Player == player {
			fmt.Println(coords, unit)
			orders = append(orders, Order{"MOVE", coords, "E"})
		}
	}

	return orders
}

func main() {
	if len(os.Args) <= 3 {
		fmt.Println("Usage: ./agent <host> <gameid> <name>")
		os.Exit(1)
	}

	host := os.Args[1]
	id := os.Args[2]
	name := os.Args[3]

	Run(host, id, name, think)
}
