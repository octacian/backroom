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
	// Run executes the adapter with the given hook and record.
	Run(hook *Hook, record *cage.Record) error
}

// ALLOWED_ADAPTERS is a map of allowed hook adapter names to their respective
// adapter implementations.
var ALLOWED_ADAPTERS = map[string]Adapter{
	"log": &LogAdapter{},
}

// InitAdapters initializes any adapters requiring dynamic configuration.
func InitAdapters() {
	if config.RC.Mail.SMTP.Host != "" {
		smtpAdapter, err := NewSMTPAdapter()
		if err != nil {
			slog.Error("Failed to create SMTP adapter", "error", err)
			os.Exit(1)
		}

		ALLOWED_ADAPTERS["smtp"] = smtpAdapter
		slog.Info("SMTP adapter enabled", "host", config.RC.Mail.SMTP.Host)
	} else {
		slog.Debug("SMTP adapter disabled", "host", config.RC.Mail.SMTP.Host)
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

// LogAdapter is an adapter that logs the created record to the console.
type LogAdapter struct{}

// Run executes the LogAdapter with the given hook and record.
func (a *LogAdapter) Run(hook *Hook, record *cage.Record) error {
	slog.Info("Caged record created", "key", record.Cage, "uuid", record.UUID)
	return nil
}
