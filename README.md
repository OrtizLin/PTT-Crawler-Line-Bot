# Line-bot (PTT beauty)

![螢幕快照](https://i.imgur.com/r4XiMa0.png)

Add Line friend !

Line-bot (PTT beauty) , user can search today's article by using keyword , or just tap the [正妹] , [神人] button to search.

If there are any articles get more 99 likes , bot will push message to you , you will not miss any hottie.



# Restrictions

In Line message API , the maximum of carousel message is 10 , so the search result can't over it.

```
for iter.Next(&result) {
			if index == 10 { //array of columns, max:10
				break
			}
```



And , LINE message - free plan , only allow 1000 message per month . So , I use LINE Notify (free service) to push article to user who has been subscribed. 

(P.S. the message will shows up at official LINE Notify account)

![螢幕快照 2018-03-19 下午3.59.26](https://i.imgur.com/l3Cdj6B.png)



# Others

Tutorial video：https://youtu.be/C9F6JESudyI

Any question are welcome.

otis@openalign.com