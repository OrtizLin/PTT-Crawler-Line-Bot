# 前言

![螢幕快照](https://i.imgur.com/r4XiMa0.png)

Line Bot 表特爆報 , 使用者可以利用關鍵字搜尋當日 PTT 表特版發文, 或使用圖文按鈕直接搜尋, 若當日有推文數>99, 亦即推爆的文章, 將會利用Line Notify 發送通知給有訂閱之使用者。



# 使用限制

Line message 裡的 template - Carousel message最大數量為10, 所以每次搜尋出來的結果最多只能10筆

```
for iter.Next(&result) {
			if index == 10 { //array of columns, max:10
				break
			}
```



免費版的Line message 每月最多只能推送1000則訊息（訊息x好友人數）, 所以當有推爆的文章需推送通知時, 使用Line Notify 推送給有訂閱的用戶. 但訊息會出現在官方的Line Notify 視窗.

![螢幕快照 2018-03-19 下午3.59.26](https://i.imgur.com/l3Cdj6B.png)



