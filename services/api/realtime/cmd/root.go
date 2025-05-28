package cmd

import (
	"fmt"
	"os"

	"github.com/TheRangiCrew/WITS/services/api/realtime/server"
	"github.com/joho/godotenv"
	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:  "realtime",
	Long: "The realtime service for WITS, handling live data feeds.",
}

var DotEnv string

var serverCmd = &cobra.Command{
	Use:   "server",
	Short: "Starts the realtime server",
	Long:  `The realtime server ingests from RabbitMQ feeds and makes the data highly available`,
	Run: func(cmd *cobra.Command, args []string) {
		if DotEnv != "" {
			godotenv.Load(DotEnv)
		}

		s := server.New()
		if s == nil {
			fmt.Fprintln(os.Stderr, "Failed to setup server")
			os.Exit(1)
		}
		s.Start()
	},
}

func init() {
	rootCmd.AddCommand(serverCmd)
	serverCmd.Flags().StringVar(&DotEnv, "loadenv", "", "The application will try and initialise a .env file for you")
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
