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
	layout := "2006-01-02"
	d1 := func(date1 string) error {
		_, err := time.Parse(layout, date1)
		if err != nil {
			return errors.New("date must be in format: 2006-02-01")
		}
		return nil
	}

	today := time.Now().Format(layout)

	date1Prompt := promptui.Prompt{
		Label:    "Enter start date",
		Validate: d1,
		Default:  today,
	}

	result, err := date1Prompt.Run()

	if err != nil {
		fmt.Println("date 1 prompt failed")
	}

	spider.Crawler.Search.DateStart = "&d1=" + result

	d2 := func(date1 string) error {
		layout := "2006-01-02"
		_, err := time.Parse(layout, date1)
		if err != nil {
			return errors.New("date must be in format: 2006-02-01")
		}
		return nil
	}

	date2Prompt := promptui.Prompt{
		Label:    "Enter end date",
		Validate: d2,
		Default:  today,
	}

	result, err = date2Prompt.Run()

	if err != nil {
		fmt.Println("date 1 prompt failed")
	}

	spider.Crawler.Search.DateEnd = "&d2=" + result

	spider.Crawler.Search.Url += spider.Crawler.Search.DateStart + spider.Crawler.Search.DateEnd

}

func GetLatestHeadlines(c func()) {
	if !spider.Crawler.Search.IsQuery {
		if spider.Crawler.Search.DateStart == "" || spider.Crawler.Search.DateEnd == "" {
			GetDateRanges()
		}
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
		spider.Crawler = spider.Crawler.Clone("news spider")
		spider.Crawler.C.AllowURLRevisit = true
		spider.Crawler.Search.Url = selectedArticle.Url
		GetLatestHeadlines(InArticleMenu)
	} else if selectedArticle.Title == "Back To Main Menu" {
		spider.Crawler = spider.Crawler.Clone("news spider")
		spider.Crawler.C.AllowURLRevisit = true
		spider.Crawler.ClearSetValues()
		InitializePrompts()
	} else {
		clone := spider.Crawler.Clone("single article crawler...")
		clone.GetArticle(selectedArticle.Url)
		c()
	}

}

func RunSearchQuery(c func()) {
	validateQuery := func(query string) error {
		if query == "" {
			return errors.New("query must not be empty")
		}
		return nil
	}

	prompt := promptui.Prompt{
		Label:    "Enter a query to being your search",
		Validate: validateQuery,
	}

	result, err := prompt.Run()

	if err != nil {
		fmt.Println("error initializing query prompt")
	}

	spider.Crawler = spider.Crawler.Clone("news spider")
	spider.Crawler.AppendQueryParam(result)
	spider.Crawler.Search.IsQuery = true
	GetLatestHeadlines(InArticleMenu)

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
		case "save":
			return nil
		default:
			return errors.New("not a valid command")
		}
	}

	prompt := promptui.Prompt{
		Label:    "[back, save]",
		Validate: validateCommand,
	}

	result, err := prompt.Run()

	if err != nil {
		fmt.Println("Article Prompt failed")
	}

	if result == "back" {
		spider.Crawler = spider.Crawler.Clone("news spider")
		if !spider.Crawler.Search.IsQuery {
			spider.Crawler.FlushQueryParam("q")
		}
		GetLatestHeadlines(InArticleMenu)
	} else if result == "save" {
		spider.Crawler.SaveToTextFile()
		InArticleMenu()
	}
}

func MainMenu() MenuLevel {

	mainMenuOptions := MenuLevel{
		Prompt: "How would you like to begin?",
		MenuOptions: []MenuOption{
			{
				Text:   "Get Latest Headlines",
				Action: GetLatestHeadlines,
			},
			{
				Text:   "Read Single Article",
				Action: ReadSingleArticle,
			},
			{
				Text:   "Search Articles",
				Action: RunSearchQuery,
			},
		},
	}

	return mainMenuOptions

}

func InitializePrompts() {

	mainMenu := MainMenu()

	begin := promptui.Select{
		Label: mainMenu.Prompt,
		Items: ReduceMenuOption(mainMenu.MenuOptions),
	}

	i, _, err := begin.Run()

	if err != nil {
		fmt.Printf("Prompt failed %v\n", err)
	}

	mainMenu.MenuOptions[i].Action(InArticleMenu)

}
