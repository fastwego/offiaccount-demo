package material

import (
	"fmt"
	"net/url"

	"github.com/fastwego/offiaccount"
	"github.com/fastwego/offiaccount/apis/material"
	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
)

var App *offiaccount.OffiAccount

func init() {
	viper.SetConfigFile(".env")
	_ = viper.ReadInConfig()

	config := offiaccount.Config{
		Appid:  viper.GetString("APPID"),
		Secret: viper.GetString("SECRET"),
	}
	App = offiaccount.New(config)
}

func ApiDemo(c *gin.Context) {

	// 上传视频
	params := url.Values{}
	params.Add("type", "video")
	fields := map[string]string{
		"description": `{"title":"Hii","introduction":"Hi"}`,
	}
	resp, err := material.AddMaterial(App, "./.data/1.mp4", params, fields)
	fmt.Println(string(resp), err)

	c.Writer.WriteString(string(resp))

	// 上传图文消息内的图片
	resp, err = material.MediaUploadImg(App, "./.data/img.jpg")
	fmt.Println(string(resp), err)

	c.Writer.WriteString(string(resp))

}
