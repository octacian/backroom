package cmd

import (
	"net/http"

	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"github.com/go-chi/cors"
	"github.com/octacian/backroom/api/cage"
	"github.com/spf13/cobra"
)

func init() {
	// Add serve command to root command
	rootCmd.AddCommand(serveCmd)
}

var serveCmd = &cobra.Command{
	Use:   "serve",
	Short: "Start the API server",
	Run:   runServeCmd,
}

// runServeCmd implements the serve command.
func runServeCmd(cmd *cobra.Command, args []string) {
	// configure chi router
	r := chi.NewRouter()

	// Basic middleware stack
	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)

	// Basic CORS
	r.Use(cors.Handler(cors.Options{
		AllowedOrigins: []string{"https://*", "http://*"},
		AllowedMethods: []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders: []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token"},
		ExposedHeaders: []string{"Link"},
		// Debug:            true,
		MaxAge: 300, // Maximum value not ignored by any of major browsers
	}))

	// Set up routes
	r.Post("/record/create/{key}", cage.HandleCreateRecord)
	r.Get("/record/{uuid}", cage.HandleGetRecord)
	r.Get("/cage/{key}", cage.HandleListRecordsByKey)
	r.Get("/keys", cage.HandleListKeys)
	r.Delete("/record/{uuid}", cage.HandleDeleteRecord)
	r.Delete("/cage/{key}", cage.HandleDeleteRecordsByKey)
	r.Get("/health", handleHealthCheck)
}

// handleHealthCheck is a simple health check endpoint.
func handleHealthCheck(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("OK"))
}
