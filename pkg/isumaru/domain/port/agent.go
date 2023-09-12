package port

import (
	"golang.org/x/net/context"
)

type AgentPort interface {
	CollectSlowQueryLog(ctx context.Context, seconds int32, path string) ([]byte, error)
}
