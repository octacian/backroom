package config

import (
	"log/slog"

	"github.com/expr-lang/expr"
	"github.com/expr-lang/expr/vm"
	"github.com/go-playground/validator/v10"
	"github.com/spf13/viper"
)

// HookEnv defines the environment variables used to evaluate hook expressions.
type HookEnv struct {
	Cage map[string]any `expr:"cage"`
}

// Hook defines a hook configuration for a cage.
type Hook struct {
	// Actions are any actions that triggers the hook.
	// Valid values are "create", "update", "delete".
	Action []string `mapstructure:"action" validate:"gt=0,dive,oneof=create update delete"`

	// If is an optional condition that must be met for the hook to run.
	// Cage data is available in the context as `cage.*`. See Expr for syntax.
	// https://expr-lang.org/docs/language-definition
	If        string `mapstructure:"if"`
	ifProgram *vm.Program

	// Cage is the identifier key for a cage.
	Cage string `mapstructure:"cage" validate:"required"`

	// Adapter is the name of the adapter to use for the hook.
	// Valid values are "log", "smtp".
	Adapter string `mapstructure:"adapter" validate:"oneof=log smtp"`

	// Target is the target log prefix or email for the hook.
	Target string `mapstructure:"target" validate:"required"`
}

// Eval returns whether or not the hook condition is met.
func (h *Hook) Eval(cageData map[string]any) (bool, error) {
	if h.ifProgram == nil {
		return true, nil // No condition, always true
	}

	// Populate the environment with the cage data
	env := &HookEnv{
		Cage: cageData,
	}

	// Evaluate the expression
	output, err := expr.Run(h.ifProgram, env)
	if err != nil {
		return false, err
	}

	// Ensure the output is a boolean
	ok := output.(bool)

	return ok, nil
}

// Config defines the expected environment variables (see .env.example.yml)
type Config struct {
	// Environment is the deployment environment.
	// Valid values are "development", and "production".
	Environment string `mapstructure:"environment" validate:"oneof=development production"`

	// AppName is the user friendly name of the deployment (default: "backroom").
	AppName string `mapstructure:"app_name"`

	// AdminEmail is the administrative email address for the deployment,
	// may receive alerts and notifications.
	AdminEmail string `mapstructure:"admin_email" validate:"email,required"`

	// APIURL is the fully qualified URL that the API is served from.
	APIURL string `mapstructure:"api_url" validate:"http_url,required"`

	// APIListen is the address and port that the API listens on.
	APIListen string `mapstructure:"api_listen" validate:"hostname_port,required"`

	Database struct {
		// User is the username to connect to the database.
		User string `mapstructure:"user" validate:"required"`
		// Password is the password to connect to the database.
		Password string `mapstructure:"password" validate:"required"`
		// Host is the hostname of the database server.
		Host string `mapstructure:"host" validate:"hostname_port,required"`
		// Name is the name of the database to connect to.
		Name string `mapstructure:"name" validate:"required"`
		// MaxConns is the maximum number of connections to the database.
		MaxConns int `mapstructure:"max_conns" validate:"omitempty,min=1,max=1000"`
	} `mapstructure:"database"`

	Mail struct {
		// FromName is the name that will be used in the "From" field of emails.
		// Defaults to the AppName if unset.
		FromName string `mapstructure:"from_name"`

		// FromAddress is the email address that will be used in the "From" field of emails.
		// Defaults to the AdminEmail if unset.
		FromAddress string `mapstructure:"from_address" validate:"omitempty,email"`

		SMTP struct {
			// Host is the hostname of the SMTP server.
			Host string `mapstructure:"host" validate:"omitempty,hostname"`
			// Port is the port of the SMTP server.
			Port int `mapstructure:"port" validate:"omitempty,min=1,max=65535"`

			// Username is the username to authenticate with the SMTP server for sending mail.
			Username string `mapstructure:"username"`
			// Password is the password to authenticate with the SMTP server for sending mail.
			Password string `mapstructure:"password"`
			// TLS is whether to use TLS when connecting to the SMTP server.
			// Defaults to true if unset.
			TLS bool `mapstructure:"tls"`
			// SkipVerify is whether to skip TLS certificate verification when connecting to the SMTP server.
			// Defaults to false if unset.
			SkipVerify bool `mapstructure:"skip_verify"`
		} `mapstructure:"smtp"`
	} `mapstructure:"mail"`

	SendGrid struct {
		// APIKey is the SendGrid API key to use for sending mail.
		APIKey string `mapstructure:"api_key"`
	} `mapstructure:"sendgrid"`

	// DeliveryMethod is the method used to send mail.
	// Valid values are "smtp", "sendgrid", "log-only".
	// Defaults to "smtp" if unset.
	DeliveryMethod string `mapstructure:"delivery_method" validate:"oneof=smtp sendgrid log-only"`

	// Hooks stores hook configuration for caged records.
	Hooks []Hook `mapstructure:"hooks" validate:"dive"`
}

// RC stores the current runtime configuration.
var RC *Config = &Config{}

// Init sets up viper and unmarshals the primary configuration file.
func Init() {
	viper.SetConfigName(".env")
	viper.SetConfigType("yaml")
	viper.SetConfigFile(".env.yml")
	if err := viper.ReadInConfig(); err != nil {
		panic(err)
	}
	if err := viper.Unmarshal(RC); err != nil {
		panic(err)
	}

	validate := validator.New(validator.WithRequiredStructEnabled())
	if err := validate.Struct(RC); err != nil {
		panic(err)
	}

	slog.Info("Viper loaded configuration", "file", viper.GetViper().ConfigFileUsed())

	// Compile any conditional hook expressions.
	for i, hook := range RC.Hooks {
		if hook.If != "" {
			program, err := expr.Compile(hook.If, expr.AsBool(), expr.Env(&HookEnv{}))
			if err != nil {
				panic(err)
			}
			RC.Hooks[i].ifProgram = program
		}
	}
}
