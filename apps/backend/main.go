package main

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"github.com/theCompanyDream/id-trials/apps/backend/cmd"
	"github.com/theCompanyDream/id-trials/apps/backend/controller"
	"github.com/theCompanyDream/id-trials/apps/backend/models"
)

var rootCmd = &cobra.Command{
	Use:   "id-benchmark",
	Short: "ID Performance Benchmark Tool",
	Long:  `A comprehensive tool for benchmarking different ID types (ULID, KSUID, UUID4, etc.) in database operations.`,
}

var serverCmd = &cobra.Command{
	Use:   "server",
	Short: "Start the web API server",
	Long:  `Starts the web server with REST API endpoints for all ID types.`,
	Run: func(command *cobra.Command, args []string) {
		controller.RunServer()
	},
}

var generateCmd = &cobra.Command{
	Use:   "generate",
	Short: "Generate test data",
	Long:  `Generate test data for all ID types in the database.`,
	Run: func(command *cobra.Command, args []string) {
		records, _ := command.Flags().GetInt("records")
		batch, _ := command.Flags().GetInt("batch")
		concurrent, _ := command.Flags().GetBool("concurrent")

		config := &models.Config{
			RecordsPerTable: records,
			BatchSize:       batch,
			Concurrent:      concurrent,
		}

		cmd.GenerateData(config)
	},
}

func init() {
	// Server command flags
	serverCmd.Flags().StringP("port", "p", "3000", "Port to run server on")
	serverCmd.Flags().StringP("host", "H", "localhost", "Host to bind server to")

	// Generate command flags
	generateCmd.Flags().IntP("records", "r", 10000, "Number of records per table")
	generateCmd.Flags().IntP("batch", "b", 1000, "Batch size for inserts")
	generateCmd.Flags().BoolP("concurrent", "c", true, "Generate data concurrently across tables")
	generateCmd.Flags().StringP("database", "d", "", "Database connection string")

	// Add commands to root
	rootCmd.AddCommand(serverCmd)
	rootCmd.AddCommand(generateCmd)
}

func main() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
