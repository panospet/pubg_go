package main

import (
	"fmt"
	"os"

	utils "github.com/pubg_go/pubg_last_id/utils"
)

func main() {
	c := make(chan string)
	vkc := make(chan string)
	playerName := os.Args[1]
	go utils.GetMatchIDs(playerName, c)
	for v := range c {
		go utils.Wrapchan(playerName, v, vkc)
	}
	for i := range vkc {
		fmt.Println(i)
	}
}
