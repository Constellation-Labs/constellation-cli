package nodemon

import (
	"bytes"
	nodegrid2 "constellation/internal/cli/nodegrid"
	"constellation/pkg/node"
	"crypto/sha256"
	"fmt"
	"github.com/go-resty/resty/v2"
	log "github.com/sirupsen/logrus"
	"io/ioutil"
	"os"
	"sort"
	"strings"
	"time"
)

type Nodemon interface {
	ExecuteNodesCheck(url node.Addr, configFile string, statusFile string, outputTheme string, operatorsFile string)
}

type nodemon struct{}

func NewNodemon() Nodemon {
	return &nodemon{}
}

func (*nodemon) ExecuteNodesCheck(addr node.Addr, configFile string, statusFile string, outputTheme string, operatorsFile string) {
	configFileBytes, err := ioutil.ReadFile(configFile)

	if err != nil {
		log.Fatal(err)
	}

	ng := nodegrid2.NewNodegrid(operatorsFile)

	nodeOps := ng.Operators()

	statusFileBytes, _ := ioutil.ReadFile(statusFile)

	webhookUrl := strings.TrimRight(string(configFileBytes), "\n")

	imageFile := fmt.Sprintf("%s/nodemon", os.TempDir())

	fmt.Println("Gathering and building cluster status")

	err, networkStatus := ng.BuildNetworkStatus(addr, true, imageFile, outputTheme, false)

	if err != nil {
		log.Printf("Error building network status")
	}

	var offlineObservations []string

	offlineNodeOperators := make(map[nodegrid2.Operator]bool)
	offlineNodes := make(map[string]bool)

	var offlineNodesObservationCount = 0
	var redownloadNodesSelfObservationCount = 0
	var slowNodes []string
	var slowNodesOperators []nodegrid2.Operator

	fmt.Println("Verifying slow nodes")

	for _, n := range networkStatus.NodesList {
		if n.AvgResponseDuration.Milliseconds() > nodegrid2.LatencyTriggerMilliseconds &&
			n.AvgResponseDuration <= 29*time.Second {

			slowNodes = append(slowNodes, n.Info.Id)
			if op, v := nodeOps[n.Info.Id]; v {
				slowNodesOperators = append(slowNodesOperators, op)
			}
		}

		if n.AvgResponseDuration > 29*time.Second {
			offlineObservations = append(offlineObservations, fmt.Sprintf("%s=%s:%s", "Nodegrid", n.Info.Id, n.Info.CardinalState()))
			offlineNodesObservationCount++
			offlineNodes[n.Info.Id] = true

			if op, v := nodeOps[n.Info.Id]; v {
				offlineNodeOperators[op] = true
			}
		}
	}

	fmt.Println("Reviewing offline nodes and redownloads")

	for observer, row := range networkStatus.NodesGrid {

		s := row[observer]

		if node.IsRedownloading(s.CardinalState()) == true {
			redownloadNodesSelfObservationCount++
		}

		for _, cell := range row {
			if node.IsOffline(cell.CardinalState()) {
				offlineObservations = append(offlineObservations, fmt.Sprintf("%s=%s:%s", observer, cell.Id, cell.CardinalState()))
				offlineNodesObservationCount++
				offlineNodes[cell.Id] = true

				if op, v := nodeOps[cell.Id]; v {
					offlineNodeOperators[op] = true
				}
			}
		}
	}

	offlineNodeMentions := make([]string, 0, len(offlineNodeOperators))
	slowNodeMentions := make([]string, 0, len(slowNodesOperators))

	for op, _ := range offlineNodeOperators {
		fmt.Printf("Notifyabout offline nodes %s - %s\n", op.DiscordId, op.Name)
		offlineNodeMentions = append(offlineNodeMentions, fmt.Sprintf("<@%s>", op.DiscordId))
	}

	for _, op := range slowNodesOperators {
		fmt.Printf("Notify about slow nodes %s - %s\n", op.DiscordId, op.Name)
		slowNodeMentions = append(slowNodeMentions, fmt.Sprintf("<@%s>", op.DiscordId))
	}

	if err == nil {
		hashCalculator := sha256.New()
		sort.Strings(offlineObservations)

		obsString := strings.Join(offlineObservations, "\n")
		currentHash := fmt.Sprintf("%x", hashCalculator.Sum([]byte(obsString)))

		oldHash := string(statusFileBytes)

		redownloadScale := 100 * float64(redownloadNodesSelfObservationCount) / float64(len(networkStatus.NodesList))

		redownloadTriggerReached := redownloadScale > 50

		// nodegrid.PrintAsciiOutput(networkStatus.NodesList, networkStatus.NodesGrid, true)

		if strings.Compare(currentHash, oldHash) != 0 || redownloadTriggerReached || len(slowNodes) > 0 {

			if strings.Compare(currentHash, oldHash) != 0 {
				fmt.Println("Network status hash is different")
			}

			var message = fmt.Sprintf("Cluster is total nodes=%d offline/partially offline nodes=%d offline observations=%d redownload=%d highLatency=%d\n",
				len(networkStatus.NodesList),
				len(offlineNodes),
				offlineNodesObservationCount,
				redownloadNodesSelfObservationCount,
				len(slowNodes))

			if redownloadTriggerReached {
				message = message + fmt.Sprintf("According to results %.2f%% of the cluster is performing a redownload. \n", redownloadScale)
			}

			if len(offlineNodeMentions) > 0 {
				message = message + fmt.Sprintf("%s - your nodes are offline or not fully reachable.", strings.Join(offlineNodeMentions, ", "))
			}

			if len(slowNodeMentions) > 0 {
				message = message + fmt.Sprintf("%s - nodegrid recorded a high network latency for your node.", strings.Join(slowNodeMentions, ", "))
			}

			fmt.Printf("Notify following operators %s, %s via webhook\n", strings.Join(offlineNodeMentions, ", "),
				strings.Join(slowNodeMentions, ", "))

			imageFileBytes, _ := ioutil.ReadFile(imageFile)

			client := resty.New()

			r, e := client.R().
				SetFormData(map[string]string{
					"username": "Nodegrid",
					"content":  message,
				}).
				SetFileReader("file", "nodegrid.png", bytes.NewReader(imageFileBytes)).Post(webhookUrl)

			if e != nil || r.StatusCode() == 400 {
				fmt.Printf("Cannot execute webhook notification, error=%s\n", e)
			} else {
				fmt.Printf("Webhook returned %d\n", r.StatusCode())
				ioutil.WriteFile(statusFile, []byte(currentHash), 0660)
			}
		}
		fmt.Printf("Network size=%d offline/partially offline nodes=%d offline observations=%d status unchanged or redownload alert trigger not met %.2f%%\nOffline summary:\n%s\nSlow nodes:\n%s\n",
			len(networkStatus.NodesList), len(offlineNodes), offlineNodesObservationCount, redownloadScale, strings.Join(offlineNodeMentions, "\n"), strings.Join(slowNodes, "\n"))

		os.Remove(imageFile)
	}
}
