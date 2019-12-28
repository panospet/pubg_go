package main

import (
	"fmt"
	"os"
	"sync"

	utils "github.com/pubg_go/pubg_last_id/utils"
)

func main() {
	wg := &sync.WaitGroup{}
	c := make(chan string)
	vkc := make(chan string)
	playerName := os.Args[1]
	go utils.GetMatchIDs(playerName, c)
	for v := range c {
		wg.Add(1)
		go utils.Wrapchan(playerName, v, vkc, wg)
	}
	go func() {
		wg.Wait()
		close(vkc)
	}()
	for i := range vkc {
		fmt.Println(i)
	}
}
