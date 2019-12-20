package utils

import (
	"encoding/json"
	"fmt"
	"github.com/joho/godotenv"
	"io/ioutil"
	"log"
	"net/http"
	"os"
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

// Load the PUBG_API_KEY environment variable
func init() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}
}

type ResultCollection []interface{}

type LogPlayerKill struct {
	KillerName string
	VictimName string
}

func keyExists(decoded map[string]interface{}, key string) bool {
	val, ok := decoded[key]
	return ok && val != nil
}

func GetKillersVictims(telUrl string) []LogPlayerKill {
	var res ResultCollection
	getTelUrlResponse := getReq(telUrl, true)
	err := json.Unmarshal([]byte(getTelUrlResponse), &res)
	if err != nil {
		panic(err)
	}

	var all []LogPlayerKill
	for i := range res {
		obj := res[i].(map[string]interface{})
		if obj["_T"] == "LogPlayerKill" {
			if keyExists(obj, "killer") && keyExists(obj, "victim") {
				killerName :=  obj["killer"].(map[string]interface{})["name"].(string)
				victimName :=  obj["victim"].(map[string]interface{})["name"].(string)
				all = append(all, LogPlayerKill{
					KillerName: killerName,
					VictimName: victimName,
				})
			}
		}
	}

	return all
}

// GetLastID fetches the last match id of a specific player along with his account id
func (p Player) GetLastID() (string, string) {
	url := "https://api.pubg.com/shards/steam/players?filter[playerNames]=meximonster"
	body := getReq(url, false)
	json.Unmarshal([]byte(body), &p)
	accid := p.Data[0].ID
	lastid := p.Data[0].Relationships.Matches.Data[0].ID
	return accid, lastid
}

type Match struct {
	Included []IncludedElement `json:"included"`
}

type IncludedElement struct {
	Type       string                 `json:"type"`
	Attributes map[string]interface{} `json:"attributes"`
}

// getTelemetry fetches the telemetry url of a certain match id provided as input
func GetTelemetryUrl(matchid string) string {
	var m Match
	var telemetryUrl string
	url := fmt.Sprintf("https://api.pubg.com/shards/steam/matches/%s", matchid)
	body := getReq(url, false)
	err := json.Unmarshal(body, &m)
	if err != nil {
		panic(err)
	}
	for i := range m.Included {
		if m.Included[i].Type == "asset" {
			telemetryUrl = m.Included[i].Attributes["URL"].(string)
		}
	}
	return telemetryUrl
}

// getReq makes the get request to an endpoint provided and given no errors, returns the body as slice of bytes
func getReq(endpoint string, useGzipHeader bool) []uint8 {
	apikey := os.Getenv("PUBG_API_KEY")
	bearer := fmt.Sprintf("Bearer %s", apikey)
	client := &http.Client{}
	req, _ := http.NewRequest("GET", endpoint, nil)
	req.Header.Set("Authorization", bearer)
	req.Header.Set("Accept", "application/vnd.api+json")
	// All telemetry URLs end with "json" and are all compressed using gzip
	if useGzipHeader {
		req.Header.Set("Accept", "Content-Encoding: gzip")
	}
	res, err := client.Do(req)
	if err != nil {
		panic(err)
	}
	defer res.Body.Close()
	body, _ := ioutil.ReadAll(res.Body)
	if resMessage := statusHandler(res.StatusCode); resMessage != "SUCCESS" {
		fmt.Print(resMessage, "\n")
		os.Exit(3)
	}
	return body
}

// Handles all possible status codes according to official PUBG API documentation
func statusHandler(statuscode int) string {
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
