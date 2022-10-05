package lb

import (
	"constellation/pkg/node"
	"encoding/json"
	log "github.com/sirupsen/logrus"
	"io"
	"net/http"
	"net/url"
)

type Loadbalancer interface {
	GetClusterNodes() (*node.Peers, error)
}

type loadbalancer struct {
	url string
}

func (n *loadbalancer) GetClusterNodes() (*node.Peers, error) {

	url, err := url.Parse(n.url)

	url.Path = "/cluster/info"

	if err != nil {
		log.Debugf("Cannot parse url=%s error=%s", n.url, err.Error())
		return nil, err
	}

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

	body, err := io.ReadAll(resp.Body)

	if err != nil {
		log.Errorf("Failed to read node response body from=%s error=%s status=%d\n", url.String(), err.Error(), resp.StatusCode)
		return nil, err
	}

	ci := node.Peers{}

	err = json.Unmarshal(body, &ci)

	if err != nil {
		log.Errorf("Failed to unmarshal node response body from=%s error=%s body=%s status=%d", url.String(), err.Error(), string(body), resp.StatusCode)

		return nil, err
	}

	log.Trace(ci)

	return &ci, nil
}

func GetClient(url string) Loadbalancer {

	return &loadbalancer{url}
}
