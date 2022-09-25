package node

import (
	"fmt"
	log "github.com/sirupsen/logrus"
	"net"
	"strconv"
	"strings"
)

type NodeState string

const (
	Initial            NodeState = "Initial"
	ReadyToJoin        NodeState = "ReadyToJoin"
	WaitingForDownload NodeState = "WaitingForDownload"
	LoadingGenesis     NodeState = "LoadingGenesis"
	GenesisReady       NodeState = "GenesisReady"
	StartingSession    NodeState = "StartingSession"
	SessionStarted     NodeState = "SessionStarted"
	Ready              NodeState = "Ready"
	Leaving            NodeState = "Leaving"
	Offline            NodeState = "Offline"
	Observing          NodeState = "Observing"

	NotSupported NodeState = "NotSupported" // Internal fallback status when not recognized
	Undefined    NodeState = "Undefined"    // Internal status when we could not obtain status for the node
)

var ValidStatuses = [...]NodeState{Initial, ReadyToJoin, WaitingForDownload, LoadingGenesis, GenesisReady, StartingSession, SessionStarted, Ready, Leaving, Offline, NotSupported, Observing, Undefined}

func IsRedownloading(status NodeState) bool {
	return status == LoadingGenesis
}

func IsOffline(status NodeState) bool {
	return status == Leaving || status == Offline || status == Initial || status == ReadyToJoin
}

func StateFromString(in string) NodeState {
	for _, v := range ValidStatuses {
		if in == fmt.Sprint(v) {
			return v
		}
	}

	log.Debugf("Status=%s is unknown to the cli tool", in)

	return NotSupported
}

type Addr struct {
	Ip   string `json:"ip"`
	Port int    `json:"publicPort"`
}

type PeerInfo struct {
	Id         string `json:"id"`
	Ip         string `json:"ip"`
	PublicPort int    `json:"publicPort"`
	P2PPort    int    `json:"p2pPort"`
	Session    string `json:"session"`
	State      string `json:"state"`

	cardinalState NodeState
}

type ClusterInfo struct {
	Id    string
	Peers *Peers
}

// TODO:  DEPRECATED, wait until State is so unreadable
func (pi *PeerInfo) CardinalState() NodeState {
	if pi.cardinalState == "" {
		pi.cardinalState = StateFromString(pi.State)
	}

	return pi.cardinalState
}

// TODO: this is a placeholder
type Metrics struct {
	PublicPort string
	Session    int
	NodeState  NodeState
}

// :9000/debug/peers
type Peers []PeerInfo

func (p *PeerInfo) Addr() Addr {
	return Addr{
		p.Ip,
		p.PublicPort,
	}
}

func (p *PeerInfo) ShortId() string {
	return p.Id[0:8]
}

func AddrOf(in string) Addr {

	host := in
	addrPort := 9000 // TODO: move to constants

	if strings.Contains(in, ":") {

		addr, port, e := net.SplitHostPort(in)
		host = addr

		if e != nil {
			panic(e)
		}

		if portNum, err := strconv.ParseUint(port, 10, 32); err != nil {
			addrPort = int(portNum)
		}
	}

	return Addr{host, addrPort}
}
