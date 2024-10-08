package main

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
	"sync"
)

var ifChannelsMapInit = false

var channelsMap sync.Map

func AddChannel(userEmail string, traceId string) {
	if !ifChannelsMapInit {
		channelsMap = sync.Map{}
		ifChannelsMapInit = true
	}
	newChannel := make(chan string)
	channelsMap.Store(userEmail+traceId, newChannel)
	log.Print("Build SSE connection for user = " + userEmail + ", trace id = " + traceId)
}

func BuildNotificationChannel(userEmail string, traceId string, c *gin.Context) {
	AddChannel(userEmail, traceId)
	c.Writer.Header().Set("Content-Type", "text/event-stream")
	c.Writer.Header().Set("Cache-Control", "no-cache")
	c.Writer.Header().Set("Connection", "keep-alive")

	// 获取http写入器并断言为flusher，让其将缓冲器的数据立即写入
	w := c.Writer
	flusher, _ := w.(http.Flusher)

	// 监听客户端通道是否被关闭
	closeNotify := c.Request.Context().Done()

	go func() {
		<-closeNotify
		channelsMap.Delete(userEmail + traceId)
		log.Print("SSE close for user = " + userEmail + ", trace id = " + traceId)
		return
	}()

	curChan, _ := channelsMap.Load(userEmail + traceId)
	for msg := range curChan.(chan string) {
		fmt.Fprintf(w, "data:%s\n\n", msg)
		flusher.Flush()
	}
}
