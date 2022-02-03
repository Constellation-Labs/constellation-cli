package nodegrid

import (
	"constellation/pkg/node"
	"github.com/jszwec/csvutil"
	log "github.com/sirupsen/logrus"
	"io/ioutil"
	"sort"
	"strings"
	"sync"
	"time"
)

type Nodegrid interface {
	// Split this as it mixes two concerns
	BuildNetworkStatus(url node.Addr, silent bool, outputImage string, outputTheme string, verbose bool) (error, *NetworkStatus)
	Operators() map[string]Operator
}

type Operator struct {
	HexId     string `csv:"id"`
	DiscordId string `csv:"discord"`
	Name      string `csv:"name"`
}

type nodegrid struct {
	operatorsFilePath string
	operators         map[string]Operator
	operatorsLoaded   bool
}

func (n *nodegrid) Operators() map[string]Operator {

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
	return &nodegrid{operatorsFile, make(map[string]Operator), false}
}

type nodePeersResult struct {
	host    string
	err     error
	peers   *node.Peers
	latency time.Duration
	addr    node.Addr
}

func peers2map(peers node.Peers) map[string]*node.PeerInfo {
	m := make(map[string]*node.PeerInfo)

	for _, peerInfo := range peers {
		m[peerInfo.Ip] = &peerInfo
	}

	return m
}

func getNodePeers(addr node.Addr) nodePeersResult {
	start := time.Now()
	peers, e := node.GetPublicClient(addr).Peers()
	duration := time.Since(start)

	if e != nil {
		log.Debugf("Cannot get peers for %s %s", addr.Ip, e)
		emptyPeers := make(node.Peers, 0)

		return nodePeersResult{
			addr.Ip,
			e,
			&emptyPeers,
			duration,
			addr,
		}
	}

	return nodePeersResult{
		addr.Ip,
		e,
		peers,
		duration,
		addr,
	}
}

func getNodePeersWorker(wg *sync.WaitGroup, cluster <-chan node.Addr, result chan<- nodePeersResult) {
	defer wg.Done()

	for addr := range cluster {
		result <- getNodePeers(addr)
	}
}

func (n *nodegrid) buildNetworkGrid(addrs *[]node.Addr) networkGrid {

	const workers = 24

	jobs := make(chan node.Addr, len(*addrs))
	results := make(chan nodePeersResult, len(*addrs))

	var wg sync.WaitGroup

	for w := 0; w < workers; w++ {
		wg.Add(1)
		go getNodePeersWorker(&wg, jobs, results)
	}

	for _, addr := range *addrs {
		jobs <- addr
	}

	close(jobs)

	wg.Wait()

	close(results)

	log.Debug("Work on results to regroup")

	clusterGrid := make(map[string]map[string]*node.PeerInfo)
	nodeLatency := make(map[string]time.Duration)

	for peersResult := range results {
		nodeLatency[peersResult.host] = time.Second * 30

		if peersResult.err == nil {
			clusterGrid[peersResult.host] = peers2map(*peersResult.peers)
			nodeLatency[peersResult.host] = peersResult.latency
		}
	}
	log.Debug("networkGrid done")

	return networkGrid{clusterGrid, nodeLatency}
}

type NodeOverview struct {
	Info                node.PeerInfo
	AvgResponseDuration time.Duration
	Operator            *Operator
	Metrics             *node.Metrics
}

type networkGrid struct {
	grid    map[string]map[string]*node.PeerInfo
	latency map[string]time.Duration
}

func (n *nodegrid) networkGridWorker(wg *sync.WaitGroup, globalClusterInfo *[]node.Addr, result chan<- networkGrid) {
	defer wg.Done()
	result <- n.buildNetworkGrid(globalClusterInfo)
}

type NetworkStatus struct {
	NodesList []NodeOverview
	NodesGrid map[string]map[string]*node.PeerInfo
}

func (n *nodegrid) BuildNetworkStatus(addr node.Addr, silent bool, outputImage string, outputTheme string, verbose bool) (error, *NetworkStatus) {

	//TODO: Until we do not have lb we will query a node
	peers, err := node.GetPublicClient(addr).Peers()

	if err != nil {
		panic(err)
	}

	addrs := make([]node.Addr, len(*peers))

	for i, v := range *peers {
		addrs[i] = v.Addr()
	}

	if err == nil {

		var wg sync.WaitGroup

		n.Operators()

		gridResults := make(chan networkGrid, 1)

		wg.Add(1)

		go n.networkGridWorker(&wg, &addrs, gridResults)

		wg.Wait()

		close(gridResults)

		networkGrid := <-gridResults

		networkOverview := make([]NodeOverview, len(networkGrid.latency))

		for i, peer := range *peers {
			networkOverview[i] = NodeOverview{peer, networkGrid.latency[addr.Ip],
				nil, nil} // TODO: replace with real values
		}

		sort.Slice(networkOverview, func(i, j int) bool {
			return strings.ToLower(networkOverview[i].Info.Ip) < strings.ToLower(networkOverview[j].Info.Ip)
		})

		if silent == false {
			PrintAsciiOutput(networkOverview, networkGrid.grid, verbose)
		}

		if outputImage != "" {
			BuildImageOutput(outputImage, networkOverview, networkGrid.grid, outputTheme)
		}

		return nil, &NetworkStatus{networkOverview, networkGrid.grid}
	} else {
		return err, nil
	}
}

const LatencyTriggerMilliseconds = 2000
