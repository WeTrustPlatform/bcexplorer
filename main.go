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

func displayBlock(
	ctx context.Context,
	w http.ResponseWriter,
	client pbeth.BeaconChainClient,
	block *pbeth.BeaconBlock) *pbeth.BeaconBlock {

	hash := hex.EncodeToString(block.StateRoot)
	nat := len(block.GetBody().Attestations)
	w.Write([]byte("<p><a href=\"/bl/" + hash + "\">" + hash + "</a> - " + fmt.Sprintf("%d", nat) + " attestations</p>"))

	resp, _ := client.ListBlocks(ctx, &pbeth.ListBlocksRequest{
		QueryFilter: &pbeth.ListBlocksRequest_Root{Root: block.ParentRoot},
		PageToken:   strconv.Itoa(0),
	}, &grpc.EmptyCallOption{})

	return resp.Blocks[0]
}

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
		block, err := beaconClient.CanonicalHead(ctx, &ptypes.Empty{})
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		w.Write([]byte("<h1>Latest blocks</h1>"))

		for i := 0; i < 20 && block != nil; i++ {
			block = displayBlock(ctx, w, ethClient, block)
		}
	})

	log.Fatal(http.ListenAndServe(":8088", nil))
}
