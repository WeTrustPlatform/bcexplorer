package main

import (
	"context"
	"encoding/hex"
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

	http.HandleFunc("/bl/", func(w http.ResponseWriter, r *http.Request) {
	})

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		head, err := beaconClient.CanonicalHead(ctx, &ptypes.Empty{})
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		w.Write([]byte("<h1>Latest blocks</h1>"))
		hash := hex.EncodeToString(head.StateRoot)
		w.Write([]byte("<a href=\"/bl/" + hash + "\">" + hash + "</a>"))
	})

	log.Fatal(http.ListenAndServe(":8088", nil))
}
