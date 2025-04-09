package hook

import (
	"fmt"

	"github.com/octacian/backroom/api/config"
)

// Hook defines a hook configuration for a cage.
type Hook config.Hook

// ListHooksByCage retrieves all hooks for a given cage from the configuration.
func ListHooksByCage(cageKey string) ([]Hook, error) {
	var hooks []Hook
	for _, hook := range config.RC.Hooks {
		if hook.Key == cageKey {
			hooks = append(hooks, Hook(hook))
		}
	}
	if len(hooks) == 0 {
		return nil, fmt.Errorf("no hooks found for cage: %s", cageKey)
	}
	return hooks, nil
}
