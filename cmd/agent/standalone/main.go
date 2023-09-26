package standalone

import (
	"fmt"

	"github.com/karamaru-alpha/isumaru/pkg/agent"
)

func Run(port int) {
	config := &agent.Config{
		Port: fmt.Sprintf("%d", port),
	}
	agent.Serve(config)
}
