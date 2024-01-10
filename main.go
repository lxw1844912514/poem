package main

import (
	"fmt"
	"gin/go-poem/gofish"
	"gin/go-poem/handle"
)

func GOScrape() {
	// Request the HTML page.
	authors := "https://so.gushiwen.cn/authors/"

	h := handle.AuthorHandle{}
	fish := gofish.NewGoFish()
	request, err := gofish.NewRequest("GET", authors, gofish.UserAgent, &h, nil)
	if err != nil {
		fmt.Println(err)
		return
	}

	fish.Request = request
	fish.Visit()

	/*res, err := http.Get(authors)
	if err != nil {
		log.Fatal(err)
	}
	defer res.Body.Close()
	if res.StatusCode != 200 {
		log.Fatalf("status code error: %d %s", res.StatusCode, res.Status)
	}

	// Load the HTML document
	doc, err := goquery.NewDocumentFromReader(res.Body)
	if err != nil {
		log.Fatal(err)
	}*/
	// Find the review items
	/*doc.Find(".left-content article .post-title").Each(func(i int, s *goquery.Selection) {
		// For each item found, get the title
		title := s.Find("a").Text()
		fmt.Printf("Review %d: %s\n", i, title)
	})*/
	/*doc.Find(".sons").Find(".cont").Find("a").Each(func(i int, s *goquery.Selection) {
		author := s.Text()
		fmt.Printf("%d author= %s \n", i, author)

		link, _ := s.Attr("href")
		fmt.Printf("%d link= %s \n", i, link)
	})*/
}

//func AppCrawler() {
//
//}

func main() {
	//web端数据少
	GOScrape()

	//	app 端
	//AppCrawler()

}
