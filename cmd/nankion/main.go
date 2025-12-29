package main

import (
	"log/slog"
	"os"

	"github.com/rodrigoTcarmo/nankion/pkg/log"
	nankionNotion "github.com/rodrigoTcarmo/nankion/pkg/notion"
	notiondatabase "github.com/rodrigoTcarmo/nankion/pkg/notion/database"
	"github.com/spf13/cobra"
)

func main() {
	log.Init()
	defer log.Close()
	rootCmd := &cobra.Command{
		Use:   "nankion",
		Short: "Nankion is a tool for managing your Notion databases",
		Run: func(cmd *cobra.Command, args []string) {
			slog.Info("Hello, Nankion!")
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

	getDatabaseCmd := &cobra.Command{
		Use:   "get [database-id]",
		Short: "Get a Notion database by ID",
		Args:  cobra.ExactArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			databaseID := args[0]
			slog.Info("Fetching database", "databaseID", databaseID)
			databaseCLI := notiondatabase.NewDatabase()
			database, err := databaseCLI.GetDatabase(databaseID)
			if err != nil {
				slog.Error("Error fetching database", "error", err.Error())
				return
			}
			slog.Info("Database info", "database", database)
		},
	}

	listDatabasePropertiesCmd := &cobra.Command{
		Use:   "get-properties [database-id]",
		Short: "Get a Notion database properties by ID",
		Args:  cobra.ExactArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			databaseID := args[0]
			slog.Info("Fetching database properties", "databaseID", databaseID)
			databaseCLI := notiondatabase.NewDatabase()
			database, err := databaseCLI.ListDatabaseProperties(databaseID)
			if err != nil {
				slog.Error("error fetching database properties", "error", err.Error())
				return
			}
			slog.Info("Database properties", "properties", database)
		},
	}

	uploadStatement := &cobra.Command{
		Use:   "upload [database-id] [statement-filepath]",
		Short: "Upload finantial statement to Notion",
		Args:  cobra.ExactArgs(2),
		Run: func(cmd *cobra.Command, args []string) {
			databaseID := args[0]
			filePath := args[1]
			slog.Info("Uploading statement...")
			nl := nankionNotion.NewNotionLoader()
			err := nl.UploadReport(databaseID, filePath)
			if err != nil {
				slog.Info("Error trying to update database: ", err)
			}
		},
	}

	databaseCmd.AddCommand(getDatabaseCmd)
	databaseCmd.AddCommand(listDatabasePropertiesCmd)
	reportCmd.AddCommand(uploadStatement)

	rootCmd.AddCommand(databaseCmd)
	rootCmd.AddCommand(reportCmd)

	if err := rootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}
