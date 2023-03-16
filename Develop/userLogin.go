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
