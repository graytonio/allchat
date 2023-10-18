package main

import (
	"bytes"
	"context"
	"html/template"
	"net/http"

	"github.com/fasthttp/websocket"
	"github.com/gin-gonic/gin"
	"github.com/graytonio/allchat/lib/chatroom"
	"github.com/graytonio/allchat/lib/config"
	"github.com/graytonio/allchat/lib/twitch"
	"github.com/sirupsen/logrus"
	"golang.org/x/exp/slices"
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

func connectToChats(chatChannel chan *chatroom.ChatMessage) {
	if slices.Contains(config.GetConfig().EnabledChats, "twitch") {
		logrus.WithFields(logrus.Fields{
			"channel": config.GetConfig().Twitch.Channel,
		}).Info("Connecting to twitch chat rooom")
		go twitch.ConnectToChat(&config.GetConfig().Twitch, chatChannel)
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
	connectToChats(chatChannel)
	r.Run()
}