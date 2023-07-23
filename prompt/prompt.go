package prompt

import (
	"errors"
	"fmt"
	"go-scraper/spider"
	"net/url"

	"github.com/manifoldco/promptui"
)

func GetLatestHeadlines() {
	fmt.Println("fetching latest headlines...")
}

func ReadSingleArticle() {
	validateInput := func(articleUrl string) error {
		_, err := url.ParseRequestURI(articleUrl)
		if err != nil {
			return errors.New("invalid url scheme - try again")
		}
		return nil
	}

	prompt := promptui.Prompt{
		Label: "Enter Article URL",
		Validate: validateInput,
	}

	result, err := prompt.Run()

	if err != nil {
		fmt.Println("Article Prompt failed")
	}

	spider.Crawler.GetArticle(result)
	
}