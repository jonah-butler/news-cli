package prompt

import "go-scraper/spider"

func ReduceMenuOption(options []MenuOption) []string {
	var reducedOptions []string
	for _, option := range options {
		reducedOptions = append(reducedOptions, option.Text)
	}
	return reducedOptions
}

func ReduceResultOption(results []spider.Result) []string {
	var reducedOptions []string
	for _, option := range results {
		reducedOptions = append(reducedOptions, option.Title)
	}
	return reducedOptions
}