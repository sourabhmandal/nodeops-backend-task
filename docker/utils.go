package docker

import (
	"encoding/json"
	"log"
	"net/http"
	"time"
)

const (
	composeFilePath      = "simple-taiko-node/docker-compose.yml"
	checkTaikoNodeApiURL = "http://localhost:8547"
	timeoutDuration      = 5 * time.Minute
)

func writeJSONResponse(w http.ResponseWriter, status int, payload interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	jsonBytes, err := json.Marshal(payload)
	if err != nil {
		log.Println("Error encoding JSON:", err)
		return
	}
	w.Write(jsonBytes)
}
