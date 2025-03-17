package cmd

import (
	"log/slog"

	"github.com/TheRangiCrew/WITS/services/parsing/awips/internal/server"
	"github.com/joho/godotenv"
	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(serverCmd)
	serverCmd.Flags().IntVar(&Log, "minlog", 0, "The minimum logging level to use")
	serverCmd.Flags().StringVar(&DotEnv, "loadenv", "", "The application will try and initialise a .env file for you")
}

var Log int
var DotEnv string

var serverCmd = &cobra.Command{
	Use:   "server",
	Short: "Starts the AWIPS server",
	Long: `The AWIPS server listens for pending products
				 and processes them`,
	Run: func(cmd *cobra.Command, args []string) {
		if DotEnv != "" {
			godotenv.Load(DotEnv)
		}

		config := server.ServerConfig{
			MinLog: Log,
		}

		s, err := server.New(config)
		if err != nil {
			slog.Error(err.Error())
		}

		s.Start()
	},
}
