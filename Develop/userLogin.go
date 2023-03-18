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

func (cc Controller) ChangePassword() {
	cc.group.POST("/user/changePassword", func(c *duang.Context) {
		getJsonData := app.GetRequest(c).JsonData
		token := app.GetRequest(c).Header["Token"]
		uid, e := duang.GetUser(token)
		if e != "" {
			c.JSON(http.StatusUnauthorized, duang.H{
				"code": 401,
				"msg":  e,
			})
			return
		}
		user := duang.Pdo_get("user", []string{}, map[string]interface{}{"id": uid})
		rePassword := app.Md5_1(app.Md5_1(getJsonData["password"].(string)) + user["salt"].(string))
		if rePassword == user["password"] {
			newPassword := getJsonData["newPassword"]
			salt := app.Md5_1(newPassword.(string))[5:20]
			newPassword = app.Md5_1(app.Md5_1(newPassword.(string)) + salt)
			d := map[string]interface{}{
				"password": newPassword,
				"salt":     salt,
			}
			result := duang.Pdo_update("user", d, map[string]interface{}{"id": uid})
			if result != 0 {
				c.JSON(http.StatusOK, duang.H{
					"code": 200,
					"msg":  "ok",
				})
				return
			}
		}
		c.JSON(http.StatusOK, duang.H{
			"code": 400,
			"msg":  "密码错误",
		})
		return
	})
}
