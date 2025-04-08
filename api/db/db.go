package db

import (
	"database/sql"
	"fmt"
	"log/slog"
	"os"

	"github.com/octacian/backroom/api/config"
	"github.com/pressly/goose/v3"
)

// SQLDB stores the current SQL database connection.
var SQLDB *sql.DB

// InitDB connects to the MySQL DB
func InitDB() {
	if SQLDB != nil && SQLDB.Ping() == nil {
		slog.Warn("Database connection already established")
		return
	}

	var err error
	db_path := fmt.Sprintf("%s:%s@tcp(%s)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		config.RC.Database.User, config.RC.Database.Password, config.RC.Database.Host, config.RC.Database.Name)

	SQLDB, err = sql.Open("mysql", db_path)
	if err != nil {
		slog.Error("Couldn't open SQL database", "user", config.RC.Database.User, "name", config.RC.Database.Name, "err", err)
		os.Exit(1)
	}

	if err := SQLDB.Ping(); err != nil {
		slog.Error("Couldn't ping SQL database", "user", config.RC.Database.User, "name", config.RC.Database.Name, "err", err)
		os.Exit(1)
	}

	if err := goose.SetDialect("mysql"); err != nil {
		slog.Error("Couldn't set database dialect for goose", "err", err)
		os.Exit(1)
	}

	SQLDB.SetMaxOpenConns(config.RC.Database.MaxConns)

	stats := SQLDB.Stats()
	slog.Info(
		"Connected to SQL database",
		"user", config.RC.Database.User,
		"name", config.RC.Database.Name,
		"maxConnections", stats.MaxOpenConnections,
		"currConnections", stats.OpenConnections,
	)

	/* if config.RC.Environment == "development" {
		initDBLogger()
	} */
}

// CloseDB closes the database connection.
func CloseDB() {
	if err := SQLDB.Close(); err != nil {
		slog.Error("Couldn't close database", "err", err)
		os.Exit(1)
	}
	slog.Info("Closed database connection")
}

// initDBLogger initializes the database statement logger
/* func initDBLogger() {
	mysql.SetQueryLogger(func(ctx context.Context, queryInfo mysql.QueryInfo) {
		_, args := queryInfo.Statement.Sql()
		slog.Debug("Executed SQL query", "args", args, "duration", queryInfo.Duration, "rows", queryInfo.RowsProcessed, "err", queryInfo.Err)

		lines := strings.Split(queryInfo.Statement.DebugSql(), "\n")
		for i, line := range lines {
			fmt.Printf("%s\t%s\n", color.CyanString(fmt.Sprintf("%03d", i)), line)
		}
	})
} */
