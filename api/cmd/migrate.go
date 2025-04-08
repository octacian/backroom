package cmd

import (
	"log/slog"
	"os"
	"strconv"

	"github.com/octacian/backroom/api/db"
	"github.com/octacian/backroom/api/migrations"
	"github.com/pressly/goose/v3"
	"github.com/spf13/cobra"
)

func init() {
	// set goose to use embedded filesystem
	goose.SetBaseFS(migrations.Migrations)

	rootCmd.AddCommand(migrateCmd)
	migrateCmd.AddCommand(migrateStatusCmd)
	migrateCmd.AddCommand(migrateCreateCmd)
	migrateCreateCmd.Flags().BoolP("sequential", "s", false, "Create a sequential migration")
	migrateCmd.AddCommand(migrateUpCmd)
	migrateCmd.AddCommand(migrateUpToCmd)
	migrateCmd.AddCommand(migrateDownCmd)
	migrateCmd.AddCommand(migrateDownToCmd)
	migrateCmd.AddCommand(migrateRedoCmd)
}

var migrateCmd = &cobra.Command{
	Use:   "migrate",
	Short: "Migrate the database",
	Run: func(cmd *cobra.Command, args []string) {
		cmd.Help()
	},
}

var migrateStatusCmd = &cobra.Command{
	Use:   "status",
	Short: "Get database migration status",
	Run: func(cmd *cobra.Command, args []string) {
		db.InitDB()
		defer db.CloseDB()

		if err := goose.Status(db.SQLDB, "."); err != nil {
			slog.Error("Couldn't get migration status", "err", err)
			os.Exit(1)
		}
	},
}

var migrateCreateCmd = &cobra.Command{
	Use:   "create [NAME] [TYPE]",
	Short: "Create a new migration",
	Args:  cobra.ExactArgs(2),
	Run: func(cmd *cobra.Command, args []string) {
		sequential, err := cmd.Flags().GetBool("sequential")
		if err != nil {
			panic(err)
		}

		if sequential {
			goose.SetSequential(sequential)
		}

		if err := goose.Create(db.SQLDB, "migrations", args[0], args[1]); err != nil {
			slog.Error("Couldn't create migration", "err", err)
			os.Exit(1)
		}

	},
}

var migrateUpCmd = &cobra.Command{
	Use:   "up",
	Short: "Apply all available database migrations",
	Run: func(cmd *cobra.Command, args []string) {
		db.InitDB()
		defer db.CloseDB()

		if err := goose.Up(db.SQLDB, "."); err != nil {
			slog.Error("Couldn't apply migrations", "err", err)
			os.Exit(1)
		}
	},
}

var migrateUpToCmd = &cobra.Command{
	Use:   "up-to [VERSION]",
	Short: "Apply all available database migrations up to a specific version",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		db.InitDB()
		defer db.CloseDB()

		version, err := strconv.ParseInt(args[0], 10, 0)
		if err != nil {
			slog.Error("Couldn't parse version", "err", err)
			os.Exit(1)
		}

		if err := goose.UpTo(db.SQLDB, ".", version); err != nil {
			slog.Error("Couldn't apply migrations to target version", "target", version, "err", err)
			os.Exit(1)
		}
	},
}

var migrateDownCmd = &cobra.Command{
	Use:   "down",
	Short: "Rollback the most recent database migration",
	Run: func(cmd *cobra.Command, args []string) {
		db.InitDB()
		defer db.CloseDB()

		if err := goose.Down(db.SQLDB, "."); err != nil {
			slog.Error("Couldn't rollback migration", "err", err)
			os.Exit(1)
		}
	},
}

var migrateDownToCmd = &cobra.Command{
	Use:   "down-to [VERSION]",
	Short: "Rollback all database migrations down to a specific version",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		db.InitDB()
		defer db.CloseDB()

		version, err := strconv.ParseInt(args[0], 10, 0)
		if err != nil {
			slog.Error("Couldn't parse version", "err", err)
			os.Exit(1)
		}

		if err := goose.DownTo(db.SQLDB, ".", version); err != nil {
			slog.Error("Couldn't rollback migrations to target version", "target", version, "err", err)
			os.Exit(1)
		}
	},
}

var migrateRedoCmd = &cobra.Command{
	Use:   "redo",
	Short: "Rollback the most recent database migration and reapply it",
	Run: func(cmd *cobra.Command, args []string) {
		db.InitDB()
		defer db.CloseDB()

		if err := goose.Redo(db.SQLDB, "."); err != nil {
			slog.Error("Couldn't redo migration", "err", err)
			os.Exit(1)
		}
	},
}
