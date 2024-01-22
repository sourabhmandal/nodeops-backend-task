package docker

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net"
	"net/http"
	"os"
	"os/exec"
	"time"
)

func StartContainer(w http.ResponseWriter, r *http.Request) {
	// Specify the path to your Docker Compose file.
	composeFilePath := "simple-taiko-node/docker-compose.yml"

	// Run the "docker-compose" command with the "up" subcommand.
	cmd := exec.Command("docker-compose", "-f", composeFilePath, "up", "-d")

	// Set the working directory for the command.
	cmd.Dir, _ = os.Getwd()

	// Redirect standard output and standard error to the Go program's standard output.
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	// Run the command.
	err := cmd.Run()
	if err != nil {
		fmt.Println("Error running Docker Compose:", err)
		payload := StartContainerResponse{
			Ok:    false,
			Error: err.Error(),
		}
		jsonBytes, err := json.Marshal(payload)
		if err != nil {
			payload.Error = err.Error()
		}

		w.WriteHeader(500)
		w.Write(jsonBytes)
		return
	}
	payload := StartContainerResponse{
		Ok: true,
	}
	jsonBytes, err := json.Marshal(payload)
	if err != nil {
		payload.Error = err.Error()
	}

	w.WriteHeader(200)
	w.Write(jsonBytes)
}

func WaitForService(w http.ResponseWriter, r *http.Request) {
	timeout := time.After(5 * time.Minute) // You can adjust the timeout as needed
	ticker := time.NewTicker(2 * time.Second)
	tickerChannel := ticker.C

	for {
		select {
		case <-timeout:
			payload := StopContainerResponse{
				Ok:    false,
				Error: "timeout waiting for service",
			}
			jsonBytes, err := json.Marshal(payload)
			if err != nil {
				payload.Error = err.Error()
			}

			w.WriteHeader(500)
			w.Write(jsonBytes)
			return
		case <-tickerChannel:
			conn, err := net.Dial("tcp", "localhost:"+"3001")
			if err == nil {
				// Service is up, close the connection and return
				conn.Close()
				fmt.Println("Service is up!")
				payload := StopContainerResponse{
					Ok: true,
				}
				jsonBytes, err := json.Marshal(payload)
				if err != nil {
					payload.Error = err.Error()
				}

				w.WriteHeader(200)
				w.Write(jsonBytes)
				ticker.Stop()
				time.Sleep(2 * time.Second)
				return
			}
		}
	}

}

func StopContainer(w http.ResponseWriter, r *http.Request) {
	// Specify the path to your Docker Compose file.
	composeFilePath := "simple-taiko-node/docker-compose.yml"

	// Run the "docker-compose" command with the "up" subcommand.
	cmd := exec.Command("docker-compose", "-f", composeFilePath, "down")

	// Set the working directory for the command.
	cmd.Dir, _ = os.Getwd()

	// Redirect standard output and standard error to the Go program's standard output.
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	// Run the command.
	err := cmd.Run()
	if err != nil {
		fmt.Println("Error running Docker Compose:", err)
		// Define the JSON payload.
		payload := StopContainerResponse{
			Ok:    false,
			Error: err.Error(),
		}
		jsonBytes, err := json.Marshal(payload)
		if err != nil {
			payload.Error = err.Error()
		}

		w.WriteHeader(500)
		w.Write(jsonBytes)
		return
	}

	payload := StopContainerResponse{
		Ok: true,
	}
	jsonBytes, err := json.Marshal(payload)
	if err != nil {
		log.Fatal(err.Error())
	}

	w.Write(jsonBytes)
}

func CheckAvailableStatus(w http.ResponseWriter, r *http.Request) {
	// Specify the API endpoint URL.
	apiURL := "http://localhost:8547" // Replace with your API URL

	// Define the JSON payload.
	payload := CheckAvailableRequest{
		Method:  "eth_blockNumber",
		Params:  []string{},
		ID:      1,
		Jsonrpc: "2.0",
	}
	jsonBytes, err := json.Marshal(payload)
	if err != nil {
		log.Fatal(err.Error())
	}
	// Create a new HTTP request.
	req, err := http.NewRequest("POST", apiURL, bytes.NewBuffer(jsonBytes))
	if err != nil {
		fmt.Println("Error creating HTTP request:", err)
		return
	}
	// Set the Content-Type header to application/json.
	req.Header.Set("Content-Type", "application/json")

	// Make the HTTP request.
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		fmt.Println("Error making POST request:", err)
		payload := CheckAvailableStatusResponse{
			Ok:      false,
			Result:  "",
			ID:      -1,
			Jsonrpc: "",
			Error:   "node is down",
		}
		jsonBytes, err := json.Marshal(payload)
		if err != nil {
			log.Fatal(err.Error())
		}
		w.WriteHeader(200)
		w.Header().Set("Content-Type", "application/json")
		w.Write(jsonBytes)
		return
	}
	defer resp.Body.Close()

	// Read the response body.
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		fmt.Println("Error reading response body:", err)
		return
	}

	var apiResponse CheckAvailableStatusResponse

	json.Unmarshal(body, &apiResponse)

	success_payload := CheckAvailableStatusResponse{
		Ok:      true,
		Result:  apiResponse.Result,
		ID:      apiResponse.ID,
		Jsonrpc: apiResponse.Jsonrpc,
	}
	jsonBytes, err = json.Marshal(success_payload)
	if err != nil {
		log.Fatal(err.Error())
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(200)
	w.Write(jsonBytes)
}
