package main

import (
	"fmt"
	"encoding/json"
	"hive-arena/common"
)

func main() {

	mapData, _ := common.LoadMap("maps/balanced.txt")

	gs := common.NewGameState(mapData, 4)

	txt, err := json.Marshal(gs)
	fmt.Printf("%+v\n", gs)
	fmt.Println(string(txt), err)
}
