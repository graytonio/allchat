package main

import (
	"bytes"
	"context"
	"html/template"
	"net/http"
	"time"

	"github.com/fasthttp/websocket"
	"github.com/gin-gonic/gin"
	"github.com/graytonio/allchat/lib/chatroom"
	"github.com/graytonio/allchat/lib/twitch"
	"github.com/sirupsen/logrus"
)

var upgrader = websocket.Upgrader{}
var componentTemplates = template.Must(template.ParseGlob("templates/components/*.html"))

func handleChatWebsocket(ctx context.Context, w http.ResponseWriter, r *http.Request, chatChannel chan *chatroom.ChatMessage) {
	c, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		logrus.WithError(err).Error("could not upgrade connection")
		return
	}
	defer c.Close()
	for {
		select {
		case <- ctx.Done():
			return
		case msg := <- chatChannel:
			var buf bytes.Buffer
			if err := componentTemplates.ExecuteTemplate(&buf, "chat-message.html", msg); err != nil {
				logrus.WithError(err).Error("could not execute chat template")
				return
			}

			err = c.WriteMessage(websocket.TextMessage, buf.Bytes())
			if err != nil {
				logrus.WithError(err).Error("could not send websocket message")
				return
			}
		}
	}
}

func sendRandomMessages(chatChannel chan *chatroom.ChatMessage) {
	for {
		chatChannel <- &chatroom.ChatMessage{
			Username: "graytonio",
			Message: "Sup",
			Source: "GO",
		}
		time.Sleep(time.Second)
	}
}

func main() {

	chatChannel := make(chan *chatroom.ChatMessage)

	r := gin.Default()
	r.LoadHTMLGlob("templates/pages/*")
	
	r.Static("/assets", "./assets/dist")

	r.GET("/", func (c *gin.Context)  {
		c.HTML(http.StatusOK, "index.html", nil)
	})

	r.GET("/connect", func(c *gin.Context) {
		handleChatWebsocket(context.TODO(), c.Writer, c.Request, chatChannel)
	})

	go twitch.ConnectToChat("graytonio", chatChannel)

	r.Run()
}