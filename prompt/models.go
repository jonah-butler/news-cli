package prompt

type MenuLevel struct {
	Prompt      string
	MenuOptions []MenuOption
}

type MenuOption struct {
	Text   string
	Action func()
}
