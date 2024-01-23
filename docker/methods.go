package docker

import (
	"bytes"
	"encoding/json"
	"io"
	"log"
	"net"
	"net/http"
	"os"
	"os/exec"
	"sync"
	"time"
)

func StartContainer(w http.ResponseWriter, r *http.Request) {
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
		log.Println("Error running Docker Compose:", err)
		payload := StartContainerResponse{
			Ok:    false,
			Error: err.Error(),
		}
		writeJSONResponse(w, 500, payload)
		return
	}
	payload := StartContainerResponse{
		Ok: true,
	}
	writeJSONResponse(w, 200, payload)

}

func WaitForService(w http.ResponseWriter, r *http.Request) {
	ticker := time.NewTicker(2 * time.Second)
	tickerChannel := ticker.C
	defer ticker.Stop()

	var wg sync.WaitGroup
	defer wg.Wait()
	wg.Add(1)

	go func() {
		defer wg.Done()
		for {
			select {
			case <-time.After(timeoutDuration):
				writeJSONResponse(w, http.StatusInternalServerError, StopContainerResponse{
					Ok:    false,
					Error: "timeout waiting for service",
				})
			case <-tickerChannel:
				conn, err := net.Dial("tcp", "localhost:"+"3001")
				if err == nil {
					// Service is up, close the connection and return
					defer conn.Close()
					log.Println("Service is up!")
					payload := StopContainerResponse{
						Ok: true,
					}
					writeJSONResponse(w, 200, payload)
					ticker.Stop()
					time.Sleep(2 * time.Second)
					return
				}
			}
		}
	}()
}

func StopContainer(w http.ResponseWriter, r *http.Request) {
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
		log.Println("Error running Docker Compose:", err)
		// Define the JSON payload.
		payload := StopContainerResponse{
			Ok:    false,
			Error: err.Error(),
		}
		writeJSONResponse(w, 500, payload)
		return
	}

	payload := StopContainerResponse{
		Ok: true,
	}
	writeJSONResponse(w, 200, payload)
}

func CheckAvailableStatus(w http.ResponseWriter, r *http.Request) {
	// Define the JSON payload.
	payload := CheckAvailableRequest{
		Method:  "eth_blockNumber",
		Params:  []string{},
		ID:      1,
		Jsonrpc: "2.0",
	}
	jsonBytes, err := json.Marshal(payload)
	if err != nil {
		log.Println(err.Error())
	}
	// Create a new HTTP request.
	req, err := http.NewRequest("POST", checkTaikoNodeApiURL, bytes.NewBuffer(jsonBytes))
	if err != nil {
		log.Println("Error creating HTTP request:", err)
		return
	}
	// Set the Content-Type header to application/json.
	req.Header.Set("Content-Type", "application/json")

	// Make the HTTP request.
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		log.Println("Error making POST request:", err)
		payload := CheckAvailableStatusResponse{
			Ok:      false,
			Result:  "",
			ID:      -1,
			Jsonrpc: "",
			Error:   "node is down",
		}
		writeJSONResponse(w, 200, payload)
		return
	}
	defer resp.Body.Close()

	// Read the response body.
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Println("Error reading response body:", err)
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
	writeJSONResponse(w, 200, success_payload)
}
