package cmd

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/fatih/color"
	"github.com/octacian/backroom/api/cage"
	"github.com/octacian/backroom/api/db"
	"github.com/octacian/backroom/api/hook"
	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(cageCmd)
	cageCmd.AddCommand(cageCreateRecordCmd)
	cageCmd.AddCommand(cageGetRecordCmd)
	cageGetRecordCmd.Flags().BoolP("clean", "c", false, "clean output suitable for machine parsing")
	cageCmd.AddCommand(cageListRecordsByKeyCmd)
	cageCmd.AddCommand(cageListKeysCmd)
	cageCmd.AddCommand(cageDeleteRecordCmd)
	cageCmd.AddCommand(cageDeleteCageCmd)
}

var cageCmd = &cobra.Command{
	Use:   "cage",
	Short: "Manage caged records",
	Run: func(cmd *cobra.Command, args []string) {
		cmd.Help()
	},
}

var cageCreateRecordCmd = &cobra.Command{
	Use:   "create [KEY] [JSON|JSON FILE|STDIN]",
	Short: "Create a new caged record",
	Args:  cobra.MinimumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		key := args[0]

		var reader io.Reader

		if len(args) < 2 {
			// args[1] doesn't exist, read from stdin
			cmd.Println("Reading JSON from stdin...")
			reader = cmd.InOrStdin()
		} else if _, err := os.Stat(args[1]); err == nil {
			// args[1] looks like a file, read from it
			file, err := os.Open(args[1])
			if err != nil {
				cmd.PrintErr("Error opening file:", err)
				return
			}
			defer file.Close()
			reader = file
			cmd.Println("Reading JSON from file:", args[1])
		} else {
			// args[1] is a JSON string
			reader = strings.NewReader(args[1])
		}

		// Read JSON from the reader
		jsonData, err := io.ReadAll(reader)
		if err != nil {
			cmd.PrintErr("Error reading JSON data:", err)
			return
		}

		// Unmarshal JSON data into a db.JSONB object
		record, err := cage.NewRecordFromString(key, string(jsonData))
		if err != nil {
			cmd.PrintErr("Error preparing JSON data:", err)
		}

		// Save a new caged record
		if err := cage.CreateRecord(record); err != nil {
			cmd.PrintErr("Error creating caged record:", err)
			return
		}

		// Run hooks after creating the record
		if err := hook.RunCreate(record); err != nil {
			cmd.PrintErr("Error running hooks:", err)
			return
		}

		cmd.Println("Caged record created with UUID:", record.UUID)
	},
}

var cageGetRecordCmd = &cobra.Command{
	Use:   "get [UUID]",
	Short: "Get a caged record by UUID",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		uuid, err := db.ParseUUID(args[0])
		if err != nil {
			cmd.PrintErr("Invalid UUID format:", err)
			return
		}

		// Retrieve the caged record
		record, err := cage.GetRecord(uuid)
		if err != nil {
			cmd.PrintErr("Error retrieving caged record:", err)
			return
		}

		// Print the caged record
		data, err := json.MarshalIndent(record.Data, "", "  ")
		if err != nil {
			cmd.PrintErr("Error marshalling caged record:", err)
			return
		}

		// Check if the output should be clean
		clean, err := cmd.Flags().GetBool("clean")
		if err != nil {
			cmd.PrintErr("Error getting clean flag:", err)
			return
		}

		if clean {
			fmt.Println(string(data))
			return
		}

		lines := strings.Split(string(data), "\n")
		for i, line := range lines {
			fmt.Printf("%s\t%s\n", color.GreenString(fmt.Sprintf("%03d", i)), line)
		}
	},
}

var cageListRecordsByKeyCmd = &cobra.Command{
	Use:   "list [KEY]",
	Short: "List all caged records by key",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		key := args[0]

		// List all caged records by key
		records, err := cage.ListRecordsByKey(key)
		if err != nil {
			cmd.PrintErr("Error listing caged records:", err)
			return
		}

		if len(records) == 0 {
			cmd.Println("No caged records found for key:", key)
			return
		}

		for _, record := range records {
			data, err := json.Marshal(record.Data)
			if err != nil {
				cmd.PrintErr(fmt.Sprintf("Error marshalling caged record %s:", record.UUID), err)
				return
			}

			fmt.Println(string(data))
		}
	},
}

var cageListKeysCmd = &cobra.Command{
	Use:   "list-keys",
	Short: "List all unique cage keys",
	Run: func(cmd *cobra.Command, args []string) {
		keys, err := cage.ListCageKeys()
		if err != nil {
			cmd.PrintErr("Error listing cage keys:", err)
			return
		}

		if len(keys) == 0 {
			cmd.Println("No cage keys found")
			return
		}

		for _, key := range keys {
			fmt.Println(key)
		}
	},
}

var cageDeleteRecordCmd = &cobra.Command{
	Use:   "delete [UUID]",
	Short: "Delete a caged record by UUID",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		uuid, err := db.ParseUUID(args[0])
		if err != nil {
			cmd.PrintErr("Invalid UUID format:", err)
			return
		}

		// Delete the caged record
		err = cage.DeleteRecord(uuid)
		if err != nil {
			cmd.PrintErr("Error deleting caged record:", err)
			return
		}

		cmd.Println("Caged record deleted with UUID:", uuid)
	},
}

var cageDeleteCageCmd = &cobra.Command{
	Use:   "delete-cage [KEY]",
	Short: "Delete all caged records by key",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		key := args[0]

		// Delete all caged records by key
		count, err := cage.DeleteCage(key)
		if err != nil {
			cmd.PrintErr("Error deleting caged records:", err)
			return
		}

		cmd.Printf("%d caged records deleted for key: %s\n", count, key)
	},
}
