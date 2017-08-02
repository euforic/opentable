package main

import (
	"encoding/json"
	"errors"
	"net/http"
	"os"
	"time"

	"github.com/bleveinc/planz/opentable"
)

func main() {

	// ensure request config data is present
	if len(os.Args) < 2 {
		writeError(errors.New("Missing config data from request"))
	}

	// get incoming request config data
	cfgJson := os.Args[1]
	cfg := map[string]string{}

	// unmarshal raw JSON data into cfg
	err := json.Unmarshal([]byte(cfgJson), &cfg)
	if err != nil {
		writeError(err)
	}

	url, ok := cfg["url"]
	if !ok {
		writeError(errors.New("A request url must be provided"))
	}

	// set custom user-agent if provided
	userAgent, ok := cfg["agent"]
	if !ok || userAgent == "" {
		userAgent = defaultUserAgent
	}

	client := &http.Client{
		Timeout: time.Second * 10,
	}

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		writeError(errors.New("ERROR: Failed to create HTTP Request"))
	}

	// set user-agent for http request
	req.Header.Set("User-Agent", userAgent)
	// make http request
	res, err := client.Do(req)
	if err != nil {
		writeError(errors.New("ERROR: Failed to crawl \"" + url + "\""))
	}
	defer res.Body.Close()

	// scrape html response for reservation results
	reservations, err := opentable.Parse(res.Body)
	if err != nil {
		writeError(err)
	}

	// encode results to JSON and write to stdout for node.js to read
	enc := json.NewEncoder(os.Stdout)
	enc.SetEscapeHTML(false)
	enc.Encode(reservations)
}

const defaultUserAgent = "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_12_6) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/59.0.3071.115 Safari/537.36"

// write error to stdout
func writeError(err error) {
	errOut := map[string]string{"error": "Go Context:" + err.Error()}
	out, err := json.Marshal(errOut)
	if err != nil {
		panic(err)
	}
	os.Stderr.Write(out)
	os.Exit(1)
}
