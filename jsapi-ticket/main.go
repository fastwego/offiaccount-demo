package main

import (
	"context"
	"crypto/sha1"
	"encoding/json"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"os"
	"os/signal"
	"strconv"
	"syscall"
	"time"

	"github.com/fastwego/offiaccount/util"

	"github.com/fastwego/offiaccount/apis/oauth"

	"github.com/fastwego/offiaccount"

	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
)

// 微信公众账号
var OffiAccount *offiaccount.OffiAccount

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
}

func main() {
	router := gin.New()
	router.Use(gin.Logger(), gin.Recovery())

	router.GET("/", func(c *gin.Context) {

		config, err := jsapiConfig(c)
		if err != nil {
			fmt.Println(err)
			return
		}

		t1, err := template.ParseFiles("index.html")
		if err != nil {
			fmt.Println(err)
			return
		}

		t1.Execute(c.Writer, config)
	})

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

func jsapiConfig(c *gin.Context) (config template.JS, err error) {

	// 优先从环缓存获取
	jsapi_ticket, err := OffiAccount.AccessToken.Cache.Fetch("jsapi_ticket:" + OffiAccount.Config.Appid)
	if len(jsapi_ticket) == 0 {
		var ttl int64
		jsapi_ticket, ttl, err = oauth.GetJSApiTicket(OffiAccount)
		if err != nil {
			return
		}

		err = OffiAccount.AccessToken.Cache.Save("jsapi_ticket:"+OffiAccount.Config.Appid, jsapi_ticket, time.Duration(ttl)*time.Second)
		if err != nil {
			return
		}
	}

	nonceStr := util.GetRandString(6)
	timestamp := strconv.FormatInt(time.Now().Unix(), 10)
	pageUrl := "http://" + c.Request.Host + c.Request.RequestURI
	plain := "jsapi_ticket=" + jsapi_ticket + "&noncestr=" + nonceStr + "&timestamp=" + timestamp + "&url=" + pageUrl

	signature := fmt.Sprintf("%x", sha1.Sum([]byte(plain)))
	fmt.Println(plain, signature)

	configMap := map[string]string{
		"url":       pageUrl,
		"nonceStr":  nonceStr,
		"appid":     OffiAccount.Config.Appid,
		"timestamp": timestamp,
		"signature": signature,
	}

	marshal, err := json.Marshal(configMap)
	if err != nil {
		return
	}

	return template.JS(marshal), nil
}
