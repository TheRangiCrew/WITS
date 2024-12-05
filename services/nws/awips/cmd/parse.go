package cmd

import (
	"fmt"

	"github.com/TheRangiCrew/WITS/services/nws/awips/internal/db"
	"github.com/TheRangiCrew/WITS/services/nws/awips/internal/server"
	"github.com/spf13/cobra"
	"github.com/surrealdb/surrealdb.go"
)

func init() {
	rootCmd.AddCommand(parseCmd)
	parseCmd.Flags().StringVarP(&Endpoint, "endpoint", "e", "", "The database endpoint to output to")
	parseCmd.Flags().StringVarP(&Username, "username", "u", "", "The database username for authentication")
	parseCmd.Flags().StringVarP(&Password, "password", "p", "", "The database password for authentication")
	parseCmd.Flags().StringVarP(&Namespace, "namespace", "n", "", "The database namespace to use")
	parseCmd.Flags().StringVarP(&Database, "database", "d", "", "The database to use")
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

		config := db.DBConfig{
			Auth: surrealdb.Auth{
				Namespace: Namespace,
				Database:  Database,
				Username:  Username,
				Password:  Password,
			},
			Endpoint: Endpoint,
			AsRoot:   RootAuth,
		}

		server.ParseText(args[0], config)
	},
}
