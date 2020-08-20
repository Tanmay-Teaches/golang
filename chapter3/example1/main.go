package main

/*
Example of only using many build-in packages in Go to reach out to a rest API to retrieve movie detail.
*/

import (
	"encoding/json"
	"errors"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"
)

//omdbapi.com API key
const APIKEY = "193ef3a"

//The structure of the return JSON from omdbapi.com
//To keep this example short, some of the value are not map into the structure
type MovieInfo struct {
	Title      string `json:"Title"`
	Year       string `json:"Year"`
	Rated      string `json:"Rated"`
	Released   string `json:"Released"`
	Runtime    string `json:"Runtime"`
	Genre      string `json:"Genre"`
	Writer     string `json:"Writer"`
	Actors     string `json:"Actors"`
	Plot       string `json:"Plot"`
	Language   string `json:"Language"`
	Country    string `json:"Country"`
	Awards     string `json:"Awards"`
	Poster     string `json:"Poster"`
	ImdbRating string `json:"imdbRating"`
	ImdbID     string `json:"imdbID"`
}

func main() {
	body, _ := SearchById("tt3896198")
	println(body.Title)
	body, _ = SearchByName("Game of")
	println(body.Title)
}

func SearchByName(name string) (*MovieInfo, error) {
	parms := url.Values{}
	parms.Set("apikey", APIKEY)
	parms.Set("t", name)
	siteURL := "http://www.omdbapi.com/?" + parms.Encode()
	body, err := sendGetRequest(siteURL)
	if err != nil {
		return nil, errors.New(err.Error() + "\nBody:" + body)
	}
	mi := &MovieInfo{}
	return mi, json.Unmarshal([]byte(body), mi)
}

func SearchById(id string) (*MovieInfo, error) {
	parms := url.Values{}
	parms.Set("apikey", APIKEY)
	parms.Set("i", id)
	siteURL := "http://www.omdbapi.com/?" + parms.Encode()
	body, err := sendGetRequest(siteURL)
	if err != nil {
		return nil, errors.New(err.Error() + "\nBody:" + body)
	}
	mi := &MovieInfo{}
	return mi, json.Unmarshal([]byte(body), mi)
}

func sendGetRequest(url string) (string, error) {
	resp, err := http.Get(url)
	if err != nil {
		return "", err
	}
	statusCode := resp.Status
	if i := strings.IndexByte(resp.Status, ' '); i != -1 {
		statusCode = resp.Status[:i]
	}

	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	if statusCode != "200" {
		return string(body), errors.New(resp.Status)
	}
	return string(body), nil
}
