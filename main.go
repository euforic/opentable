package main

import (
	"encoding/base64"
	"encoding/json"
	"os"

	"github.com/bleveinc/planz/otscraper"
)

func main() {
	seedUrlBase64 := os.Args[1]

	seedUrl, err := base64.StdEncoding.DecodeString(seedUrlBase64)
	if err != nil {
		panic(err)
	}
	data, _ := otscraper.FetchData(string(seedUrl), "")

	reservations, _ := otscraper.Scrape(data)

	enc := json.NewEncoder(os.Stdout)
	enc.SetEscapeHTML(false)
	enc.SetIndent("", "  ")
	enc.Encode(reservations)
}
