package nodemon

import (
	"bytes"
	"constellation_cli/nodegrid"
	"constellation_cli/pkg/node"
	"crypto/sha256"
	"fmt"
	"github.com/go-resty/resty/v2"
	"io/ioutil"
	"log"
	"os"
	"strings"
)

type Nodemon interface {
	ExecuteNodesCheck(url string, configFile string, statusFile string, outputTheme string, operatorsFile string)
}

type nodemon struct {}

func NewNodemon() Nodemon {
	return & nodemon {}
}

func (* nodemon) ExecuteNodesCheck(url string, configFile string, statusFile string, outputTheme string, operatorsFile string) {
	configFileBytes, err := ioutil.ReadFile(configFile)

	if err != nil {
		log.Fatal(err)
	}

	ng := nodegrid.NewNodegrid(operatorsFile)

 	nodeOps := ng.Operators()

	statusFileBytes, _ := ioutil.ReadFile(statusFile)

	webhookUrl := strings.TrimRight(string(configFileBytes), "\n")

	imageFile := fmt.Sprintf("%s/nodemon", os.TempDir())

	err, networkStatus := ng.BuildNetworkStatus(url, true, imageFile, outputTheme, false)

	var offlineObservations []string

	offlineNodeOperators := make(map[nodegrid.Operator]bool)

	var offlineNodesObservationCount = 0
	var redownloadNodesSelfObservationCount = 0

	for observer, row := range networkStatus.NodesGrid {

		s := row[observer]

		if node.IsRedownloading(s.Status) == true{
			redownloadNodesSelfObservationCount++
		}

		for _, cell := range row {
			if node.IsOffline(cell.Status) {
				offlineObservations = append(offlineObservations, fmt.Sprintf("%s=%s;%s", observer, cell.Id.Hex,cell.Status))
				offlineNodesObservationCount++

				if op, v := nodeOps[cell.Id.Hex]; v{
					offlineNodeOperators[op] = true
				}
			}
		}
	}

	offlineNodeMentions := make([]string, 0, len(offlineNodeOperators))

	for op, _ := range offlineNodeOperators {

		fmt.Printf("Notify %s - %s\n", op.DiscordId, op.Name)
		offlineNodeMentions = append(offlineNodeMentions, fmt.Sprintf("<@%s>", op.DiscordId))
	}

	if err == nil {
		hashCalculator := sha256.New()
		obsString:= strings.Join(offlineObservations, "\n")

		currentHash := fmt.Sprintf("%x", hashCalculator.Sum([]byte(obsString)))
		oldHash := string(statusFileBytes)

		redownloadScale := 100 * float64(redownloadNodesSelfObservationCount) / float64(len(networkStatus.NodesList))

		redownloadTriggerReached := redownloadScale > 50

		if strings.Compare(currentHash, oldHash) != 0 || redownloadTriggerReached {

			var message = ""

			if redownloadTriggerReached {
				message = fmt.Sprintf("According to results %.2f%% of the cluster is performing a redownload. \n", redownloadScale)
			}

			if len(offlineNodeMentions) > 0 {
				message = message + fmt.Sprintf("Operators %s, we need you since your nodes are offline or marked as offline.", strings.Join(offlineNodeMentions, ", "))
			}

			fmt.Printf("Notify following operators %s via webhook\n", strings.Join(offlineNodeMentions, ", "))

			imageFileBytes, _ := ioutil.ReadFile(imageFile)

			var content = ""
			if len(offlineNodeMentions) > 0 {
				content = message
			}

			client := resty.New()

			r, e := client.R().
				SetFormData(map[string]string{
					"content": content,
				}).
				SetFileReader("file", "nodegrid.png", bytes.NewReader(imageFileBytes)).Post(webhookUrl)

			if e != nil || r.StatusCode() == 400 {
				fmt.Printf("Cannot execute webhook notification, error=%s\n", e)
			} else {
				fmt.Printf("Webhook returned %d\n", r.StatusCode())
			}

			ioutil.WriteFile(statusFile, []byte(currentHash), 0660)
		} else {
			fmt.Printf("Network offline status unchanged or redownload alert trigger not met %.2f%%\n", redownloadScale)
		}

		os.Remove(imageFile)
	}
}