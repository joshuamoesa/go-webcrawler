package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"os"

	"github.com/anaskhan96/soup"
)

type Omdbapiretrieve struct {
	Title    string `json:"Title"`
	Year     string `json:"Year"`
	Rated    string `json:"Rated"`
	Released string `json:"Released"`
	Runtime  string `json:"Runtime"`
	Genre    string `json:"Genre"`
	Director string `json:"Director"`
	Writer   string `json:"Writer"`
	Actors   string `json:"Actors"`
	Plot     string `json:"Plot"`
	Language string `json:"Language"`
	Country  string `json:"Country"`
	Awards   string `json:"Awards"`
	Poster   string `json:"Poster"`
	Ratings  []struct {
		Source string `json:"Source"`
		Value  string `json:"Value"`
	} `json:"Ratings"`
	Metascore  string `json:"Metascore"`
	ImdbRating string `json:"imdbRating"`
	ImdbVotes  string `json:"imdbVotes"`
	ImdbID     string `json:"imdbID"`
	Type       string `json:"Type"`
	DVD        string `json:"DVD"`
	BoxOffice  string `json:"BoxOffice"`
	Production string `json:"Production"`
	Website    string `json:"Website"`
	Response   string `json:"Response"`
}

func main() {
	resp, err := soup.Get("https://www.pathe-thuis.nl/films/collectie/81/nieuw")

	if err != nil {
		os.Exit(1)
	}
	doc := soup.HTMLParse(resp)
	list := doc.Find("body").FindAll("li", "class", "vertical-poster-list__item")

	fmt.Println("Debug: starting for-loop...")

	for _, i := range list {

		movie := i.Find("a")
		name := movie.Attrs()["data-product-name"]
		id := movie.Attrs()["data-product-id"]
		fmt.Println(name)
		fmt.Println("imdb: " + getImdbRating(name))
		//		fmt.Println("imdb: " + getImdbRating(title) + " | willem: " + willemRating.Text())
		//fmt.Println("imdb: " + getImdbRating(title))
		fmt.Println("https://www.pathe-thuis.nl/film/" + id)
		fmt.Println(" ")
	}
}

func getImdbRating(movie string) string {

	// QueryEscape escapes the phone string so
	// it can be safely placed inside a URL query
	safeMovie := url.QueryEscape(movie)

	url := fmt.Sprintf("http://www.omdbapi.com/?apikey=e6d94a21&t=%s&r=json", safeMovie)

	// Build the request
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		log.Fatal("NewRequest: ", err)
		//return
		os.Exit(1)
	}

	// For control over HTTP client headers,
	// redirect policy, and other settings,
	// create a Client
	// A Client is an HTTP client
	client := &http.Client{}

	// Send the request via a client
	// Do sends an HTTP request and
	// returns an HTTP response
	resp, err := client.Do(req)
	if err != nil {
		log.Fatal("Do: ", err)
		//return
		os.Exit(1)
	}

	// Callers should close resp.Body
	// when done reading from it
	// Defer the closing of the body
	defer resp.Body.Close()

	// Fill the record with the data from the JSON
	var record Omdbapiretrieve

	// Use json.Decode for reading streams of JSON data
	if err := json.NewDecoder(resp.Body).Decode(&record); err != nil {
		log.Println(err)
	}

	//fmt.Println("ImdbRating. = ", record.ImdbRating)

	if len(record.ImdbRating) < 1 {
		return "N/A"
	}

	return record.ImdbRating

}
