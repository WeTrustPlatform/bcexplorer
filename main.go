package main

import (
	"context"
	"log"
	"net/http"

	ptypes "github.com/gogo/protobuf/types"
	pb "github.com/prysmaticlabs/prysm/proto/beacon/rpc/v1"
	"google.golang.org/grpc"
)

func main() {
	ctx := context.Background()

	conn, err := grpc.DialContext(ctx, "127.0.0.1:4000", grpc.WithInsecure())
	if err != nil {
		panic(err)
	}

	beaconClient := pb.NewBeaconServiceClient(conn)

	head, err := beaconClient.CanonicalHead(ctx, &ptypes.Empty{})
	if err != nil {
		panic(err)
	}

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("<h1>Latest blocks</h1>"))
		w.Write(head.RandaoReveal)
	})

	log.Fatal(http.ListenAndServe(":8088", nil))
}
