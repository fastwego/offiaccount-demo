package oauth

import (
	"encoding/json"
	"fmt"

	"github.com/fastwego/offiaccount"
	"github.com/fastwego/offiaccount/apis/oauth"
	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
)

var Ctx *offiaccount.OffiAccount

func init() {
	viper.SetConfigFile(".env")
	_ = viper.ReadInConfig()

	config := offiaccount.Config{
		Appid:  viper.GetString("APPID"),
		Secret: viper.GetString("SECRET"),
	}
	Ctx = offiaccount.New(config)
}

// Oauth 演示
func Oauth(c *gin.Context) {

	code := c.Request.URL.Query().Get("code")
	if code != "" {
		// code 换取 accessToken
		accessToken, err := oauth.GetAccessToken(Ctx.Config.Appid, Ctx.Config.Secret, code)
		if err != nil {
			c.Writer.WriteString(err.Error())
			return
		}

		// 拉取用户信息
		userInfo, err := oauth.GetUserInfo(accessToken.AccessToken, accessToken.Openid, oauth.LANG_zh_CN)
		if err != nil {
			c.Writer.WriteString(err.Error())
			return
		}

		info, err := json.Marshal(userInfo)
		if err != nil {
			c.Writer.WriteString(err.Error())
			return
		}

		c.Writer.Write(info)

		// 判断 accesstoken 是否有效
		isValid, err := oauth.Auth(accessToken.AccessToken, accessToken.Openid)
		fmt.Println(isValid, err)

		// 刷新 AccessToken
		oauthAccessToken, err := oauth.RefreshToken(Ctx.Config.Appid, accessToken.RefreshToken)
		if err != nil {
			c.Writer.WriteString(err.Error())
			return
		}
		fmt.Println(oauthAccessToken)

		return
	}

	// 获取授权跳转链接
	link := oauth.GetAuthorizeUrl(Ctx.Config.Appid, "http:/127.0.0.1/api/weixin/oauth", oauth.ScopeSnsapiUserinfo, "STATE")
	c.Header("Content-Type", "text/html; charset=utf-8")
	c.Writer.WriteString(fmt.Sprintf("在微信中访问:<br/> <a href='%s'>%s</a>", link, link))
}
