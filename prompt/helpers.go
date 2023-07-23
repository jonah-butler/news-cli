package prompt

func ReduceMenuOption(options []MenuOption) []string {
	var reducedOptions []string
	for _, option := range options {
		reducedOptions = append(reducedOptions, option.Text)
	}
	return reducedOptions
}