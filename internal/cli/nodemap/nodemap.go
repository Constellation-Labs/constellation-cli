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
	DiscoverNetwork(url node.Addr, silent bool, verbose bool, outputImage string, outputTheme string) (error, *NetworkStatus)
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

func peers2map(peers node.Peers) map[string]node.PeerInfo {
	m := make(map[string]node.PeerInfo)

	log.Debugf("Making map of %d list", len(peers))

	for _, peerInfo := range peers {
		m[peerInfo.Id] = peerInfo
	}

	log.Debugf("Made map of %d list", len(m))

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

	log.Debugf("Return info with %d peers %s", len(*clusterInfo.Peers), clusterInfo.Id)

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

func (n *nodemap) buildNetworkmap(addrs *[]node.Addr) networkGrid {

	log.Debug("Fetch from addrs", addrs)

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

	clustermap := make(map[string]map[string]node.PeerInfo)
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
	grid    map[string]map[string]node.PeerInfo
	latency map[string]time.Duration
}

func (n *nodemap) networkmapWorker(wg *sync.WaitGroup, globalClusterInfo *[]node.Addr, result chan<- networkGrid) {
	defer wg.Done()
	result <- n.buildNetworkmap(globalClusterInfo)
}

type NetworkStatus struct {
	DiscoveredNodes []*node.PeerInfo
	NodesMap        map[string]map[string]node.PeerInfo
}

type ClusterNode struct {
	Addr     node.Addr
	Id       string
	SelfInfo *node.PeerInfo
}

func (p ClusterNode) ShortId() string {
	return p.Id[0:8]
}

func (n *nodemap) fetchPool(freshAddrPool []node.Addr) map[string]map[string]node.PeerInfo {
	var wg sync.WaitGroup

	mapResults := make(chan networkGrid, 1)
	wg.Add(1)
	go n.networkmapWorker(&wg, &freshAddrPool, mapResults)
	wg.Wait()
	close(mapResults)

	partialResult := <-mapResults

	return partialResult.grid
}

func (n *nodemap) DiscoverNetwork(addr node.Addr, silent bool, verbose bool, outputImage string, outputTheme string) (error, *NetworkStatus) {

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

	freshAddrPool := addrs

	networkGridAccumulator := networkGrid{
		make(map[string]map[string]node.PeerInfo),
		make(map[string]time.Duration),
	}

	for len(freshAddrPool) > 0 {
		log.Debugf("Discovery in progress for %d peers", len(freshAddrPool))

		partialGrid := n.fetchPool(freshAddrPool)
		freshAddrPool = nil

		for nodeid, nodesgrid := range partialGrid {

			log.Trace("Size of grid for %s is %d", nodeid, len(nodesgrid))
			log.Debug(nodesgrid)

			networkGridAccumulator.grid[nodeid] = nodesgrid

			for kk, pinfo := range nodesgrid {

				log.Trace("Checking %s if contains %s %s", kk, pinfo.Ip, pinfo.Id)

				if !slices.Contains(addrs, pinfo.Addr()) {

					log.Infof("Discovered new peer from %s->%s %s", nodeid[0:8], pinfo.Id[0:8], pinfo.Ip)
					freshAddrPool = append(freshAddrPool, pinfo.Addr())
					ids = append(ids, pinfo.Id)
					addrs = append(addrs, pinfo.Addr())
				}
			}
		}
	}

	log.Infof("discovered %d grid addrs %d ids %d", len(networkGridAccumulator.grid), len(addrs), len(ids))

	clusterOverview := make([]ClusterNode, len(addrs))

	for i, peer := range ids {

		if selfInfo, ok := networkGridAccumulator.grid[peer][peer]; ok {
			clusterOverview[i] = ClusterNode{
				addrs[i],
				ids[i],
				&selfInfo,
			}
		} else {
			clusterOverview[i] = ClusterNode{
				addrs[i],
				ids[i],
				nil,
			}
		}
	}

	if err == nil {

		if silent == false {
			PrintAsciiOutput(clusterOverview, networkGridAccumulator.grid, verbose)
		}

		if outputImage != "" {
			BuildImageOutput(outputImage, clusterOverview, networkGridAccumulator.grid, outputTheme)
		}

	} else {
		return err, nil
	}

	return nil, nil
}
