package opentable

import (
	"bytes"
	"encoding/json"
	"flag"
	"io/ioutil"
	"net/http"
	"testing"
	"time"

	"github.com/google/go-cmp/cmp"
)

var update = flag.Bool("update", false, "update golden files")

func TestParse(t *testing.T) {
	tests := []struct {
		name    string
		want    map[string]Resturant
		wantErr bool
	}{
		{"basic", map[string]Resturant{}, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if *update {
				updateGoldenFiles(t, tt.name)
			}
			reader := bytes.NewReader(readGoldenFile(t, tt.name+".html"))
			got, err := Parse(reader)
			if (err != nil) != tt.wantErr {
				t.Errorf("Parse() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			wantFile := readGoldenFile(t, tt.name+".json")
			if err := json.Unmarshal(wantFile, &tt.want); err != nil {
				t.Fatalf("err: %s", err)
			}

			if diff := cmp.Diff(got, tt.want); diff != "" {
				t.Errorf("%s: result different: (-got +want)\n%s", tt.name, diff)
			}
		})
	}
}

// Test Helpers

func updateGoldenFiles(t *testing.T, name string) {
	t.Helper()
	var err error

	date := time.Now()
	d := date.AddDate(0, 1, 0).Format("2006-01-02")

	urls := map[string]string{
		"basic": "https://www.opentable.com/s/?covers=2&dateTime=" + d + "%2019%3A00&latitude=33.611674&longitude=-117.74882&metroId=496&term=italian&enableSimpleCuisines=true&freetext%5BLimit%5D=200&pageType=0",
	}

	resp, err := http.Get(urls[name])
	if err != nil {
		t.Fatalf("err: %s", err)
	}
	defer resp.Body.Close()

	data, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		t.Fatalf("err: %s", err)
	}

	writeGoldenFile(t, name+".html", data)

	res, err := Parse(bytes.NewReader(data))
	if err != nil {
		t.Fatalf("err: %s", err)
	}
	out, err := json.MarshalIndent(res, "", "  ")
	if err != nil {
		t.Fatalf("err: %s", err)
	}
	writeGoldenFile(t, name+".json", out)
	return

}

func writeGoldenFile(t *testing.T, name string, data []byte) {
	t.Helper()
	err := ioutil.WriteFile("test-fixtures/"+name+".golden", data, 0644)
	if err != nil {
		t.Fatalf("err: %s", err)
	}
}

func readGoldenFile(t *testing.T, name string) []byte {
	t.Helper()
	file, err := ioutil.ReadFile("test-fixtures/" + name + ".golden")
	if err != nil {
		t.Fatalf("err: %s", err)
	}

	return file
}
