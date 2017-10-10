package cmd

import (
	"context"
	"encoding/json"
	"log"
	"os"
	"strings"
	"time"

	"github.com/euforic/opentable/otpb"
	"github.com/euforic/opentable/otserver"
	"github.com/spf13/cobra"
)

var searchOpts = struct {
	Term       string
	Latitude   string
	Longitude  string
	People     string
	DateTime   string
	Sort       string
	Agent      string
	DirectCall bool
	Opts       []string
}{}

var searchCmd = &cobra.Command{
	Use:   "search",
	Short: "Search for reservations",
	Run: func(cmd *cobra.Command, args []string) {
		t := time.Now()

		req := &otpb.SearchReq{
			People:    searchOpts.People,
			Time:      &t,
			Latitude:  searchOpts.Latitude,
			Longitude: searchOpts.Longitude,
			Term:      searchOpts.Term,
		}

		opts := map[string]string{}

		// convert other args passed in as a `[]string{"KEY=VALUE",...}`
		// into `map[string]string{"KEY": "VALUE", ...}`
		for _, o := range searchOpts.Opts {
			op := strings.Split(o, "=")
			if len(opts) < 2 {
				continue
			}
			opts[op[0]] = op[1]
		}

		req.Opts = opts

		server := otserver.New()

		var r *otpb.SearchRes
		var err error

		if searchOpts.DirectCall {
			r, err = server.Search(context.Background(), req)
		} else {
			r, err = client().Search(context.Background(), req)
		}

		if err != nil {
			log.Fatal(r, err)
		}
		enc := json.NewEncoder(os.Stdout)
		enc.SetEscapeHTML(false)
		enc.SetIndent("", "  ")
		enc.Encode(r)
	},
}

func init() {
	clientCmd.AddCommand(searchCmd)
	searchCmd.Flags().StringVar(&searchOpts.People, "people", "2", "Number of people attending")
	searchCmd.Flags().StringVar(&searchOpts.Term, "term", "", "Term to search for")
	searchCmd.Flags().StringVar(&searchOpts.DateTime, "date_time", "", "Date / Time for reservation")
	searchCmd.Flags().StringVar(&searchOpts.Latitude, "latitude", "", "Latitude to search around")
	searchCmd.Flags().StringVar(&searchOpts.Longitude, "longitude", "", "Longitude to search around")
	searchCmd.Flags().StringVar(&searchOpts.Agent, "agent", "", "User-Agent string to send request as")
	searchCmd.Flags().StringVar(&searchOpts.Sort, "sort", "", "Sort by RATING, DISTANCE, POPULARITY or NAME. Default is RATING")
	searchCmd.Flags().BoolVar(&searchOpts.DirectCall, "direct", false, "Call the server function directly wo/ starting a server")
	searchCmd.Flags().StringSliceVar(&searchOpts.Opts, "opts", []string{}, "Additional call options")
}
