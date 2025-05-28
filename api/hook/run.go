package hook

import (
	"log/slog"
	"slices"

	"github.com/octacian/backroom/api/cage"
)

type Action string

const (
	ActionCreate Action = "create"
	ActionUpdate Action = "update"
	ActionDelete Action = "delete"
)

// RunHooksByAction runs all hooks for a particular action.
func RunHooksByAction(act Action, record *cage.Record) error {
	// Get all hooks for the record's cage
	hooks, err := ListHooksByCage(record.Cage)
	if err != nil {
		return err
	}
	slog.Debug("Running hooks", "action", act, "cage", record.Cage, "hooks", hooks)

	for _, hook := range hooks {
		// Make sure the hook action matches the action we're running
		if !slices.Contains(hook.Action, string(act)) {
			continue // Skip this hook if the action does not match
		}

		// Check if the hook condition is met
		data := record.Data.ToMap()
		ok, err := hook.Eval(data)
		if err != nil {
			slog.Error("Failed to evaluate hook condition", "hook", hook, "error", err)
			return err
		}

		if !ok {
			slog.Debug("Hook condition not met, skipping", "hook", hook)
			continue // Skip this hook if the condition is not met
		}

		// Get the adapter for the hook
		adapter, err := GetAdapter(hook.Adapter)
		if err != nil {
			return err
		}

		// Run the adapter with the hook and record
		if err := adapter.Run(act, &hook, record); err != nil {
			slog.Error("Failed to run create hook", "hook", hook, "error", err)
			return err
		}

		slog.Info("Hook executed successfully", "hook", hook, "cage", record.Cage, "uuid", record.UUID)
	}

	return nil
}
