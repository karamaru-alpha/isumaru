package port

import (
	"bytes"
	"compress/gzip"
	"context"
	"encoding/json"
	"io"
	"net/http"
	"strings"

	"github.com/karamaru-alpha/isumaru/pkg/isumaru/domain/port"
)

type agentPort struct{}

func NewAgentPort() port.AgentPort {
	return &agentPort{}
}

type AgentCollectSlowQueryLogRequest struct {
	Seconds int32  `json:"seconds"`
	Path    string `json:"path"`
}

func (p *agentPort) CollectSlowQueryLog(ctx context.Context, seconds int32, path string) ([]byte, error) {
	requestBody, err := json.Marshal(&AgentCollectSlowQueryLogRequest{
		Seconds: seconds,
		Path:    path,
	})
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, "http://localhost:19000/mysql/collect", bytes.NewBuffer(requestBody))
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept-Encoding", "gzip")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var bodyReader io.Reader = resp.Body
	if strings.Contains(resp.Header.Get("Content-Encoding"), "gzip") {
		gzipReader, err := gzip.NewReader(resp.Body)
		if err != nil {
			return nil, err
		}
		defer gzipReader.Close()

		bodyReader = gzipReader
	}

	bodyContent, err := io.ReadAll(bodyReader)
	if err != nil {
		return nil, err
	}
	return bodyContent, nil
}
