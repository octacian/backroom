package hook

import (
	"github.com/octacian/backroom/api/config"
)

// Hook defines a hook configuration for a cage.
type Hook = config.Hook

// ListHooksByCage retrieves all hooks for a given cage from the configuration.
func ListHooksByCage(cageKey string) ([]Hook, error) {
	hooks := make([]Hook, 0)
	for _, hook := range config.RC.Hooks {
		if hook.Cage == cageKey {
			hooks = append(hooks, Hook(hook))
		}
	}
	return hooks, nil
}
