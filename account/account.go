//package account 账号管理
package account

import (
	"github.com/fastwego/offiaccount"
	"github.com/fastwego/offiaccount/apis/account"
	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
)

var OffiAccounts = map[string]*offiaccount.OffiAccount{}

func init() {
	viper.SetConfigFile(".env")
	_ = viper.ReadInConfig()

	config := offiaccount.OffiAccountConfig{
		Appid:  viper.GetString("APPID"),
		Secret: viper.GetString("SECRET"),
	}
	OffiAccounts["account"] = offiaccount.New(config)
}

func Account(c *gin.Context) {
	action := c.Request.URL.Query().Get("action")
	switch action {
	case "created/qrcode":
		payload := []byte(`{
				"expire_seconds": 604800, 
				"action_name": "QR_SCENE", 
				"action_info": {
					"scene": {"scene_id": 123
				}
			}
		}`)
		resp, err := account.CreateQRCode(OffiAccounts["account"], payload)
		if err != nil {
			c.Writer.WriteString(err.Error())
			return
		}
		c.Writer.WriteString(string(resp))
	case "create/qrcode/limit_scene":
		payload := []byte(`{
			"action_name":"QR_LIMIT_SCENE",
			"action_info":{
				"scene":{
					"scene_id":123
				}
			}
		}`)
		resp, err := account.CreateQRCode(OffiAccounts["account"], payload)
		if err != nil {
			c.Writer.WriteString(err.Error())
			return
		}
		c.Writer.WriteString(string(resp))
	case "short/url":
		payload := []byte(`{
			"action":"long2short",
			"long_url":"https://github.com/fastwego/offiaccount"
		}`)
		resp, err := account.ShortUrl(OffiAccounts["account"], payload)
		if err != nil {
			c.Writer.WriteString(err.Error())
			return
		}
		c.Writer.WriteString(string(resp))
	default:
		listen := viper.GetString("LISTEN")
		c.Writer.WriteString(action + " eg: //" + listen + "/api/weixin/account?action=created/qrcode/scene")

	}
}
