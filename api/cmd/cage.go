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
	rootCmd.AddCommand(recordCmd)
	recordCmd.AddCommand(recordCreateCmd)
	recordCmd.AddCommand(recordGetCmd)
	recordGetCmd.Flags().BoolP("clean", "c", false, "clean output suitable for machine parsing")
	recordCmd.AddCommand(recordListByCageCmd)
	recordCmd.AddCommand(recordListCagesCmd)
	recordCmd.AddCommand(recordUpdateCmd)
	recordCmd.AddCommand(recordDeleteCmd)
	recordCmd.AddCommand(recordDeleteCageCmd)
}

var recordCmd = &cobra.Command{
	Use:   "record",
	Short: "Manage backroom records",
	Run: func(cmd *cobra.Command, args []string) {
		cmd.Help()
	},
}

var recordCreateCmd = &cobra.Command{
	Use:   "create [CAGE] [JSON|JSON FILE|STDIN]",
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
		if err := hook.RunHooksByAction(hook.ActionCreate, record); err != nil {
			cmd.PrintErr("Error running hooks:", err)
			return
		}

		cmd.Println("Caged record created with UUID:", record.UUID)
	},
}

var recordGetCmd = &cobra.Command{
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

var recordListByCageCmd = &cobra.Command{
	Use:   "list [CAGE]",
	Short: "List all records belonging to a common cage",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		cageKey := args[0]

		// List all caged records by key
		records, err := cage.ListRecordsByCage(cageKey)
		if err != nil {
			cmd.PrintErr("Error listing caged records:", err)
			return
		}

		if len(records) == 0 {
			cmd.Println("No caged records found for key:", cageKey)
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

var recordListCagesCmd = &cobra.Command{
	Use:   "list-cages",
	Short: "List all unique cages",
	Run: func(cmd *cobra.Command, args []string) {
		keys, err := cage.ListCages()
		if err != nil {
			cmd.PrintErr("Error listing cages:", err)
			return
		}

		if len(keys) == 0 {
			cmd.Println("No cages found")
			return
		}

		for _, key := range keys {
			fmt.Println(key)
		}
	},
}

var recordUpdateCmd = &cobra.Command{
	Use:   "update [UUID] [JSON|JSON FILE|STDIN]",
	Short: "Update an existing caged record",
	Args:  cobra.MinimumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		uuid, err := db.ParseUUID(args[0])
		if err != nil {
			cmd.PrintErr("Invalid UUID format:", err)
			return
		}

		// Retrieve the existing caged record
		record, err := cage.GetRecord(uuid)
		if err != nil {
			cmd.PrintErr("Error retrieving caged record:", err)
			return
		}

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
		var data db.JSONB
		if err := json.Unmarshal(jsonData, &data); err != nil {
			cmd.PrintErr("Error unmarshalling JSON data:", err)
			return
		}
		record.Data = data

		// Update the record
		if err := cage.UpdateRecord(record); err != nil {
			cmd.PrintErr("Error updating caged record:", err)
			return
		}

		// Run hooks after updating the record
		if err := hook.RunHooksByAction(hook.ActionUpdate, record); err != nil {
			cmd.PrintErr("Error running hooks:", err)
			return
		}

		cmd.Println("Caged record updated with UUID:", record.UUID)
	},
}

var recordDeleteCmd = &cobra.Command{
	Use:   "delete [UUID]",
	Short: "Delete a record by UUID",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		uuid, err := db.ParseUUID(args[0])
		if err != nil {
			cmd.PrintErr("Invalid UUID format:", err)
			return
		}

		record, err := cage.GetRecord(uuid)
		if err != nil {
			cmd.PrintErr("Error retrieving caged record:", err)
			return
		}

		// Run hooks before deleting the record
		if err := hook.RunHooksByAction(hook.ActionDelete, record); err != nil {
			cmd.PrintErr("Error running hooks:", err)
			return
		}

		// Delete the caged record
		err = cage.DeleteRecord(uuid)
		if err != nil {
			cmd.PrintErr("Error deleting record:", err)
			return
		}

		cmd.Println("Record deleted with UUID:", uuid)
	},
}

var recordDeleteCageCmd = &cobra.Command{
	Use:   "delete-cage [CAGE]",
	Short: "Delete all records belonging to a common cage",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		cageKey := args[0]

		// Delete all caged records by cageKey
		count, err := cage.DeleteCage(cageKey)
		if err != nil {
			cmd.PrintErr("Error deleting caged records:", err)
			return
		}

		cmd.Printf("%d caged records deleted for cage: %s\n", count, cageKey)
	},
}
