package main

import (
	"encoding/json"
	"fmt"
	"github.com/teomiscia/hexbattle/internal/game"
	"github.com/teomiscia/hexbattle/internal/model"
)

func main() {
	gs := game.NewGameState("id", model.RoomSettings{}, model.PlayerState{}, model.PlayerState{}, 1)
	b, err := json.Marshal(gs)
	if err != nil {
		fmt.Println("ERROR:", err)
	} else {
		fmt.Println("OK:", string(b))
	}
}
