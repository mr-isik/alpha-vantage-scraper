package cmd

import (
	"encoding/csv"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"sort"
	"time"

	"github.com/spf13/cobra"
)

var stocksCmd = &cobra.Command{
	Use:   "stocks",
	Short: "Fetch stock market data",
	Run:   fetchStocks,
}

func init() {
	stocksCmd.Flags().String("symbol", "", "Stock symbol to fetch (e.g., IBM)")
	stocksCmd.MarkFlagRequired("symbol")
	stocksCmd.Flags().String("start-date", "", "Start date in YYYY-MM-DD format")
	stocksCmd.Flags().String("end-date", "", "End date in YYYY-MM-DD format")
	stocksCmd.Flags().String("output", "output.csv", "Output CSV file name")
}

type TimeSeries map[string]struct {
	Open   string `json:"1. open"`
	High   string `json:"2. high"`
	Low    string `json:"3. low"`
	Close  string `json:"4. close"`
	Volume string `json:"5. volume"`
}

type StockData struct {
	MetaData   map[string]string `json:"Meta Data"`
	TimeSeries TimeSeries        `json:"Time Series (Daily)"`
}

func fetchStocks(cmd *cobra.Command, args []string) {
	symbol, _ := cmd.Flags().GetString("symbol")
	startDateStr, _ := cmd.Flags().GetString("start-date")
	endDateStr, _ := cmd.Flags().GetString("end-date")
	output, _ := cmd.Flags().GetString("output")
	apiKey := os.Getenv("API_KEY")

	if apiKey == "" {
		er("API_KEY not found in .env file")
	}

	url := fmt.Sprintf("https://www.alphavantage.co/query?function=TIME_SERIES_DAILY&symbol=%s&apikey=%s&outputsize=full", symbol, apiKey)

	resp, err := http.Get(url)
	if err != nil {
		er(err)
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		er(err)
	}

	var stockData StockData
	err = json.Unmarshal(body, &stockData)
	if err != nil {
		// Check for error message from API
		var apiError map[string]string
		json.Unmarshal(body, &apiError)
		if msg, ok := apiError["Error Message"]; ok {
			er(msg)
		}
		er(err)
	}

	if stockData.TimeSeries == nil {
		er("Could not retrieve stock data. Check your symbol or API key.")
	}

	file, err := os.Create(output)
	if err != nil {
		er(err)
	}
	defer file.Close()

	writer := csv.NewWriter(file)
	defer writer.Flush()

	writer.Write([]string{"date", "open", "high", "low", "close", "volume"})

	var startDate, endDate time.Time
	var dateErr error

	if startDateStr != "" {
		startDate, dateErr = time.Parse("2006-01-02", startDateStr)
		if dateErr != nil {
			er("Invalid start date format. Use YYYY-MM-DD.")
		}
	}
	if endDateStr != "" {
		endDate, dateErr = time.Parse("2006-01-02", endDateStr)
		if dateErr != nil {
			er("Invalid end date format. Use YYYY-MM-DD.")
		}
	}

	// Sort dates
	dates := make([]string, 0, len(stockData.TimeSeries))
	for date := range stockData.TimeSeries {
		dates = append(dates, date)
	}
	sort.Strings(dates)

	for _, dateStr := range dates {
		data := stockData.TimeSeries[dateStr]
		currentDate, _ := time.Parse("2006-01-02", dateStr)

		if (startDateStr == "" || !currentDate.Before(startDate)) && (endDateStr == "" || !currentDate.After(endDate)) {
			row := []string{dateStr, data.Open, data.High, data.Low, data.Close, data.Volume}
			if err := writer.Write(row); err != nil {
				er(err)
			}
		}
	}

	fmt.Println("Successfully saved stock data to", output)
}
