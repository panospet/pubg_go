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
	utils.Wrap(playerName)
	t := time.Now()
	elapsed := t.Sub(start)
	fmt.Printf("Execution took %v\n", elapsed)
}
