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
	endpoint := "https://richmond.com/search/?nsa=eedition&app=editorial&d1=2023-07-17&d2=2023-07-18&s=start_time&sd=asc&l=25&t=article&nfl=ap"
	spider.Crawler.GetArticleLinks(endpoint)

	template := &promptui.SelectTemplates{
		Label: "{{ . }}",
		Active: "\U0001F4F0 {{ .Title | green }}",
		Inactive: "  {{ .Title | cyan }}",
		Selected: "\U0001F4F0 {{ .Title | green }}",
	}

	searcher := func(input string, index int) bool {
		article := spider.Crawler.Search.Results[index]
		name := strings.Replace(strings.ToLower(article.Title), " ", "", -1)
		input = strings.Replace(strings.ToLower(input), " ", "", -1)

		return strings.Contains(name, input)
	}

	articlePrompt := promptui.Select {
		Label: "Article Search",
		Items: spider.Crawler.Search.Results,
		Templates: template,
		Size: 10,
		Searcher: searcher,
	}

	i, _, err := articlePrompt.Run()

	if err != nil {
		fmt.Println("failed to initialize article prompt...", err)
	}

	clone := spider.Crawler.Clone("single article crawler...", spider.Crawler.Search.Results[i].Url)

	clone.GetArticle(spider.Crawler.Search.Results[i].Url)

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