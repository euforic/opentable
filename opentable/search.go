package opentable

import (
	"time"
)

// SearchOpts is a struct that holds the values required
// to form the url for an opentable search request
type SearchOpts struct {
	UserAgent string
	People    string
	Time      time.Time
	Latitude  string
	Longitude string
	Term      string
	Sort      string
	Opts      map[string]string
}

// String takes the values from the SearchOpts and encodes them
// into a usable opentable search url string
func (s SearchOpts) String() string {
	url := baseURL + "s/?Limit=200"
	url += "&term=" + s.Term
	url += "&covers=" + s.People
	url += "&longitude=" + s.Longitude
	url += "&latitude=" + s.Latitude
	url += "&dateTime=" + s.Time.Format("2006-01-02+15:00")

	for k, v := range s.Opts {
		url += "&" + k + "=" + v
	}
	return url
}

// Search performs a search against opentable.com with the given SearchOpts
// to be encoded as query params and then parses the returned HTML response
// scraping out the resturant reservation data
func Search(opts SearchOpts) (map[string]Resturant, error) {
	var result map[string]Resturant

	res, err := fetchHTML(opts)
	if err != nil {
		return result, err
	}
	defer res.Body.Close()

	result, err = Parse(res.Body)
	if err != nil {
		return result, err
	}

	return result, nil
}
