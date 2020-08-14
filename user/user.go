package user

import (
	"net/url"

	"github.com/fastwego/offiaccount"
	"github.com/fastwego/offiaccount/apis/user"
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

	App.SetLogger(nil)
}

func ApiDemo(c *gin.Context) {
	action := c.Request.URL.Query().Get("action")
	switch action {
	case "/user/info":

		params := url.Values{}
		params.Add("openid", "o8jDwjrgxfOcQZ2_7V_iy_ZSIcok")
		params.Add("lang", "zh_CN")

		resp, err := user.GetUserInfo(App, params)
		if err != nil {
			c.Writer.WriteString(err.Error())
			return
		}
		c.Writer.WriteString(string(resp))
	default:
		listen := viper.GetString("LISTEN")
		c.Writer.WriteString(action + " eg: //" + listen + "/api/weixin/user?action=/user/info")
	}
}

//获取帐号的关注者列表,第一页
func GetUserList(c *gin.Context) {
	params := url.Values{}
	params.Add("next_openid", "")
	resp, err := user.Get(App, params)
	if err != nil {
		c.Writer.WriteString(err.Error())
		return
	}
	c.Writer.Write(resp)
	return
}
