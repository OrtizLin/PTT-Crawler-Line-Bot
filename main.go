package main

import (
	"github.com/line/line-bot-sdk-go/linebot"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"

	"log"
	"net/http"
	"os"
)

type Article struct {
	Title     string
	LikeCount int
	Link      string
	Date      string
}

func main() {
	app, err := NewLineBot(
		os.Getenv("ChannelSecret"),
		os.Getenv("ChannelAccessToken"),
		os.Getenv("APP_BASE_URL"),
	)
	if err != nil {
		log.Fatal(err)
	}

	http.HandleFunc("/callback", app.Callback)
	if err := http.ListenAndServe(":"+os.Getenv("PORT"), nil); err != nil {
		log.Fatal(err)
	}
}

// line_bot app
type LineBot struct {
	bot         *linebot.Client
	appBaseURL  string
	downloadDir string
}

// NewLineBot function
func NewLineBot(channelSecret, channelToken, appBaseURL string) (*LineBot, error) {
	bot, err := linebot.New(
		channelSecret,
		channelToken,
	)
	if err != nil {
		return nil, err
	}

	return &LineBot{
		bot:         bot,
		appBaseURL:  appBaseURL,
		downloadDir: "test",
	}, nil
}
func (app *LineBot) Callback(w http.ResponseWriter, r *http.Request) {
	events, err := app.bot.ParseRequest(r)

	if err != nil {
		if err == linebot.ErrInvalidSignature {
			w.WriteHeader(400)
		} else {
			w.WriteHeader(500)
		}
		return
	}
	for _, event := range events {
		log.Printf("Got event %v", event)
		switch event.Type {
		case linebot.EventTypeMessage:
			switch message := event.Message.(type) {
			case *linebot.TextMessage:
				if err := app.handleText(message, event.ReplyToken, event.Source); err != nil {
					log.Print(err)
				}
			default:
				log.Printf("Unknown message: %v", message)
			}
		default:
			log.Printf("Unknown event: %v", event)
		}
	}
}

func (app *LineBot) handleText(message *linebot.TextMessage, replyToken string, source *linebot.EventSource) error {
	switch message.Text {
	case "tonygrr":
		log.Printf("Echo message to %s: %s", replyToken, message.Text)
		if _, err := app.bot.ReplyMessage(
			replyToken,
			linebot.NewTextMessage("http://www.jav777.cc/"),
		).Do(); err != nil {
			return err
		}
	default:

		session, errs := mgo.Dial(os.Getenv("DBURL"))
		if errs != nil {
			panic(errs)
		}
		defer session.Close()
		c := session.DB("xtest").C("xtest")
		result := Article{}
		var columns []*linebot.CarouselColumn
		iter := c.Find(bson.M{"title": bson.M{"$regex": message.Text}}).Iter()
		var index = 0
		for iter.Next(&result) {
			if index == 5 {
				break
			}
			thumbnailImageUrl := "https://www.atanews.net/upload_edit/images/201605/20160524174909_79517e8e.jpg"
			column := linebot.NewCarouselColumn(
				thumbnailImageUrl, result.Date, result.Title,
				linebot.NewURITemplateAction("點我查看更多", result.Link),
			)
			columns = append(columns, column)
			index++
		}
		//if serch result is null
		if index == 0 {
			log.Printf("Echo message to %s: %s", replyToken, message.Text)
			if _, err := app.bot.ReplyMessage(
				replyToken,
				linebot.NewTextMessage("抱歉！目前無相關文章"),
			).Do(); err != nil {
				return err
			}
		} else {
			//reply carousel message if search result exist
			template := linebot.NewCarouselTemplate(columns...)
			if _, err := app.bot.ReplyMessage(
				replyToken,
				linebot.NewTemplateMessage("正妹來囉！", template),
			).Do(); err != nil {
				return err
			}
		}
	}
	return nil

}
