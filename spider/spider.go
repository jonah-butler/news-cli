package spider

import (
	"errors"
	"fmt"
	"net/url"
	"os"
	"strconv"
	"time"

	"github.com/gocolly/colly"
)

var SEARCH_URL string

const REQUEST_TIMEOUT = 120 * time.Second

const SAVED_ARTICLES_DIR = "saved_files"

type Article struct {
	Title string
	Body  string
	Url   string
}

type Spider struct {
	Name      string
	ActiveUrl string
	Data      *Article
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
	IsQuery   bool
}

var (
	Crawler *Spider
)

func InitSpider(name string, elements Elements) {
	spider := Spider{
		name,
		"",
		&Article{},
		colly.NewCollector(),
		elements,
		SearchResults{},
	}
	spider.Search.Url = SEARCH_URL
	Crawler = &spider
}

func (s Spider) Clone(name string) *Spider {
	spider := Spider{
		name,
		"",
		&Article{},
		s.C.Clone(),
		s.Html,
		s.Search,
	}
	return &spider
}

func(s *Spider) SaveToTextFile() {
	var title string
	if s.Data.Title != "" {
		title = s.Data.Title
	} else {
		// convert current unix timestamp to string for
		// non-overlapping file titles
		title = strconv.FormatInt(time.Now().Unix(), 10)
	}

	// make directory if it does not exist
	if _, err := os.Stat(SAVED_ARTICLES_DIR); errors.Is(err, os.ErrNotExist) {
		err := os.Mkdir(SAVED_ARTICLES_DIR, os.ModePerm)
		if err != nil {
			fmt.Println("Failed to make directory: ", err)
		}
	}

	f, err := os.Create(SAVED_ARTICLES_DIR + "/" + title + ".txt")
	if err != nil {
		fmt.Printf("Failed to save file: %s", title)
	}

	defer f.Close()

	articleBody := fmt.Sprintf("%s\n\n\nurl: %s", s.Data.Body, s.Data.Url)

	_, err = f.WriteString(articleBody)
	if err != nil {
		fmt.Printf("Failed to write article body to file: %s", title)
	} else {
		fmt.Printf("Write to file %s was successful", title + ".txt")
	}
}

func (s *Spider) ClearSetValues() {
	s.Search.IsQuery = false
	// eventually put these values in env
	// to allow for dynamic usage
	s.FlushQueryParam("q")
	s.FlushQueryParam("d1")
	s.FlushQueryParam("d2")
	s.FlushQueryParam("o")
	s.Search.DateEnd = ""
	s.Search.DateStart = ""
}

func (s *Spider) FlushQueryParam(param string) {
	u, err := url.Parse(s.Search.Url)
	if err != nil {
		fmt.Println("failed to parse query param in FlushQueryParam")
	}

	q := u.Query()
	q.Del(param)
	u.RawQuery = q.Encode()
	s.Search.Url = u.String()
}

func (s *Spider) AppendQueryParam(query string) {
	u, err := url.Parse(s.Search.Url)
	if err != nil {
		fmt.Println("failed to parse query param in AppendQueryParam")
	}

	vals := u.Query()

	vals.Add("q", query)

	u.RawQuery = vals.Encode()

	s.Search.Url = s.Html.BaseUrl + "/search/?" + u.RawQuery
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
		
		Crawler.Data = &article

	})

	s.C.OnError(func(r *colly.Response, e error) {
		fmt.Println("COLLY ERROR - endpoint: ", endpoint, "\nERROR: ", e)
	})

	// release the spider
	s.C.Visit(s.ActiveUrl)
}

func (s *Spider) GetArticleLinks(endpoint string) {

	s.C.SetRequestTimeout(REQUEST_TIMEOUT)

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
		previous := e1.ChildAttr("li.prev a", "href")
		if previous != "" {
			s.BuildAndStoreResultsLink(previous, "Previous")
		}

		// get next previous button href
		next := e1.ChildAttr("li.next a", "href")
		if next != "" {
			s.BuildAndStoreResultsLink(next, "Next")
		}

		// add default option to go back to main menu
		s.BuildAndStoreResultsLink("/", "Back To Main Menu")

	})

	s.C.OnError(func(r *colly.Response, e error) {
		fmt.Println("COLLY ERROR - endpoint: ", s.Search.Url, "\nERROR: ", e)
	})
	err := s.C.Visit(s.Search.Url)

	if err != nil {
		fmt.Println("ERROR RETURNED FROM VISIT: ", err)
	}
}

func (s *Spider) ClearStoredArticleLinks() {
	if len(s.Search.Results) != 0 {
		s.Search.Results = []Result{}
	}
}
