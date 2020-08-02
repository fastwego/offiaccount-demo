package main

import (
	"context"
	"encoding/xml"
	menu "github.com/fastwego/offiaccount-demo/menu"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"os/signal"
	"path"
	"strconv"
	"syscall"
	"time"

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

	OffiAccounts["account1"] = offiaccount.New(offiaccount.OffiAccountConfig{
		Appid:          viper.GetString("APPID"),
		Secret:         viper.GetString("SECRET"),
		Token:          viper.GetString("TOKEN"),
		EncodingAESKey: viper.GetString("EncodingAESKey"),
	})

	OffiAccounts["account2"] = offiaccount.New(offiaccount.OffiAccountConfig{
		Appid:          viper.GetString("APPID2"),
		Secret:         viper.GetString("SECRET2"),
		Token:          viper.GetString("TOKEN2"),
		EncodingAESKey: viper.GetString("EncodingAESKey2"),
	})
}

func EchoStr(c *gin.Context) {
	// 区分不同账号
	account := path.Base(c.Request.URL.Path)

	// 调用相应公众号服务
	OffiAccounts[account].Server.EchoStr(c.Writer, c.Request)
}

func HandleMessage(c *gin.Context) {
	// 区分不同账号
	account := path.Base(c.Request.URL.Path)

	// 调用相应公众号服务
	body, _ := ioutil.ReadAll(c.Request.Body)
	log.Println(string(body))

	message, err := OffiAccounts[account].Server.ParseMessage(body)
	if err != nil {
		log.Println(err)
	}

	switch message.(type) {
	case type_message.MessageText: // 文本 消息
		msg := message.(type_message.MessageText)
		replyMsg := type_message.ReplyMessageText{
			ReplyMessage: type_message.ReplyMessage{
				ToUserName:   type_message.CDATA(msg.FromUserName),
				FromUserName: type_message.CDATA(msg.ToUserName),
				CreateTime:   strconv.FormatInt(time.Now().Unix(), 10),
				MsgType:      type_message.ReplyMsgTypeText,
			},
			Content: type_message.CDATA(msg.Content),
		}

		data, err := xml.Marshal(replyMsg)
		if err != nil {
			return
		}
		OffiAccounts[account].Server.Response(c.Writer, c.Request, data)
	}

	if !c.Writer.Written() {
		log.Println("default")
		// 兜底响应 success 告知微信服务器
		c.Writer.WriteString(offiaccount.SUCCESS)
	}
}

func main() {

	router := gin.New()
	router.Use(gin.Logger(), gin.Recovery())

	// 账号 1 服务
	router.GET("/api/weixin/account1", EchoStr)
	router.POST("/api/weixin/account1", HandleMessage)

	// 账号 2 服务
	router.GET("/api/weixin/account2", EchoStr)
	router.POST("/api/weixin/account2", HandleMessage)

	// 接口演示
	router.GET("/api/weixin/menu", menu.ApiDemo)

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
