package node

import (
	"encoding/json"
	"fmt"
	log "github.com/sirupsen/logrus"
	"io/ioutil"
	"net"
	"net/http"
	"net/url"
)

type PublicApi interface {
	Peers() (*Peers, error)
	ClusterInfo() (*ClusterInfo, error)
}

type node struct {
	addr Addr
}

func (n *node) Peers() (*Peers, error) {

	url := url.URL{Scheme: "http", Host: net.JoinHostPort(n.addr.Ip, fmt.Sprintf("%d", n.addr.Port)), Path: "/debug/peers"}

	log.Debugf("Make request for peers on %s", url.String())

	resp, err := http.Get(url.String())

	if err != nil {
		log.Debugf("Node %s is not operational at the moment, returned error=%s", url.String(), err.Error())

		return nil, err
	}

	if resp.StatusCode != 200 {
		log.Debugf("Node %s returned status code=%d", url.String(), resp.StatusCode)

		return nil, err
	}

	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)

	if err != nil {
		log.Errorf("Failed to read node response body from=%s error=%s status=%d\n", url.String(), err.Error(), resp.StatusCode)
		return nil, err
	}

	id := resp.Header.Get("X-Id")

	if id == "" {
		return nil, fmt.Errorf("no id in response from %s:%d", n.addr.Ip, n.addr.Port)
	}

	ci := Peers{}

	err = json.Unmarshal(body, &ci)

	if err != nil {
		log.Errorf("Failed to unmarshal node response body from=%s error=%s body=%s status=%d", url.String(), err.Error(), string(body), resp.StatusCode)

		return nil, err
	}

	log.Debugf("Got peers %d from %s %s", len(ci), n.addr.Ip, id)

	return &ci, nil
}

func (n *node) ClusterInfo() (*ClusterInfo, error) {

	url := url.URL{Scheme: "http", Host: net.JoinHostPort(n.addr.Ip, fmt.Sprintf("%d", n.addr.Port)), Path: "/cluster/info"}

	log.Debugf("Make request for peers on %s", url.String())

	resp, err := http.Get(url.String())

	if err != nil {
		log.Debugf("Node %s is not operational at the moment, returned error=%s", url.String(), err.Error())

		return nil, err
	}

	if resp.StatusCode != 200 {
		log.Debugf("Node %s returned status code=%d", url.String(), resp.StatusCode)

		return nil, err
	}

	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)

	if err != nil {
		log.Errorf("Failed to read node response body error=%s status=%d\n", err.Error(), resp.StatusCode)
		return nil, err
	}

	id := resp.Header.Get("X-Id")

	if id == "" {
		return nil, fmt.Errorf("no id in response from %s:%d", n.addr.Ip, n.addr.Port)
	}

	ci := Peers{}

	err = json.Unmarshal(body, &ci)

	if err != nil {
		log.Errorf("Failed to unmarshal node response body error=%s body=%s status=%d", err.Error(), string(body), resp.StatusCode)

		return nil, err
	}

	log.Debugf("Got peers %d from %s %s", len(ci), n.addr.Ip, id)

	return &ClusterInfo{id, &ci}, nil
}

func GetPublicClient(addr Addr) PublicApi {

	return &node{addr}
}
