package db

import (
	"github.com/utahta/go-linenotify"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
	"log"
	"os"
	"strings"
)

type User struct {
	UserToken string
}

type Article struct {
	Title           string
	LikeCount       int
	Link            string
	Date            string
	ImageLink       string
	LikeCountString string
}

func getDB() *mgo.Database {
	session, err := mgo.Dial(os.Getenv("DBURL"))
	if err != nil {
		panic(err)
	}

	session.SetMode(mgo.Monotonic, true)
	db := session.DB("xtest")
	return db
}

func SaveToken(token string) bool {

	c := getDB().C("tokendb")
	user := User{}
	user.UserToken = token
	errs := c.Insert(&User{user.UserToken})
	if errs != nil {
		log.Fatal(errs)
		return false
	} else {
		connect := linenotify.New()
		connect.NotifyWithImageURL(user.UserToken, "恭喜您已與表特爆報連動 , 若表特版有精彩文章將會立即通知您。", "https://i.imgur.com/wIdGRrU.jpg", "https://i.imgur.com/wIdGRrU.jpg")
		return true
	}

}

func InsertArticle(title string, likeCount int, link string, date string, imageLink string, likeCountString string) {

	c := getDB().C("xtest")
	c2 := getDB().C("alreadysent")
	c3 := getDB().C("tokendb")
	errs := c.Insert(&Article{title, likeCount, link, date, imageLink, likeCountString})
	if errs != nil {
		log.Fatal(errs)
	} else {
		if likeCountString == "爆" && strings.Contains(title, "帥哥") == false && strings.Contains(title, "創作") == false {
			result := Article{}
			err := c2.Find(bson.M{"link": link}).One(&result) //check if article already send
			if err != nil {
				err3 := c2.Insert(&Article{title, likeCount, link, date, imageLink, likeCountString})
				if err3 != nil {
					log.Fatal(err3)
				}

				users := User{}
				iter := c3.Find(nil).Iter()
				for iter.Next(&users) {
					connect := linenotify.New()
					content := "\n" + link
					connect.NotifyWithImageURL(users.UserToken, content, imageLink, imageLink)
				}

			}

		}
	}
}

func SearchArticle(message string) (article []Article) {
	
	var articles []Article

	c := getDB().C("xtest")
	result := Article{}
	iter := c.Find(bson.M{"title": bson.M{"$regex": message}}).Iter()
	count := 0
	for iter.Next(&result) {
		if count == 10 {
			break
		}
		articles = append(articles, result)
		count++
	}
	return articles
}

func RemoveALL() {

	c := getDB().C("xtest")
	//Clean DB
	c.RemoveAll(nil)
}
