package lb

import (
	"constellation/pkg/node"
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"strings"
)


type Loadbalancer interface {
	GetClusterInfo() (*node.ClusterInfo, error)
}

type loadbalancer struct {
	url string
}

// lb.constellationnetwork.io:9000
func (lb *loadbalancer) GetClusterInfo() (*node.ClusterInfo, error) {
	endpointUrl := strings.TrimRight(lb.url, "/") + "/cluster/info"

	resp, err := http.Get(endpointUrl)

	if err != nil {
		log.Fatalf("Cannot execute request err=%s", err.Error())

		return nil, err
	}

	if (resp.StatusCode != 200) { // TODO: Find constants for http status codes in go std lib
		log.Fatal("Loadbalancer is not operational at the moment, returned http code=503")
		return nil, nil
	}

	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)

	if err != nil {
		log.Fatalf("Cannot read response body, lb resp code=%d err=%s", resp.StatusCode, err.Error())
		return nil, err
	}

	ci := node.ClusterInfo{}

	err = json.Unmarshal(body, &ci)

	if err != nil {
		log.Fatalf("Cannot unmarshal JSON response error=%s body=%s", err.Error(), string(body))
		return nil, err
	}

	return &ci, nil
}

func GetClient(url string) Loadbalancer {

	return & loadbalancer { url }
}