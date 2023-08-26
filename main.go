package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"strconv"
	"time"
)

type Node struct {
	ID      string `json:"ID"`
	Name    string `json:"Name"`
	Address string `json:"Address"`
}

var (
	apiURL   = ""
	interval = 300
	nodePort = "4646"
	filePath = "nomad-nodes.yaml"
)

func main() {

	err := validate()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	fmt.Printf("Writing to %s every %d seconds...\n", filePath, interval)

	err = run()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func validate() (err error) {
	nodePortInput := os.Getenv("NOMAD_NODE_PORT")
	if nodePortInput != "" {
		nodePort = nodePortInput
	}

	apiURLInput := os.Getenv("NOMAD_API_URL")
	if apiURLInput == "" {
		return fmt.Errorf("NOMAD_API_URL environment variable not set.")
	} else {
		apiURL = apiURLInput
	}

	intervalInput := os.Getenv("REFRESH_INTERVAL")
	if intervalInput != "" {
		interval, err = strconv.Atoi(intervalInput)
		if err != nil {
			return fmt.Errorf("Invalid REFRESH_INTERVAL. %s", err)
		}
	}

	filePathInput := os.Getenv("OUTPUT_FILE_PATH")
	if filePathInput != "" {
		filePath = filePathInput
	}

	return nil
}

func run() error {
	for {
		// Make an HTTP GET request to the Nomad API
		resp, err := http.Get(fmt.Sprintf("%s/v1/nodes", apiURL))
		if err != nil {
			return fmt.Errorf("Error making the request: %s", err)
		}
		defer resp.Body.Close()

		// Read the response body
		var nodes []Node
		if err := json.NewDecoder(resp.Body).Decode(&nodes); err != nil {
			return fmt.Errorf("Error decoding response: %s", err)
		}

		// Prepare the YAML content
		var yamlContent string
		for _, node := range nodes {
			yamlContent += fmt.Sprintf(
				"- targets: ['%s:%s']\n  labels:\n    node_id: '%s'\n    node_name: '%s'\n",
				node.Address, nodePort, node.ID, node.Name,
			)
		}

		// Write the YAML content to the specified file
		err = os.WriteFile(filePath, []byte(yamlContent), 0644)
		if err != nil {
			return fmt.Errorf("Error writing to file: %s", err)
		}

		// Sleep for 5 minutes before the next update
		time.Sleep(time.Duration(interval) * time.Second)
	}
}
