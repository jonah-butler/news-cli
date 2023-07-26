package menus

import (
	"errors"
	"fmt"
	"go-scraper/prompt"
	"go-scraper/spider"

	"github.com/manifoldco/promptui"
)

func InitializePrompts() {

	mainMenu := MainMenu()

	begin := promptui.Select{
		Label: mainMenu.Prompt,
		Items: prompt.ReduceMenuOption(mainMenu.MenuOptions),
	}

	i, _, err := begin.Run()

	if err != nil {
		fmt.Printf("Prompt failed %v\n", err)
	}

	mainMenu.MenuOptions[i].Action(InArticleMenu)

}

func MainMenu() prompt.MenuLevel {

	mainMenuOptions := prompt.MenuLevel{
		Prompt: "How would you like to begin?",
		MenuOptions: []prompt.MenuOption{
			{
				Text:   "Get Latest Headlines",
				Action: prompt.GetLatestHeadlines,
			},
			{
				Text:   "Read Single Article",
				Action: prompt.ReadSingleArticle,
			},
		},
	}

	return mainMenuOptions

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

	backPrompt := promptui.Prompt{
		Label:    "[back]",
		Validate: validateCommand,
	}

	result, err := backPrompt.Run()

	if err != nil {
		fmt.Println("Article Prompt failed")
	}
	
	if result == "back" {
		if spider.Crawler.Search.DateStart == "" {
			InitializePrompts()
		} else {
			spider.Crawler = spider.Crawler.Clone("news spider", spider.Crawler.Search.Url)
			prompt.GetLatestHeadlines(InArticleMenu)
		}
	}
}
