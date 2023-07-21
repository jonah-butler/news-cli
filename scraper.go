package main

import (
	"fmt"
	"log"
	"os"
	"strconv"
	"time"

	"github.com/gocolly/colly"
	"github.com/joho/godotenv"
)

// https://richmond.com/search/?nsa=eedition&app=editorial&d1=2023-07-17&d2=2023-07-18&s=start_time&sd=asc&l=25&t=article&nfl=ap",

type Article struct {
	Title string
	Body  string
	Url   string
}

type Spider struct {
	// name of crawler
	Name string
	// stores scraped article data
	Data []Article
	// the spider isntance - can be cloned
	C *colly.Collector
	// current page if crawling a results page
	Page int
	// limit results for results page - query param
	// also used in results offset if more than one
	// page of results exists
	Limit int
	// active url
	ActiveUrl string
	// previously visited url
	PreviousUrl string
	// temp storage for article body
	Body string
	// all required html elements for
	// navigating result and single article pages
	Html Elements
}

type Elements struct {
	// base url for navigation
	BaseUrl string
	// article body container
	ArticleBody string
	// article title inside ArticleBody
	ArticleTitle string
	// article text inside ArticleBody
	ArticleText string
	// search results container
	ResultsContainer string
	// link inside results container
	ResultsLink string
}

const REQUEST_TIMEOUT = 120

func LoadEnv() {
	err := godotenv.Load(".env")
	if err != nil {
		log.Fatalf("failed to load env file - can't run scraper without html declarations")
	}
}

func initSpider(url string, elements Elements) *Spider {
	spider := Spider{
		"news crawler",
		make([]Article, 0),
		colly.NewCollector(),
		1,
		20,
		url,
		"",
		"",
		elements,
	}
	return &spider
}

func (s Spider) Clone(name string, url string) *Spider {
	spider := Spider{
		name,
		make([]Article, 0),
		s.C.Clone(),
		s.Page,
		s.Limit,
		url,
		"",
		"",
		s.Html,
	}
	return &spider
}

func main() {

	LoadEnv()

	// links on search page are relative, so need base URL
	// to do - delete this
	BASE_URL := os.Getenv("BASE_URL")

	// all html tags needed to get related article info
	// to do - delete these
	ARTICLE_BODY := os.Getenv("ARTICLE_BODY")
	ARTICLE_TITLE := os.Getenv("ARTICLE_TITLE")
	ARTICLE_TEXT := os.Getenv("ARTICLE_TEXT")

	// all html tags needed to get all related search data
	// to do - delete these
	RESULTS_CONTAINER := os.Getenv("RESULTS_CONTAINER")
	RESULTS_LINK := os.Getenv("RESULTS_LINK")

	htmlElements := Elements{
		// links on search page are relative, so need base URL
		os.Getenv("BASE_URL"),
		// all html tags needed to get related article info
		os.Getenv("ARTICLE_BODY"),
		os.Getenv("ARTICLE_TITLE"),
		os.Getenv("ARTICLE_TEXT"),
		// all html tags needed to get all related search data
		os.Getenv("RESULTS_CONTAINER"),
		os.Getenv("RESULTS_LINK"),
	}

	spider := initSpider("news crawler", htmlElements)

	spider.C.SetRequestTimeout(120 * time.Second)

	spider.C.OnRequest(func(r *colly.Request) {

		fmt.Println("visiting site: ", r.URL)

	})

	spider.C.OnHTML(RESULTS_CONTAINER, func(e *colly.HTMLElement) {

		clone := spider.C.Clone()

		e.ForEach(RESULTS_LINK, func(_ int, h *colly.HTMLElement) {

			url := h.Attr("href")
			article := BASE_URL + url

			clone.SetRequestTimeout(120 * time.Second)

			clone.OnHTML(ARTICLE_BODY, func(e *colly.HTMLElement) {

				article := Article{}

				article.Url = url
				article.Title = e.ChildText(ARTICLE_TITLE)

				tempInt := 0

				e.ForEach(ARTICLE_TEXT, func(i int, h2 *colly.HTMLElement) {

					tempInt = i

					article.Body += h2.Text
					article.Body += "\n"

				})

				fileName := ""

				if article.Title == "" {
					fileName = strconv.Itoa(tempInt) + ".txt"
				} else {
					fileName = article.Title + ".txt"
				}

				f, err := os.Create(fileName)

				if err != nil {
					panic("failed to write article")
				}

				defer f.Close()

				_, err = f.WriteString(article.Body)

				if err != nil {
					panic("failed to write string data to article.txt")
				}

			})

			clone.Visit(article)

		})

	})

	// get next page
	// spider.C.OnHTML("a.page-numbers.next", func(e *colly.HTMLElement) {

	// 	spider.PreviousUrl = spider.ActiveUrl
	// 	spider.ActiveUrl = e.Attr("href")
	// 	spider.Page = spider.Page + 1

	// })

	spider.C.OnResponse(func(r *colly.Response) {

		fmt.Println("got a response from: ", r.Request)

	})

	spider.C.OnError(func(r *colly.Response, e error) {

		fmt.Println("received an error from colly: ", e)

	})

	// spider.C.OnScraped(func(r *colly.Response) {
	// 	fmt.Println("scrape finshed: ", r.StatusCode)
	// 	if spider.Page < 3 {

	// 		spider.C.Visit(spider.ActiveUrl)

	// 	}
	// })

	spider.C.Visit(spider.ActiveUrl)
}
