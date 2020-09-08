package menu

import (
	"github.com/fastwego/offiaccount"
	"github.com/fastwego/offiaccount/apis/menu"
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
	action := c.Request.URL.Query().Get("action")
	switch action {
	case "/menu/create":
		payload := []byte(`
		{
		  "button": [
			{
			  "name": "发图",
			  "sub_button": [
				{
				  "type": "pic_sysphoto",
				  "name": "系统拍照发图",
				  "key": "rselfmenu_1_0",
				  "sub_button": []
				},
				{
				  "type": "pic_photo_or_album",
				  "name": "拍照或者相册发图",
				  "key": "rselfmenu_1_1",
				  "sub_button": []
				},
				{
				  "type": "pic_weixin",
				  "name": "微信相册发图",
				  "key": "rselfmenu_1_2",
				  "sub_button": []
				}
			  ]
			},
			{
			  "name": "发送位置",
			  "type": "location_select",
			  "key": "rselfmenu_2_0"
			}
		  ]
		}`)
		resp, err := menu.Create(App, payload)
		if err != nil {
			c.Writer.WriteString(err.Error())
			return
		}
		c.Writer.WriteString(string(resp))
	case "/menu/get":
		resp, err := menu.Get(App)
		if err != nil {
			c.Writer.WriteString(err.Error())
			return
		}
		c.Writer.WriteString(string(resp))
	default:
		listen := viper.GetString("LISTEN")
		c.Writer.WriteString(action + " eg: //" + listen + "/api/weixin/menu?action=/menu/create")
	}
}
