package main

import (
	"encoding/json"
	"fmt"
	"os"
	"time"

	"github.com/bleveinc/planz/otscraper"
)

func main() {
	seedUrl := os.Args[1]

	startTime := time.Now()
	data, _ := otscraper.FetchData(seedUrl)
	fetchTime := time.Since(startTime).String()

	startScraper := time.Now()
	reservations, _ := otscraper.Crawl(data)

	enc := json.NewEncoder(os.Stdout)
	enc.SetEscapeHTML(false)
	enc.SetIndent("", "  ")
	enc.Encode(reservations)

	fmt.Println("Fetch Time: " + fetchTime)
	fmt.Println("Scrape Time: " + time.Since(startScraper).String())
	fmt.Println("Total Time: " + time.Since(startTime).String())
}
