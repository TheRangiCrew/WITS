package cmd

import (
	"fmt"

	"github.com/TheRangiCrew/WITS/services/parsing/awips/internal/server"
	"github.com/joho/godotenv"
	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(parseCmd)
	parseCmd.Flags().IntVar(&Log, "minlog", 0, "The minimum logging level to use")
	parseCmd.Flags().StringVar(&DotEnv, "loadenv", "", "The application will try and initialise a .env file for you")

}

var parseCmd = &cobra.Command{
	Use:   "parse",
	Short: "Parses provided AWIPS products",
	Long: `Given a file or directory, the AWIPS products will attempt to be parsed.
	This is mainly for testing.`,
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) < 1 {
			fmt.Printf("Need 1 argument. %d were provided. Please provide a filename.\n", len(args))
			return
		}

		if len(args) > 1 {
			fmt.Printf("Need 1 argument. %d were provided. You can only provide a filename.\n", len(args))
			return
		}

		if DotEnv != "" {
			godotenv.Load(DotEnv)
		}

		server.ParseText(args[0], Log)
	},
}
