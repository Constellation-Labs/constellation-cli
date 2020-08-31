package reporter

import (
	"constellation_cli/pkg/lb"
	"constellation_cli/pkg/node"
	"sort"
	"strings"
	"sync"
	"time"
)

type Nodegrid interface {
	BuildNetworkStatus(url string, silent bool, outputImage string)
}

type nodegrid struct {}

func NewNodegrid() Nodegrid {
	return & nodegrid {}
}

type nodeResult struct {
	host string
	err error
	info *node.ClusterInfo
}

func nodeInfoMap(info node.ClusterInfo) map[string]node.NodeInfo{
	m := make(map[string]node.NodeInfo)

	for _, nodeInfo := range info {
		m[nodeInfo.Ip.Host] = nodeInfo
	}

	return m
}

func clusterInfo(addr node.NodeAddr) nodeResult {
	ci, e := node.GetClient(addr).GetClusterInfo()

	return nodeResult {
		 addr.Host,
		  e,
		  ci,
	}
}

func queryNodeForClusterInfoWorker(wg *sync.WaitGroup, cluster <-chan node.NodeAddr, result chan<- nodeResult) {
	defer wg.Done()

	for addr := range cluster {
		ci := clusterInfo(addr)
		result <- ci
	}
}

func buildNetworkGrid(ci *node.ClusterInfo) map[string]map[string]node.NodeInfo{

	const workers = 24

	jobs := make(chan node.NodeAddr, len(*ci))
	results := make(chan nodeResult, len(*ci))

	var wg sync.WaitGroup

	for w := 0; w < workers; w++ {
		wg.Add(1)
		go queryNodeForClusterInfoWorker(&wg, jobs, results)
	}

	for _, nodeInfo := range *ci {
		jobs <- nodeInfo.Ip
	}

	close(jobs)

	wg.Wait()

	close(results)

	clusterGrid := make(map[string]map[string]node.NodeInfo)

	for cir := range results {
		if cir.err == nil {
			clusterGrid[cir.host] = nodeInfoMap(*cir.info)
		}
	}

	return clusterGrid
}

type NodeOverview struct {
	info node.NodeInfo
	metrics *node.Metrics
	metricsResponseDuration time.Duration
}

func buildNodeOverviewWorker(wg *sync.WaitGroup, nodes <-chan node.NodeInfo, result chan<- NodeOverview) {
	defer wg.Done()

	for nodeInfo := range nodes {
		start := time.Now()
		m, _ := node.GetClient(nodeInfo.Ip).GetNodeMetrics()
		elapsed := time.Since(start)
		result <- NodeOverview { nodeInfo, m, elapsed}
	}
}

func buildClusterOverview(globalClusterInfo *node.ClusterInfo) []NodeOverview {

	const workers = 24

	jobs := make(chan node.NodeInfo, len(*globalClusterInfo))
	results := make(chan NodeOverview, len(*globalClusterInfo))

	var wg sync.WaitGroup

	for w := 0; w < workers; w++ {
		wg.Add(1)
		go buildNodeOverviewWorker(&wg, jobs, results)
	}

	for _, nodeInfo := range *globalClusterInfo {
		jobs <- nodeInfo
	}

	close(jobs)

	wg.Wait()

	close(results)

	var clusterOverview []NodeOverview

	for nodeOverview := range results {
		clusterOverview = append(clusterOverview, nodeOverview)
	}

	return clusterOverview
}

func networkOverviewWorker(wg *sync.WaitGroup, globalClusterInfo *node.ClusterInfo, result chan<- []NodeOverview ) {
	defer wg.Done()
	result <- buildClusterOverview(globalClusterInfo)
}

func networkGridWorker(wg *sync.WaitGroup, globalClusterInfo *node.ClusterInfo, result chan<- map[string]map[string]node.NodeInfo ) {
	defer wg.Done()
	result <- buildNetworkGrid(globalClusterInfo)
}


func (n *nodegrid) BuildNetworkStatus(url string, silent bool, outputImage string) {

	globalClusterInfo, err := lb.GetClient(url).GetClusterInfo()

	if err == nil {

		var wg sync.WaitGroup

		nodeResults := make(chan []NodeOverview, 1)
		gridResults := make(chan map[string]map[string]node.NodeInfo, 1)

		wg.Add(2)

		go networkOverviewWorker(&wg, globalClusterInfo, nodeResults)
		go networkGridWorker(&wg, globalClusterInfo, gridResults)

		wg.Wait()

		close(nodeResults)
		close(gridResults)

		networkOverview, networkGrid := <-nodeResults, <-gridResults

		sort.Slice(networkOverview, func(i, j int) bool {
			return strings.ToLower(networkOverview[i].info.Alias) < strings.ToLower(networkOverview[j].info.Alias)
		})

		if silent == false {
			PrintAsciiOutput(networkOverview, networkGrid)
		}

		if outputImage != ""  {
			BuildImageOutput(outputImage, networkOverview, networkGrid)
		}
	} else {
		println("error")
		println(err.Error())
	}
}