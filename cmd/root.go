package cmd

import (
	"fmt"
	"os"

	"github.com/joho/godotenv"
	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "av-scraper",
	Short: "A CLI to scrape data from Alpha Vantage",
	Long:  `A CLI application to fetch stock and news data from the Alpha Vantage API and save it to a CSV file.`,
	PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
		return godotenv.Load()
	},
}

func Execute() error {
	return rootCmd.Execute()
}

func init() {
	rootCmd.AddCommand(stocksCmd)
	rootCmd.AddCommand(newsCmd)
}

func er(msg interface{}) {
	fmt.Println("Error:", msg)
	os.Exit(1)
}
