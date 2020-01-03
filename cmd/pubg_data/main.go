package main

import (
	"fmt"
	"os"
	"sync"
	"time"

	utils "github.com/pubg_go/pubg_last_id/utils"
)

func main() {
	start := time.Now()
	wg := &sync.WaitGroup{}
	c := make(chan string)
	vkc := make(chan utils.Player)
	playerName := os.Args[1]
	go utils.GetMatchIDs(playerName, c)
	for v := range c {
		wg.Add(1)
		go utils.Wrapchan(playerName, v, vkc, wg)
	}
	go utils.Wait(wg, vkc)
	for i := range vkc {
		fmt.Println(i)
	}
	t := time.Now()
	elapsed := t.Sub(start)
	fmt.Printf("Execution took %v seconds", elapsed)
}
