package cmd

import (
	"github.com/spf13/cobra"
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
		// if len(args) < 1 {
		// 	fmt.Printf("Need 1 argument. %d were provided. Please provide a filename.\n", len(args))
		// }

		// if len(args) > 1 {
		// 	fmt.Printf("Need 1 argument. %d were provided. You can only provide a filename.\n", len(args))
		// }

		// data, err := os.ReadFile(args[0])
		// if err != nil {
		// 	slog.Error(err.Error())
		// 	return
		// }

		// text := string(data)

		// config := server.ServerConfig{
		// 	DB: db.DBConfig{
		// 		Auth: surrealdb.Auth{
		// 			Namespace: Namespace,
		// 			Database:  Database,
		// 			Username:  Username,
		// 			Password:  Password,
		// 		},
		// 		Endpoint: Endpoint,
		// 		AsRoot:   RootAuth,
		// 	},
		// }

		// slog.Info(fmt.Sprintf("Parsing %s", args[0]))

		// h, err := handler.New(config.DB)
		// if err != nil {
		// 	slog.Error(err.Error())
		// 	return
		// }

		// err = h.Handle(text, time.Now())
		// if err != nil {
		// 	slog.Error(err.Error())
		// }
	},
}
