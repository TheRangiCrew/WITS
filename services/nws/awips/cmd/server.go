package cmd

import (
	"github.com/TheRangiCrew/WITS/services/nws/awips/internal/db"
	"github.com/TheRangiCrew/WITS/services/nws/awips/internal/server"
	"github.com/spf13/cobra"
	"github.com/surrealdb/surrealdb.go"
)

func init() {
	rootCmd.AddCommand(serverCmd)
	serverCmd.Flags().StringVarP(&Endpoint, "endpoint", "e", "", "The database endpoint to output to")
	serverCmd.Flags().StringVarP(&Username, "username", "u", "", "The database username for authentication")
	serverCmd.Flags().StringVarP(&Password, "password", "p", "", "The database password for authentication")
	serverCmd.Flags().StringVarP(&Namespace, "namespace", "n", "", "The database namespace to use")
	serverCmd.Flags().StringVarP(&Database, "database", "d", "", "The database to use")
	serverCmd.Flags().BoolVar(&RootAuth, "as-root", false, "Connect to the database using root credentials")

}

var Endpoint string
var Username string
var Password string
var Namespace string
var Database string
var RootAuth bool

var serverCmd = &cobra.Command{
	Use:   "server",
	Short: "Starts the main AWIPS server",
	Long: `The AWIPS server listens for pending products
				on the database and processes them`,
	Run: func(cmd *cobra.Command, args []string) {
		auth := server.ServerConfig{
			DB: db.DBConfig{
				Auth: surrealdb.Auth{
					Namespace: Namespace,
					Database:  Database,
					Username:  Username,
					Password:  Password,
				},
				Endpoint: Endpoint,
				AsRoot:   RootAuth,
			},
		}
		server.Start(auth)
	},
}
