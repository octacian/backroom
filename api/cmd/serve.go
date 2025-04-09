package cmd

import (
	"log/slog"
	"net/http"
	"os"

	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"github.com/go-chi/cors"
	"github.com/octacian/backroom/api/config"
	"github.com/octacian/backroom/api/httphandle"
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
	r.Post("/record/create", httphandle.HandleCreateRecord)
	r.Get("/record/{uuid}", httphandle.HandleGetRecord)
	r.Get("/cage/{key}", httphandle.HandleListRecordsByCage)
	r.Get("/cages", httphandle.HandleListCages)
	r.Delete("/record/{uuid}", httphandle.HandleDeleteRecord)
	r.Delete("/cage/{key}", httphandle.HandleDeleteRecordsByKey)
	r.Get("/health", handleHealthCheck)

	// Start the server
	slog.Info("Listening", "address", config.RC.APIListen, "url", config.RC.APIURL)

	err := http.ListenAndServe(config.RC.APIListen, r)
	if err != nil {
		slog.Error("Failed to start server", "error", err)
		os.Exit(1)
	}
}

// handleHealthCheck is a simple health check endpoint.
func handleHealthCheck(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("OK"))
}
