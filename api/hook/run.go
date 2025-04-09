package hook

import (
	"log/slog"

	"github.com/octacian/backroom/api/cage"
)

// RunCreate runs all create hooks for a new caged record.
func RunCreate(cage *cage.Cage) error {
	// get the list of hooks for the key
	hooks, err := ListHooksByKey(cage.Key)
	if err != nil {
		return err
	}

	// run each hook
	for _, hook := range hooks {
		adapter, err := GetAdapter(hook.Adapter)
		if err != nil {
			return err
		}

		if err := adapter.Run(&hook, cage); err != nil {
			slog.Error("Failed to run create hook", "hook", hook, "error", err)
			return err
		}
	}

	return nil
}
