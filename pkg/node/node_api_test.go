package node

import (
	"encoding/json"
	"testing"
)

func TestPeers(t *testing.T) {

	jsonBytes := []byte(`[
		 {
			"id": "32b7be1864798fba069e500900d1db122c6a51bc45aee391a3c781b96af7b6bf10d732cdc77c79021a7b3e2f22bed2f5ae3544d904a44d995576055e8c1c8aad",
			"ip": "128.199.7.18",
			"publicPort": 9000,
			"p2pPort": 9001,
			"session": "9c3bbfac-c3dd-4546-9593-c94e39a2726e",
			"state": "SessionStarted"
		  }
		]`)

	d := Peers{}

	e := json.Unmarshal(jsonBytes, &d)

	if e != nil {
		t.Error(e)
	}
}
