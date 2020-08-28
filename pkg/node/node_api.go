package node

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"io/ioutil"
	"strings"
)

type Node interface {
	GetClusterInfo() (*ClusterInfo, error)
	GetNodeMetrics() (*Metrics, error)
}

type node struct {
	url string
}

func (n *node) GetNodeMetrics() (*Metrics, error) {
	url := strings.TrimRight(n.url, "/") + "/metrics"
	resp, err := http.Get(url)

	if err != nil {

		return nil, err
	}

	if resp.StatusCode != 200 {
		log.Fatalf("Node %s returned status code=%d", url, resp.StatusCode)

		return nil, err
	}

	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)

	if err != nil {
		log.Fatalf("Failed to read node response body error=%s status=%d", err.Error(), resp.StatusCode)
		return nil, err
	}

	e := MetricsEnvelope{}

	err = json.Unmarshal(body, &e)

	if err != nil {
		log.Fatalf("Failed to unmarshal node response body error=%s body=%s status=%d", err.Error(), string(body), resp.StatusCode)

		return nil, err
	}

	return &e.Metrics, nil
}

func (n *node) GetClusterInfo() (*ClusterInfo, error) {
	url := strings.TrimRight(n.url, "/") + "/cluster/info"
	resp, err := http.Get(url)

	if err != nil {
		//log.Fatalf("Node %s is not operational at the moment, returned error=%s", url, err.Error())

		return nil, err
	}

	if resp.StatusCode != 200 {
		//log.Fatalf("Node %s returned status code=%d", url, resp.StatusCode)

		return nil, err
	}

	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)

	if err != nil {
		//log.Fatalf("Failed to read node response body error=%s status=%d", err.Error(), resp.StatusCode)
		return nil, err
	}

	ci := ClusterInfo{}

	err = json.Unmarshal(body, &ci)

	if err != nil {
		log.Fatalf("Failed to unmarshal node response body error=%s body=%s status=%d", err.Error(), string(body), resp.StatusCode)

		return nil, err
	}

	return &ci, nil
}

func GetClient(addr NodeAddr) Node {

	return & node { fmt.Sprintf("http://%s:%d", addr.Host, addr.Port - 1) }
}