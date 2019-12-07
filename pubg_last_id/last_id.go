package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"

	"github.com/joho/godotenv"
)

// Player object
type Player struct {
	Data []struct {
		Relationships struct {
			Matches struct {
				Data []struct {
					ID string `json:"id"`
				} `json:"data"`
			} `json:"matches"`
		} `json:"relationships"`
	} `json:"data"`
}

func init() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}
}

func main() {
	apikey := os.Getenv("PUBG_API_KEY")
	bearer := fmt.Sprintf("Bearer %s", apikey)
	p := Player{}
	lastmatchid := p.Getplayer(bearer)
	fmt.Print(lastmatchid)
}

// Getplayer func
func (p Player) Getplayer(bearer string) string {
	client := &http.Client{}
	req, _ := http.NewRequest("GET", "https://api.pubg.com/shards/steam/players?filter[playerNames]=meximonster", nil)
	req.Header.Set("Authorization", bearer)
	req.Header.Set("Accept", "application/vnd.api+json")
	res, err := client.Do(req)
	if err != nil {
		panic(err)
	}
	defer res.Body.Close()
	//s := res.StatusCode
	body, _ := ioutil.ReadAll(res.Body)
	json.Unmarshal([]byte(body), &p)
	lastid := p.Data[0].Relationships.Matches.Data[0].ID
	return lastid

	// ---- IN CASE ALL IDS ARE NEEDED ----
	//ids := []string{}
	//for i := range p.Data[0].Relationships.Matches.Data {
	//	ids = append(ids, p.Data[0].Relationships.Matches.Data[i].ID)
	//}

}
