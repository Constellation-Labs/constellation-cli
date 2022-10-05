package status

import (
	"constellation/pkg/lb"
	"constellation/pkg/node"
	"fmt"
	"golang.org/x/exp/slices"
)

type Status interface {
	ProvideAsciiStatus(nodeId string)
}

type status struct {
	lb string
}

func NewStatus(url string) Status {
	return &status{url}
}

func (s *status) ProvideAsciiStatus(nodeId string) {

	nodes, err := lb.GetClient(s.lb).GetClusterNodes()

	if err != nil {

		fmt.Println("🖤")
		return
	}

	nodesArray := make([]node.PeerInfo, len(*nodes))

	for i, n := range *nodes {
		nodesArray[i] = n
	}

	nodeIndex := slices.IndexFunc(nodesArray, func(p node.PeerInfo) bool { return p.Id == nodeId })

	if nodeIndex == -1 {
		fmt.Println("💔")
		fmt.Println("---")
	} else {
		switch nodeState := nodesArray[nodeIndex].State; nodeState {
		case "Ready":
			fmt.Println("💚")
		case "Observing":
			fmt.Println("💜")
		default:
			fmt.Println("💙")
		}
		fmt.Println("---")

		fmt.Printf("Your node is %s\n", nodesArray[nodeIndex].State)
	}

	fmt.Printf("%d L0 nodes\n", len(*nodes))
}
