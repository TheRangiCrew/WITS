package cmd

// import (
// 	"log/slog"
// 	"time"

// 	"github.com/TheRangiCrew/WITS/services/parsing/awips/internal/server"
// 	"github.com/joho/godotenv"
// 	"github.com/spf13/cobra"
// )

// func init() {
// 	rootCmd.AddCommand(iemCmd)
// 	iemCmd.Flags().IntVar(&Log, "minlog", 0, "The minimum logging level to use")
// 	iemCmd.Flags().StringVar(&DotEnv, "loadenv", "", "The application will try and initialise a .env file for you")
// 	iemCmd.Flags().StringVarP(&StartDate, "startDate", "s", "", "The date to start pulling from the IEM. YYYY-MM-DD")
// 	iemCmd.Flags().StringVarP(&EndDate, "endDate", "e", "", "The date to end pulling from the IEM. If left empty, the start date will be used as the end date.")
// 	iemCmd.Flags().StringVarP(&Product, "product", "p", "", "The AWIPS product to search for. WAR for warnings (IEM shortcut).")
// 	iemCmd.Flags().StringVarP(&WFO, "wfo", "w", "", "The WFO who owns the product. ALL to get all (using the database).")
// }

// var StartDate string
// var EndDate string
// var Product string
// var WFO string

// var iemCmd = &cobra.Command{
// 	Use:   "iem",
// 	Short: "Fetches IEM text data and archives it. Thanks Daryl! :)",
// 	Run: func(cmd *cobra.Command, args []string) {
// 		if len(StartDate) == 0 {
// 			slog.Error("start date argument is missing.")
// 			return
// 		}

// 		startDate, err := time.Parse("2006-01-02", StartDate)
// 		if err != nil {
// 			slog.Error("error parsing start date: " + err.Error())
// 		}

// 		var endDate *time.Time
// 		if len(EndDate) > 0 {
// 			t, err := time.Parse("2006-01-02", EndDate)
// 			if err != nil {
// 				slog.Error("error parsing end date: " + err.Error())
// 			}
// 			endDate = &t
// 		}

// 		if len(WFO) != 3 {
// 			slog.Error("invalid wfo length. Must be 3 characters.")
// 			return
// 		}

// 		if len(Product) < 1 || len(Product) > 3 {
// 			slog.Error("invalid AWIPS product identifier length. Must be between 1 and 3 characters.")
// 			return
// 		}

// 		if DotEnv != "" {
// 			err := godotenv.Load(DotEnv)
// 			if err != nil {
// 				slog.Error("failed to load env: " + err.Error())
// 				return
// 			}
// 		}

// 		config := server.IEMConfig{
// 			StartDate:      startDate,
// 			EndDate:        endDate,
// 			Office:         WFO,
// 			Product:        Product,
// 			MaxConcurrency: 8,
// 		}

// 		server.IEM(config, Log)
// 	},
// }
