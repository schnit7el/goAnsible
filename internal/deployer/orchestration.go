package deployer

import (
	"fmt"
	"time"

	"github.com/schnit7el/goAnsible/internal/config"
)

func RunDeployment(cfg *config.Config) error {
	for _, node := range cfg.Nodes {
		client, err := NewSSHClient(node)
		if err != nil {
			return fmt.Errorf("node connection failed: %w", err)
		}
		defer client.Close()

		for _, task := range cfg.Tasks {
			handler, exists := taskHandlers[task.Type]
			if !exists {
				return fmt.Errorf("unknown task type: %s", task.Type)
			}

			if err := handler.Validate(task); err != nil {
				return fmt.Errorf("task validation failed: %w", err)
			}

			var lastErr error
			for attempt := 0; attempt <= task.Retries; attempt++ {
				if attempt > 0 {
					time.Sleep(task.Delay)
				}

				if err := handler.Execute(client, task); err == nil {
					lastErr = nil
					break
				} else {
					lastErr = err
				}
			}

			if lastErr != nil && !task.IgnoreError {
				return fmt.Errorf("task failed: %w", lastErr)
			}
		}
	}
	return nil
}
