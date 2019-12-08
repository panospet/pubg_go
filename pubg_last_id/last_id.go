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
	p := Player{}
	lastid := p.GetLastID()
	fmt.Print(lastid)
}

// GetLastID fetches the last match id of a specific player
func (p Player) GetLastID() string {
	url := "https://api.pubg.com/shards/steam/players?filter[playerNames]=meximonster"
	body := getreq(url)
	json.Unmarshal([]byte(body), &p)
	lastid := p.Data[0].Relationships.Matches.Data[0].ID
	return lastid
}

func getreq(endpoint string) []uint8 {
	apikey := os.Getenv("PUBG_API_KEY")
	bearer := fmt.Sprintf("Bearer %s", apikey)
	client := &http.Client{}
	req, _ := http.NewRequest("GET", endpoint, nil)
	req.Header.Set("Authorization", bearer)
	req.Header.Set("Accept", "application/vnd.api+json")
	res, err := client.Do(req)
	if err != nil {
		panic(err)
	}
	defer res.Body.Close()
	body, _ := ioutil.ReadAll(res.Body)
	if resMessage := statushandler(res.StatusCode); resMessage != "SUCCESS" {
		fmt.Print(resMessage, "\n")
		os.Exit(3)
	}
	return body
}

// Handles all possible status codes according to official PUBG API documentation
func statushandler(statuscode int) string {
	var result string
	switch s := statuscode; s {
	case 401:
		result = "API key invalid or missing."
	case 404:
		result = "The specified resource was not found."
	case 415:
		result = "Content type incorrect or not specified."
	case 429:
		result = "Too many requests."
	default:
		result = "SUCCESS"
	}
	return result
}
