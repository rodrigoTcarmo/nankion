package main

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"github.com/rodrigoTcarmo/nankion/pkg/cli/databases"
)

func main() {
	rootCmd := &cobra.Command{
		Use:   "nankion",
		Short: "Nankion is a tool for managing your Notion databases",
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Println("Hello, Nankion!")
		},
	}

	databaseCmd := &cobra.Command{
		Use:   "database",
		Short: "Manage your Notion databases",
	}

	getDatabaseCmd := &cobra.Command{
		Use:   "get [database-id]",
		Short: "Get a Notion database by ID",
		Args:  cobra.ExactArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			databaseID := args[0]
			fmt.Printf("Fetching database: %s\n", databaseID)
			databases.GetDatabase(databaseID)
		},
	}

	listDatabasePropertiesCmd := &cobra.Command {
		Use: "get-properties [database-id]",
		Short: "Get a Notion database properties by ID",
		Args: cobra.ExactArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			databaseID := args[0]
			fmt.Println("Fetching database properties: %s\n", databaseID)
			databases.ListDatabaseProperties(databaseID)
		},
	}

	databaseCmd.AddCommand(getDatabaseCmd)
	databaseCmd.AddCommand(listDatabasePropertiesCmd)
	rootCmd.AddCommand(databaseCmd)

	if err := rootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}
