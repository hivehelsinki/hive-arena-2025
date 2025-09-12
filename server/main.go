package main

import (
	"fmt"
	"encoding/json"
	"hive-arena/common"
)

func main() {

	mapData, _ := common.LoadMap("maps/balanced.txt")

	gs := common.NewGameState(mapData, 4)

	txt, _ := json.Marshal(gs)
	fmt.Println(string(txt))
}
