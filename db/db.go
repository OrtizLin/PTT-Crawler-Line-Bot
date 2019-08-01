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

type HotBorads struct {
	Board []string
}

type Article struct {
	Title           string
	LikeCount       int
	Link            string
	Date            string
	ImageLink       string
	LikeCountString string
	Board 			string
}

func SaveToken(token string) bool {

	session, errs := mgo.Dial(os.Getenv("DBURL"))
	if errs != nil {
		panic(errs)
	}
	defer session.Close()
	collect := session.DB("xtest").C("tokendb")
	user := User{}
	user.UserToken = token
	errs = collect.Insert(&User{user.UserToken})
	if errs != nil {
		log.Fatal(errs)
		return false
	} else {
		connect := linenotify.New()
		connect.NotifyWithImageURL(user.UserToken, "恭喜您已與表特爆報連動 , 若表特版有精彩文章將會立即通知您。", "https://i.imgur.com/wIdGRrU.jpg", "https://i.imgur.com/wIdGRrU.jpg")
		return true
	}

}

func AllHotBoards() (hotboardlist []string) {
	var results []HotBorads
	var list []string
	session, errs := mgo.Dial(os.Getenv("DBURL"))
	if errs != nil {
		panic(errs)
	}
	defer session.Close()
	c := session.DB("xtest").C("hotboard")
	err := c.Find(nil).All(&results)
	if err != nil {
		panic(err)
	}
	for i := 0; i < len(results); i++ {
		list = results[i].Board
	}

	return list

}

func InsertHotBoard(boards []string) {

	session, errs := mgo.Dial(os.Getenv("DBURL"))
	if errs != nil {
		panic(errs)
	}
	defer session.Close()
	c := session.DB("xtest").C("hotboard")
	errs = c.Insert(&HotBorads{boards})
	if errs != nil {
		log.Fatal(errs)
	} 

}

func InsertArticle(title string, likeCount int, link string, date string, imageLink string, likeCountString string, board string) {
	session, errs := mgo.Dial(os.Getenv("DBURL"))
	if errs != nil {
		panic(errs)
	}
	defer session.Close()
	c := session.DB("xtest").C("xtest")
	c2 := session.DB("xtest").C("alreadysent")
	c3 := session.DB("xtest").C("tokendb")
	errs = c.Insert(&Article{title, likeCount, link, date, imageLink, likeCountString, board})
	if errs != nil {
		log.Fatal(errs)
	} else {
		if likeCountString == "爆" && strings.Contains(title, "帥哥") == false && strings.Contains(title, "創作") == false {
			result := Article{}
			err := c2.Find(bson.M{"link": link}).One(&result) //check if article already send
			if err != nil {
				err3 := c2.Insert(&Article{title, likeCount, link, date, imageLink, likeCountString, board})
				if err3 != nil {
					log.Fatal(err3)
				}

				users := User{}
				iter := c3.Find(nil).Iter()
				for iter.Next(&users) {
					connect := linenotify.New()
					content := board + "版 - " + title + "\n" + link
						connect.NotifyWithImageURL(users.UserToken, content, imageLink, imageLink)
				}

			}

		}
	}
}

func SearchArticle(message string) (article []Article) {
	var articles []Article
	session, errs := mgo.Dial(os.Getenv("DBURL"))
	if errs != nil {
		panic(errs)
	}
	defer session.Close()
	c := session.DB("xtest").C("xtest")
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

func RemoveALL(collect string) {
	session, errs := mgo.Dial(os.Getenv("DBURL"))
	if errs != nil {
		panic(errs)
	}
	defer session.Close()
	c := session.DB("xtest").C(collect)
	//Clean DB
	c.RemoveAll(nil)

}
