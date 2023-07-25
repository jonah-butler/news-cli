package menus

import (
	"fmt"
	"go-scraper/prompt"

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

	mainMenu.MenuOptions[i].Action()

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
