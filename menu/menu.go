package menu

import (
	"github.com/fastwego/offiaccount"
	"github.com/fastwego/offiaccount/apis/menu"
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

func ApiDemo(c *gin.Context) {
	action := c.Request.URL.Query().Get("action")
	switch action {
	case "/menu/create":
		payload := []byte(`
		{
			 "button":[
			 {
				   "name":"菜单",
				   "sub_button":[
				   {	
					   "type":"view",
					   "name":"搜索",
					   "url":"http://www.soso.com/"
					}]
			   }]
		}`)
		resp, err := menu.Create(OffiAccounts["account"], payload)
		if err != nil {
			c.Writer.WriteString(err.Error())
			return
		}
		c.Writer.WriteString(string(resp))
	case "/menu/get":
		resp, err := menu.Get(OffiAccounts["account"])
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
