package node

import (
	"encoding/json"
	"fmt"
	log "github.com/sirupsen/logrus"
	"io/ioutil"
	"net/http"
	"strings"
)

type PublicApi interface {
	Peers() (*Peers, error)
}

type node struct {
	url string
}

func (n *node) Peers() (*Peers, error) {
	url := strings.TrimRight(n.url, "/") + "/debug/peers"

	log.Debugf("Make request for peers on %s", url)

	resp, err := http.Get(url)

	if err != nil {
		log.Debugf("Node %s is not operational at the moment, returned error=%s", url, err.Error())

		return nil, err
	}

	if resp.StatusCode != 200 {
		log.Debugf("Node %s returned status code=%d", url, resp.StatusCode)

		return nil, err
	}

	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)

	if err != nil {
		log.Errorf("Failed to read node response body error=%s status=%d\n", err.Error(), resp.StatusCode)
		return nil, err
	}

	ci := Peers{}

	err = json.Unmarshal(body, &ci)

	if err != nil {
		log.Errorf("Failed to unmarshal node response body error=%s body=%s status=%d", err.Error(), string(body), resp.StatusCode)

		return nil, err
	}

	log.Debugf("Got peers %d from %s", len(ci), n.url)

	return &ci, nil
}

func GetPublicClient(addr Addr) PublicApi {

	return &node{fmt.Sprintf("http://%s:%d", addr.Ip, addr.Port)}
}
