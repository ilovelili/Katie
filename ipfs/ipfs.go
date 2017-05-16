package ipfs

import (
	"context"
	"fmt"
	"os/user"
	"path"
	"runtime"

	"github.com/ipfs/go-ipfs/core"
	"github.com/ipfs/go-ipfs/repo/fsrepo"
)

func usr() *user.User {
	usr, err := user.Current()
	if err != nil {
		panic(err)
	}

	return usr
}

func defaultPath() string {
	if runtime.GOOS == "windows" {
		return path.Join(usr().HomeDir, ".ipfs")
	}
	return "~/.ipfs"
}

// StartNode Start IPFS Node
func StartNode() (*core.IpfsNode, error) {
	// Assume the user has run 'ipfs init'
	repo := defaultPath()
	r, err := fsrepo.Open(repo)
	if err != nil {
		return nil, err
	}

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	cfg := &core.BuildCfg{
		Repo:   r,
		Online: true,
	}

	return core.NewNode(ctx, cfg)
}

// GetStrings get strings of IpfsNode
func GetStrings(node *core.IpfsNode, name string) (stringArr []string, err error) {
	stringArr = make([]string, 0)
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	path, err := node.Namesys.Resolve(ctx, name)
	if err != nil {
		return stringArr, err
	}

	nodeGetter := node.DAG
	cid, err := core.ResolveToCid(ctx, node, path)
	fmt.Println("cid is", cid)

	if err != nil {
		return stringArr, err
	}

	nd, err := nodeGetter.Get(ctx, cid)
	fmt.Println("the node is", nd)
	if err != nil {
		fmt.Println(err)
		return
	}

	fmt.Printf("bout to crash")
	fmt.Printf("%s", nd.String())
	fmt.Println("not crashed ")

	for {
		var err error
		if len(nd.Links()) == 0 {
			break
		}

		nd, err = nd.Links()[0].GetNode(ctx, nodeGetter)
		if err != nil {
			fmt.Println(err)
			break
		}

		data := nd.String()
		fmt.Println(data)
		stringArr = append(stringArr, data)
	}

	return stringArr, nil
}

/*
// AddString add input string to ipfs node
func AddString(node *core.IpfsNode, inputString string) (*cid.Cid, error) {
	pointsTo, err := node.Namesys.Resolve(node.Context(), node.Identity.Pretty())
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	//If there is an error, user is new and hasn't yet created a DAG.
	if err != nil {
		newProtoNode := makeStringNode(inputString)
		cid, err := node.DAG.Add(newProtoNode)
		if err != nil {
			return nil, err
		}

		err = node.Namesys.Publish(ctx, node.PrivateKey, path.FromCid(cid))
		if err != nil {
			return nil, err
		}

		return cid, nil
	}

	// Else user has already creatd a DAG
	newProtoNode := makeStringNode(inputString)
	cid, err := core.ResolveToCid(ctx, node, pointsTo)
	if err != nil {
		return nil, err
	}

	oldProtoNode, err := node.DAG.Get(ctx, cid)
	if err != nil {
		return nil, err
	}

	err = newProtoNode.AddNodeLink("next", oldProtoNode)
	if err != nil {
		return nil, err
	}

	node.DAG.Add(newProtoNode)
	err = node.Namesys.Publish(ctx, node.PrivateKey, pointsTo)
	if err != nil {
		return nil, err
	}

	return cid, nil
}

func makeStringNode(s string) *merkledag.ProtoNode {
	nd := new(merkledag.ProtoNode)
	nd.SetData([]byte(s))
	return nd
}
*/
