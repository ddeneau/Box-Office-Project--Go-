package main

import (
	"encoding/csv"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strings"
)

//tt0103285 //IMDB ID - Sample
// My Key : 6110282

var titles [25][1]string
var index int

func findTitles() {
	currentID := ""
	batchSize := 20
	a := [25]string{"tt0117008", "tt4801232", "tt0117008", "tt0091939", "tt1247692",
		"tt0088847", "tt0090305", "tt0092718", "tt0329774", "tt2015381",
		"tt0088763", "tt0351283", "tt0180093", "tt2584384", "tt10954652",
		"tt1375670", "tt0099685", "tt0114369", "tt4846232", "tt3281548",
		"tt0081237", "tt0097523", "tt0112642", "tt0090830", "tt0090830"}

	for i := 0; i < batchSize; i++ {
		//n0 := rand.Int(4)
		//currentID = fmt.Sprintf("tt%d%d%d%d%d%d%d", n0, n1, n2, n3, n4, n5, n6)
		currentID = a[i]
		if index < 0 {
			index = 0
		}
		if i < 0 {
			i = 0
		}

		if connectAndCollect(currentID) == false {
			i--
			index--
			continue
		}
	}

	fmt.Printf(titles[0][0])
}

func connectAndCollect(id string) bool {
	movieMap := make(map[string]string) // Create Map

	url := fmt.Sprintf("http://www.omdbapi.com/?i=%s&apikey=f27c112f&type=movie", id) // Establish URL

	req, _ := http.NewRequest("GET", url, nil) // Establish Request.

	res, err := http.DefaultClient.Do(req) // Aquire Response (or error)

	if err != nil { // Check for error.
		fmt.Println("Error in func connectAndCollect")
		return false
	}

	defer res.Body.Close() // Close the resonse.

	body, _ := ioutil.ReadAll(res.Body) // Convert HTML Reqest Body into string.

	json.Unmarshal([]byte(body), &movieMap) // Decode JSON format into our map.

	string_ := formatData(movieMap["Title"], movieMap["Released"], movieMap["Runtime"],
		movieMap["Metascore"], movieMap["BoxOffice"], movieMap["Actors"], movieMap["Director"], movieMap["Writer"])

	titles[index][0] = string_
	index++
	return checkForFields(movieMap)
}

// Checks that we have the properites we need for each movie.
func checkForFields(_map map[string]string) bool {
	if _map["Type"] == "episode" {
		return false
	}

	if _map["Title"] == "N/A" || _map["Released"] == "N/A" ||
		_map["Runtime"] == "N/A" || _map["imdbRating"] == "N/A" || strings.Contains(_map["BoxOffice"], "United States") {
		return false
	} else {
		return true
	}
}

func formatData(title string, year string, runtime string, score string, boxOffice string,
	cast string, director string, writer string) string {
	return fmt.Sprintf("%s,%s,%s,%s,%s,%s,%s,%s", title, year, runtime, score,
		boxOffice, strings.Replace(cast, ",", " ", -1), director, strings.Replace(cast, ",", " ", -1))
}

func main() {
	findTitles()

	csvFile, err := os.Create(fmt.Sprintf("%s.csv", "films"))

	if err != nil {
		log.Fatalf("failed creating file: %s", err)
	}

	fileWriter := csv.NewWriter(csvFile)
	fileWriter.Comma = '\t'

	defer csvFile.Close()

	//i := 1; i < 5; i++
	fileWriter.Write([]string{"Title", "Release", "Runtime", "Metascore", "Box Office", "Cast", "Director", "Writer"})
	for i := 0; i < len(titles); i++ {
		a := titles[i]
		fileWriter.Write(a[:])
	}

	fileWriter.Flush()
}
