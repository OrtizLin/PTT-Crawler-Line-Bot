![release-badge](https://img.shields.io/badge/Release-Ver2.0.0-blue.svg)![line-bot](https://img.shields.io/badge/WebCrawler-orange.svg)![web-crawler](https://img.shields.io/badge/Golang-green.svg)![golang](https://img.shields.io/badge/Bot-Line-brown.svg)

# Line-bot (PTT Beauty, Sex)

![螢幕快照](https://i.imgur.com/r4XiMa0.png)

Add Line friend !

This is a Line-bot which user can search today's article in PTT Beauty and Sex forum. 

And if there are any articles get more then 99 likes (push) , bot will send notify to you , you will not miss any hottie.



# Restrictions

First , by using Line message API , the maximum of carousel message is 10 , so the search result can't be more then the restriction.

```
for iter.Next(&result) {
			if index == 10 { //array of columns, max:10
				break
			}
```



Second , LINE message - free plan , only allow 1000 message per month , So I use LINE Notify (free service) to push article to users whom has been subscribed. 

(P.S. the message will shows up at official LINE Notify account)

![螢幕快照 2018-03-19 下午3.59.26](https://i.imgur.com/l3Cdj6B.png)

# Others

Tutorial video：https://youtu.be/C9F6JESudyI

Any question are welcome.

chenghsien852@gmail.com