package hook

import (
	"log/slog"

	"github.com/octacian/backroom/api/cage"
)

// RunCreate runs all create hooks for a new record.
func RunCreate(record *cage.Record) error {
	// get the list of hooks for the key
	hooks, err := ListHooksByCage(record.Cage)
	if err != nil {
		return err
	}

	// run each hook
	for _, hook := range hooks {
		adapter, err := GetAdapter(hook.Adapter)
		if err != nil {
			return err
		}

		if err := adapter.Run(&hook, record); err != nil {
			slog.Error("Failed to run create hook", "hook", hook, "error", err)
			return err
		}
	}

	return nil
}
