package main

import (
	"encoding/json"
	"github.com/gorilla/mux"
	"github.com/pubg_go/pubg_last_id/utils"
	"net/http"
)

type ComparisonResponse struct {
	Player1      string
	Player1Stats utils.PlayerSeasonStats
	Player2      string
	Player2Stats utils.PlayerSeasonStats
}

func compare(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)

	player1 := vars["player1"]
	player2 := vars["player2"]

	acc1 := utils.GetAccid(player1)
	acc2 := utils.GetAccid(player2)
	stats1, stats2 := utils.GetSeasonStats(acc1, acc2)
	resp := ComparisonResponse{
		Player1:      player1,
		Player1Stats: stats1,
		Player2:      player2,
		Player2Stats: stats2,
	}

	js, err := json.Marshal(resp)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/json")
	w.Write(js)
}

func main() {
	r := mux.NewRouter()

	r.HandleFunc("/compare/{player1}/{player2}", compare)

	err := http.ListenAndServe(":4444", r)
	if err != nil {
		panic(err)
	}
}
