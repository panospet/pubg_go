package main

import (
	"fmt"
	"time"

	utils "github.com/pubg_go/pubg_last_id/utils_v2"
)

func main() {
	start := time.Now()
	p := utils.Player{}
	_, lastid := p.GetLastID()
	//fmt.Printf("Account id: %v\nLast match id: %v\n", accid, lastid)
	telURL := utils.GetTelemetryURL(lastid)
	all := utils.GetKillersVictims(telURL)
	for i := range all {
		if all[i].KillerName == "meximonster" || all[i].VictimName == "meximonster" {
			fmt.Println(all[i])
		}
	}
	t := time.Now()
	elapsed := t.Sub(start)
	fmt.Printf("Execution took %v\n", elapsed)
}
