package cmd

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net"
	"net/http"
	"strings"
	"time"

	"github.com/euforic/opentable/otpb"
	"github.com/euforic/opentable/otserver"
	"github.com/spf13/cobra"
	"google.golang.org/grpc"
)

var serverOpts = struct {
	Port    string
	UseRest bool
}{}

var serverCmd = &cobra.Command{
	Use:   "server",
	Short: "Run the gRPC server",
	Run: func(cmd *cobra.Command, args []string) {
		if serverOpts.UseRest {
			startRest()
			return
		}
		startGRPC()
	},
}

func init() {
	RootCmd.AddCommand(serverCmd)
	serverCmd.Flags().StringVarP(&serverOpts.Port, "port", "p", "8080", "app server port")
	serverCmd.Flags().BoolVarP(&serverOpts.UseRest, "rest", "r", false, "Start Rest API server instead of gRPC Server")
}

func startRest() {
	server := otserver.New()

	http.HandleFunc("/search", func(w http.ResponseWriter, r *http.Request) {
		query := r.URL.Query()

		qTime := query.Get("date_time")
		// hacky to fix `+` getting removed from URL when not encoded
		qTime = strings.Replace(qTime, " ", "+", 1)
		if qTime == "" {
			qTime = time.Now().Format("2006-01-02+15:00")
		}

		t, err := time.Parse("2006-01-02+15:00", qTime)
		if err != nil {
			resErr(w, err)
			return
		}

		optsStr := query.Get("opts")
		opts := map[string]string{}

		if optsStr != "" {
			a := strings.Split(optsStr, ",")
			for _, kvs := range a {
				kv := strings.Split(kvs, "=")
				if len(kv) != 2 {
					continue
				}
				opts[kv[0]] = kv[1]
			}
		}

		req := otpb.SearchReq{
			People:    query.Get("people"),
			Time:      &t,
			Latitude:  query.Get("latitude"),
			Longitude: query.Get("longitude"),
			Term:      query.Get("term"),
			Opts:      opts,
		}

		sort := strings.ToUpper(query.Get("sort"))
		if val, ok := otpb.SearchReq_Sort_value[sort]; ok {
			req.Sort = otpb.SearchReq_Sort(val)
		}

		res, err := server.Search(context.Background(), &req)
		if err != nil {
			resErr(w, err)
			return
		}

		resJSON, err := json.Marshal(res)
		if err != nil {
			resErr(w, err)
			return
		}
		w.WriteHeader(http.StatusOK)
		w.Write(resJSON)
		return
	})

	fmt.Println("Starting OTServer REST Server on :", serverOpts.Port)
	log.Fatal(http.ListenAndServe(":"+serverOpts.Port, nil))
}

func startGRPC() {
	lis, err := net.Listen("tcp", fmt.Sprintf(":%s", serverOpts.Port))
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	var opts []grpc.ServerOption
	grpcServer := grpc.NewServer(opts...)
	otpb.RegisterOTServiceServer(grpcServer, otserver.New())
	fmt.Println("Starting OTServer gRPC Server on :", serverOpts.Port)
	grpcServer.Serve(lis)

}

func resErr(w http.ResponseWriter, err error) {
	w.WriteHeader(http.StatusInternalServerError)
	w.Write([]byte(err.Error()))
}
