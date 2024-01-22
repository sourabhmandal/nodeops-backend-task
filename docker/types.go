package docker

type CheckAvailableRequest struct {
	Method  string   `json:"method"`
	Params  []string `json:"params"`
	ID      int      `json:"id"`
	Jsonrpc string   `json:"jsonrpc"`
}

type StartContainerResponse struct {
	Ok    bool   `json:"ok"`
	Error string `json:"error,omitempty"`
}

type StopContainerResponse struct {
	Ok    bool   `json:"ok"`
	Error string `json:"error,omitempty"`
}

type WaitContainerResponse struct {
	Loaded bool   `json:"loaded"`
	Error  string `json:"error,omitempty"`
}

type CheckAvailableStatusResponse struct {
	Ok      bool   `json:"ok"`
	Jsonrpc string `json:"jsonrpc,omitempty"`
	ID      int    `json:"id,omitempty"`
	Result  string `json:"result,omitempty"`
	Error   string `json:"error,omitempty"`
}
