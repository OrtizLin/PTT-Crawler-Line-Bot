package linebot

import (
	"github.com/line/line-bot-sdk-go/linebot"
	"line_bot_final/db"
	"line_bot_final/linenotify"
	"log"
	"net/http"
)

type LineBotStruct struct {
	bot         *linebot.Client
	appBaseURL  string
	downloadDir string
}

func NewLineBot(channelSecret, channelToken, appBaseURL string) (*LineBotStruct, error) {
	bots, err := linebot.New(
		channelSecret,
		channelToken,
	)
	if err != nil {
		return nil, err
	}
	return &LineBotStruct{
		bot:         bots,
		appBaseURL:  appBaseURL,
		downloadDir: "testing",
	}, nil
}

func (app *LineBotStruct) Callback(w http.ResponseWriter, r *http.Request) {
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
		case linebot.EventTypeFollow:
			if err := app.handleFollow(event.ReplyToken, event.Source); err != nil {
				log.Print(err)
			}
		default:
			log.Printf("Unknown event: %v", event)
		}
	}
}
func (app *LineBotStruct) handleText(message *linebot.TextMessage, replyToken string, source *linebot.EventSource) error {
	switch message.Text {
	case "tonygrr":
		if _, err := app.bot.ReplyMessage(
			replyToken,
			linebot.NewTextMessage("https://18comic.org/"),
		).Do(); err != nil {
			return err
		}
	default:

		result := db.SearchArticle(message.Text)

		if len(result) == 0 {
			log.Printf("Echo message to %s: %s", replyToken, message.Text)
			if _, err := app.bot.ReplyMessage(
				replyToken,
				linebot.NewTextMessage("抱歉！目前無相關文章"),
			).Do(); err != nil {
				return err
			}
		} else {
			count := 0
			var columns []*linebot.CarouselColumn
			for i := 0; i < len(result); i++ {
				if count == 9 {
					break
				}
				thumbnailImageUrl := result[i].ImageLink
				column := linebot.NewCarouselColumn(
					thumbnailImageUrl, result[i].Date, result[i].Title,
					linebot.NewURITemplateAction("點我查看更多", result[i].Link),
				)
				columns = append(columns, column)
				count++
			}

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

func (app *LineBotStruct) handleFollow(replyToken string, source *linebot.EventSource) error {
	//GET USER PROFILE
	profile, err := app.bot.GetProfile(source.UserID).Do()
	if err != nil {
		log.Print(err)
	}
	//send notify to me when someone follow this robot.
	linenotify.SomeOneFollow(profile.DisplayName, profile.PictureURL)

	return nil
}
