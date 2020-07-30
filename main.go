package main

import (
	"context"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"

	"github.com/fastwego/offiaccount"
	"github.com/fastwego/offiaccount/server"
	"github.com/spf13/viper"

	"github.com/gin-gonic/gin"
)

func init() {
	// 加载配置文件
	viper.SetConfigFile(".env")
	_ = viper.ReadInConfig()

	// 初始化 公众号配置
	{
		offiaccount.Appid = viper.GetString("APPID")   // 公众号 appid
		offiaccount.Secret = viper.GetString("SECRET") // 公众号 secret

		server.Token = viper.GetString("TOKEN")                   // 接口校验 Token
		server.EncodingAESKey = viper.GetString("EncodingAESKey") // 消息 aes 加密 key

		log.Println(strings.Join([]string{offiaccount.Appid, offiaccount.Secret, server.Token, server.EncodingAESKey}, "/"))
		if offiaccount.Appid == "" || offiaccount.Secret == "" {
			panic("APPID/SECRET not found")
		}
	}
}

func EchoStr(c *gin.Context) {
	server.EchoStr(c.Writer, c.Request)
	if !c.Writer.Written() {
		// 兜底响应
		c.Writer.WriteString(server.SUCCESS)
	}
}

func HandleMessage(c *gin.Context) {

	body, _ := ioutil.ReadAll(c.Request.Body)
	log.Println(string(body))

	message, _ := server.ParseMessage(body)

	switch message.(type) {
	case server.MessageText: // 文本 消息
		HandleTextMessage(c.Writer, message.(server.MessageText))
	case server.MessageImage: // 图片 消息
		HandleImageMessage(c.Writer, message.(server.MessageImage))
	case server.MessageVoice:
		// TODO
	case server.MessageVideo:
		// TODO
	case server.MessageShortVideo:
		// TODO
	case server.MessageLink:
		// TODO
	case server.MessageLocation:
		// TODO
	case server.MessageFile:
		// TODO
	case server.MessageEvent: // 事件 处理
		event, _ := server.ParseEvent(body)
		switch event.(type) {
		case server.EventSubscribe: // 关注
			// TODO
		case server.EventUnsubscribe: // 取关
			// TODO
		case server.EventScan: // 已关注 扫码
			// TODO
		case server.EventLocation: // 位置
			// TODO
		case server.EventClick: // 点击菜单
			// TODO
		case server.EventView: // 点击菜单链接
			// TODO
		}
	}

	if !c.Writer.Written() {
		// 兜底响应 success 告知微信服务器
		c.Writer.WriteString(server.SUCCESS)
	}
}

func main() {

	router := gin.New()
	router.Use(gin.Logger(), gin.Recovery())

	// 服务器校验
	router.GET("/api/weixin", EchoStr)

	// 消息响应
	router.POST("/api/weixin", HandleMessage)

	// 接口演示
	router.GET("/api/weixin/demo", ApiDemo)

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
