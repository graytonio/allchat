package twitch

import (
	"fmt"
	"log"
	"net/url"
	"strings"

	"github.com/fasthttp/websocket"
	"github.com/graytonio/allchat/lib/chatroom"
	"github.com/sirupsen/logrus"
)

func ConnectToChat(channelName string, chatChannel chan *chatroom.ChatMessage) {
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

	err = c.WriteMessage(websocket.TextMessage, []byte(fmt.Sprintf("JOIN #%s", channelName)))
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

		log.Printf("recv: %s", message)
	}
}

func parseTwitchMessage(raw []byte) *chatroom.ChatMessage {
	if !strings.Contains(string(raw), "PRIVMSG") {
		return nil
	}
	
	user := strings.Split(string(raw), "!")[0]

	content := strings.Split(string(raw), ":")[2]

	// inUser := false
	// inCMD := false
	// inMessage := true
	// for _, c := range raw {
		
	// }

	return &chatroom.ChatMessage{
		Username: user,
		Message: content,
		Source: "twitch",
	}
}