package opentable

import (
	"flag"
	"io"
	"os"
	"path/filepath"
	"reflect"
	"testing"
)

var update = flag.Bool("update", false, "update golden files")

func TestScrape(t *testing.T) {
	tests := []struct {
		name    string
		data    io.Reader
		want    map[string]Resturant
		wantErr bool
	}{
		{"basic", nil, map[string]Resturant{}, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := Scrape(tt.data)
			if (err != nil) != tt.wantErr {
				t.Errorf("Scrape() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Scrape() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestFetchData(t *testing.T) {
	urls := map[string]string{
		"basic": "https://www.opentable.com/s/?covers=2&dateTime=2017-08-01%2019%3A00&latitude=33.611674&longitude=-117.74882&metroId=496&term=italian&enableSimpleCuisines=true&freetext%5BLimit%5D=200&pageType=0",
	}
	type args struct {
		requestURL string
		userAgent  string
	}
	tests := []struct {
		name    string
		args    args
		want    io.ReadCloser
		wantErr bool
	}{
		{"basic", args{urls["basic"], ""}, nil, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := FetchData(tt.args.requestURL, tt.args.userAgent)
			defer got.Close()
			if *update {
				f, err := os.Create(filepath.Join("test-fixtures", tt.name+".golden"))
				if err != nil {
					t.Fail()
				}
				defer f.Close()
				io.Copy(f, got)
				return
			}

			if (err != nil) != tt.wantErr {
				t.Errorf("FetchData() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("FetchData() = %v, want %v", got, tt.want)
			}
		})
	}
}
