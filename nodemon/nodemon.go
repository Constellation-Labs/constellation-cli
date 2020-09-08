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
			if node.IsRedownloading(cell.Status) {
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

		if strings.Compare(currentHash, oldHash) != 0 || redownloadNodesSelfObservationCount > len(networkStatus.NodesList)/2 {

			fmt.Printf("Notify following operators %s via webhook\n", strings.Join(offlineNodeMentions, ", "))

			imageFileBytes, _ := ioutil.ReadFile(imageFile)

			var content = ""
			if len(offlineNodeMentions) > 0 {
				content = fmt.Sprintf("Operators %s, we need you since your nodes are offline.", strings.Join(offlineNodeMentions, ", "))
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
			fmt.Printf("Network offline status unchanged and no need to notify\n")
		}

		os.Remove(imageFile)
	}
}