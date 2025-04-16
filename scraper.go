package main

import (
	"go-scraper/helpers"
	"go-scraper/prompt"
	"go-scraper/spider"
	"log"
	"os"
	"path/filepath"

	"github.com/joho/godotenv"
)

func LoadEnv() {
	exePath, err := os.Executable()
	if err != nil {
		log.Fatal("could not determine executable path:", err)
	}
	exeDir := filepath.Dir(exePath)

	envPath := filepath.Join(exeDir, ".env")

	err = godotenv.Load(envPath)
	if err != nil {
		log.Fatalf("failed to load env file - can't run scraper without html declarations")
	}
}

func main() {

	LoadEnv()

	helpers.InitRandSource()

	htmlElements := spider.Elements{
		// links on search page are relative, so need base URL
		BaseUrl: os.Getenv("BASE_URL"),
		// all html tags needed to get related article info
		ArticleBody:  os.Getenv("ARTICLE_BODY"),
		ArticleTitle: os.Getenv("ARTICLE_TITLE"),
		ArticleText:  os.Getenv("ARTICLE_TEXT"),
		// all html tags needed to get all related search data
		ResultsContainer: os.Getenv("RESULTS_CONTAINER"),
		ResultsLink:      os.Getenv("RESULTS_LINK"),
	}
	spider.SEARCH_URL = os.Getenv("SEARCH_URL")

	spider.InitSpider("new scraper", htmlElements)

	// menus.InitializePrompts()
	prompt.InitializePrompts()

}
