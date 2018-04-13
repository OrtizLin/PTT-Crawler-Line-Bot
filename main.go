package main

import (
	"line_bot_final/crawler"
	"line_bot_final/linebot"
	"line_bot_final/linenotify"
	"log"
	"net/http"
	"os"
)

func main() {
	app, err := linebot.NewLineBot(
		os.Getenv("ChannelSecret"),
		os.Getenv("ChannelAccessToken"),
		os.Getenv("APP_BASE_URL"),
	)
	if err != nil {
		log.Fatal(err)
	}

	http.HandleFunc("/callback", app.Callback)
	http.HandleFunc("/auth", linenotify.Auth)
	http.HandleFunc("/pushnotify", linenotify.Token)
	http.HandleFunc("/wakeup", crawler.Start)
	if err := http.ListenAndServe(":"+os.Getenv("PORT"), nil); err != nil {
		log.Fatal(err)
	}
}
