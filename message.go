package main

import (
	"encoding/xml"
	"net/http"
	"strconv"
	"time"

	"github.com/fastwego/offiaccount/server"
)

func HandleTextMessage(writer http.ResponseWriter, msg server.MessageText) {
	replyMsg := server.ReplyMessageText{
		ReplyMessage: server.ReplyMessage{
			ToUserName:   server.CDATA(msg.FromUserName),
			FromUserName: server.CDATA(msg.ToUserName),
			CreateTime:   strconv.FormatInt(time.Now().Unix(), 10),
			MsgType:      server.ReplyMsgTypeText,
		},
		Content: server.CDATA(msg.Content),
	}

	data, err := xml.Marshal(replyMsg)
	if err != nil {
		return
	}
	_, err = writer.Write(data)
	if err != nil {
		return
	}
}

func HandleImageMessage(writer http.ResponseWriter, msg server.MessageImage) {
	replyMsg := server.ReplyMessageImage{
		ReplyMessage: server.ReplyMessage{
			ToUserName:   server.CDATA(msg.FromUserName),
			FromUserName: server.CDATA(msg.ToUserName),
			CreateTime:   strconv.FormatInt(time.Now().Unix(), 10),
			MsgType:      server.ReplyMsgTypeImage,
		},
		Image: struct{ MediaId server.CDATA }{
			MediaId: server.CDATA(msg.MediaId),
		},
	}

	data, err := xml.Marshal(replyMsg)
	if err != nil {
		return
	}
	_, err = writer.Write(data)
	if err != nil {
		return
	}
}
