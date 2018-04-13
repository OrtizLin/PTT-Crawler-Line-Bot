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
	getAllArticles()
}

func getAllArticles() {

	var BOOL = true
	var exist = true
	var url string = ""  //default url
	var href string = "" //next page url
	crawlerCount := 0

	//today's date
	loc, _ := time.LoadLocation("Asia/Chongqing")
	time := time.Now().In(loc)

	for BOOL {
		if href == "" {
			url = BasePttAddress + "/bbs/Beauty/index.html"
		} else {
			url = BasePttAddress + href
		}

		doc, err := goquery.NewDocument(url)
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
			article.ImageLink = "https://i.imgur.com/aQjMlmV.jpg"
			article.LikeCountString = s.Find(".nrec span").Text()
			if article.Date != time.Format("1/02") {
				if crawlerCount > 0 {
					BOOL = false
				}
			}

			if article.Date == time.Format("1/02") && article.Link != BasePttAddress {
				//search image link in article
				doc, err := goquery.NewDocument(article.Link)
				if err != nil {
					log.Fatal(err)
				}

				//article-metaline
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
				db.InsertArticle(article.Title, article.LikeCount, article.Link, article.Date, article.ImageLink, article.LikeCountString)
			}
		})
		crawlerCount = crawlerCount + 1
	}
}
