package standalone

import (
	"context"
	"fmt"

	"github.com/karamaru-alpha/isumaru/pkg/agent"
)

func Run(ctx context.Context, port int) {
	config := &agent.Config{
		Port: fmt.Sprintf("%d", port),
	}
	agent.Serve(ctx, config)
}
