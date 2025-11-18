package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"github.com/hajimehoshi/ebiten/v2"
	"io"
	"net/http"
)

import . "hive-arena/common"

type Viewer struct {
}

func (viewer *Viewer) Update() error {
	return nil
}

func (viewer *Viewer) Draw(screen *ebiten.Image) {

}

func (viewer *Viewer) Layout(outsideWidth, outsideHeight int) (screenWidth, screenHeight int) {
	return outsideWidth, outsideHeight
}

func GetURL(url string) *PersistedGame {
	res, err := http.Get(url)
	if err != nil {
		fmt.Println(err)
		return nil
	}

	body, err := io.ReadAll(res.Body)
	res.Body.Close()

	if res.StatusCode > 299 {
		fmt.Printf("Response failed with status code: %d and\nbody: %s\n", res.StatusCode, body)
		return nil
	}

	if err != nil {
		fmt.Println(err)
		return nil
	}

	var game PersistedGame
	err = json.Unmarshal(body, &game)

	if err != nil {
		fmt.Println(err)
		return nil
	}

	return &game
}

func main() {
	url := flag.String("url", "", "URL of the history file to view")
	flag.Parse()

	if *url == "" {
		flag.PrintDefaults()
		return
	}

	game := GetURL(*url)
	if game == nil {
		return
	}

	fmt.Println(game)

	ebiten.SetWindowSize(1024, 768)
	ebiten.SetWindowTitle("Hive Arena Viewer")
	ebiten.SetWindowResizingMode(ebiten.WindowResizingModeEnabled)

	viewer := &Viewer{}
	err := ebiten.RunGame(viewer)

	if err != nil {
		fmt.Println(err)
	}
}
