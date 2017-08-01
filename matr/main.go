package main

import (
	"context"
	"errors"
	"log"
	"os"

	"github.com/euforic/matr"
	"github.com/euforic/matr/tlkn"
)

var platforms = map[string]bool{
	"linux":  true,
	"darwin": true,
}

func main() {
	m := matr.New()

	m.Handle("build", buildHandler)

	if _, err := m.Run(context.Background(), os.Args[1:]...); err != nil {
		log.Fatal(err)
	}
}

func buildHandler(ctx context.Context) (context.Context, error) {

	c, ok := matr.ContextFrom(ctx)
	if !ok {
		return ctx, errors.New("No matr context found")
	}

	p := c.Args[0]
	if _, ok := platforms[p]; !ok {
		return ctx, errors.New("Invalid platform '" + p + "' for build. Acceptable platforms are 'linux' or 'darwin'")
	}

	tlkn.Bash("GOOS=" + p + " CGO_ENABLED=0 go build -a -ldflags '-s' -installsuffix cgo -o ./bin/opentable ./opentable/cmd/.")
	return ctx, nil
}
