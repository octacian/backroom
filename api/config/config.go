package config

import (
	"log/slog"

	"github.com/go-playground/validator/v10"
	"github.com/spf13/viper"
)

// Config defines the expected environment variables (see .env.example.yml)
type Config struct {
	// Environment is the deployment environment.
	// Valid values are "development", and "production".
	Environment string `validate:"oneof=development production"`

	// AppName is the user friendly name of the deployment (default: "backroom").
	AppName string

	// AdminEmail is the administrative email address for the deployment,
	// may receive alerts and notifications.
	AdminEmail string `validate:"email,required"`

	// APIURL is the fully qualified URL that the API is served from.
	APIURL string `validate:"http_url,required"`

	// APIListen is the address and port that the API listens on.
	APIListen string `validate:"hostname_port,required"`

	Database struct {
		// User is the username to connect to the database.
		User string `validate:"required"`
		// Password is the password to connect to the database.
		Password string `validate:"required"`
		// Host is the hostname of the database server.
		Host string `validate:"hostname_port,required"`
		// Name is the name of the database to connect to.
		Name string `validate:"required"`
		// MaxConns is the maximum number of connections to the database.
		MaxConns int `validate:"min=1,max=1000"`
	}

	Mail struct {
		// FromName is the name that will be used in the "From" field of emails.
		// Defaults to the AppName if unset.
		FromName string

		// FromAddress is the email address that will be used in the "From" field of emails.
		// Defaults to the AdminEmail if unset.
		FromAddress string `validate:"email"`

		SMTP struct {
			// Host is the hostname of the SMTP server.
			Host string `validate:"hostname,required"`
			// Port is the port of the SMTP server.
			Port string `validate:"min=1,max=65535,required"`

			// Username is the username to authenticate with the SMTP server for sending mail.
			Username string
			// Password is the password to authenticate with the SMTP server for sending mail.
			Password string
			// TLS is whether to use TLS when connecting to the SMTP server.
			// Defaults to true if unset.
			TLS bool
			// SkipVerify is whether to skip TLS certificate verification when connecting to the SMTP server.
			// Defaults to false if unset.
			SkipVerify bool
		}
	}

	SendGrid struct {
		// APIKey is the SendGrid API key to use for sending mail.
		APIKey string
	}

	// DeliveryMethod is the method used to send mail.
	// Valid values are "smtp" and "sendgrid".
	// Defaults to "smtp" if unset.
	DeliveryMethod string `validate:"oneof=smtp sendgrid"`
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
}
