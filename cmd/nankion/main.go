package main

import (
	"fmt"
	"os"

	nankionNotion "github.com/rodrigoTcarmo/nankion/pkg/notion"
	notiondatabase "github.com/rodrigoTcarmo/nankion/pkg/notion/database"
	"github.com/rodrigoTcarmo/nankion/pkg/statement"
	"github.com/spf13/cobra"
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

	reportCmd := &cobra.Command{
		Use:   "report",
		Short: "Manage your reports",
	}

	printReportCmd := &cobra.Command{
		Use:   "print [ofx-filepath]",
		Short: "Print a report from an OFX file",
		Args:  cobra.ExactArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			filepath := args[0]
			fmt.Printf("Printing report from: %s\n", filepath)
			report, err := statement.LoadReport(filepath)
			if err != nil {
				fmt.Printf("Error loading report: %s\n", err)
				return
			}
			fmt.Println("Statements\n", report.Statements)
		},
	}

	getDatabaseCmd := &cobra.Command{
		Use:   "get [database-id]",
		Short: "Get a Notion database by ID",
		Args:  cobra.ExactArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			databaseID := args[0]
			fmt.Printf("Fetching database: %s\n", databaseID)
			databaseCLI := notiondatabase.NewDatabase()
			database, err := databaseCLI.GetDatabase(databaseID)
			if err != nil {
				fmt.Printf("Error fetching database: %s\n", err)
				return
			}
			fmt.Println("Database\n", database)
		},
	}

	listDatabasePropertiesCmd := &cobra.Command{
		Use:   "get-properties [database-id]",
		Short: "Get a Notion database properties by ID",
		Args:  cobra.ExactArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			databaseID := args[0]
			fmt.Printf("Fetching database properties: %s\n", databaseID)
			databaseCLI := notiondatabase.NewDatabase()
			database, err := databaseCLI.ListDatabaseProperties(databaseID)
			if err != nil {
				fmt.Printf("Error fetching database properties: %s\n", err)
				return
			}
			fmt.Println("Database properties\n", database)
		},
	}

	uploadStatement := &cobra.Command{
		Use:   "upload-statement [database-id] [statement-filepath]",
		Short: "Upload finantial statement to Notion",
		Args:  cobra.ExactArgs(2),
		Run: func(cmd *cobra.Command, args []string) {
			databaseID := args[0]
			filePath := args[1]
			fmt.Println("Uploading statement...")
			nl := nankionNotion.NewNotionLoader()
			err := nl.UploadReport(databaseID, filePath)
			if err != nil {
				fmt.Println("Error trying to update database: ", err)
			}
		},
	}

	databaseCmd.AddCommand(getDatabaseCmd)
	databaseCmd.AddCommand(listDatabasePropertiesCmd)
	databaseCmd.AddCommand(uploadStatement)

	reportCmd.AddCommand(printReportCmd)

	rootCmd.AddCommand(databaseCmd)
	rootCmd.AddCommand(reportCmd)

	if err := rootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}
