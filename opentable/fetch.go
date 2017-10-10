package opentable

import (
	"net/http"
	"time"
)

const defaultUserAgent = "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_12_6) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/59.0.3071.115 Safari/537.36"
const baseURL = "https://www.opentable.com/"

// fetchHTML forms the request url based on the given options
// then makes a request and returns a pointer to the http.Response
func fetchHTML(opts SearchOpts) (*http.Response, error) {
	// set custom user-agent if provided
	userAgent := opts.UserAgent

	if userAgent == "" {
		userAgent = defaultUserAgent
	}

	client := &http.Client{Timeout: time.Second * 10}

	req, err := http.NewRequest("GET", opts.String(), nil)
	if err != nil {
		return nil, err
	}

	// set user-agent for http request
	req.Header.Set("User-Agent", userAgent)
	// make http request
	res, err := client.Do(req)
	return res, err
}
