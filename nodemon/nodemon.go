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
	ExecuteNodesCheck(url string, configFile string, statusFile string, outputTheme string)
}

type nodemon struct {}

func NewNodemon() Nodemon {
	return & nodemon {}
}

func (* nodemon) ExecuteNodesCheck(url string, configFile string, statusFile string, outputTheme string) {
	configFileBytes, err := ioutil.ReadFile(configFile)

	if err != nil {
		log.Fatal(err)
	}

	statusFileBytes, _ := ioutil.ReadFile(statusFile)

	webhookUrl := strings.TrimRight(string(configFileBytes), "\n")

	ng := nodegrid.NewNodegrid()

	imageFile := fmt.Sprintf("%s/nodemon", os.TempDir())

	err, networkStatus := ng.BuildNetworkStatus(url, true, imageFile, outputTheme)

	var importantObservations []string

	for _, row := range networkStatus.NodesGrid {
		for _, cell := range row {
			if cell.Status != node.Ready && cell.Status != node.SnapshotCreation{
				importantObservations = append(importantObservations, fmt.Sprintf("%s", cell.Ip.Host))
			}
		}
	}

	if err == nil {
		hashCalculator := sha256.New()

		currentHash := fmt.Sprintf("%x", hashCalculator.Sum([]byte(fmt.Sprintf("%v", importantObservations))))
		oldHash := string(statusFileBytes)

		if strings.Compare(currentHash, oldHash) != 0{


			fmt.Printf("Notify webook %s\n", webhookUrl)

			imageFileBytes, _ := ioutil.ReadFile(imageFile)

			client := resty.New()
			r, e := client.R().
				SetFileReader("file", "nodegrid.png", bytes.NewReader(imageFileBytes)).Post(webhookUrl)

			if e != nil || r.StatusCode() == 400 {
				fmt.Printf("Cannot execute webhook notification, error=%s\n", e)
			} else {
				fmt.Printf("Webhook returned %d\n", r.StatusCode())
			}

			ioutil.WriteFile(statusFile, []byte(currentHash), 0660)
		} else {
			fmt.Printf("Network status unchanged\n")
		}

		os.Remove(imageFile)
	}
}