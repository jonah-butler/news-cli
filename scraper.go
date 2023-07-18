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

// "https://richmond.com/search/?nsa=eedition&app=editorial&d1=2023-07-17&d2=2023-07-18&s=start_time&sd=asc&l=100&t=article&nfl=ap"
// "https://richmond.com/news/local/weather/how-did-the-james-river-rise-so-fast-our-chief-meteorologist-explains/article_63a94ada-24e1-11ee-973f-6fc456f72867.html#tracking-source=home-top-story",
// const BASE_URL = "https://richmond.com"
// const ARTICLE_BODY = "section#main-page-container"
// const ARTICLE_TITLE = "h1.headline span"
// const ARTICLE_TEXT = ".lee-article-text"
// const RESULTS_CONTAINER = "div#results-col"
// const RESULTS_LINK = "h3.tnt-headline a"

func loadEnv() {
	err := godotenv.Load(".env")
	if err != nil {
		log.Fatalf("failed to load env file - can't run scraper without html declarations")
	}
}

func initSpider() *Spider {
	spider := Spider{
		"news crawler",
		make([]Article, 0),
		colly.NewCollector(),
		1,
		"https://richmond.com/search/?nsa=eedition&app=editorial&d1=2023-07-17&d2=2023-07-18&s=start_time&sd=asc&l=25&t=article&nfl=ap",
		"",
		"",
	}
	return &spider
}

func (s Spider) Test() int {
	return 0
}

func main() {

	loadEnv()

	// links on search page are relative, so need base URL
	BASE_URL := os.Getenv("BASE_URL")

	// all html tags needed to get related article info
	ARTICLE_BODY := os.Getenv("ARTICLE_BODY")
	ARTICLE_TITLE := os.Getenv("ARTICLE_TITLE")
	ARTICLE_TEXT := os.Getenv("ARTICLE_TEXT")

	// all html tags needed to get all related search data
	RESULTS_CONTAINER := os.Getenv("RESULTS_CONTAINER")
	RESULTS_LINK := os.Getenv("RESULTS_LINK")

	spider := initSpider()

	spider.C.SetRequestTimeout(120 * time.Second)

	spider.C.OnRequest(func(r *colly.Request) {

		fmt.Println("visiting site: ", r.URL)

	})

	spider.C.OnHTML(RESULTS_CONTAINER, func(e *colly.HTMLElement) {

		// title := e.ChildText(ARTICLE_TITLE)
		// fmt.Println("article title: ", title)
		// data := ""

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

		// f, err := os.Create("article1.txt")

		// if err != nil {
		// 	panic("failed to write article")
		// }

		// defer f.Close()

		// _, err = f.WriteString(data)

		// if err != nil {
		// 	panic("failed to write string data to article.txt")
		// }

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
