package port

import (
	"bytes"
	"compress/gzip"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/karamaru-alpha/isumaru/pkg/isumaru/domain/port"
)

type agentPort struct{}

func NewAgentPort() port.AgentPort {
	return &agentPort{}
}

type AgentCollectRequest struct {
	Seconds int32  `json:"seconds"`
	Path    string `json:"path"`
}

func (p *agentPort) CollectLog(_ context.Context, agentURL, path string, duration time.Duration) (io.ReadCloser, error) {
	// ブラウザのリロードでcancelされないように別のctxを使用する
	ctx := context.Background()
	// リクエスト作成
	requestBody, err := json.Marshal(&AgentCollectRequest{
		Seconds: int32(duration.Seconds()),
		Path:    path,
	})
	if err != nil {
		return nil, err
	}
	req, err := http.NewRequestWithContext(
		ctx,
		http.MethodPost,
		fmt.Sprintf("%s/debug/collect", agentURL),
		bytes.NewBuffer(requestBody))
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept-Encoding", "gzip")

	// リクエスト
	client := &http.Client{
		Transport: &http.Transport{
			DisableCompression: false,
		},
	}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	pr, pw := io.Pipe()
	go func() {
		defer resp.Body.Close()
		defer pw.Close()

		var reader io.ReadCloser
		switch resp.Header.Get("Content-Encoding") {
		case "gzip":
			reader, err = gzip.NewReader(resp.Body)
			if err != nil {
				pw.CloseWithError(err)
				return
			}
			defer reader.Close()
		default:
			reader = resp.Body
		}

		if _, err := io.Copy(pw, reader); err != nil {
			pw.CloseWithError(err)
		}
	}()
	if resp.StatusCode != http.StatusOK {
		return nil, errors.New(fmt.Sprintf("fail to collect log. err=%+v", resp.Body))
	}
	if resp.ContentLength == 0 {
		return nil, errors.New("log is empty")
	}

	return pr, nil
}
