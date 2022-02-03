package node

import (
	"fmt"
	"net"
	"strconv"
	"strings"
)

type NodeId struct {
	Hex string `json:"hex"`
}

type NodeState string

const (
	PendingDownload                   NodeState = "PendingDownload"
	SessionStarted                    NodeState = "SessionStarted"
	ReadyForDownload                  NodeState = "ReadyForDownload"
	DownloadInProgress                NodeState = "DownloadInProgress"
	DownloadCompleteAwaitingFinalSync NodeState = "DownloadCompleteAwaitingFinalSync"
	SnapshotCreation                  NodeState = "SnapshotCreation"
	Ready                             NodeState = "Ready"
	Leaving                           NodeState = "Leaving"
	Offline                           NodeState = "Offline"
	Undefined                         NodeState = "Undefined"
)

var ValidStatuses = [...]NodeState{SessionStarted, PendingDownload, ReadyForDownload, DownloadInProgress, DownloadCompleteAwaitingFinalSync, SnapshotCreation, Ready, Leaving, Offline}

func IsRedownloading(status NodeState) bool {
	return status == PendingDownload || status == ReadyForDownload || status == DownloadInProgress || status == DownloadCompleteAwaitingFinalSync
}

func IsOffline(status NodeState) bool {
	return status == Leaving || status == Offline
}

func StateFromString(in string) NodeState {
	for _, v := range ValidStatuses {
		if in == fmt.Sprint(v) {
			return v
		}
	}

	panic(fmt.Sprintf("Status unknown is %s", in))
}

type Addr struct {
	Ip   string `json:"ip"`
	Port int    `json:"publicPort"`
}

type PeerInfo struct {
	Id         string                 `json:"id"`
	Ip         string                 `json:"ip"`
	PublicPort int                    `json:"publicPort"`
	P2PPort    int                    `json:"p2pPort"`
	Session    string                 `json:"session"`
	State      map[string]interface{} `json:"state"`

	cardinalState NodeState
}

// TODO:  DEPRECATED, wait until State is so unreadable
func (pi *PeerInfo) CardinalState() NodeState {
	if pi == nil {
		return Undefined
	}

	if pi.cardinalState == "" {
		s := ""
		for k, _ := range pi.State {
			s = k
			break
		}
		pi.cardinalState = StateFromString(s)
	}

	return pi.cardinalState
}

//TODO: this is a placeholder
type Metrics struct {
	Version            string
	LastSnapshotHeight string
	NodeState          NodeState
	Alias              string
}

// :9000/debug/peers
type Peers []PeerInfo

func (p *PeerInfo) Addr() Addr {
	return Addr{
		p.Ip,
		p.PublicPort,
	}
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
