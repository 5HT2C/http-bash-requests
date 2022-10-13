package httpBashRequests

import (
	"bytes"
	"fmt"
	"io"
	"net/http"
	"sync"
)

var client *Client

type Client struct {
	Mutex      sync.Mutex
	Addr       string
	HttpClient *http.Client
}

// Setup will configure the client to use for Run
func Setup(clientIn *Client) {
	client = clientIn
}

// Run will run a bash command through HBR, and return the result of said command
func Run(str string) ([]byte, error) {
	return RunBinary(str, "", "", false)
}

// RunBinary will run binPath with binArg through HBR, and return the result of said command
func RunBinary(str, binPath, binArg string, splitBody bool) ([]byte, error) {
	client.Mutex.Lock()
	defer client.Mutex.Unlock()

	req, err := http.NewRequest("GET", client.Addr, bytes.NewBuffer([]byte(str)))
	if err != nil {
		return nil, err
	}

	req.Header.Set("X-Bin-Path", binPath)
	req.Header.Set("X-Bin-Arg", binArg)
	req.Header.Set("X-Body-Split", fmt.Sprintf("%v", splitBody))

	res, err := client.HttpClient.Do(req)
	if err != nil {
		return nil, err
	}

	if res.Body != nil {
		defer res.Body.Close()
	}

	resBody, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}

	return resBody, nil
}
