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

func (n *nodegrid) buildNetworkGrid(ci *node.ClusterInfo) map[string]map[string]node.NodeInfo{

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
	operator *Operator
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

func (n *nodegrid) networkGridWorker(wg *sync.WaitGroup, globalClusterInfo *node.ClusterInfo, result chan<- map[string]map[string]node.NodeInfo ) {
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
		gridResults := make(chan map[string]map[string]node.NodeInfo, 1)

		wg.Add(2)

		go n.networkOverviewWorker(&wg, globalClusterInfo, nodeResults)
		go n.networkGridWorker(&wg, globalClusterInfo, gridResults)

		wg.Wait()

		close(nodeResults)
		close(gridResults)

		networkOverview, networkGrid := <-nodeResults, <-gridResults

		sort.Slice(networkOverview, func(i, j int) bool {
			return strings.ToLower(networkOverview[i].info.Alias) < strings.ToLower(networkOverview[j].info.Alias)
		})

		if silent == false {
			PrintAsciiOutput(networkOverview, networkGrid, verbose)
		}

		if outputImage != ""  {
			BuildImageOutput(outputImage, networkOverview, networkGrid, outputTheme)
		}

		return nil, &NetworkStatus{networkOverview, networkGrid}
	} else {
		return err, nil
	}
}