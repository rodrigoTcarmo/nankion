package main

import (
	"log/slog"
	"os"

	"github.com/rodrigoTcarmo/nankion/pkg/log"
	nankionNotion "github.com/rodrigoTcarmo/nankion/pkg/notion"
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
			nl := nankionNotion.NewNotionLoader()
			database, err := nl.Database.GetDatabase(databaseID)
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
			nl := nankionNotion.NewNotionLoader()
			database, err := nl.Database.ListDatabaseProperties(databaseID)
			if err != nil {
				slog.Error("error fetching database properties", "error", err.Error())
				return
			}
			slog.Info("Database properties", "properties", database)
		},
	}

	uploadSingleStatement := &cobra.Command{
		Use:   "upload [database-id] [statement-filepath]",
		Short: "Upload financial statement to Notion",
		Args:  cobra.ExactArgs(2),
		Run: func(cmd *cobra.Command, args []string) {
			databaseID := args[0]
			filePath := args[1]
			slog.Info("Uploading statement...")
			nl := nankionNotion.NewNotionLoader()
			err := nl.UploadStatement(databaseID, filePath)
			if err != nil {
				slog.Info("error trying to update database: ", "error", err.Error())
			}
		},
	}

	uploadMultipleStatements := &cobra.Command{
		Use:   "upload-statements [database-id] [statement-folderpath]",
		Short: "Upload multiple .ofx financial statements to Notion",
		Args:  cobra.ExactArgs(2),
		Run: func(cmd *cobra.Command, args []string) {
			databaseID := args[0]
			folderPath := args[1]
			slog.Info("Uploading statements...")
			nl := nankionNotion.NewNotionLoader()
			err := nl.UploadStatements(databaseID, folderPath)
			if err != nil {
				slog.Info("error trying to update database: ", "error", err.Error())
			}
		},
	}

	databaseCmd.AddCommand(getDatabaseCmd)
	databaseCmd.AddCommand(listDatabasePropertiesCmd)
	reportCmd.AddCommand(uploadSingleStatement)
	reportCmd.AddCommand(uploadMultipleStatements)

	rootCmd.AddCommand(databaseCmd)
	rootCmd.AddCommand(reportCmd)

	if err := rootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}
