package main

import (
	"context"
	"encoding/hex"
	"fmt"
	"log"
	"net/http"
	"strconv"

	ptypes "github.com/gogo/protobuf/types"
	pb "github.com/prysmaticlabs/prysm/proto/beacon/rpc/v1"
	pbeth "github.com/prysmaticlabs/prysm/proto/eth/v1alpha1"
	"google.golang.org/grpc"
)

func main() {
	ctx := context.Background()

	conn, err := grpc.DialContext(ctx, "127.0.0.1:4000", grpc.WithInsecure())
	if err != nil {
		panic(err)
	}

	beaconClient := pb.NewBeaconServiceClient(conn)
	ethClient := pbeth.NewBeaconChainClient(conn)

	http.HandleFunc("/bl/", func(w http.ResponseWriter, r *http.Request) {
	})

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		head, err := beaconClient.CanonicalHead(ctx, &ptypes.Empty{})
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		hash := hex.EncodeToString(head.StateRoot)
		nat := len(head.GetBody().Attestations)
		w.Write([]byte("<p><a href=\"/bl/" + hash + "\">" + hash + "</a> - " + fmt.Sprintf("%d", nat) + " attestations</p>"))

		blockResp, err := ethClient.ListBlocks(ctx, &pbeth.ListBlocksRequest{
			QueryFilter: &pbeth.ListBlocksRequest_Root{Root: head.ParentRoot},
			PageToken:   strconv.Itoa(0),
		}, &grpc.EmptyCallOption{})
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		block := blockResp.Blocks[0]
		hash = hex.EncodeToString(block.StateRoot)
		nat = len(block.GetBody().Attestations)
		w.Write([]byte("<p><a href=\"/bl/" + hash + "\">" + hash + "</a> - " + fmt.Sprintf("%d", nat) + " attestations</p>"))
	})

	log.Fatal(http.ListenAndServe(":8088", nil))
}
