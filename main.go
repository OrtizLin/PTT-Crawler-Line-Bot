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
	case "carousel":
		log.Printf("Echo message to %s: %s", replyToken, message.Text)
		if _, err := app.bot.ReplyMessage(
			replyToken,
			linebot.NewTextMessage(message.Text),
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
		for iter.Next(&result) {
			colum := linbot.NewCarouselColumn(
				thumbnailImageUrl, result.Title, "3/12",
				linebot.NewURITemplateAction("這裡放URL", result.Link),
			)
			columns = append(columns, column)
		}
		template := linebot.NewCarouselTemplate(columns...)
		// thumbnailImageUrl := "https://www.atanews.net/upload_edit/images/201605/20160524174909_79517e8e.jpg"
		// template := linebot.NewCarouselTemplate(
		// 	linebot.NewCarouselColumn(
		// 		thumbnailImageUrl, "[正妹]超星拳婦", "3/12",
		// 		linebot.NewURITemplateAction("這裡放URL", "https://tw.yahoo.com/"),
		// 	),
		// 	linebot.NewCarouselColumn(
		// 		thumbnailImageUrl, "[正妹]冰冰", "3/12",
		// 		linebot.NewURITemplateAction("這裡放URL", "https://tw.yahoo.com/"),
		// 	),
		// )
		if _, err := app.bot.ReplyMessage(
			replyToken,
			linebot.NewTemplateMessage("Carousel alt text", template),
		).Do(); err != nil {
			return err
		}
	}
	return nil

}
