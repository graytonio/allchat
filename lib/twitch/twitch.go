package twitch

import (
	"fmt"
	"net/url"
	"strings"

	"github.com/fasthttp/websocket"
	"github.com/graytonio/allchat/lib/chatroom"
	"github.com/graytonio/allchat/lib/config"
	"github.com/sirupsen/logrus"
)

func ConnectToChat(conf *config.TwitchConfig, chatChannel chan *chatroom.ChatMessage) {
	u := url.URL{Scheme: "ws", Host: "irc-ws.chat.twitch.tv:80"}
	c, _, err := websocket.DefaultDialer.Dial(u.String(), nil)
	if err != nil {
		logrus.WithError(err).Error("could not connect to twitch")
		return
	}
	defer c.Close()

	err = c.WriteMessage(websocket.TextMessage, []byte("NICK justinfan3456"))
	if err != nil {
		logrus.WithError(err).Error("could not join channel")
		return
	}

	err = c.WriteMessage(websocket.TextMessage, []byte(fmt.Sprintf("JOIN #%s", conf.Channel)))
	if err != nil {
		logrus.WithError(err).Error("could not join channel")
		return
	}

	for {
		_, message, err := c.ReadMessage()
		if err != nil {
			logrus.WithError(err).Error("could not receiving message")
			return
		}

		parsedMessage := parseTwitchMessage(message)
		if parsedMessage != nil {
			chatChannel <- parsedMessage
		}

		logrus.WithField("message", string(message)).Debugf("received socket message")
	}
}

func parseTwitchMessage(raw []byte) *chatroom.ChatMessage {
	if !strings.Contains(string(raw), "PRIVMSG") {
		return nil
	}
	
	user := strings.Split(string(raw), "!")[0]
	user = strings.TrimPrefix(user, ":")
	content := strings.Split(string(raw), ":")[2]

	return &chatroom.ChatMessage{
		Username: user,
		Message: content,
		Source: "twitch",
	}
}