package handle

import (
	"fmt"
	"gin/go-poem/db"
	"github.com/PuerkitoBio/goquery"
	"io"
	"log"
	"strings"
)

var poemHome = "https://so.gushiwen.cn/shiwens/default.aspx?page=1&tstr=&astr=authorStr&cstr=&xstr="
var mingJuHome = "https://so.gushiwen.cn/mingjus/default.aspx?page=1&tstr=&astr=authorStr&cstr=&xstr="

func getUrls(url string, pageSize int, author string) []string {
	urls := make([]string, 0)
	//替换页码
	urlTmp := strings.Replace(url, "page=1", "page=%d", 1)

	//替换作者
	urlTmp = strings.Replace(urlTmp, "astr=authorStr", "astr="+author, 1)

	for i := 1; i <= pageSize; i++ {
		urls = append(urls, fmt.Sprintf(urlTmp, i))

		fmt.Println(fmt.Sprintf(urlTmp, i))
	}
	return urls
}

type PoemInfoHandle struct {
}

func (h *PoemInfoHandle) Worker(body io.Reader, url string) {
	doc, err := goquery.NewDocumentFromReader(body)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("元素数：",doc.Find("#leftZhankai").Find(".sons").Find(".cont").Length())
	//fmt.Println("最后一个元素：",doc.Find("#leftZhankai").Find(".sons").Find(".cont").Last().Text())

	doc.Find("#leftZhankai").Find(".sons").Find(".cont").Each(func(i int, s *goquery.Selection) {
		author := ""
		title := ""
		dynsty := ""
		content := ""

		//获取 作者朝代标题
		title = strings.TrimSpace(s.Find("p").Find("b").Text())
		authoranddynsty := strings.TrimSpace(s.Find(".source").Text())
		authoranddynstySlice := strings.Split(authoranddynsty, "〔")
		if len(authoranddynstySlice) == 2 {
			author = authoranddynstySlice[0]
			dynsty = strings.Trim(authoranddynstySlice[1], "〔〕")
		}
		fmt.Printf("第%d条：,作者：%s, 朝代：%s,标题：%s,请求地址: %s  \n",i, author, dynsty, title, url)

		//获取内容
		s.Find(".contson").Each(func(i int, s *goquery.Selection) {
			content = strings.TrimSpace(s.Text())
			fmt.Printf("内容：%s \n", content)
		})

		//
		if author != "" && dynsty != "" && title != "" && content != "" {
			p := db.Poem{}
			p.Author = author
			p.Dynasty = dynsty
			p.Title = title
			p.Content = content
			p.Save()
		}
	})
}
