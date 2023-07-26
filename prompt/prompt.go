package prompt

import (
	"errors"
	"fmt"
	"go-scraper/spider"
	"strings"
	"time"

	"net/url"

	"github.com/manifoldco/promptui"
)

func GetDateRanges() {
	d1 := func(date1 string) error {
		layout := "2006-01-02"
		_, err := time.Parse(layout, date1)
		if err != nil {
			return errors.New("date must be in format: 2006-02-01")
		}
		return nil
	}

	date1Prompt := promptui.Prompt {
		Label: "Enter start date",
		Validate: d1,
	}

	result, err := date1Prompt.Run()

	if err != nil {
		fmt.Println("date 1 prompt failed")
	}

	spider.Crawler.Search.DateStart = "&d1="+result


	d2 := func(date1 string) error {
		layout := "2006-01-02"
		_, err := time.Parse(layout, date1)
		if err != nil {
			return errors.New("date must be in format: 2006-02-01")
		}
		return nil
	}

	date2Prompt := promptui.Prompt {
		Label: "Enter end date",
		Validate: d2,
	}

	result, err = date2Prompt.Run()

	if err != nil {
		fmt.Println("date 1 prompt failed")
	}

	spider.Crawler.Search.DateEnd = "&d2="+result

	spider.Crawler.Search.Url += spider.Crawler.Search.DateStart + spider.Crawler.Search.DateEnd

}

func GetLatestHeadlines(c func()) {
	if spider.Crawler.Search.DateStart == "" || spider.Crawler.Search.DateEnd == "" {
		GetDateRanges()
	}

	spider.Crawler.GetArticleLinks(spider.Crawler.Search.Url)

	template := &promptui.SelectTemplates{
		Label:    "{{ . }}",
		Active:   "\U0001F50E {{ .Title | green }}",
		Inactive: "  {{ .Title | cyan }}",
		Selected: "\U0001F50E {{ .Title | green }}",
	}

	searcher := func(input string, index int) bool {
		article := spider.Crawler.Search.Results[index]
		name := strings.Replace(strings.ToLower(article.Title), " ", "", -1)
		input = strings.Replace(strings.ToLower(input), " ", "", -1)

		return strings.Contains(name, input)
	}

	articlePrompt := promptui.Select{
		Label:     "Article Search",
		Items:     spider.Crawler.Search.Results,
		Templates: template,
		Size:      10,
		Searcher:  searcher,
	}

	i, _, err := articlePrompt.Run()

	if err != nil {
		fmt.Println("failed to initialize article prompt...", err)
	}

	selectedArticle := spider.Crawler.Search.Results[i]

	if selectedArticle.Title == "Next" || selectedArticle.Title == "Previous" {
		spider.Crawler = spider.Crawler.Clone("news spider", selectedArticle.Url)
		spider.Crawler.Search.Url = selectedArticle.Url
		GetLatestHeadlines(InArticleMenu)
	} else {
		clone := spider.Crawler.Clone("single article crawler...", spider.Crawler.Search.Results[i].Url)
		clone.GetArticle(spider.Crawler.Search.Results[i].Url)
	}

	c()

}

func ReadSingleArticle(c func()) {
	validateInput := func(articleUrl string) error {
		_, err := url.ParseRequestURI(articleUrl)
		if err != nil {
			return errors.New("invalid url scheme - try again")
		}
		return nil
	}

	prompt := promptui.Prompt{
		Label:    "Enter Article URL",
		Validate: validateInput,
	}

	result, err := prompt.Run()

	if err != nil {
		fmt.Println("Article Prompt failed")
	}

	spider.Crawler.GetArticle(result)
	c()

}

func InArticleMenu() {
	validateCommand := func(command string) error {
		switch command {
		case "back":
			return nil
		default:
			return errors.New("not a valid command")
		}
	}

	prompt := promptui.Prompt{
		Label:    "[back]",
		Validate: validateCommand,
	}

	result, err := prompt.Run()

	if err != nil {
		fmt.Println("Article Prompt failed")
	}
	
	if result == "back" {
		if spider.Crawler.Search.DateStart == "" {
			// menus.InitializePrompts()
		} else {
			spider.Crawler = spider.Crawler.Clone("news spider", spider.Crawler.Search.Url)
			GetLatestHeadlines(InArticleMenu)
		}
	}
}
