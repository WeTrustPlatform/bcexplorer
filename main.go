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

func getBlock(ctx context.Context, client pbeth.BeaconChainClient, root []byte) *pbeth.BeaconBlock {
	resp, err := client.ListBlocks(ctx, &pbeth.ListBlocksRequest{
		QueryFilter: &pbeth.ListBlocksRequest_Root{Root: root},
		PageToken:   strconv.Itoa(0),
	}, &grpc.EmptyCallOption{})
	if err != nil {
		panic(err)
	}

	return resp.Blocks[0]
}

func displayBlock(w http.ResponseWriter, block *pbeth.BeaconBlock) {
	root := hex.EncodeToString(block.StateRoot)
	nat := len(block.GetBody().Attestations)
	w.Write([]byte("<p><a href=\"/block?root=" + root + "\">" + root + "</a> - " + fmt.Sprintf("%d", nat) + " attestations</p>"))
}

func main() {
	ctx := context.Background()

	conn, err := grpc.DialContext(ctx, "127.0.0.1:4000", grpc.WithInsecure())
	if err != nil {
		panic(err)
	}

	beaconClient := pb.NewBeaconServiceClient(conn)
	ethClient := pbeth.NewBeaconChainClient(conn)

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		block, err := beaconClient.CanonicalHead(ctx, &ptypes.Empty{})
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		w.Write([]byte("<h1>Latest blocks</h1>"))

		for i := 0; i < 20 && block != nil; i++ {
			displayBlock(w, block)
			fmt.Println(block.StateRoot, hex.EncodeToString(block.StateRoot))
			block = getBlock(ctx, ethClient, block.ParentRoot)
		}
	})

	http.HandleFunc("/block", func(w http.ResponseWriter, r *http.Request) {
		rootStrings, _ := r.URL.Query()["root"]
		w.Write([]byte("<h1>" + rootStrings[0] + "</h1>"))

		root, _ := hex.DecodeString(rootStrings[0])

		fmt.Println(">>", root, rootStrings[0])
		block := getBlock(ctx, ethClient, root)
		fmt.Println(block)

		// w.Write([]byte("<p>Signature: " + hex.EncodeToString(block.Signature) + "</p>"))
		return
	})

	log.Fatal(http.ListenAndServe(":8088", nil))
}
