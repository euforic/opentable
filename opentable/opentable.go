package opentable

import (
	"errors"
	"io"
	"net/http"
	"net/url"
	"strings"
	"time"

	"golang.org/x/net/html"
)

const defaultUserAgent = "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_12_6) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/59.0.3071.115 Safari/537.36"

// Resturant struct that holds resturant data and reservation list
type Resturant struct {
	ID           string
	Name         string
	Reservations []Reservation
}

// Reservation struct that holds reservation data
type Reservation struct {
	Time time.Time
	URL  string
}

// Scrape takes in html content via an io.ReadCloser and parses out the reservation data
func Scrape(data io.ReadCloser) (map[string]Resturant, error) {
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
							Time: dt,
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

// FetchData makes a request to opentable and returns the results html
func FetchData(requestURL string, userAgent string) (io.ReadCloser, error) {
	// set custom user-agent if provided
	if userAgent == "" {
		userAgent = defaultUserAgent
	}

	client := &http.Client{
		Timeout: time.Second * 10,
	}

	req, err := http.NewRequest("GET", requestURL, nil)
	if err != nil {
		return nil, errors.New("ERROR: Failed to create HTTP Request")
	}

	// set user-agent for http request
	req.Header.Set("User-Agent", userAgent)
	// make http request
	resp, err := client.Do(req)

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
