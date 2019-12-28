package utils

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"sync"

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

// GetMatchIDs fetches all match ids for the last 2 weeks
func GetMatchIDs(playerName string, c chan string) {
	url := "https://api.pubg.com/shards/steam/players?filter[playerNames]=" + playerName
	body := getReq(url, true, false)
	var p fastjson.Parser
	v, err := p.ParseBytes([]byte(body))
	if err != nil {
		log.Fatal(err)
	}
	vv := v.GetArray("data", "0", "relationships", "matches", "data")
	for i := 0; i < 10; i++ {
		c <- string(vv[i].GetStringBytes("id"))
	}
	close(c)
}

// GetTelemetryURL fetches the telemetry url of a certain match id provided as input
func GetTelemetryURL(matchid string) string {
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
	return telemetryURL
}

// GetKillersVictims fetches the killers and victims of a match
func GetKillersVictims(playerName string, telURL string) ([]string, string) {
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
	return victims, killer
}

// getReq makes the get request to an endpoint provided and given no errors, returns the body as slice of bytes
func getReq(endpoint string, needAuth bool, useGzipHeader bool) []uint8 {
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
	return body
}

// Check if get request fails
func statusHandler(endpoint string, statuscode int) {
	if statuscode != 200 {
		log.Fatalf("Get request to %v failed with status code %v", endpoint, statuscode)
	}
}

// Handleresults manages the output
func Handleresults(v []string, k string, vkc chan string) {
	if len(v) != 0 {
		for i := range v {
			vkc <- v[i] + ".victim"
		}
	}
	if k != "" {
		vkc <- k + ".killer"
	}
}

// Wrapchan sums all the above
func Wrapchan(playerName, lastid string, vkc chan string, wg *sync.WaitGroup) {
	defer wg.Done()
	telURL := GetTelemetryURL(lastid)
	v, k := GetKillersVictims(playerName, telURL)
	Handleresults(v, k, vkc)
}
