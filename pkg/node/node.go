package node

import (
	"fmt"
	"net"
	"strconv"
	"strings"
)

type NodeState string

const (
	Initial         NodeState = "Initial"
	ReadyToJoin     NodeState = "ReadyToJoin"
	LoadingGenesis  NodeState = "LoadingGenesis"
	GenesisReady    NodeState = "GenesisReady"
	StartingSession NodeState = "StartingSession"
	SessionStarted  NodeState = "SessionStarted"
	Ready           NodeState = "Ready"
	Leaving         NodeState = "Leaving"
	Offline         NodeState = "Offline"
	Unknown         NodeState = "Unknown"
	Undefined       NodeState = "Undefined"
)

var ValidStatuses = [...]NodeState{Initial, ReadyToJoin, LoadingGenesis, GenesisReady, StartingSession, SessionStarted, Ready, Leaving, Offline, Unknown, Undefined}

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

	panic(fmt.Sprintf("Status unknown is %s", in))
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
	if pi == nil {
		return Undefined
	}

	if pi.cardinalState == "" {
		pi.cardinalState = StateFromString(pi.State)
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
