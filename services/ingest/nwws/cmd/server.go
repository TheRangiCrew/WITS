package cmd

import (
	"github.com/TheRangiCrew/WITS/services/ingest/nwws/internal"
	"github.com/joho/godotenv"
	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(serverCmd)
	serverCmd.Flags().StringVar(&DotEnv, "loadenv", "", "The application will try and initialise a .env file for you")
}

var Log int
var DotEnv string

var serverCmd = &cobra.Command{
	Use:   "server",
	Short: "Starts the main NWWS client",
	Long:  `The NWWS client listens for messages and queues them for distribution`,
	Run: func(cmd *cobra.Command, args []string) {
		if DotEnv != "" {
			godotenv.Load(DotEnv)
		}

		internal.Run()
	},
}
