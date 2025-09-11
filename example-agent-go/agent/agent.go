package agent

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/gorilla/websocket"
	"io"
	"net/http"
	"os"
)

func request(url string) string {

	resp, err := http.Get(url)
	if err != nil {
		fmt.Println("Could not get " + url)
		os.Exit(1)
	}
	defer resp.Body.Close()

	bodyBytes, _ := io.ReadAll(resp.Body)
	body := string(bodyBytes)

	if resp.StatusCode != 200 {
		fmt.Println("Error:", body)
		os.Exit(1)
	}

	return body
}

type JoinResponse struct {
	Id    uint
	Token string
}

type WebSocketMessage struct {
	Turn     uint
	GameOver bool
}

func joinGame(host string, id uint, name string) JoinResponse {

	url := "http://" + host + fmt.Sprintf("/join?id=%d&name=%s", id, name)
	body := request(url)

	var response JoinResponse
	json.Unmarshal([]byte(body), &response)

	fmt.Printf("Joined game %d as player %d\n", id, response.Id)

	return response
}

func startWebSocket(host string, id uint) *websocket.Conn {

	url := "ws://" + host + fmt.Sprintf("/ws?id=%d", id)

	ws, _, err := websocket.DefaultDialer.Dial(url, nil)

	if err != nil {
		fmt.Println("Websocket error: ", err)
		os.Exit(1)
	}

	return ws
}

func getState(host string, id uint, token string) GameState {

	url := "http://" + host + fmt.Sprintf("/game?id=%d&token=%s", id, token)
	body := request(url)

	var response GameState
	json.Unmarshal([]byte(body), &response)

	return response
}

func sendOrders(host string, id uint, token string, orders []Order) {
	url := "http://" + host + fmt.Sprintf("/orders?id=%d&token=%s", id, token)
	payload, err := json.Marshal(orders)

	resp, err := http.Post(url, "application/json", bytes.NewReader(payload))
	if err != nil {
		fmt.Println("Could not post to " + url)
		os.Exit(1)
	}
	defer resp.Body.Close()

	bodyBytes, _ := io.ReadAll(resp.Body)
	body := string(bodyBytes)

	if resp.StatusCode != 200 {
		fmt.Println("Error:", body)
		os.Exit(1)
	}
}

func Run(host string, id uint, name string, callback func(*GameState, uint) []Order ) {

	playerInfo := joinGame(host, id, name)
	ws := startWebSocket(host, id)

	for {
		var message WebSocketMessage
		err := ws.ReadJSON(&message)
		if err != nil {
			fmt.Println("Error:", err)
			os.Exit(1)
		}

		if message.GameOver {
			fmt.Println("Game is over")
			break
		} else {
			fmt.Printf("Starting turn %d\n", message.Turn)
		}

		state := getState(host, id, playerInfo.Token)
		orders := callback(&state, playerInfo.Id)

		sendOrders(host, id, playerInfo.Token, orders)
	}
}
