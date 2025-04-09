package hook

import (
	"fmt"

	"github.com/octacian/backroom/api/config"
)

// Hook defines a hook configuration for a cage.
type Hook config.Hook

// ListHooksByKey retrieves all hooks for a given key from the configuration.
func ListHooksByKey(key string) ([]Hook, error) {
	var hooks []Hook
	for _, hook := range config.RC.Hooks {
		if hook.Key == key {
			hooks = append(hooks, Hook(hook))
		}
	}
	if len(hooks) == 0 {
		return nil, fmt.Errorf("no hooks found for key: %s", key)
	}
	return hooks, nil
}
