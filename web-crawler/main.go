package main

import (
	"fmt"
	"net/http"
	"os"
	"strings"

	log "github.com/llimllib/loglevel"
	"golang.org/x/net/html"
	"golang.org/x/net/html/atom"
)

// MaxDepth ...
var MaxDepth = 2

// Link ...
type Link struct {
	url   string
	text  string
	depth int
}

// HTTPError ...
type HTTPError struct {
	original string
}

// LinkReader ...
func LinkReader(res *http.Response, depth int) []Link {
	page := html.NewTokenizer(res.Body)
	links := []Link{}

	var start *html.Token
	var tokenData string

	for {
		_ = page.Next()
		token := page.Token()

		if token.Type == html.ErrorToken {
			break
		}

		if start != nil && token.Type == html.TextToken {
			tokenData = fmt.Sprintf("%s%s", tokenData, token.Data)
		}

		if token.DataAtom == atom.A {
			switch token.Type {
			case html.StartTagToken:
				if len(token.Attr) > 0 {
					start = &token
				}
			case html.EndTagToken:
				if start == nil {
					log.Warnf("Link End found without Start: %s", tokenData)
					continue
				}
				link := NewLink(*start, tokenData, depth)
				if link.isValid() {
					links = append(links, link)
					log.Debugf("Link Found %v", link)
				}

				start = nil
				tokenData = ""
			}
		}
	}

	log.Debug(links)

	return links
}

// NewLink ...
func NewLink(tag html.Token, text string, depth int) Link {
	link := Link{text: strings.TrimSpace(text), depth: depth}

	for i := range tag.Attr {
		if tag.Attr[i].Key == "href" {
			link.url = strings.TrimSpace(tag.Attr[i].Val)
		}
	}

	return link
}

func (link Link) String() string {
	spacer := strings.Repeat("\t", link.depth)
	formatted := fmt.Sprintf("%s%s (%d) - %s", spacer, link.text, link.depth, link.url)

	return formatted
}

func (link Link) isValid() bool {
	if link.depth >= MaxDepth {
		return false
	}

	if len(link.text) == 0 {
		return false
	}

	if len(link.url) == 0 || strings.Contains(strings.ToLower(link.url), "javascript") {
		return false
	}

	return true
}

func (error HTTPError) Error() string {
	return error.original
}

func recurDownloader(url string, depth int) {
	log.Infoln("\n")
	page, err := downloader(url)

	if err != nil {
		log.Error(err)
		return
	}

	links := LinkReader(page, depth)

	for _, link := range links {
		fmt.Println(link)
		if depth+1 < MaxDepth {
			recurDownloader(link.url, depth+1)
		}
	}
}

func downloader(url string) (resp *http.Response, err error) {
	log.Debugf("Downloading: %s", url)
	resp, err = http.Get(url)
	if err != nil {
		log.Debugf("Error: %s", err)
		return
	}

	if resp.StatusCode > 299 {
		err = HTTPError{fmt.Sprintf("Error (%d): %s", resp.StatusCode, url)}
		log.Debug(err)
		return
	}

	return
}

func main() {
	log.SetPriorityString("info")
	log.SetPrefix("crawler ")

	log.Debug(os.Args)

	if len(os.Args) < 2 {
		log.Fatalln("Missing `Url` argument!")
	}

	recurDownloader(os.Args[1], 0)
}
