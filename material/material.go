// Copyright 2020 FastWeGo
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//      http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

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
	// 新增临时素材
	params := url.Values{}
	params.Add("type", "image")
	resp, err := material.MediaUpload(App, "./.data/img.jpg", params)
	fmt.Println(string(resp), err)

	c.Writer.WriteString(string(resp))

	// 新增永久素材
	params = url.Values{}
	params.Add("type", "video")
	fields := map[string]string{
		"description": `{"title":"Hii","introduction":"Hi"}`,
	}
	resp, err = material.AddMaterial(App, "./.data/1.mp4", params, fields)
	fmt.Println(string(resp), err)

	c.Writer.WriteString(string(resp))

	// 上传图文消息内的图片
	resp, err = material.MediaUploadImg(App, "./.data/img.jpg")
	fmt.Println(string(resp), err)

	c.Writer.WriteString(string(resp))

}
