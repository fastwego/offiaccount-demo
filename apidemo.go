package main

import (
	"github.com/fastwego/offiaccount/apis/menu"
	"github.com/gin-gonic/gin"
)

func ApiDemo(c *gin.Context) {
	api := c.Request.URL.Query().Get("api")
	switch api {
	case "/menu/create":
		payload := []byte(`
		{
			 "button":[
			 {	
				  "type":"click",
				  "name":"今日歌曲",
				  "key":"V1001_TODAY_MUSIC"
			  },
			  {
				   "name":"菜单",
				   "sub_button":[
				   {	
					   "type":"view",
					   "name":"搜索",
					   "url":"http://www.soso.com/"
					}
					{
					   "type":"click",
					   "name":"赞一下~",
					   "key":"V1001_GOOD"
					}]
			   }]
		}`)
		resp, err := menu.Create(OffiAccounts["account2"], payload)
		if err != nil {
			c.Writer.WriteString(err.Error())
			return
		}
		c.Writer.WriteString(string(resp))
	case "/menu/get":
		resp, err := menu.Get(OffiAccounts["account2"])
		if err != nil {
			c.Writer.WriteString(err.Error())
			return
		}
		c.Writer.WriteString(string(resp))
	default:
		c.Writer.WriteString(api + " eg: http://localhost/api/weixin/demo?api=/menu/create")
	}
}
