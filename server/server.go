package main

import (
	"encoding/json"
	"fmt"
	"log"
	"maps"
	"net/http"
	"os"
	"slices"
	"strconv"
	"strings"
	"sync"
	"time"
)

import . "hive-arena/common"

const MapDir = "maps"
const GameStartTimeout = 5 * time.Minute

type Server struct {
	mutex sync.Mutex

	Maps  map[string]MapData
	Games map[string]*GameSession
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

func logRoute(r *http.Request) {
	log.Printf("%s %v %v", r.Method, r.URL, r.RemoteAddr)
}

func writeJson(w http.ResponseWriter, payload any, status int) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(payload)
}

func (server *Server) handleNewGame(w http.ResponseWriter, r *http.Request) {

	logRoute(r)

	mapname := r.URL.Query().Get("map")
	mapdata, mapfound := server.Maps[mapname]
	if !mapfound {
		writeJson(w, "Map not found: "+mapname, http.StatusBadRequest)
		return
	}

	playerStr := r.URL.Query().Get("players")
	players, ok := strconv.Atoi(playerStr)
	if ok != nil || !IsValidNumPlayers(players) {
		writeJson(w, "Invalid number of players: "+playerStr, http.StatusBadRequest)
		return
	}

	server.mutex.Lock()
	id := GenerateUniqueID(server.Games)
	game := NewGameSession(id, players, mapname, mapdata)
	server.Games[id] = game
	server.mutex.Unlock()

	time.AfterFunc(GameStartTimeout, func() { server.removeIfNotStarted(id) })

	log.Printf("Created game %s (%s, %d players)", id, mapname, players)

	writeJson(w, map[string]any{
		"id":         game.ID,
		"adminToken": game.AdminToken,
	}, http.StatusOK)
}

func (server *Server) removeIfNotStarted(id string) {
	server.mutex.Lock()
	defer server.mutex.Unlock()

	game := server.Games[id]
	if game != nil && !game.IsFull() {
		delete(server.Games, id)
		log.Printf("Removed game %s because of timeout", id)
	}
}

func (server *Server) handleStatus(w http.ResponseWriter, r *http.Request) {
	logRoute(r)

	server.mutex.Lock()
	defer server.mutex.Unlock()

	var statuses []map[string]any
	for _, game := range server.Games {
		statuses = append(statuses, map[string]any{
			"id": game.ID,
			"createdDate": game.CreatedDate,
			"numPlayers": game.State.NumPlayers,
			"playersJoined": len(game.Players),
			"map": game.Map,
		})
	}

	writeJson(w, statuses, http.StatusOK)
}

func RunServer(port int) {

	server := Server{
		Maps:  loadMaps(),
		Games: make(map[string]*GameSession),
	}

	http.HandleFunc("GET /newgame", server.handleNewGame)
	http.HandleFunc("GET /status", server.handleStatus)

	log.Printf("Listening on port %d", port)

	err := http.ListenAndServe(fmt.Sprintf(":%d", port), nil)
	fmt.Println(err)
}
