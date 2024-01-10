package handle

import (
	"fmt"
	"gin/go-poem/gofish"
	"github.com/PuerkitoBio/goquery"
	"io"
	"log"
	"strings"
)

var Baseurl = "https://www.gushiwen.cn"
var SOBaseurl = "https://so.gushiwen.cn"
var BaseurlShiWen = "https://so.gushiwen.cn/shiwens/default.aspx?astr="
var BaseurlMingJu = "https://so.gushiwen.cn/mingjus/default.aspx?astr="

type AuthorHandle struct {
}

func (h *AuthorHandle) Worker(body io.Reader, url string) {
	//1.获取作者列表地址link
	doc, err := goquery.NewDocumentFromReader(body)
	if err != nil {
		log.Fatal(err)
	}

	doc.Find(".sons").Find(".cont").Find("a").Each(func(i int, s *goquery.Selection) {
		author := s.Text()
		fmt.Printf("%d author= %s \n", i, author)

		link, _ := s.Attr("href")
		fmt.Printf("%d link= %s \n", i, link)

		//2.获取作者生平介绍
		h := PoemHomeHandle{}
		fish := gofish.NewGoFish()
		request, err := gofish.NewRequest("GET", SOBaseurl+link, gofish.UserAgent, &h, nil)
		if err != nil {
			fmt.Println(err)
			return
		}

		//fmt.Println(request.Url)

		fish.Request = request
		fish.Visit()

	})
}

type PoemHomeHandle struct {
}

func (h *PoemHomeHandle) Worker(body io.Reader, url string) {
	doc, err := goquery.NewDocumentFromReader(body)
	if err != nil {
		log.Fatal(err)
	}

	//获取第一个匹配的元素:a标签
	//doc.Find(".sonspic").Find(".cont").Find("p").Find("a").First().Each(func(i int, s *goquery.Selection) {
	doc.Find(".sonspic").Find(".cont").Find("p").Find("a").Each(func(i int, s *goquery.Selection) {
		link, _ := s.Attr("href")
		shiwenUrl := ""
		mingjuUrl := ""
		if strings.Contains(link, "shiwens") {
			shiwenUrl = SOBaseurl + link
			fmt.Println("诗文=", shiwenUrl)
		} else {
			mingjuUrl = SOBaseurl + link
			fmt.Println("名句=", mingjuUrl)
		}

		//fmt.Println("诗文Or名句=", SOBaseurl+link)

		h := PoemInfoHandle{}
		fish := gofish.NewGoFish()
		request, err := gofish.NewRequest("GET", shiwenUrl, gofish.UserAgent, &h, nil)
		if err != nil {
			fmt.Println(err)
			return
		}
		fish.Request = request
		fish.Visit()
	})

}
