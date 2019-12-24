package utils

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/joho/godotenv"
)

// Player object
type Player struct {
	Data []struct {
		ID            string `json:"id"`
		Relationships struct {
			Matches struct {
				Data []struct {
					ID string `json:"id"`
				} `json:"data"`
			} `json:"matches"`
		} `json:"relationships"`
	} `json:"data"`
}

// Match object
type Match struct {
	Included []IncludedElement `json:"included"`
}

// IncludedElement is needed for retrieving TelemetryUrl
type IncludedElement struct {
	Type       string                 `json:"type"`
	Attributes map[string]interface{} `json:"attributes"`
}

// Events object
type Events []interface{}

// LogPlayerKill event
type LogPlayerKill struct {
	KillerName string
	VictimName string
}

// Load the PUBG_API_KEY environment variable
func init() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}
}

// keyExists returns true if key exists in map, else false
func keyExists(decoded map[string]interface{}, key string) bool {
	val, ok := decoded[key]
	return ok && val != nil
}

// GetLastID fetches the last match id of a specific player along with his account id
func (p Player) GetLastID() (string, string) {
	start := time.Now()
	url := "https://api.pubg.com/shards/steam/players?filter[playerNames]=" + os.Args[1]
	body := getReq(url, false)
	err := json.Unmarshal([]byte(body), &p)
	if err != nil {
		panic(err)
	}
	accid := p.Data[0].ID
	lastid := p.Data[0].Relationships.Matches.Data[0].ID
	t := time.Now()
	elapsed := t.Sub(start)
	fmt.Printf("GetLastID took %v\n", elapsed)
	return accid, lastid
}

// GetTelemetryURL fetches the telemetry url of a certain match id provided as input
func GetTelemetryURL(matchid string) string {
	start := time.Now()
	var m Match
	var telemetryURL string
	url := "https://api.pubg.com/shards/steam/matches/" + matchid
	body := getReq(url, false)
	err := json.Unmarshal(body, &m)
	if err != nil {
		panic(err)
	}
	for i := range m.Included {
		if m.Included[i].Type == "asset" {
			telemetryURL = m.Included[i].Attributes["URL"].(string)
		}
	}
	t := time.Now()
	elapsed := t.Sub(start)
	fmt.Printf("GetTelemetryURL took %v\n", elapsed)
	return telemetryURL
}

// GetKillersVictims fetches the killers and victims of a match
func GetKillersVictims(telURL string) []LogPlayerKill {
	start := time.Now()
	var res Events
	gettelURLResponse := getReq(telURL, true)
	err := json.Unmarshal([]byte(gettelURLResponse), &res)
	if err != nil {
		panic(err)
	}

	var all []LogPlayerKill
	for i := range res {
		obj := res[i].(map[string]interface{})
		if obj["_T"] == "LogPlayerKill" {
			if keyExists(obj, "killer") && keyExists(obj, "victim") {
				killerName := obj["killer"].(map[string]interface{})["name"].(string)
				victimName := obj["victim"].(map[string]interface{})["name"].(string)
				all = append(all, LogPlayerKill{
					KillerName: killerName,
					VictimName: victimName,
				})
			}
		}
	}
	t := time.Now()
	elapsed := t.Sub(start)
	fmt.Printf("GetKillersVictims took %v\n", elapsed)
	return all
}

// getReq makes the get request to an endpoint provided and given no errors, returns the body as slice of bytes
func getReq(endpoint string, useGzipHeader bool) []uint8 {
	apikey := os.Getenv("PUBG_API_KEY")
	bearer := fmt.Sprintf("Bearer %s", apikey)
	client := &http.Client{}
	req, _ := http.NewRequest("GET", endpoint, nil)
	req.Header.Set("Authorization", bearer)
	req.Header.Set("Accept", "application/vnd.api+json")
	// All telemetry URLs are all compressed using gzip
	if useGzipHeader {
		req.Header.Set("Accept", "Content-Encoding: gzip")
	}
	res, err := client.Do(req)
	if err != nil {
		panic(err)
	}
	defer res.Body.Close()
	body, _ := ioutil.ReadAll(res.Body)
	statusHandler(endpoint, res.StatusCode)
	return body
}

// Check if get request fails
func statusHandler(endpoint string, statuscode int) {
	if statuscode != 200 {
		log.Fatalf("Get request to %v failed with status code %v", endpoint, statuscode)
	}
}
