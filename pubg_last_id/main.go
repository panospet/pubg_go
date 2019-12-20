package main

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

func main() {
	p := Player{}
	_, lastid := p.GetLastID()
	//fmt.Printf("Account id: %v\nLast match id: %v\n", accid, lastid)
	telURL := getTelemetryUrl(lastid)
	getKillersVictims(telURL)
}

type ResultCollection []Result
type Result struct {
	Type string `json:"_T"`
}

func getKillersVictims(telUrl string) {
	var res ResultCollection
	getTelUrlResponse := getreq(telUrl, true)
	err := json.Unmarshal([]byte(getTelUrlResponse), &res)
	if err != nil {
		panic(err)
	}

	for i := range res {
		fmt.Println(res[i].Type)
	}
}

// GetLastID fetches the last match id of a specific player along with his account id
func (p Player) GetLastID() (string, string) {
	url := "https://api.pubg.com/shards/steam/players?filter[playerNames]=meximonster"
	body := getreq(url, false)
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
func getTelemetryUrl(matchid string) string {
	var m Match
	var telemetryUrl string
	url := fmt.Sprintf("https://api.pubg.com/shards/steam/matches/%s", matchid)
	body := getreq(url, false)
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

// getreq makes the get request to an endpoint provided and given no errors, returns the body as slice of bytes
func getreq(endpoint string, useGzipHeader bool) []uint8 {
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
	if resMessage := statushandler(res.StatusCode); resMessage != "SUCCESS" {
		//fmt.Print(resMessage, "\n")
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

func getSample() string {
	SAMPLE := `[
    {
        "MatchId": "match.bro.official.pc-2018-05.steam.squad-fpp.eu.2019.12.19.22.d1624c59-2192-486a-851f-27062ce1988a",
        "PingQuality": "low",
        "SeasonState": "progress",
        "_D": "2019-12-19T22:10:34.7332693Z",
        "_T": "LogMatchDefinition"
    },
    {
        "_D": "2019-12-19T22:09:13.235Z",
        "_T": "LogPlayerLogin",
        "accountId": "account.7db33b515d9343ac9a4e65a5c22770cf",
        "common": {
            "isGame": 0
        }
    },
    {
        "_D": "2019-12-19T22:09:13.264Z",
        "_T": "LogPlayerCreate",
        "character": {
            "accountId": "account.7db33b515d9343ac9a4e65a5c22770cf",
            "health": 100,
            "isInBlueZone": false,
            "isInRedZone": false,
            "location": {
                "x": 795740.625,
                "y": 21110.17578125,
                "z": 547.231201171875
            },
            "name": "Zolasky_",
            "ranking": 0,
            "teamId": 3,
            "zone": []
        },
        "common": {
            "isGame": 0
        }
    },
    {
        "_D": "2019-12-19T22:09:13.265Z",
        "_T": "LogPlayerLogin",
        "accountId": "account.3c31556102d5431fb51f4ee8c71e68e5",
        "common": {
            "isGame": 0
        }
    },
    {
        "_D": "2019-12-19T22:09:13.301Z",
        "_T": "LogPlayerCreate",
        "character": {
            "accountId": "account.3c31556102d5431fb51f4ee8c71e68e5",
            "health": 100,
            "isInBlueZone": false,
            "isInRedZone": false,
            "location": {
                "x": 344180.03125,
                "y": 170009.15625,
                "z": 1540.6497802734375
            },
            "name": "HanK_HulaY",
            "ranking": 0,
            "teamId": 7,
            "zone": []
        },
        "common": {
            "isGame": 0
        }
    },
    {
        "_D": "2019-12-19T22:09:13.302Z",
        "_T": "LogPlayerLogin",
        "accountId": "account.ee503bca56e1409abdc544810a58c8f0",
        "common": {
            "isGame": 0
        }
    }]
`

	return SAMPLE
}
