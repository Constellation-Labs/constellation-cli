package nodegrid

import (
	"constellation/pkg/lb"
	"constellation/pkg/node"
	log "github.com/sirupsen/logrus"
	"sort"
	"strings"
	"sync"
	"time"
)

type Nodegrid interface {
	BuildNetworkStatus(url string, silent bool, outputImage string, outputTheme string, verbose bool) (error, *NetworkStatus)
}

type nodegrid struct {
}

func NewNodegrid(operatorsFile string) Nodegrid {
	return &nodegrid{}
}

type nodePeersResult struct {
	nodeId  string
	err     error
	peers   *node.Peers
	latency time.Duration
	addr    node.Addr
}

func peers2map(peers node.Peers) map[string]node.PeerInfo {
	m := make(map[string]node.PeerInfo)

	for _, peerInfo := range peers {
		m[peerInfo.Id] = peerInfo
	}

	return m
}

func getNodePeers(addr node.Addr) nodePeersResult {
	start := time.Now()
	clusterInfo, e := node.GetPublicClient(addr).ClusterInfo()
	duration := time.Since(start)

	if e != nil {
		log.Debugf("Cannot get peers for %s %s", addr.Ip, e)
		emptyPeers := make(node.Peers, 0)

		return nodePeersResult{
			"",
			e,
			&emptyPeers,
			duration,
			addr,
		}
	}

	return nodePeersResult{
		clusterInfo.Id,
		e,
		clusterInfo.Peers,
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

	clusterGrid := make(map[string]map[string]node.PeerInfo)
	nodeLatency := make(map[string]time.Duration)

	for peersResult := range results {
		nodeLatency[peersResult.nodeId] = time.Second * 30

		if peersResult.err == nil {
			clusterGrid[peersResult.nodeId] = peers2map(*peersResult.peers)
			nodeLatency[peersResult.nodeId] = peersResult.latency
		}
	}
	log.Debug("networkGrid done")

	return networkGrid{clusterGrid, nodeLatency}
}

type NodeOverview struct {
	Id         string
	Ip         string
	PublicPort int

	LbInfo   node.PeerInfo
	SelfInfo *node.PeerInfo
}

type networkGrid struct {
	grid    map[string]map[string]node.PeerInfo
	latency map[string]time.Duration
}

func (n *nodegrid) networkGridWorker(wg *sync.WaitGroup, globalClusterInfo *[]node.Addr, result chan<- networkGrid) {
	defer wg.Done()
	result <- n.buildNetworkGrid(globalClusterInfo)
}

type NetworkStatus struct {
	NodesList []NodeOverview
	NodesGrid map[string]map[string]node.PeerInfo
}

func (n *nodegrid) BuildNetworkStatus(url string, silent bool, outputImage string, outputTheme string, verbose bool) (error, *NetworkStatus) {

	clusterNodes, err := lb.GetClient(url).GetClusterNodes()

	if err != nil {
		panic(err)
	}

	addrs := make([]node.Addr, len(*clusterNodes))

	for i, v := range *clusterNodes {
		addrs[i] = v.Addr()
	}

	if err == nil {

		var wg sync.WaitGroup

		gridResults := make(chan networkGrid, 1)

		wg.Add(1)

		go n.networkGridWorker(&wg, &addrs, gridResults)

		wg.Wait()

		close(gridResults)

		networkGrid := <-gridResults

		networkOverview := make([]NodeOverview, len(*clusterNodes))

		for i, peer := range *clusterNodes {

			selfInfo, exists := networkGrid.grid[peer.Id][peer.Id]

			if exists {
				networkOverview[i] = NodeOverview{
					peer.Id,
					peer.Ip,
					peer.PublicPort,
					peer,
					&selfInfo,
				}
			} else {
				networkOverview[i] = NodeOverview{
					peer.Id,
					peer.Ip,
					peer.PublicPort,
					peer,
					nil,
				}
			}
		}

		sort.Slice(networkOverview, func(i, j int) bool {
			return strings.ToLower(networkOverview[i].Ip) < strings.ToLower(networkOverview[j].Ip)
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
