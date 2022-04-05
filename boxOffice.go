package main

import (
	"encoding/csv"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"math/rand"
	"net/http"
	"os"
)

//tt0103285 //IMDB ID - Sample
// My Key : 6110282

var titles [100][1]string

func findTitles() {
	currentID := ""
	batchSize := 100
	
	for i := 0; i < batchSize; i++ {
		n1, n2 := rand.Intn(9), rand.Intn(9)
		n3, n4 := rand.Intn(9), rand.Intn(9)
		n5 := rand.Intn(9)	
		currentID = fmt.Sprintf("tt0%d0%d%d%d%d", n1, n2, n3, n4, n5)
		connectAndCollect(currentID, i)
	}

	fmt.Printf(titles[0][0])
}

func connectAndCollect(id string, index int) {
	movieMap := make(map[string]string) // Create Map

	url := fmt.Sprintf("http://www.omdbapi.com/?i=%s&apikey=6110282", id) // Establish URL

	req, _ := http.NewRequest("GET", url, nil) // Establish Request.

	res, err := http.DefaultClient.Do(req) // Aquire Response (or error)

	if err != nil { // Check for error.
		fmt.Println(err)
	}

	defer res.Body.Close() // Close the resonse.

	body, _ := ioutil.ReadAll(res.Body) // Convert HTML Reqest Body into string.

	json.Unmarshal([]byte(body), &movieMap) // Decode JSON format into our map.

	string_ := formatData(movieMap["Title"], movieMap["Year"], movieMap["Director"])

	titles[index][0] = string_
}

func formatData(title string, year string, director string) string{
	return fmt.Sprintf("%s,%s,%s", title, year, director)
}

func writeTo(data [][]string, filename string) {
	csvFile, err := os.Create(fmt.Sprintf("%s.csv", filename))
	defer csvFile.Close()

	if err != nil {
		log.Fatalf("failed creating file: %s", err)
	}

	fileWriter := csv.NewWriter(csvFile)
	fileWriter.Comma = '|'

	fileWriter.WriteAll(data)

	fileWriter.Flush()
}

func main() {
	findTitles()
}
