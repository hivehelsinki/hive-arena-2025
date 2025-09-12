package main

import (
	"fmt"
	"hive-arena/common"
)

func main() {

	foo, err := common.LoadMap("maps/balanced.txt")

	fmt.Printf("%+v \n %v", foo, err)
}
