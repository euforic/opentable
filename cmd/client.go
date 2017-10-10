package cmd

import (
	"log"

	"github.com/euforic/opentable/otpb"
	"github.com/spf13/cobra"
	"google.golang.org/grpc"
)

var clientHost string

var clientCmd = &cobra.Command{
	Use:   "client",
	Short: "Run the gRPC client",
}

func init() {
	RootCmd.AddCommand(clientCmd)
	clientCmd.Flags().StringVar(&clientHost, "host", ":8080", "app server host")
}

func client() otpb.OTServiceClient {
	var opts []grpc.DialOption
	opts = append(opts, grpc.WithInsecure())

	conn, err := grpc.Dial(clientHost, opts...)
	if err != nil {
		log.Fatal(err)
	}
	return otpb.NewOTServiceClient(conn)
}
