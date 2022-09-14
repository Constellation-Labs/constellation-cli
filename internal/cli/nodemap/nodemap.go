package nodemap

import (
	"constellation/pkg/node"
	log "github.com/sirupsen/logrus"
	"golang.org/x/exp/slices"
	"sync"
	"time"
)

type Nodemap interface {
	// Split this as it mixes two concerns
	DiscoverNetwork(url node.Addr, silent bool) (error, *NetworkStatus)
}

type nodemap struct {
}

func NewNodemap() Nodemap {
	return &nodemap{}
}

type nodePeersResult struct {
	nodeId  string
	err     error
	peers   *node.Peers
	latency time.Duration
	addr    node.Addr
}

func peers2map(peers node.Peers) map[string]*node.PeerInfo {
	m := make(map[string]*node.PeerInfo)

	for _, peerInfo := range peers {
		m[peerInfo.Id] = &peerInfo
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
		log.Debug("Get peers", addr)
		result <- getNodePeers(addr)
	}
}

func (n *nodemap) buildNetworkmap(addrs *[]node.Addr) networkGrid {

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

	clustermap := make(map[string]map[string]*node.PeerInfo)
	nodeLatency := make(map[string]time.Duration)

	for peersResult := range results {
		nodeLatency[peersResult.nodeId] = time.Second * 30

		if peersResult.err == nil {
			clustermap[peersResult.nodeId] = peers2map(*peersResult.peers)
			nodeLatency[peersResult.nodeId] = peersResult.latency
		}
	}
	log.Debug("networkmap done")

	return networkGrid{clustermap, nodeLatency}
}

type networkGrid struct {
	grid    map[string]map[string]*node.PeerInfo
	latency map[string]time.Duration
}

func (n *nodemap) networkmapWorker(wg *sync.WaitGroup, globalClusterInfo *[]node.Addr, result chan<- networkGrid) {
	defer wg.Done()
	result <- n.buildNetworkmap(globalClusterInfo)
}

type NetworkStatus struct {
	DiscoveredNodes []*node.PeerInfo
	NodesMap        map[string]map[string]*node.PeerInfo
}

type ClusterNode struct {
	Addr     node.Addr
	Id       string
	SelfInfo *node.PeerInfo
}

func (n *nodemap) DiscoverNetwork(addr node.Addr, verbose bool) (error, *NetworkStatus) {

	clusterInfo, err := node.GetPublicClient(addr).ClusterInfo()

	if err != nil {
		panic(err)
	}

	addrs := make([]node.Addr, len(*clusterInfo.Peers))
	ids := make([]string, len(*clusterInfo.Peers))

	for i, v := range *clusterInfo.Peers {
		addrs[i] = v.Addr()
		ids[i] = v.Id
	}

	newAddrs := addrs

	networkGridAccumulator := networkGrid{
		make(map[string]map[string]*node.PeerInfo),
		make(map[string]time.Duration),
	}

	for len(newAddrs) > 0 {
		log.Debugf("Discovery in progress for %d peers", len(newAddrs))

		var wg sync.WaitGroup

		mapResults := make(chan networkGrid, 1)
		wg.Add(1)
		go n.networkmapWorker(&wg, &newAddrs, mapResults)
		wg.Wait()
		close(mapResults)
		newAddrs = nil

		partialResult := <-mapResults

		for k, v := range partialResult.grid {

			networkGridAccumulator.grid[k] = v

			for _, pinfo := range v {

				if !slices.Contains(addrs, pinfo.Addr()) {
					newAddrs = append(newAddrs, pinfo.Addr())
					ids = append(ids, pinfo.Id)
					addrs = append(addrs, pinfo.Addr())
				}
			}
		}
	}

	log.Infof("discovered %d grid addrs %d ids %d", len(networkGridAccumulator.grid), len(addrs), len(ids))

	clusterOverview := make([]ClusterNode, len(addrs))

	for i, peer := range ids {

		clusterOverview[i] = ClusterNode{
			addrs[i],
			ids[i],
			networkGridAccumulator.grid[peer][peer],
		}
	}

	if err == nil {
		PrintAsciiOutput(clusterOverview, networkGridAccumulator.grid, verbose)

	} else {
		return err, nil
	}

	return nil, nil
}
