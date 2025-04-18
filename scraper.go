package main

import (
	"flag"
	"fmt"
	"go-scraper/helpers"
	"go-scraper/prompt"
	"go-scraper/spider"
	"log"
	"os"
	"path/filepath"

	"github.com/joho/godotenv"
)

func LoadEnv() {
	// set default env
	mode := flag.String("env", "dev", "Defines the application state: (prod/dev)")
	flag.Parse()

	if *mode == "dev" {
		err := godotenv.Load(".env")
		if err != nil {
			log.Fatalf("Failed to load env in dev environment: %s", err)
		}
	} else if *mode == "prod" {
		exePath, err := os.Executable()
		if err != nil {
			log.Fatal("could not determine executable path:", err)
		}
		exeDir := filepath.Dir(exePath)

		fmt.Println(exeDir)
		envPath := filepath.Join(exeDir, ".env")

		err = godotenv.Load(envPath)
		if err != nil {
			log.Fatalf("Failed to load env in prod environment: %s", err)
		}
	} else {
		log.Fatalf("Invalid environment flag: <prod> or <dev> allowed")
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
