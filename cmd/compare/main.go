package main

import (
	"fmt"
	"os"

	utils "github.com/pubg_go/pubg_last_id/utils"
)

func main() {
	player1 := os.Args[1]
	player2 := os.Args[2]
	acc1 := utils.GetAccid(player1)
	acc2 := utils.GetAccid(player2)
	stats1, stats2 := utils.GetSeasonStats(acc1, acc2)
	fmt.Printf("MatchesPlayed: %v %v\n", stats1.RoundsPlayed, stats2.RoundsPlayed)
	fmt.Printf("Wins: %v %v\n", stats1.Wins, stats2.Wins)
	fmt.Printf("Losses: %v %v\n", stats1.Losses, stats2.Losses)
	fmt.Printf("Top10S: %v %v\n", stats1.Top10S, stats2.Top10S)
	fmt.Printf("Kills: %v %v\n", stats1.Kills, stats2.Kills)
	fmt.Printf("DamageDealt: %v %v\n", stats1.DamageDealt, stats2.DamageDealt)
	fmt.Printf("Assists: %v %v\n", stats1.Assists, stats2.Assists)
	fmt.Printf("DBNOs: %v %v\n", stats1.DBNOs, stats2.DBNOs)
	fmt.Printf("HeadshotKills: %v %v\n", stats1.HeadshotKills, stats2.HeadshotKills)
	fmt.Printf("LongestKill: %v %v\n", stats1.LongestKill, stats2.LongestKill)
	fmt.Printf("MaxKillStreaks: %v %v\n", stats1.MaxKillStreaks, stats2.MaxKillStreaks)
	fmt.Printf("Revives: %v %v\n", stats1.Revives, stats2.Revives)
	fmt.Printf("RoundMostKills: %v %v\n", stats1.RoundMostKills, stats2.RoundMostKills)
	fmt.Printf("Suicides: %v %v\n", stats1.Suicides, stats2.Suicides)
	fmt.Printf("TeamKills: %v %v\n", stats1.TeamKills, stats2.TeamKills)
}
