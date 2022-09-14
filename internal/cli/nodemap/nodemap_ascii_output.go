package nodemap

import (
	"constellation/pkg/node"
	"fmt"
)

const (
	OperationalColor = "\033[1;92m%s\033[0m"
	WarningColor     = "\033[1;33m%s\033[0m"
	OfflineColor     = "\033[1;31m%s\033[0m"
	WorkingColor     = "\033[1;36m%s\033[0m"
	UnknownColor     = "\033[1;34m%s\033[0m"
	UndefinedColor   = "\033[1;31m%s\033[0m"
)

func statusColorFmt(status node.NodeState) string {

	switch status {
	case node.Initial:
		return OfflineColor
	case node.ReadyToJoin:
		return WarningColor
	case node.LoadingGenesis:
		return WarningColor
	case node.GenesisReady:
		return WarningColor
	case node.StartingSession:
		return WarningColor
	case node.SessionStarted:
		return OperationalColor
	case node.Ready:
		return OperationalColor
	case node.Leaving:
		return OfflineColor
	case node.Offline:
		return OfflineColor
	case node.Unknown:
		return UnknownColor
	case node.Undefined:
		return UndefinedColor
	default:
		return UnknownColor
	}
}

func statusSymbol(status node.NodeState) string {
	switch status {
	case node.Initial:
		return `==`
	case node.ReadyToJoin:
		return `==`
	case node.LoadingGenesis:
		return `∎∎`
	case node.GenesisReady:
		return `∎∎`
	case node.StartingSession:
		return `∎∎`
	case node.SessionStarted:
		return `■■`
	case node.Ready:
		return `■■`
	case node.Leaving:
		return `==`
	case node.Offline:
		return `--`
	case node.Unknown:
		return `∎∎`
	case node.Undefined:
		return `~~`
	default:
		return `~~`
	}
}

func symbol(status node.NodeState) string {
	return fmt.Sprintf(statusColorFmt(status), statusSymbol(status))
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
			fmt.Sprintf(statusColorFmt(selfState), selfState))
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
