package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
)

// tt0103285 IMDB ID - Sample
// My Key : 6110282

func connectAndCollect() {
	movieMap := make(map[string]string) // Create Map

	url := "http://www.omdbapi.com/?t=joker&apikey=6110282" // Establish URL

	req, _ := http.NewRequest("GET", url, nil) // Establish Request.

	res, err := http.DefaultClient.Do(req) // Aquire Response (or error)
	
	if(err != nil) { 					// Check for error.
		fmt.Println(err)
	}

	defer res.Body.Close() 				// Close the resonse.

	body, _ := ioutil.ReadAll(res.Body) 	// Convert HTML Reqest Body into string.

	json.Unmarshal([]byte(body), &movieMap) // Decode JSON format into our map.

	fmt.Println(string(movieMap["Title"])) // Test by printing out the title.
}

func writeTo(_map map[string]string, filename string) {
	
}

func main() {
	connectAndCollect()
}
