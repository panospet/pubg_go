package utils

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"

	"github.com/joho/godotenv"
	"github.com/valyala/fastjson"
)

// Load the PUBG_API_KEY environment variable
func init() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}
}

// GetLastID fetches the last match id of a specific player along with his account id
func GetLastID(playerName string) string {
	//start := time.Now()
	url := "https://api.pubg.com/shards/steam/players?filter[playerNames]=" + playerName
	body := getReq(url, true, false)
	lastid := fastjson.GetString([]byte(body), "data", "0", "relationships", "matches", "data", "0", "id")
	// t := time.Now()
	// elapsed := t.Sub(start)
	// fmt.Printf("GetLastID took %v\n", elapsed)
	return lastid
}

// GetTelemetryURL fetches the telemetry url of a certain match id provided as input
func GetTelemetryURL(matchid string) string {
	//start := time.Now()
	var telemetryURL string
	url := "https://api.pubg.com/shards/steam/matches/" + matchid
	body := getReq(url, false, false)
	var p fastjson.Parser
	v, err := p.ParseBytes([]byte(body))
	if err != nil {
		log.Fatal(err)
	}
	vv := v.GetArray("included")
	for i := range vv {
		if vv[i].Exists("attributes", "URL") {
			telemetryURL = string(vv[i].GetStringBytes("attributes", "URL"))
			break
		}
	}
	// t := time.Now()
	// elapsed := t.Sub(start)
	// fmt.Printf("GetTelemetryURL took %v\n", elapsed)
	return telemetryURL
}

// GetKillersVictims fetches the killers and victims of a match
func GetKillersVictims(playerName string, telURL string) ([]string, string) {
	// start := time.Now()
	gettelURLResponse := getReq(telURL, true, true)
	victims := []string{}
	var killer string
	var p fastjson.Parser
	v, err := p.ParseBytes([]byte(gettelURLResponse))
	if err != nil {
		log.Fatal(err)
	}
	vv := v.GetArray()
	for i := range vv {
		if string(vv[i].GetStringBytes("_T")) == "LogPlayerKill" {
			if string(vv[i].GetStringBytes("killer", "name")) == playerName {
				victims = append(victims, string(vv[i].GetStringBytes("victim", "name")))
				continue
			}
			if string(vv[i].GetStringBytes("victim", "name")) == playerName {
				killer = string(vv[i].GetStringBytes("killer", "name"))
			}
		}
	}
	// t := time.Now()
	// elapsed := t.Sub(start)
	// fmt.Printf("GetKillersVictims took %v\n", elapsed)
	return victims, killer
}

// getReq makes the get request to an endpoint provided and given no errors, returns the body as slice of bytes
func getReq(endpoint string, needAuth bool, useGzipHeader bool) []uint8 {
	// start := time.Now()
	client := &http.Client{}
	req, _ := http.NewRequest("GET", endpoint, nil)
	req.Header.Set("Accept", "application/vnd.api+json")
	if needAuth {
		apikey := os.Getenv("PUBG_API_KEY")
		bearer := fmt.Sprintf("Bearer %s", apikey)
		req.Header.Set("Authorization", bearer)
	}
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
	// t := time.Now()
	// elapsed := t.Sub(start)
	// fmt.Printf("getReq of %v took %v\n", endpoint, elapsed)
	return body
}

// Check if get request fails
func statusHandler(endpoint string, statuscode int) {
	if statuscode != 200 {
		log.Fatalf("Get request to %v failed with status code %v", endpoint, statuscode)
	}
}

// PrintResults manages the output
func PrintResults(v []string, k string) {
	fmt.Print("Victims : ")
	if len(v) != 0 {
		for i := range v {
			fmt.Print(v[i], ", ")
		}
	} else {
		fmt.Print("None!")
	}
	fmt.Print("\n")
	if k != "" {
		fmt.Println("Killer : ", k)
	} else {
		fmt.Println("You either survived or deathtype not by player.")
	}
}

// Wrap sums all the above
func Wrap(playerName string) {
	lastid := GetLastID(playerName)
	telURL := GetTelemetryURL(lastid)
	v, k := GetKillersVictims(playerName, telURL)
	PrintResults(v, k)
}
