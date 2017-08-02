package opentable

import (
	"io"
	"net/url"
	"strings"
	"time"

	"golang.org/x/net/html"
)

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

// Parse takes in html content via an io.Reader and parses out the reservation data
func Parse(data io.Reader) (map[string]Resturant, error) {
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
