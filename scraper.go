package main

import (
	"go-scraper/menus"
	"go-scraper/spider"
	"log"
	"os"

	"github.com/joho/godotenv"
)

const REQUEST_TIMEOUT = 120

func LoadEnv() {
	err := godotenv.Load(".env")
	if err != nil {
		log.Fatalf("failed to load env file - can't run scraper without html declarations")
	}
}

func main() {

	LoadEnv()

	htmlElements := spider.Elements {
		// links on search page are relative, so need base URL
		BaseUrl: os.Getenv("BASE_URL"),
		// all html tags needed to get related article info
		ArticleBody: os.Getenv("ARTICLE_BODY"),
		ArticleTitle: os.Getenv("ARTICLE_TITLE"),
		ArticleText: os.Getenv("ARTICLE_TEXT"),
		// all html tags needed to get all related search data
		ResultsContainer: os.Getenv("RESULTS_CONTAINER"),
		ResultsLink: os.Getenv("RESULTS_LINK"),
	}

	spider.InitSpider("new scraper", htmlElements)

	menus.InitializePrompts()

}
