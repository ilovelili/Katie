package main

import (
	"fmt"
	"html/template"
	"katie/ipfs"
	"log"
	"net/http"
	"path"

	"github.com/ipfs/go-ipfs/core"
	"github.com/julienschmidt/httprouter"
)

// IPFSHandler IPFSHandler
type IPFSHandler struct {
	node *core.IpfsNode
}

// Page view model passing to front-end
type Page struct {
	Title  string
	Tweet  []string
	isMine bool
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

	log.Fatal(http.ListenAndServe(":8080", router))
}

// TextInput Called on all profile pages. Fills the profile page with tweets for the relevant user.
func TextInput(node *core.IpfsNode) httprouter.Handle {
	return func(w http.ResponseWriter, r *http.Request, params httprouter.Params) {
		userID := params.ByName("name")
		var page Page

		if userID == "" {
			// [1] If its your home profile page
			tweetArray, err := ipfs.GetStrings(node, node.Identity.Pretty())
			if err != nil {
				panic(err)
			}

			fmt.Printf("the tweet array is %s\n", tweetArray)
			//[1A] If no tweets, send nil
			if tweetArray == nil {
				fmt.Println("tweetarray is nil")
				page = Page{"Decentralized Twitter", nil, true}
			} else {
				fmt.Println("tweetarray is not nil")
				page = Page{"Decentralized Twitter", tweetArray, true}
			}

			fp := path.Join("templates", "index.html")
			tmpl, err := template.ParseFiles(fp)

			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}

			if err := tmpl.Execute(w, page); err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
			}
		} else {
			// //[2] If its another profile Pull from IPNS
			tweetArray, err := ipfs.GetStrings(node, userID)
			if err != nil {
				fmt.Println("tweeterarray is nil")
				page = Page{"Decentralized Twitter", nil, false}
			} else {
				fmt.Println("RESOLVED")
				if err != nil {
					panic(err)
				}

				page = Page{"Decentralized Twitter", tweetArray, false}
			}

			fp := path.Join("templates", "index.html")
			tmpl, err := template.ParseFiles(fp)
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
			if err := tmpl.Execute(w, page); err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
			}
		}
	}
}
