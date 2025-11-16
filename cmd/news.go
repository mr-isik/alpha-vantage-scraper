package cmd

import (
	"encoding/csv"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"

	"github.com/spf13/cobra"
)

var newsCmd = &cobra.Command{
	Use:   "news",
	Short: "Fetch news and sentiment data",
	Run:   fetchNews,
}

func init() {
	newsCmd.Flags().String("tickers", "", "Comma-separated list of stock tickers (e.g., IBM,AAPL)")
	newsCmd.Flags().String("topics", "", "Comma-separated list of topics (e.g., technology,finance)")
	newsCmd.Flags().String("start-date", "", "Start date in YYYYMMDDTHHMM format")
	newsCmd.Flags().String("end-date", "", "End date in YYYYMMDDTHHMM format")
	newsCmd.Flags().Int("limit", 50, "Number of results to return")
	newsCmd.Flags().String("output", "output.csv", "Output CSV file name")
}

type NewsArticle struct {
	Title            string `json:"title"`
	URL              string `json:"url"`
	TimePublished    string `json:"time_published"`
	Summary          string `json:"summary"`
	BannerImage      string `json:"banner_image"`
	Source           string `json:"source"`
	Category         string `json:"category_within_source"`
	SourceDomain     string `json:"source_domain"`
	OverallSentiment string `json:"overall_sentiment_label"`
}

type NewsData struct {
	Feed []NewsArticle `json:"feed"`
}

func fetchNews(cmd *cobra.Command, args []string) {
	tickers, _ := cmd.Flags().GetString("tickers")
	topics, _ := cmd.Flags().GetString("topics")
	startDate, _ := cmd.Flags().GetString("start-date")
	endDate, _ := cmd.Flags().GetString("end-date")
	limit, _ := cmd.Flags().GetInt("limit")
	output, _ := cmd.Flags().GetString("output")
	apiKey := os.Getenv("API_KEY")

	if apiKey == "" {
		er("API_KEY not found in .env file")
	}

	baseURL := "https://www.alphavantage.co/query"
	params := url.Values{}
	params.Add("function", "NEWS_SENTIMENT")
	params.Add("apikey", apiKey)

	if tickers != "" {
		params.Add("tickers", tickers)
	}
	if topics != "" {
		params.Add("topics", topics)
	}
	if startDate != "" {
		params.Add("time_from", startDate)
	}
	if endDate != "" {
		params.Add("time_to", endDate)
	}
	params.Add("limit", fmt.Sprintf("%d", limit))

	fullURL := baseURL + "?" + params.Encode()

	resp, err := http.Get(fullURL)
	if err != nil {
		er(err)
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		er(err)
	}

	var newsData NewsData
	err = json.Unmarshal(body, &newsData)
	if err != nil {
		var apiError map[string]string
		json.Unmarshal(body, &apiError)
		if msg, ok := apiError["Information"]; ok {
			er(msg)
		}
		er(err)
	}

	if len(newsData.Feed) == 0 {
		er("No news found for the given criteria.")
	}

	file, err := os.Create(output)
	if err != nil {
		er(err)
	}
	defer file.Close()

	writer := csv.NewWriter(file)
	defer writer.Flush()

	writer.Write([]string{"title", "url", "time_published", "summary", "banner_image", "source", "category", "source_domain", "overall_sentiment"})

	for _, article := range newsData.Feed {
		row := []string{
			article.Title,
			article.URL,
			article.TimePublished,
			article.Summary,
			article.BannerImage,
			article.Source,
			article.Category,
			article.SourceDomain,
			article.OverallSentiment,
		}
		if err := writer.Write(row); err != nil {
			er(err)
		}
	}

	fmt.Println("Successfully saved news data to", output)
}
