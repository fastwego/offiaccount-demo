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
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"os/signal"
	"strconv"
	"syscall"
	"time"

	"github.com/fastwego/dingding"
	"github.com/fastwego/dingding/apis/ai"

	"github.com/fastwego/offiaccount/type/type_message"

	"github.com/fastwego/offiaccount"
	"github.com/spf13/viper"

	"github.com/gin-gonic/gin"
)

// 微信公众账号
var OffiAccount *offiaccount.OffiAccount

// 钉钉 App 实例
var DingApp *dingding.App

func init() {
	// 加载配置文件
	viper.SetConfigFile(".env")
	_ = viper.ReadInConfig()

	// 创建公众号实例
	OffiAccount = offiaccount.New(offiaccount.Config{
		Appid:          viper.GetString("APPID"),
		Secret:         viper.GetString("SECRET"),
		Token:          viper.GetString("TOKEN"),
		EncodingAESKey: viper.GetString("EncodingAESKey"),
	})

	// 创建钉钉应用实例
	DingApp = dingding.NewApp(dingding.AppConfig{
		AppKey:    viper.GetString("AppKey"),
		AppSecret: viper.GetString("AppSecret"),
	})
}

func HandleMessage(c *gin.Context) {

	// 获取公众号消息
	body, _ := ioutil.ReadAll(c.Request.Body)
	log.Println(string(body))

	message, err := OffiAccount.Server.ParseXML(body)
	if err != nil {
		log.Println(err)
	}

	var output interface{}
	switch message.(type) {
	case type_message.MessageText: // 文本 消息
		msg := message.(type_message.MessageText)

		// 调用 钉钉 翻译服务
		params := struct {
			Query          string `json:"query"`
			TargetLanguage string `json:"target_language"`
			SourceLanguage string `json:"source_language"`
		}{}

		params.Query = msg.Content
		params.SourceLanguage = "zh"
		params.TargetLanguage = "fr"

		data, err := json.Marshal(params)
		if err != nil {
			fmt.Println(string(data), err)
			return
		}

		// 翻译接口
		resp, err := ai.Translate(DingApp, data)
		fmt.Println(string(resp), err)

		if err != nil {
			return
		}

		// 翻译结果
		result := struct {
			Errcode int64  `json:"errcode"`
			Errmsg  string `json:"errmsg"`
			Result  string `json:"result"`
		}{}
		err = json.Unmarshal(resp, &result)
		fmt.Println(result, err)
		if err != nil {
			return
		}

		// 回复公众号 翻译结果文本消息
		output = type_message.ReplyMessageText{
			ReplyMessage: type_message.ReplyMessage{
				ToUserName:   type_message.CDATA(msg.FromUserName),
				FromUserName: type_message.CDATA(msg.ToUserName),
				CreateTime:   strconv.FormatInt(time.Now().Unix(), 10),
				MsgType:      type_message.ReplyMsgTypeText,
			},
			Content: type_message.CDATA(result.Result),
		}
	}

	OffiAccount.Server.Response(c.Writer, c.Request, output)
}

func main() {

	router := gin.New()
	router.Use(gin.Logger(), gin.Recovery())

	// 公众号 服务 url 校验
	router.GET("/api/weixin/dingding", func(c *gin.Context) {
		OffiAccount.Server.EchoStr(c.Writer, c.Request)
	})

	// 公众号 服务
	router.POST("/api/weixin/dingding", HandleMessage)

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
