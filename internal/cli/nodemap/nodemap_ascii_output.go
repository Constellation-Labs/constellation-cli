package nodemap

import (
	"constellation/internal/cli/nodegrid"
	"constellation/pkg/node"
	"fmt"
)

func symbol(status node.NodeState) string {
	return fmt.Sprintf(nodegrid.StatusColorFmt(status), nodegrid.StatusSymbol(status))
}

func PrintAsciiOutput(clusterOverview []ClusterNode, grid map[string]map[string]*node.PeerInfo, verbose bool) {

	fmt.Printf("Cluster discovery result nodes [%d]\n", len(clusterOverview))

	fmt.Printf("\u001B[1;35m##  %-132s %-21s %s\u001B[0m\n", "Id", "Address", "Status")

	for i, nodeOverview := range clusterOverview {

		selfState := node.Undefined
		if nodeOverview.SelfInfo != nil {
			selfState = nodeOverview.SelfInfo.CardinalState()
		}

		fmt.Printf("\u001B[1;36m%02d\u001B[0m  %-132s %-21s %-10s\n",
			i,
			nodeOverview.Id,
			fmt.Sprintf("%s:%d", nodeOverview.Addr.Ip, nodeOverview.Addr.Port),
			fmt.Sprintf(nodegrid.StatusColorFmt(selfState), selfState))
	}

	fmt.Printf("\n\nLegend\n   ")
	for i, status := range node.ValidStatuses {
		fmt.Printf("%s %-35s   ", symbol(status), status)
		if (i+1)%3 == 0 {
			fmt.Print("\n   ")
		}
	}

	fmt.Printf("\n\n  ")

	for i, _ := range clusterOverview {
		fmt.Printf(" %02d", i)
	}

	fmt.Println()

	for i, rowNode := range clusterOverview {
		fmt.Printf("%02d", i)

		rowMap := grid[rowNode.Id]

		for _, colNode := range clusterOverview {

			cardinalState := node.Undefined
			if cell := rowMap[colNode.Id]; cell != nil {
				cardinalState = cell.CardinalState()
			}

			fmt.Printf(" %s", symbol(cardinalState))
		}

		fmt.Printf("\n")
	}
}
