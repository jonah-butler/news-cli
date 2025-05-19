package prompt

import (
	"errors"
	"fmt"
	"go-scraper/mail"
	"go-scraper/spider"
	"os"
	"strings"
	"time"

	"net/url"

	"github.com/manifoldco/promptui"
)

func SendEmail() {

	mail.SetupSMTPAuth()

	// get recipient address
	v1 := func(emailAddress string) error {
		if emailAddress == "" {
			return errors.New("email address can not be empty")
		}
		return nil
	}

	toPrompt := promptui.Prompt{
		Label:    "Recipient address:",
		Validate: v1,
	}

	recipient, err := toPrompt.Run()
	if err != nil {
		fmt.Printf("Failed to run recipient prompt: %s", err.Error())
	}

	// get email subject
	v2 := func(subject string) error {
		if subject == "" {
			return errors.New("subject line can not be empty")
		}
		return nil
	}

	subjectPrompt := promptui.Prompt{
		Label:    "Email Subject:",
		Validate: v2,
	}

	subject, err := subjectPrompt.Run()
	if err != nil {
		fmt.Printf("Failed to run subject prompt: %s", err.Error())
	}

	// get additional body
	v3 := func(body string) error {
		return nil
	}

	bodyPrompt := promptui.Prompt{
		Label:    "Enter additional body to email[optional]:",
		Validate: v3,
	}

	body, err := bodyPrompt.Run()
	if err != nil {
		fmt.Printf("Failed to run additional body prompt: %s", err.Error())
	}

	err = mail.SendMail(body, recipient, subject, spider.Crawler.Data.Title, spider.Crawler.Data.Body)
	if err != nil {
		fmt.Printf("Error sending email to: %s - %s", recipient, err.Error())
	}

	InArticleMenu()

}

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
		spider.Crawler.C.AllowURLRevisit = true
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

func ExitApplication(c func()) {
	os.Exit(0)
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
		case "email":
			return nil
		case "exit":
			return nil
		default:
			return errors.New("not a valid command")
		}
	}

	prompt := promptui.Prompt{
		Label:    "[back, save, email, exit]",
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
	} else if result == "email" {
		SendEmail()
	} else if result == "exit" {
		os.Exit(0)
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
			{
				Text:   "Exit CLI",
				Action: ExitApplication,
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
