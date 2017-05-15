package main

import (
	"katie/ipfs"
	"net/http"

	"github.com/ipfs/go-ipfs/core"
	"github.com/julienschmidt/httprouter"
)

// Page view model passing to front-end
type Page struct {
	Title   string
	Author  string
	Tweet   []string
	isMine  bool
	Balance float64
}

// PeerList List of all peers in IPFS
type PeerList struct {
	AllPeers []string
}

func main() {
	node, err := ipfs.StartNode()
	if err != nil {
		panic(err)
	}

	router := httprouter.New()
	router.GET("/", TextInput(node))
}

// TextInput Called on all profile pages. Fills the profile page with tweets for the relevant user.
func TextInput(node *core.IpfsNode) httprouter.Handle {
	return func(w http.ResponseWriter, r *http.Request, params httprouter.Params) {
		// userID := params.ByName("name")
		// if userID == "" {
		// 	pointsTo, err := katja.GetDAG(node)
		// }
	}
}
