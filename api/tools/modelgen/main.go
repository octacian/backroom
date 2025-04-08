package main

import (
	"fmt"
	"log/slog"
	"os"
	"strconv"
	"strings"

	"github.com/go-jet/jet/v2/generator/metadata"
	"github.com/go-jet/jet/v2/generator/postgres"
	"github.com/go-jet/jet/v2/generator/template"
	postgres2 "github.com/go-jet/jet/v2/postgres"
	_ "github.com/lib/pq"
	"github.com/octacian/backroom/api/config"
	"github.com/octacian/backroom/api/db"
	"github.com/spf13/cobra"
)

func main() {
	execute()
}

var rootCmd = &cobra.Command{
	Use:   "modelgen",
	Short: "Generate JET model files from a database scheme",
	Run: func(cmd *cobra.Command, args []string) {
		// initialize configuration
		config.Init()

		split := strings.Split(config.RC.Database.Host, ":")
		if len(split) != 2 {
			fmt.Printf("Invalid database path: %s\n", config.RC.Database.Host)
			return
		}

		port, err := strconv.Atoi(split[1])
		if err != nil {
			fmt.Printf("Invalid database path, missing valid port: %s\n", config.RC.Database.Host)
			return
		}

		conn := postgres.DBConnection{
			Host:       split[0],
			Port:       port,
			User:       config.RC.Database.User,
			Password:   config.RC.Database.Password,
			DBName:     config.RC.Database.Name,
			SchemaName: "public",
			SslMode:    "disable",
		}

		tmpl := template.Default(postgres2.Dialect).
			UseSchema(func(schema metadata.Schema) template.Schema {
				return template.DefaultSchema(schema).
					UseModel(template.DefaultModel().
						UseTable(func(table metadata.Table) template.TableModel {
							return template.DefaultTableModel(table).
								UseField(func(column metadata.Column) template.TableModelField {
									defaultTableModelField := template.DefaultTableModelField(column)

									if column.Name == "uuid" {
										defaultTableModelField.Type = template.NewType(db.UUID{})
									} else if column.DataType.Name == "jsonb" {
										defaultTableModelField.Type = template.NewType(db.JSONB{})
									}

									return defaultTableModelField
								})
						}))
			})

		err = postgres.Generate(
			"./.gen",
			conn,
			tmpl,
		)

		if err != nil {
			slog.Error("Couldn't generate jet model files", "err", err)
		}
	},
}

func execute() {
	if err := rootCmd.Execute(); err != nil {
		slog.Error("Couldn't execute root command", "err", err)
		os.Exit(1)
	}
}
