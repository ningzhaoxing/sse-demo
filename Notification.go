package main

import (
	"encoding/json"
	"log"
	"strings"
	"time"
)

type NotificationLog struct {
	MessageBody string    `json:"message_body"`
	UserEmail   string    `json:"user_email"`
	Type        string    `json:"type"`
	Status      string    `json:"status"`
	CreatTime   time.Time `json:"creat_time"`
}

func SendNotification(userEmail string, messageBody string, actionType string) {
	log.Print("Send notification to user = " + userEmail)
	var msg = NotificationLog{
		MessageBody: messageBody,
		UserEmail:   userEmail,
		Type:        actionType,
		Status:      "UNREAD",
		CreatTime:   time.Now(),
	}
	msgBytes, _ := json.Marshal(msg)
	channelsMap.Range(func(key, value any) bool {
		k := key.(string)
		if strings.Contains(k, userEmail) {
			channel := value.(chan string)
			channel <- string(msgBytes)
		}
		return true
	})
}
