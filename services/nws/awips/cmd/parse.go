package cmd

import (
	"fmt"
	"os"

	"github.com/TheRangiCrew/WITS/services/nws/awips/internal/db"
	"github.com/TheRangiCrew/WITS/services/nws/awips/internal/server"
	"github.com/joho/godotenv"
	"github.com/spf13/cobra"
	"github.com/surrealdb/surrealdb.go"
)

func init() {
	rootCmd.AddCommand(parseCmd)
	parseCmd.Flags().IntVar(&Log, "minlog", 0, "The minimum logging level to use")
	parseCmd.Flags().StringVar(&DotEnv, "loadenv", "", "The application will try and initialise a .env file for you")
	// parseCmd.Flags().StringVarP(&Endpoint, "endpoint", "e", "", "The database endpoint to output to")
	// parseCmd.Flags().StringVarP(&Username, "username", "u", "", "The database username for authentication")
	// parseCmd.Flags().StringVarP(&Password, "password", "p", "", "The database password for authentication")
	// parseCmd.Flags().StringVarP(&Namespace, "namespace", "n", "", "The database namespace to use")
	// parseCmd.Flags().StringVarP(&Database, "database", "d", "", "The database to use")
	parseCmd.Flags().BoolVar(&RootAuth, "as-root", false, "Connect to the database using root credentials")

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

		config := db.DBConfig{
			Auth: surrealdb.Auth{
				Namespace: os.Getenv("DB_NS"),
				Database:  os.Getenv("DB_DB"),
				Username:  os.Getenv("DB_USER"),
				Password:  os.Getenv("DB_PASS"),
			},
			Endpoint: os.Getenv("DB_URL"),
			AsRoot:   RootAuth,
		}

		server.ParseText(args[0], config, Log)
	},
}
