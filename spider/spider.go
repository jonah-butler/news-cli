package spider

import (
	"github.com/gocolly/colly"
)

type Article struct {
	Title string
	Body  string
	Url   string
}

type Spider struct {
	Name        string
	Data        []Article
	C           *colly.Collector
	Page        int
	ActiveUrl   string
	PreviousUrl string
	Body        string
}

func (s Spider) Clone(name string, url string) *Spider {
	spider := Spider {
		name,
		make([]Article, 0),
		colly.NewCollector(),
		1,
		url,
		"",
		"",
	}
	return &spider
}

func (s Spider) GetArticle() {

}