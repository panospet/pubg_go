package main

import (
	"fmt"

	utils "github.com/pubg_go/pubg_last_id"
)

func main() {
	p := utils.Player{}
	_, lastid := p.GetLastID()
	//fmt.Printf("Account id: %v\nLast match id: %v\n", accid, lastid)
	telURL := utils.GetTelemetryURL(lastid)
	all := utils.GetKillersVictims(telURL)
	for i := range all {
		fmt.Println(all[i])
	}
}
