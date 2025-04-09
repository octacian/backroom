package hook

import (
	"errors"
	"log/slog"
	"os"

	"github.com/octacian/backroom/api/cage"
	"github.com/octacian/backroom/api/config"
)

var ErrBadAdapter = errors.New("bad adapter")

// Adapter defines expected execution methods for a hook adapter.
type Adapter interface {
	// Run executes the adapter with the given hook and cage.
	Run(hook *Hook, cage *cage.Cage) error
}

// ALLOWED_ADAPTERS is a map of allowed hook adapter names to their respective
// adapter implementations.
var ALLOWED_ADAPTERS = map[string]Adapter{
	"log": &LogAdapter{},
}

func init() {
	if config.RC.Mail.SMTP.Host != "" {
		smtpAdapter, err := NewSMTPAdapter()
		if err != nil {
			slog.Error("Failed to create SMTP adapter", "error", err)
			os.Exit(1)
		}

		ALLOWED_ADAPTERS["smtp"] = smtpAdapter
		slog.Info("SMTP adapter enabled", "host", config.RC.Mail.SMTP.Host)
	}
}

// GetAdapter returns an adapter by name if it exists in ALLOWED_ADAPTERS.
// Returns ErrBadAdapter if the adapter is not found.
func GetAdapter(name string) (Adapter, error) {
	adapter, ok := ALLOWED_ADAPTERS[name]
	if !ok {
		return nil, ErrBadAdapter
	}
	return adapter, nil
}

// LogAdapter is an adapter that logs the created cage to the console.
type LogAdapter struct{}

// Run executes the LogAdapter with the given hook and cage.
func (a *LogAdapter) Run(hook *Hook, cage *cage.Cage) error {
	slog.Info("Caged record created", "key", cage.Key, "uuid", cage.UUID)
	return nil
}
