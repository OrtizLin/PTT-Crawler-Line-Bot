package crawler

import (
	"github.com/PuerkitoBio/goquery"
	"line_bot_final/db"
	"log"
	"net/http"
	"strconv"
	"strings"
	"time"
)

const BasePttAddress = "https://www.ptt.cc"

type Article struct {
	Title           string
	LikeCount       int
	Link            string
	Date            string
	ImageLink       string
	LikeCountString string
}

func Start(w http.ResponseWriter, r *http.Request) {
	db.RemoveALL()
	getAllArticles("Beauty")
	getAllArticles("Sex")
}

func getAllArticles(fourm string) {

	var BOOL = true
	var exist = true
	var url string = ""  // default url
	var href string = "" // next page url
	var crawlerCount = 0

	// today's date
	loc, _ := time.LoadLocation("Asia/Chongqing")
	time := time.Now().In(loc)

	// 開始爬蟲
	for BOOL {

		if href == "" {
			url = BasePttAddress + "/bbs/" + fourm + "/index.html" // 首頁
		} else {
			url = BasePttAddress + href // 翻至下一頁
		}

	// 設定 header 以及 滿18歲cookie
	client:=&http.Client{}
	req, err := http.NewRequest("GET", url, nil)
	req.Header.Add("User-Agent", "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_12_5) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/58.0.3029.110 Safari/537.36")
	req.Header.Add("Referer", url)
	cookie := http.Cookie {
		Name: "over18",
		Value: "1",
	}
	req.AddCookie(&cookie)
	res, err := client.Do(req)
	defer res.Body.Close()

	// 最後直接把res傳给goquery就可以來解析網頁
		doc, err := goquery.NewDocumentFromResponse(res)
		if err != nil {
			log.Fatal(err)
		}

		//Find previous link
		doc.Find(".btn-group a").Each(func(i int, s *goquery.Selection) {
			if strings.Contains(s.Text(), "上頁") {
				href, exist = s.Attr("href")
			}
		})

		doc.Find(".r-ent").Each(func(i int, s *goquery.Selection) {
			article := Article{}
			article.Title = strings.TrimSpace(s.Find(".title").Text())
			article.LikeCount, _ = strconv.Atoi(s.Find(".nrec span").Text())
			hrefs, _ := s.Find(".title a").Attr("href")
			article.Link = BasePttAddress + hrefs
			article.Date = strings.TrimSpace(s.Find(".meta").Find(".date").Text())
			article.ImageLink = "https://i.imgur.com/aQjMlmV.jpg" // 先塞入預設imageLink
			article.LikeCountString = s.Find(".nrec span").Text()
			if article.Date != time.Format("1/02") {
				if crawlerCount > 0 {
					BOOL = false // 爬不到今日文章後 停止爬蟲
				}
			}

			// 今日文章且未被刪除（被刪除文章url會變成BasePttAddress)
			// 若文章內含有https及.jpg 的字串, 儲存為article.ImageLink.
			if article.Date == time.Format("1/02") && article.Link != BasePttAddress {
				//search image link in article
				doc, err := goquery.NewDocument(article.Link)
				if err != nil {
					log.Fatal(err)
				}

				doc.Find("#main-content > a").EachWithBreak(func(i int, s *goquery.Selection) bool {
					imgLink := s.Text()
					if strings.Contains(imgLink, ".jpg") {
						if strings.Contains(imgLink, "https") {
							article.ImageLink = imgLink
							return false
						}
					}
					return true
				})
				log.Println(article.Date + fourm + " 版 - " + " 標題: (" + article.LikeCountString + ")" + article.Title)
				db.InsertArticle(article.Title, article.LikeCount, article.Link, article.Date, article.ImageLink, article.LikeCountString)
			}
		})
		crawlerCount = crawlerCount + 1
	}

}
