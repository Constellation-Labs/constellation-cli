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

func printableNodeStatus(metrics *node.Metrics) string {
	if metrics == nil {
		return fmt.Sprintf("/"+statusColorFmt(node.Offline), node.Offline)
	}

	return fmt.Sprintf("/"+statusColorFmt(metrics.NodeState), metrics.NodeState)
}

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
	case node.Observing:
		return ObservingColor
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
	case node.Observing:
		return `oo`
	default:
		return `~~`
	}
}

func symbol(status node.NodeState) string {
	return fmt.Sprintf(statusColorFmt(status), statusSymbol(status))
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

func operatorName(ni NodeOverview) string {
	if ni.Operator != nil {
		return fmt.Sprintf("%s <%s>", ni.Operator.Name, ni.Operator.DiscordId)
	}

	return ""
}

func PrintAsciiOutput(clusterOverview []NodeOverview, grid map[string]map[string]*node.PeerInfo, verbose bool) {

	fmt.Printf("Constellation Hypergraph Network nodes [%d], majority status\n", len(clusterOverview))

	if verbose {
		fmt.Printf("\u001B[1;35m##  %-129s %-20s %-40s %-21s %-10s %-10s %-10s %s\u001B[0m\n", "Id", "Alias", "Operator", "Address", "Version", "Snapshot", "Latency", "Status Lb/Node")
	} else {
		fmt.Printf("\u001B[1;35m##  %-20s %-21s %-10s %-10s %-10s %s\u001B[0m\n", "Alias", "Address", "Version", "Snapshot", "Latency", "Status Lb/Node")
	}

	for i, nodeOverview := range clusterOverview {
		var version = "?"
		var snap = "?"
		var latency = "♾"

		if nodeOverview.Metrics != nil {
			version = nodeOverview.Metrics.Version
			snap = nodeOverview.Metrics.LastSnapshotHeight
		}

		latency = fmtLatencyAscii(nodeOverview.AvgResponseDuration)

		if verbose {
			fmt.Printf("\u001B[1;36m%02d\u001B[0m  %-129s %-20s %-40s %-21s %-10s %-10s %-10s %s%s\n",
				i,
				nodeOverview.SelfInfo.Id,
				nodeOverview.SelfInfo.Ip, // TODO: replace with alias if available
				operatorName(nodeOverview),
				fmt.Sprintf("%s:%d", nodeOverview.SelfInfo.Ip, nodeOverview.SelfInfo.PublicPort),
				version,
				snap,
				latency,
				fmt.Sprintf(statusColorFmt(nodeOverview.SelfInfo.CardinalState()), nodeOverview.SelfInfo.CardinalState()), // TODO: no status in peer info
				printableNodeStatus(nodeOverview.Metrics))
		} else {
			selfState := node.Undefined
			if nodeOverview.SelfInfo != nil {
				selfState = nodeOverview.SelfInfo.CardinalState()
			}

			fmt.Printf("\u001B[1;36m%02d\u001B[0m  %-20s %-21s %-10s %-10s %-10s %s%s\n",
				i,
				nodeOverview.Ip,
				fmt.Sprintf("%s:%d", nodeOverview.Ip, nodeOverview.PublicPort),
				version,
				snap,
				latency,
				fmt.Sprintf(statusColorFmt(selfState), selfState),
				printableNodeStatus(nodeOverview.Metrics))
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
			if cell := rowMap[colNode.Id]; cell != nil {
				cardinalState = cell.CardinalState()
			}

			fmt.Printf(" %s", symbol(cardinalState))
		}

		fmt.Printf("\n")
	}
}
