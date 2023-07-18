package main

import (
	"fmt"
	"os"
	"time"

	"github.com/gocolly/colly"
)

type Blog struct {
	Title   string
	Url     string
}

type Spider struct {
	Name string
	Data []Blog
	C *colly.Collector
	Page int
	ActiveUrl string
	PreviousUrl string
	Body string
}

func initSpider() *Spider {
	spider := Spider{
		"rvalibrary.org crawler",
		make([]Blog, 0),
		colly.NewCollector(),
		1,
		"https://richmond.com/news/local/weather/how-did-the-james-river-rise-so-fast-our-chief-meteorologist-explains/article_63a94ada-24e1-11ee-973f-6fc456f72867.html#tracking-source=home-top-story",
		"",
		"",
	}
	return &spider
}

// func (s Spider) Test() int {
// 	fmt.Println("hello from spider method")
// 	return 0
// }


func main() {

	// spider := Spider{
	// 	"rvalibrary.org crawler",
	// 	make([]Blog, 0),
	// 	colly.NewCollector(),
	// 	1,
	// 	"https://rvalibrary.org/shelf-respect/",
	// }

	spider := initSpider()

	// c := colly.NewCollector()

	spider.C.SetRequestTimeout(120 * time.Second)

	spider.C.OnRequest(func(r *colly.Request) {

		fmt.Println("visiting site: ", r.URL)

	})

	// get all related blog data
	// "div#article-body"
	spider.C.OnHTML("div#article-body", func(e *colly.HTMLElement) {

		data := ""

		e.ForEach(".lee-article-text", func(i int, h *colly.HTMLElement) {

			data += h.Text
			data += "\n"

		})


		f, err := os.Create("article1.txt")

		if err != nil {
			panic("failed to write article")
		}

		defer f.Close()

		_, err = f.WriteString(data)

		if err != nil {
			panic("failed to write string data to article.txt")
		}

		// if(len(articleText) > 0) {

		// 	f, err := os.Create("article.txt")
			
		// 	if err != nil {
		// 		panic("failed to create article.txt")
		// 	}

		// 	defer f.Close()

		// 	_, err = f.WriteString(articleText)
			

		// 	if err != nil {
		// 		panic("failed to write article contents...")
		// 	}

			// f2, _ := os.Open("article.txt")

			// scanner := bufio.NewScanner(f2)

			// for scanner.Scan() {
			// 	line := scanner.Text()
			// 	fmt.Println("text file values: ", line)
			// }

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