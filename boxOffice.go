package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strings"

	"golang.org/x/net/html"
)

/* I'm not even sure yet how comments work in the Go community. But this is my first project.
It persists data in a MySQl server, which is aggregated from BoxOfficeMojo.com.

Author: Daniel Deneau

Eventually, other sites should be added to build a more robust database for ML models.
(scrape multiple sites in different threads)

*/

var siteBody = "" /* Just a practice global variable to store the siteBody text.
Not needed, we can just pass around in function arguments.*/

// (found on zetcode.com/golang) cycles through attributes of an HTML element to check if the desired one exists.
func getAttribute(n *html.Node, key string) (string, bool) {
	for _, attr := range n.Attr {
		if attr.Val == key {
			return attr.Val, true
		}
	}

	return "", false
}

/* (found on zetcode.com/golang) If the passed in HTML node is an element, it
will be tested to see if it contains the desired attribute */
func checkId(n *html.Node, id string) bool {
	if n.Type == html.ElementNode {

		stringOut, isAttribute := getAttribute(n, id)

		if isAttribute && stringOut == id {
			return false
		}
	}

	return true
}

// (from zetcode.net/golang)
// Learning comment: Takes in a pointer to a html.Node, a string, and returns a pointer to an html.Node.
// Production comment: Filters the site for only data wanted by the program.
func collectData(n *html.Node, id string) *html.Node {
	if checkId(n, id) {
		return n
	}

	for c := n.FirstChild; c != nil; c = c.NextSibling {
		current := collectData(c, id)

		if current != nil {
			return current
		}
	}
	return nil // we aren't returning anything if this function doesn't find anything.
}

/* (mostly from zetcode also)
connects to the website desired. */
func connectAndCollect() {
	resp, err := http.Get("https://boxofficemojo.com")
	if err != nil {
		log.Fatalln(err)
	}

	//We Read the response body on the line below.
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatalln(err)
	}

	//Convert the body to type string
	sb := string(body)
	siteBody = sb

	doc, err := html.Parse(strings.NewReader(sb))

	if err != nil {
		fmt.Println("HTML Parse Error.")
	}

	collectData(doc, "a-link-normal")
}

/* (once again, thank you zetcode.net/golang)
 */
func parseCollectedData() {
	myTokenizer := html.NewTokenizer(strings.NewReader(siteBody))
	textOut := ""

	for {
		tkn_iter := myTokenizer.Next()
		switch {

		case tkn_iter == html.ErrorToken:
			textOut += "error token"

		case tkn_iter == html.TextToken:
			tkn := myTokenizer.Token()
			siteBody = tkn.Data
		}
	}
}

func main() {
	connectAndCollect()
	parseCollectedData()
	fmt.Println(siteBody)
}
