package otscraper

import (
	"errors"
	"io"
	"net/http"
	"net/url"
	"strings"
	"time"

	"golang.org/x/net/html"
)

type Resturant struct {
	ID           string
	Name         string
	Reservations []Reservation
}

type Reservation struct {
	Time string
	URL  string
}

// Extract all http** links from a given webpage
func Crawl(data io.ReadCloser) (map[string]Resturant, error) {
	defer data.Close()
	vals := map[string]Resturant{}

	var rr Resturant

	z := html.NewTokenizer(data)

	for {
		tt := z.Next()

		switch tt {
		case html.StartTagToken:
			t := z.Token()

			// Check if we entered a new resturant row
			if t.Data == "div" {
				id := getAttr(t.Attr, "data-rid")
				if id == "" {
					continue
				}
				if rr.ID != id {
					if rr.ID != "" {
						vals[rr.ID] = rr
					}
					rr = Resturant{
						ID:           id,
						Reservations: []Reservation{},
					}
					continue
				}
				rr.ID = id
			}

			// Get the resturant name
			if rr.ID != "" && t.Data == "span" && hasClass(t.Attr, "rest-row-name-text") {
				z.Next()
				rr.Name = z.Token().String()
			}

			// Get the resturant available reservations
			if rr.ID != "" && t.Data == "a" {
				for _, a := range t.Attr {
					if (a.Key == "href" || a.Key == "data-href") && strings.HasPrefix(a.Val, "/book/") {
						qs, err := url.ParseQuery(a.Val)
						if err != nil {
							return vals, err
						}

						dt, err := time.Parse("2006-01-02 15:04", qs.Get("sd"))
						if err != nil {
							return vals, err
						}

						rr.Reservations = append(rr.Reservations, Reservation{
							Time: dt.Format("03:04PM"),
							URL:  "https://opentable.com" + a.Val,
						})
					}
				}
			}
		case html.ErrorToken:
			// End of the document, we're done
			return vals, nil
		}
	}
}

func FetchData(requestURL string) (io.ReadCloser, error) {

	resp, err := http.Get(requestURL)

	if err != nil {
		return nil, errors.New("ERROR: Failed to crawl \"" + requestURL + "\"")
	}

	b := resp.Body
	return b, nil
}

func hasClass(attrs []html.Attribute, className string) bool {
	for _, a := range attrs {
		if a.Key == "class" && a.Val == className {
			return true
		}
	}
	return false
}

func getAttr(attrs []html.Attribute, attrName string) string {
	for _, a := range attrs {
		if a.Key == attrName {
			return a.Val
		}
	}
	return ""
}
