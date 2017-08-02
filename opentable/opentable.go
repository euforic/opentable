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
	URL          string
	Recommended  string
	Reservations []Reservation
}

// Reservation struct that holds reservation data
type Reservation struct {
	Time time.Time
	URL  string
}

// Scrape takes in html content via an io.ReadCloser and parses out the reservation data
func Scrape(data io.Reader) (map[string]Resturant, error) {
	vals := map[string]Resturant{}

	var rr Resturant

	z := html.NewTokenizer(data)

	for {
		tt := z.Next()

		switch tt {
		case html.StartTagToken:
			t := z.Token()

			attrMap := attrToMap(t.Attr)

			// Check if we entered a new resturant row
			if t.Data == "div" {
				id, ok := attrMap["data-rid"]
				if !ok || id == "" {
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

			if rr.ID != "" && t.Data == "span" {
				// Get the resturant name
				if hasClass(attrMap, "rest-row-name-text") {
					z.Next()
					rr.Name = z.Token().String()
					continue
				}

				// get resturant recommended percentage
				if hasClass(attrMap, "recommended-small") {
					z.Next()
					rr.Recommended = z.Token().String()
					continue
				}
			}

			// Get the resturant available reservations
			if rr.ID != "" && t.Data == "a" {

				v := attrMap["href"]

				if val, ok := attrMap["data-href"]; ok && val != "" {
					v = val
				}

				if hasClass(attrMap, "rest-row-name") {
					rr.URL = "https://opentable.com" + v
					continue
				}

				// check if the url is a booking url
				if strings.HasPrefix(v, "/book/") {
					qs, err := url.ParseQuery(v)
					if err != nil {
						return vals, err
					}

					dt, err := time.Parse("2006-01-02 15:04", qs.Get("sd"))
					if err != nil {
						return vals, err
					}

					rr.Reservations = append(rr.Reservations, Reservation{
						Time: dt,
						URL:  "https://opentable.com" + v,
					})
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

func hasClass(attrMap map[string]string, name string) bool {
	classNames, ok := attrMap["class"]
	if !ok {
		return false
	}
	c := strings.Split(classNames, " ")
	for _, n := range c {
		if n == name {
			return true
		}
	}
	return false
}

func attrToMap(attrs []html.Attribute) map[string]string {
	m := map[string]string{}
	for _, a := range attrs {
		m[a.Key] = a.Val
	}
	return m
}
