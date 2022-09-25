package nodegrid

import (
	"constellation/pkg/node"
	"fmt"
	"time"
)

const (
	OperationalColor = "\033[1;92m%s\033[0m"
	WarningColor     = "\033[1;33m%s\033[0m"
	OfflineColor     = "\033[1;31m%s\033[0m"
	UnknownColor     = "\033[1;34m%s\033[0m"
	UndefinedColor   = "\033[1;31m%s\033[0m"
	ObservingColor   = "\033[1;36m%s\033[0m"
)

// TODO: Verify if this is based on lb or not
func printableNodeStatus(no NodeOverview) string {
	//if no.SelfInfo == nil {
	//	return fmt.Sprintf("/"+StatusColorFmt(node.Offline), node.Offline)
	//}

	return fmt.Sprintf("/"+StatusColorFmt(no.LbInfo.CardinalState()), no.LbInfo.CardinalState())
}

func StatusColorFmt(status node.NodeState) string {

	switch status {
	case node.Initial:
		return OfflineColor
	case node.ReadyToJoin:
		return WarningColor
	case node.WaitingForDownload:
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
	case node.Observing:
		return ObservingColor
	case node.Undefined:
		return UndefinedColor
	case node.NotSupported:
		return UnknownColor
	default:
		return UnknownColor
	}
}

func StatusSymbol(status node.NodeState) string {
	switch status {
	case node.Observing:
		return `oo`
	case node.Initial:
		return `==`
	case node.ReadyToJoin:
		return `==`
	case node.WaitingForDownload:
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
	case node.NotSupported:
		return `∎∎`
	case node.Undefined:
		return `~~`
	default:
		return `~~`
	}
}

func symbol(status node.NodeState) string {
	return fmt.Sprintf(StatusColorFmt(status), StatusSymbol(status))
}

func fmtLatency(d time.Duration) string {

	if d.Seconds() >= 1 {
		return fmt.Sprintf("%.3f[s]", d.Seconds())
	}

	return fmt.Sprintf("%d[ms]", d.Milliseconds())
}

func fmtLatencyAscii(d time.Duration) string {
	var lat = fmtLatency(d)

	if d.Milliseconds() >= LatencyTriggerMilliseconds {
		lat = fmt.Sprintf("\033[1;33m%-10s\033[0m", lat)
	} else {
		lat = fmt.Sprintf("\033[0;37m%-10s\033[0m", lat)
	}

	return lat
}

func PrintAsciiOutput(clusterOverview []NodeOverview, grid map[string]map[string]node.PeerInfo, verbose bool) {

	fmt.Printf("Constellation Hypergraph Network nodes [%d], majority status\n", len(clusterOverview))

	if verbose {
		fmt.Printf("\u001B[1;35m## %9s  %-20s %-10s \u001B[0m\n", "Id", "Address", "Status Lb/Node")
	} else {
		fmt.Printf("\u001B[1;35m## %9s %-20s %s\u001B[0m\n", "Id", "Address", "Status Lb/Node")
	}

	for i, nodeOverview := range clusterOverview {

		if verbose {
			fmt.Printf("\u001B[1;36m%02d\u001B[0m  %-9s %-20s %-21s %s%s\n",
				i,
				nodeOverview.LbInfo.ShortId(),
				nodeOverview.SelfInfo.Ip, // TODO: replace with alias if available
				fmt.Sprintf("%s:%d", nodeOverview.SelfInfo.Ip, nodeOverview.SelfInfo.PublicPort),

				fmt.Sprintf(StatusColorFmt(nodeOverview.SelfInfo.CardinalState()), nodeOverview.SelfInfo.CardinalState()), // TODO: no status in peer info
				printableNodeStatus(nodeOverview))
		} else {
			selfState := node.Undefined
			if nodeOverview.SelfInfo != nil {
				selfState = nodeOverview.SelfInfo.CardinalState()
			}

			fmt.Printf("\u001B[1;36m%02d\u001B[0m  %-9s %-21s %s%s\n",
				i,
				nodeOverview.LbInfo.ShortId(),
				fmt.Sprintf("%s:%d", nodeOverview.Ip, nodeOverview.PublicPort),
				fmt.Sprintf(StatusColorFmt(selfState), selfState),
				printableNodeStatus(nodeOverview))
		}
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
			if cell, ok := rowMap[colNode.Id]; ok {
				cardinalState = cell.CardinalState()
			}

			fmt.Printf(" %s", symbol(cardinalState))
		}

		fmt.Printf("\n")
	}
}
