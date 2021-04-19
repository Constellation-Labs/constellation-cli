package nodegrid

import (
	"constellation_cli/pkg/lb"
	"constellation_cli/pkg/node"
	"github.com/jszwec/csvutil"
	"io/ioutil"
	"sort"
	"strings"
	"sync"
	"time"
)

type Nodegrid interface {
	// Split this as it mixes two concerns
	BuildNetworkStatus(url string, silent bool, outputImage string, outputTheme string, verbose bool)(error, *NetworkStatus)
	Operators() map[string]Operator
}

type Operator struct {
	HexId string `csv:"id"`
	DiscordId string `csv:"discord"`
	Name string `csv:"name"`
}

type nodegrid struct {
	operatorsFilePath string
	operators map[string]Operator
	operatorsLoaded bool
}

func (n *nodegrid) Operators() map[string]Operator{

	if !n.operatorsLoaded {
		var operators []Operator

		operatorsFileBytes, _ := ioutil.ReadFile(n.operatorsFilePath)

		csvutil.Unmarshal(operatorsFileBytes, &operators)

		for _, o := range operators {
			n.operators[o.HexId] = o
		}

		n.operatorsLoaded = true
	}

	return n.operators
}

func NewNodegrid(operatorsFile string) Nodegrid {
	return & nodegrid {operatorsFile, make(map[string]Operator), false}
}

type nodeResult struct {
	host string
	err error
	info *node.ClusterInfo
	latency time.Duration
}

func nodeInfoMap(info node.ClusterInfo) map[string]node.NodeInfo{
	m := make(map[string]node.NodeInfo)

	for _, nodeInfo := range info {
		m[nodeInfo.Ip.Host] = nodeInfo
	}

	return m
}

func clusterInfo(addr node.NodeAddr) nodeResult {
	start := time.Now()
	ci, e := node.GetClient(addr).GetClusterInfo()
	duration := time.Since(start)

	return nodeResult {
		 addr.Host,
		  e,
		  ci,
		duration,
	}
}

func queryNodeForClusterInfoWorker(wg *sync.WaitGroup, cluster <-chan node.NodeAddr, result chan<- nodeResult) {
	defer wg.Done()

	for addr := range cluster {
		ci := clusterInfo(addr)
		result <- ci
	}
}

func (n *nodegrid) buildNetworkGrid(ci *node.ClusterInfo) networkGrid{

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
	nodeLatency := make(map[string]time.Duration)

	for cir := range results {
		nodeLatency[cir.host] = time.Second * 30

		if cir.err == nil {
			clusterGrid[cir.host] = nodeInfoMap(*cir.info)
			nodeLatency[cir.host] = cir.latency
		}
	}

	return networkGrid {clusterGrid, nodeLatency }
}

type NodeOverview struct {
	Info                node.NodeInfo
	Metrics             *node.Metrics
	AvgResponseDuration time.Duration
	Operator            *Operator
}

func (n *nodegrid) buildNodeOverviewWorker(wg *sync.WaitGroup, nodes <-chan node.NodeInfo, result chan<- NodeOverview) {
	defer wg.Done()

	for nodeInfo := range nodes {
		var op *Operator = nil

		o, e := n.Operators()[nodeInfo.Id.Hex]
		if e {
			op = &o
		}

		start := time.Now()
		m, _ := node.GetClient(nodeInfo.Ip).GetNodeMetrics()
		elapsed := time.Since(start)

		result <- NodeOverview { nodeInfo, m, elapsed, op}
	}
}

func (n *nodegrid) buildClusterOverview(globalClusterInfo *node.ClusterInfo) []NodeOverview {

	const workers = 24

	jobs := make(chan node.NodeInfo, len(*globalClusterInfo))
	results := make(chan NodeOverview, len(*globalClusterInfo))

	var wg sync.WaitGroup

	for w := 0; w < workers; w++ {
		wg.Add(1)
		go n.buildNodeOverviewWorker(&wg, jobs, results)
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

func (n *nodegrid) networkOverviewWorker(wg *sync.WaitGroup, globalClusterInfo *node.ClusterInfo, result chan<- []NodeOverview ) {
	defer wg.Done()
	result <- n.buildClusterOverview(globalClusterInfo)
}

type networkGrid struct {
	grid map[string]map[string]node.NodeInfo
	latency map[string]time.Duration
}

func (n *nodegrid) networkGridWorker(wg *sync.WaitGroup, globalClusterInfo *node.ClusterInfo, result chan<- networkGrid ) {
	defer wg.Done()
	result <- n.buildNetworkGrid(globalClusterInfo)
}

type NetworkStatus struct {
	NodesList []NodeOverview
	NodesGrid map[string]map[string]node.NodeInfo
}

func (n *nodegrid) BuildNetworkStatus(url string, silent bool, outputImage string, outputTheme string, verbose bool) (error, *NetworkStatus) {

	globalClusterInfo, err := lb.GetClient(url).GetClusterInfo()

	if err == nil {

		var wg sync.WaitGroup

		n.Operators()

		nodeResults := make(chan []NodeOverview, 1)
		gridResults := make(chan networkGrid, 1)

		wg.Add(2)

		go n.networkOverviewWorker(&wg, globalClusterInfo, nodeResults)
		go n.networkGridWorker(&wg, globalClusterInfo, gridResults)

		wg.Wait()

		close(nodeResults)
		close(gridResults)

		networkOverview, networkGrid := <-nodeResults, <-gridResults

		for _, n := range networkOverview {
			n.AvgResponseDuration = (n.AvgResponseDuration + networkGrid.latency[n.Info.Ip.Host])/2
		}

		sort.Slice(networkOverview, func(i, j int) bool {
			return strings.ToLower(networkOverview[i].Info.Alias) < strings.ToLower(networkOverview[j].Info.Alias)
		})

		if silent == false {
			PrintAsciiOutput(networkOverview, networkGrid.grid, verbose)
		}

		if outputImage != ""  {
			BuildImageOutput(outputImage, networkOverview, networkGrid.grid, outputTheme)
		}

		return nil, &NetworkStatus{networkOverview, networkGrid.grid}
	} else {
		return err, nil
	}
}

const LatencyTriggerMilliseconds = 2000