package main

import (
	"fmt"
	"os"
	"time"

	utils "github.com/pubg_go/pubg_last_id/utils"
)

func main() {
	start := time.Now()
	playerName := os.Args[1]
	wrap(playerName)
	// lastid := utils.GetLastID(playerName)
	// telURL := utils.GetTelemetryURL(lastid)
	// v, k := utils.GetKillersVictims(playerName, telURL)
	// utils.PrintResults(v, k)
	t := time.Now()
	elapsed := t.Sub(start)
	fmt.Printf("Execution took %v\n", elapsed)
}

func wrap(playerName string) {
	lastid := utils.GetLastID(playerName)
	telURL := utils.GetTelemetryURL(lastid)
	v, k := utils.GetKillersVictims(playerName, telURL)
	utils.PrintResults(v, k)
}
