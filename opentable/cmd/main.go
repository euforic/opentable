package main

import (
	"encoding/json"
	"errors"
	"os"

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

	// make request to opentable
	data, err := opentable.FetchData(cfg["url"], cfg["agent"])
	if err != nil {
		writeError(err)
	}

	// scrape html response for reservation results
	reservations, err := opentable.Scrape(data)
	if err != nil {
		writeError(err)
	}

	// encode results to JSON and write to stdout for node.js to read
	enc := json.NewEncoder(os.Stdout)
	enc.SetEscapeHTML(false)
	enc.Encode(reservations)
}

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
