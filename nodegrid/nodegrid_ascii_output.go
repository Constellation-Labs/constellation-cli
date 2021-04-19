package nodegrid

import (
	"constellation_cli/pkg/node"
	"fmt"
	"time"
)

const (
	OperationalColor = "\033[1;92m%s\033[0m"
	WarningColor = "\033[1;33m%s\033[0m"
	OfflineColor = "\033[1;31m%s\033[0m"
	WorkingColor = "\033[1;36m%s\033[0m"
	UnknownColor = "\033[1;37m%s\033[0m"
)

func statusColorFmt(status node.NodeStatus) string {
	switch status {
	case node.DownloadCompleteAwaitingFinalSync:
		return WarningColor
	case node.ReadyForDownload:
		return WarningColor
	case node.DownloadInProgress:
		return WarningColor
	case node.PendingDownload:
		return WarningColor
	case node.Leaving:
		return OfflineColor
	case node.Offline:
		return OfflineColor
	case node.SnapshotCreation:
		return WorkingColor
	case node.Ready:
		return OperationalColor
	default:
		return UnknownColor
	}
}

func statusSymbol(status node.NodeStatus) string {
	switch status {
	case node.DownloadCompleteAwaitingFinalSync:
		return `∎∎`
	case node.ReadyForDownload:
		return `∎∎`
	case node.DownloadInProgress:
		return `∎∎`
	case node.PendingDownload:
		return `∎∎`
	case node.Ready:
		return `■■`
	case node.Leaving:
		return `==`
	case node.SnapshotCreation:
		return `■■`
	case node.Offline:
		return `--`
	default:
		return `~~`
	}
}

func symbol(status node.NodeStatus) string {
	return fmt.Sprintf(statusColorFmt(status), statusSymbol(status))
}

func printableNodeStatus(metrics *node.Metrics) string {
	if metrics == nil {
		return fmt.Sprintf("/"+statusColorFmt(node.Offline), node.Offline)
	}

	return fmt.Sprintf("/"+statusColorFmt(metrics.NodeState), metrics.NodeState)
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

func PrintAsciiOutput(clusterOverview []NodeOverview, grid map[string]map[string]node.NodeInfo, verbose bool) {

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
			latency = fmtLatencyAscii(nodeOverview.AvgResponseDuration)
		}

		if verbose {
			fmt.Printf("\u001B[1;36m%02d\u001B[0m  %-129s %-20s %-40s %-21s %-10s %-10s %-10s %s%s\n",
				i,
				nodeOverview.Info.Id.Hex,
				nodeOverview.Info.Alias,
				operatorName(nodeOverview),
				fmt.Sprintf("%s:%d", nodeOverview.Info.Ip.Host, nodeOverview.Info.Ip.Port),
				version,
				snap,
				latency,
				fmt.Sprintf(statusColorFmt(nodeOverview.Info.Status), nodeOverview.Info.Status),
				printableNodeStatus(nodeOverview.Metrics))
		} else {
			fmt.Printf("\u001B[1;36m%02d\u001B[0m  %-20s %-21s %-10s %-10s %-10s %s%s\n",
				i,
				nodeOverview.Info.Alias,
				fmt.Sprintf("%s:%d", nodeOverview.Info.Ip.Host, nodeOverview.Info.Ip.Port),
				version,
				snap,
				latency,
				fmt.Sprintf(statusColorFmt(nodeOverview.Info.Status), nodeOverview.Info.Status),
				printableNodeStatus(nodeOverview.Metrics))
		}
	}

	fmt.Println()
	fmt.Println()

	fmt.Println("Legend")
	fmt.Print("   ")

	for i, status := range node.ValidStatuses {
		fmt.Printf("%s %-35s   ", symbol(status), status)
		if (i+1)%3 == 0 {
			fmt.Print("\n   ")
		}
	}

	fmt.Println()
	fmt.Println()

	fmt.Print("  ")
	for i, _ := range clusterOverview {
		fmt.Printf(" %02d", i)
	}

	fmt.Println()

	for i, rowNode := range clusterOverview {
		fmt.Printf("%02d", i)

		rowMap := grid[rowNode.Info.Ip.Host]

		for _, colNode := range clusterOverview {
			cell := rowMap[colNode.Info.Ip.Host]
			fmt.Printf(" %s", symbol(cell.Status))
		}

		fmt.Printf("\n")
	}
}