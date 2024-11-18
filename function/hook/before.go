package hook

import (
	"chatbot/config"
	"chatbot/logger"
	"chatbot/msg"
	my_mongo "chatbot/storage/mongo"
	"context"
	"encoding/base64"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"io"
	"net/http"
	"sync"
	"time"
)

var (
	msgCol   *mongo.Collection
	imageCol *mongo.Collection
	mu       sync.Mutex
)

var OnlyWhiteList Hook = func(message *msg.ReceiveMessage) bool {
	whiteGroup := config.Get().Server.WhiteGroup
	if whiteGroup != nil && len(whiteGroup) > 0 {
		for _, group := range whiteGroup {
			if message.GroupId == group {
				return true
			}
		}
	}
	logger.Infof(context.Background(), "msg source %v is not on whitelist", message.GroupId)
	return false
}

var SaveMessage Hook = func(message *msg.ReceiveMessage) bool {
	if msgCol == nil && mu.TryLock() {
		defer mu.Unlock()
		msgCol = my_mongo.Get("default").Database("bot").Collection("msg")
	}

	go func() {
		ctx, cancel := context.WithTimeout(context.Background(), time.Second*1)
		defer cancel()
		msgCol.InsertOne(ctx, message)
	}()
	return true
}

var GetImage Hook = func(message *msg.ReceiveMessage) bool {
	if imageCol == nil && mu.TryLock() {
		defer mu.Unlock()
		imageCol = my_mongo.Get("default").Database("bot").Collection("images")
	}

	if len(message.Message) > 0 && message.Message[0].Type == msg.SubImageMsg {
		data := message.Message[0].Data
		name := data["file"].(string)
		url := data["url"].(string)
		if res, err := http.Get(url); err == nil {
			defer res.Body.Close()
			if res.StatusCode != 200 {
				return true
			}
			body, _ := io.ReadAll(res.Body)
			encoded := base64.StdEncoding.EncodeToString(body)

			go func() {
				ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
				defer cancel()
				imageCol.InsertOne(ctx, bson.M{
					"name":       name,
					"data":       encoded,
					"time":       time.Now().Unix(),
					"message_id": message.MessageId,
					"group_id":   message.GroupId,
					"user_id":    message.UserId,
				})
			}()
		}
	}
	return true
}
