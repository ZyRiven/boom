/**
 *
 * @company: Co.预见（天津）智能科技有限公司
 * @Author:  ZhaoYi
 * @Date:    2023/2/2 17:20
 */

package Develop

import (
	"gohello/app"
	"gohello/duang"
	"net/http"
	"time"
)

func (cc Controller) Login() {
	cc.group.POST("/user/login", func(c *duang.Context) {
		getJsonData := app.GetRequest(c).JsonData
		if getJsonData["mobile"] == nil || getJsonData["password"] == nil {
			c.JSON(http.StatusBadRequest, duang.H{
				"code": 400,
				"msg":  "用户名或密码不为空",
			})
			return
		}
		w := map[string]interface{}{
			"mobile": getJsonData["mobile"],
		}
		user := duang.Pdo_get("user", []string{}, w)
		if user == nil {
			c.JSON(http.StatusBadRequest, duang.H{
				"code": 400,
				"msg":  "用户名或密码不正确",
			})
			return
		}
		rePassword := app.Md5_1(app.Md5_1(getJsonData["password"].(string)) + user["salt"].(string))
		if rePassword == user["password"] {
			user["token"], _ = duang.EnToken(user["id"].(string))
			delete(user, "password")
			delete(user, "salt")
			c.JSON(http.StatusOK, duang.H{
				"code": 200,
				"msg":  "ok",
				"data": user,
			})
			return
		}
		c.JSON(http.StatusBadRequest, duang.H{
			"code": 400,
			"msg":  "用户名或密码不正确",
		})
		return
	})
}

// Register 注册用户
func (cc Controller) Register() {
	cc.group.POST("/user/register", func(c *duang.Context) {
		data := app.GetRequest(c).JsonData
		if data["mobile"] == nil || data["password"] == nil {
			c.JSON(http.StatusBadRequest, duang.H{
				"code": 400,
				"msg":  "用户名或密码不为空",
			})
			return
		}
		w := map[string]interface{}{
			"mobile": data["mobile"],
		}
		user := duang.Pdo_count("user", w)
		if user != 0 {
			c.JSON(http.StatusBadRequest, duang.H{
				"code": 400,
				"msg":  "用户名已存在",
			})
			return
		}
		salt := app.Md5_1(data["password"].(string))[5:20]
		password := app.Md5_1(app.Md5_1(data["password"].(string)) + salt)
		data["password"] = password
		data["salt"] = salt
		data["jointime"] = time.Now().Unix()
		result := duang.Pdo_insert("user", data)
		if result != 0 {
			c.JSON(http.StatusOK, duang.H{
				"code": 200,
				"msg":  "ok",
			})
			return
		}
		c.JSON(http.StatusBadRequest, duang.H{
			"code": 400,
			"msg":  "注册失败",
		})
		return
	})
}

// UpdateUser 修改用户信息
func (cc Controller) UpdateUser() {
	cc.group.POST("/user/updateUser", func(c *duang.Context) {
		data := app.GetRequest(c).JsonData
		token := app.GetRequest(c).Header["Token"]
		userId, err := duang.GetUser(token)
		if err != "" {
			c.JSON(http.StatusUnauthorized, duang.H{
				"code": 401,
				"msg":  err,
			})
			return
		}
		if data["password"] != nil {
			salt := app.Md5_1(data["password"].(string))[5:20]
			password := app.Md5_1(app.Md5_1(data["password"].(string)) + salt)
			data["password"] = password
			data["salt"] = salt
		}
		result := duang.Pdo_update("user", data, map[string]interface{}{"id": userId})
		if result != 0 {
			c.JSON(http.StatusOK, duang.H{
				"code": 200,
				"msg":  "ok",
			})
			return
		}
		c.JSON(http.StatusBadRequest, duang.H{
			"code": 400,
			"msg":  "error",
		})
		return
	})
}
