package port

import (
	"context"
	"io"
	"time"
)

type AgentPort interface {
	CollectLog(ctx context.Context, agentURL, path string, duration time.Duration) (io.ReadCloser, error)
	PProf(ctx context.Context, agentURL string, duration time.Duration) (io.ReadCloser, error)
}
