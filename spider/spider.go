package spider

import (
	"fmt"
	"strconv"
	"time"

	"github.com/gocolly/colly"
)

const REQUEST_TIMEOUT = 120 * time.Second
const DUMMY_URL = "https://richmond.com/search/?nsa=eedition&app=editorial&d1=2023-07-17&d2=2023-07-18&s=start_time&sd=asc&l=25&t=article&nfl=ap"

type Article struct {
	Title string
	Body  string
	Url   string
}

type Spider struct {
	Name      string
	ActiveUrl string
	Data      []Article
	C         *colly.Collector
	Html      Elements
	Search    SearchResults
}

type Elements struct {
	BaseUrl          string
	ArticleBody      string
	ArticleTitle     string
	ArticleText      string
	ResultsContainer string
	ResultsLink      string
}

type Result struct {
	Title string
	Url   string
}

type SearchResults struct {
	Url       string
	DateStart string
	DateEnd   string
	Results   []Result
	Limit     int
}

var (
	Crawler *Spider
)

func InitSpider(name string, elements Elements) {
	spider := Spider{
		name,
		"",
		make([]Article, 0),
		colly.NewCollector(),
		elements,
		SearchResults{},
	}
	spider.Search.Url = DUMMY_URL
	Crawler = &spider
}

func (s Spider) Clone(name string, url string) *Spider {
	spider := Spider{
		name,
		"",
		make([]Article, 0),
		s.C.Clone(),
		s.Html,
		s.Search,
	}
	return &spider
}

/**
* appends base url and relative article link
* to form full path, usable by spider to scrape
* article contents
**/
func (s *Spider) BuildAndStoreResultsLink(relativeUrl string, title string) {
	result := Result{
		title,
		s.Html.BaseUrl + relativeUrl,
	}
	s.Search.Results = append(s.Search.Results, result)
}

/**
* SETS UP REQUEST FOR A PARTICULAR ARTICLE AND BUILDS ARTICLE
* within a given article endpoint
* - store article url
* - store article title
* - loop through article text paragraphs and build article body
**/
func (s *Spider) GetArticle(endpoint string) {
	s.ActiveUrl = endpoint
	s.C.SetRequestTimeout(REQUEST_TIMEOUT)

	s.C.OnRequest(func(r *colly.Request) {
		fmt.Println("visiting endpoint: ", endpoint)
	})

	// GRAB NECESSARY ARTICLE FILEDS ONCE HTML IS LOADED
	s.C.OnHTML(s.Html.ArticleBody, func(e1 *colly.HTMLElement) {

		article := Article{}

		article.Url = endpoint
		article.Title = e1.ChildText(s.Html.ArticleTitle)

		if article.Title == "" {
			article.Title = strconv.FormatInt(time.Now().UTC().UnixMilli(), 10)
		}

		article.Body += "\n\n"

		e1.ForEach(s.Html.ArticleText, func(i int, e2 *colly.HTMLElement) {

			article.Body += e2.Text
			article.Body += "\n\n"

		})

		fmt.Println(article.Body)

	})

	s.C.OnError(func(r *colly.Response, e error) {
		fmt.Println("COLLY ERROR - endpoint: ", endpoint, "\nERROR: ", e)
	})

	// release the spider
	s.C.Visit(s.ActiveUrl)
}

func (s *Spider) GetArticleLinks(endpoint string) {
	s.C.SetRequestTimeout(REQUEST_TIMEOUT)

	s.C.OnRequest(func(r *colly.Request) {
		fmt.Println("scraping results page: ", s.Search.Url)
	})

	s.C.OnHTML(s.Html.ResultsContainer, func(e1 *colly.HTMLElement) {

		// new search results - so clear old search result links
		s.ClearStoredArticleLinks()

		e1.ForEach(s.Html.ResultsLink, func(_ int, e2 *colly.HTMLElement) {

			relativeUrl := e2.Attr("href")
			title := e2.Attr("aria-label")
			if title != "" {
				s.BuildAndStoreResultsLink(relativeUrl, title)
			}

		})

		// get previous button href
		previous := e1.ChildAttr("li.previous a", "href")
		if previous != "" {
			s.BuildAndStoreResultsLink(previous, "Previous")
		}

		// get next previous button href
		next := e1.ChildAttr("li.next a", "href")
		if next != "" {
			s.BuildAndStoreResultsLink(next, "Next")
		}

	})

	s.C.OnError(func(r *colly.Response, e error) {
		fmt.Println("COLLY ERROR - endpoint: ", s.Search.Url, "\nERROR: ", e)
	})

	s.C.Visit(endpoint)
}

func (s *Spider) ClearStoredArticleLinks() {
	if len(s.Search.Results) != 0 {
		s.Search.Results = []Result{}
	}
}
