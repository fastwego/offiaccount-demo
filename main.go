// Copyright 2020 FastWeGo
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//      http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package main

import (
	"context"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"os/signal"
	"path"
	"strconv"
	"syscall"
	"time"

	"github.com/fastwego/offiaccount-demo/material"

	"github.com/fastwego/offiaccount-demo/user"

	account "github.com/fastwego/offiaccount-demo/account"
	menu "github.com/fastwego/offiaccount-demo/menu"
	"github.com/fastwego/offiaccount-demo/oauth"

	"github.com/fastwego/offiaccount/type/type_message"

	"github.com/fastwego/offiaccount"
	"github.com/spf13/viper"

	"github.com/gin-gonic/gin"
)

// 微信公众账号池
var OffiAccounts = map[string]*offiaccount.OffiAccount{}

func init() {
	// 加载配置文件
	viper.SetConfigFile(".env")
	_ = viper.ReadInConfig()

	OffiAccounts["account1"] = offiaccount.New(offiaccount.Config{
		Appid:          viper.GetString("APPID"),
		Secret:         viper.GetString("SECRET"),
		Token:          viper.GetString("TOKEN"),
		EncodingAESKey: viper.GetString("EncodingAESKey"),
	})

	OffiAccounts["account2"] = offiaccount.New(offiaccount.Config{
		Appid:          viper.GetString("APPID2"),
		Secret:         viper.GetString("SECRET2"),
		Token:          viper.GetString("TOKEN2"),
		EncodingAESKey: viper.GetString("EncodingAESKey2"),
	})
}

func HandleMessage(c *gin.Context) {
	// 区分不同账号
	account := path.Base(c.Request.URL.Path)

	// 调用相应公众号服务
	body, _ := ioutil.ReadAll(c.Request.Body)
	log.Println(string(body))

	message, err := OffiAccounts[account].Server.ParseXML(body)
	if err != nil {
		log.Println(err)
	}

	var output interface{}
	switch message.(type) {
	case type_message.MessageText: // 文本 消息
		msg := message.(type_message.MessageText)

		// 回复文本消息
		output = type_message.ReplyMessageText{
			ReplyMessage: type_message.ReplyMessage{
				ToUserName:   type_message.CDATA(msg.FromUserName),
				FromUserName: type_message.CDATA(msg.ToUserName),
				CreateTime:   strconv.FormatInt(time.Now().Unix(), 10),
				MsgType:      type_message.ReplyMsgTypeText,
			},
			Content: type_message.CDATA(msg.Content),
		}
	}

	OffiAccounts[account].Server.Response(c.Writer, c.Request, output)
}

func main() {

	router := gin.New()
	router.Use(gin.Logger(), gin.Recovery())

	// 账号 1 服务
	router.GET("/api/weixin/account1", func(c *gin.Context) {
		OffiAccounts["account1"].Server.EchoStr(c.Writer, c.Request)
	})
	router.POST("/api/weixin/account1", HandleMessage)

	// 账号 2 服务
	router.GET("/api/weixin/account2", func(c *gin.Context) {
		OffiAccounts["account2"].Server.EchoStr(c.Writer, c.Request)
	})
	router.POST("/api/weixin/account2", HandleMessage)

	// 接口演示
	router.GET("/api/weixin/menu", menu.ApiDemo)
	router.GET("/api/weixin/account", account.ApiDemo)
	router.GET("/api/weixin/oauth", oauth.ApiDemo)
	router.GET("/api/weixin/user", user.ApiDemo)
	router.GET("/api/weixin/material", material.ApiDemo)

	svr := &http.Server{
		Addr:    viper.GetString("LISTEN"),
		Handler: router,
	}

	go func() {
		err := svr.ListenAndServe()
		if err != nil && err != http.ErrServerClosed {
			log.Fatalln(err)
		}
	}()

	quit := make(chan os.Signal)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	timeout := time.Duration(5) * time.Second
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	if err := svr.Shutdown(ctx); err != nil {
		log.Fatalln(err)
	}
}
