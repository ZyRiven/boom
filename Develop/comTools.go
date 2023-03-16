/**
 *
 * @company: Co.预见（天津）智能科技有限公司
 * @Author:  ZhaoYi
 * @Date:    2023/2/7 14:24
 */

package Develop

import (
	"bufio"
	"fmt"
	"gohello/app"
	"gohello/duang"
	"io"
	"net/http"
	"os"
	"strings"
)

func (cc Controller) List() {
	cc.group.GET("/list", func(c *duang.Context) {
		c.JSON(200, duang.H{
			"code": 200,
			"data": "11",
		})
	})
}

// Upload 上传图片
func (cc Controller) Upload() {
	cc.group.POST("/upload", func(c *duang.Context) {
		token := app.GetRequest(c).Header["Token"]
		_, err := app.GetUser(token)
		if err != "" {
			c.JSON(http.StatusOK, duang.H{
				"code": 400,
				"msg":  err,
			})
			return
		}

		// ParseMultipartForm将请求的主体作为multipart/form-data解析。
		// 请求的整个主体都会被解析，得到的文件记录最多maxMemery字节保存在内存，其余部分保存在硬盘的temp文件里
		//r.ParseMultipartForm(32 << 20)
		file, handler, er := c.Req.FormFile("file")
		if er != nil {
			c.JSON(http.StatusOK, duang.H{
				"code": 400,
				"msg":  "验证失败",
			})
			return
		}
		defer file.Close()
		boolFileList := duang.CheckFileIsExist("./public/img/")
		if boolFileList == false {
			e := os.Mkdir("./public/img/", os.ModePerm) //创建文件夹
			if e != nil {
				c.JSON(http.StatusOK, duang.H{
					"code": 400,
					"msg":  e,
				})
				return
			}
			fmt.Println("public/img/，创建成功")
		}
		f, e := os.OpenFile("./public/img/"+handler.Filename, os.O_WRONLY|os.O_CREATE, 0666)
		if e != nil {
			c.JSON(http.StatusOK, duang.H{
				"code": 400,
				"msg":  e,
			})
			return
		}
		defer f.Close()
		io.Copy(f, file)
		c.JSON(http.StatusOK, duang.H{
			"code":     200,
			"fileName": handler.Filename,
			"url":      "/public/img/" + handler.Filename,
			"curl":     c.Req.Host + "/web/handleFile/" + handler.Filename,
		})
		return
	})
}

// HandleFile 获取上传图片文件流
func (cc Controller) HandleFile() {
	cc.group.GET("/handleFile/*", func(c *duang.Context) {
		path := strings.Split(c.Path, "/")
		file, er := os.Open("./public/img/" + path[3])
		if er != nil {
			c.String(http.StatusNotFound, "404 NOT FOUND: %s\n", c.Path)
			return
		}
		defer file.Close()
		reader := bufio.NewReader(file)
		_, er = io.Copy(c.Writer, reader)
		if er != nil {
			c.String(http.StatusNotFound, "404 NOT FOUND: %s\n", c.Path)
			return
		}
		return
	})
}
