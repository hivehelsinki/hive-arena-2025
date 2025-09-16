package main

import (
	"net/http"
	"fmt"
	"os"
	"log"
	"strings"
	"maps"
	"slices"
)

import . "hive-arena/common"

const MapDir = "maps"

type Server struct {
	Maps map[string]MapData
}

func (server *Server) handleIndex(w http.ResponseWriter, r *http.Request) {
	log.Printf("Maps %+v", server.Maps)

	body := "Hello world " + GenerateID()

	w.Write([]byte(body))
}

func loadMaps() map[string]MapData {

	data := make(map[string]MapData)

	entries, err := os.ReadDir(MapDir)
	if err != nil {
		log.Fatalf("Could not find maps directory")
	}

	for _, entry := range entries {
		name := entry.Name()
		path := MapDir + "/" + name
		mapdata, err := LoadMap(path)
		if err != nil {
			log.Fatalf("Could not load map %s: %s", name, err)
		}

		name = strings.ReplaceAll(name, ".txt", "")
		data[name] = mapdata
	}

	log.Printf("Loaded maps: %s", strings.Join(slices.Collect(maps.Keys(data)), ", "))

	return data
}

func RunServer(port int) {

	server := Server{
		Maps: loadMaps(),
	}

	http.HandleFunc("/", server.handleIndex)

	log.Printf("Listening on port %d", port)
	err := http.ListenAndServe(fmt.Sprintf(":%d", port), nil)
	fmt.Println(err)
}
