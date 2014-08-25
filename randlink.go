package main

import (
    "fmt"
	"io/ioutil"
	"log"
	"math/rand"
	"net/http"
	"os"
	"regexp"
	"strings"
	"time"

	"code.google.com/p/go.net/html"
)

func FilterLinks(links []string, filter string) []string {
	var filteredLinks []string
	absoluteUrl := regexp.MustCompile(`^(http|https)://`)
	domain := regexp.MustCompile("^(http|https)://"+filter)

	for _, link := range links {
		if absoluteUrl.MatchString(link) &&
			domain.MatchString(link) == false {
			filteredLinks = append(filteredLinks, link)
		}
	}

	return filteredLinks
}

func GetLinks(url string) []string {
	res, err := http.Get(url)
	if err != nil {
		log.Fatal(err)
	}

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		log.Fatal(err)
	}

	doc, err := html.Parse(strings.NewReader(string(body)))
	if err != nil {
		log.Fatal(err)
	}

	var links []string
	var f func(*html.Node)
	f = func(n *html.Node) {
		if n.Type == html.ElementNode && n.Data == "a" {
			for _, a := range n.Attr {
				if a.Key == "href" {
					links = append(links, a.Val)
					break
				}
			}
		}
		for c := n.FirstChild; c != nil; c = c.NextSibling {
			f(c)
		}
	}
	f(doc)

	return links
}

func Index(w http.ResponseWriter, req *http.Request) {
    places := `<a href="hackernews">hacker news</a>
        <br/><a href="lobsters">lobsters</a>`

    fmt.Fprintf(w, places)
}

func RandomLobsters(w http.ResponseWriter, req *http.Request) {
	links := GetLinks("http://lobste.rs")
	nonLobstersLinks := FilterLinks(links, "lobste.rs")

	rand.Seed(time.Now().UnixNano())
	randomLobstersLink := nonLobstersLinks[rand.Intn(len(nonLobstersLinks))]

	http.Redirect(w, req, randomLobstersLink, 307)
}

func RandomHackerNews(w http.ResponseWriter, req *http.Request) {
	links := GetLinks("http://news.ycombinator.com")
	nonHNLinks := FilterLinks(links, "news.ycombinator.com")

	rand.Seed(time.Now().UnixNano())
	randomHNLink := nonHNLinks[rand.Intn(len(nonHNLinks))]

	http.Redirect(w, req, randomHNLink, 307)
}

func main() {

	http.HandleFunc("/", Index)
	http.HandleFunc("/hackernews", RandomHackerNews)
	http.HandleFunc("/lobsters", RandomLobsters)
	err := http.ListenAndServe(":"+os.Getenv("PORT"), nil)
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}

}
