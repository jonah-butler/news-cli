package prompt

import (
	"errors"
	"fmt"
	"go-scraper/spider"
	"strings"

	"net/url"

	"github.com/manifoldco/promptui"
)

func GetLatestHeadlines() {
	spider.Crawler.GetArticleLinks(spider.Crawler.Search.Url)

	template := &promptui.SelectTemplates{
		Label:    "{{ . }}",
		Active:   "\U0001F4F0 {{ .Title | green }}",
		Inactive: "  {{ .Title | cyan }}",
		Selected: "\U0001F4F0 {{ .Title | green }}",
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
		GetLatestHeadlines()
	} else {
		clone := spider.Crawler.Clone("single article crawler...", spider.Crawler.Search.Results[i].Url)
		clone.GetArticle(spider.Crawler.Search.Results[i].Url)
	}

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
		Label:    "Enter Article URL",
		Validate: validateInput,
	}

	result, err := prompt.Run()

	if err != nil {
		fmt.Println("Article Prompt failed")
	}

	spider.Crawler.GetArticle(result)

}
