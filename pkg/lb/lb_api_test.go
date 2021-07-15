package lb

import (
	"constellation/pkg/node"
	"testing"
	"encoding/json"
	)


func Test(t *testing.T) {

	test := []byte("[{\"alias\":\"Mel-B\",\"id\":{\"hex\":\"4d7ce3292306c8fb15a4d8a6a91136c2369b061c634b2374c23e3c06a2850926a026ef679f7a6e3a576c1c94a5bfa34f999c968b47a89a374d394cd6e0fff119\"},\"ip\":{\"host\":\"142.93.177.241\",\"port\":9001},\"status\":\"Ready\",\"reputation\":0},{\"alias\":\"Mintaka\",\"id\":{\"hex\":\"710b3dc521b805aea7a798d61f5d4dae39601124f1f34fac9738a78047adeff60931ba522250226b87a2194d3b7d39da8d2cbffa35d6502c70f1a7e97132a4b0\"},\"ip\":{\"host\":\"54.177.239.156\",\"port\":9001},\"status\":\"Ready\",\"reputation\":0},{\"alias\":\"My Little Pony\",\"id\":{\"hex\":\"0619ffcd96eaf4d82195d021dbc38f3899c0dfc97e8818867d0f50cc3d00df29bb7dbb3332678e52f40e1978a77b248478af12b8b5500226bc68736c060ba236\"},\"ip\":{\"host\":\"167.99.132.22\",\"port\":9001},\"status\":\"Ready\",\"reputation\":0}]")

	ci := node.ClusterInfo{}

	e := json.Unmarshal(test, &ci)

	if e != nil {
		t.Error(e)
	}
}