package port

import (
	"context"
	"io"
	"time"
)

type AgentPort interface {
	CollectSlowQueryLog(ctx context.Context, agentURL, slowQueryLogPath string, duration time.Duration) (io.ReadCloser, error)
}
