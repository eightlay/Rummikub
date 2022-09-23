package server

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/eightlay/rummikub-server/iternal/server/hub"
	"github.com/gorilla/mux"
)

func StartServer() {
	m := hub.CreateManager()
	r := mux.NewRouter()

	r.HandleFunc("/api/health", func(w http.ResponseWriter, r *http.Request) {
		json.NewEncoder(w).Encode(map[string]bool{"ok": true})
	})

	r.HandleFunc("/api/join/{roomSize}", m.AddPlayer)
	r.HandleFunc("/api/room-ready/{roomID}", m.RoomIsReady)
	r.HandleFunc("/api/start-room/{roomID}", m.StartRoom)
	r.HandleFunc("/api/room-started/{roomID}", m.RoomStarted)
	r.HandleFunc("/api/game-state/{roomID}/{playerID}", m.CurrentGameState)
	r.HandleFunc("/api/make-action/{roomID}", m.ReceiveInput)

	log.Fatal(http.ListenAndServe(":8080", r))
}
