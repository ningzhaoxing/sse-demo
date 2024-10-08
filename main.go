package main

import "github.com/gin-gonic/gin"

func main() {
	r := gin.Default()
	r.Use(CORSMiddleware())
	r.GET("/notification/socket-connection", SocketConnection) // 建立sse连接
	r.GET("/notification/export-excel", ExportExcel)           // 触发通知，发送消息

	r.Run("localhost:8080")
}

func SocketConnection(c *gin.Context) {
	in := struct {
		UserEmail string `json:"user_email"`
		TraceId   string `json:"trace_id"`
	}{}

	err := c.ShouldBindJSON(&in)
	if err != nil {
		return
	}

	BuildNotificationChannel(in.UserEmail, in.TraceId, c)
}

func ExportExcel(c *gin.Context) {
	in := struct {
		UserEmail   string `json:"user_email"`
		MessageBody string `json:"message_body"`
		ActionType  string `json:"action_type"`
	}{}

	err := c.ShouldBindJSON(&in)
	if err != nil {
		return
	}

	SendNotification(in.UserEmail, in.MessageBody, in.ActionType)
}
