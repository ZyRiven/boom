/**
 *
 * @company: Co.预见（天津）智能科技有限公司
 * @Author:  ZhaoYi
 * @Date:    2023/3/6 14:59
 */

package admin

import (
	"fmt"
	"gohello/app"
	"gohello/duang"
	"net/http"
	"strings"
)

func (m Memu) AmdinMemu() {
	m.group.GET("/memu", func(c *duang.Context) {
		list, _, _ := duang.Pdo_getall("auth_rule", []string{}, map[string]interface{}{})
		for _, v := range list {
			v["meta"] = map[string]string{
				"icon":  v["icon"].(string),
				"title": v["title"].(string),
			}
			delete(v, "icon")
			delete(v, "title")
		}
		tree := app.ListToTree(&list, 0)
		c.JSON(http.StatusOK, duang.H{
			"code": 200,
			"data": tree,
		})
	})
}

func (m Memu) RoleData() {
	m.group.POST("/role", func(c *duang.Context) {
		data := app.GetRequest(c).JsonData
		w := map[string]interface{}{
			"Limit": [2]int{int(data["page"].(float64)), int(data["pageSize"].(float64))},
			//"Order": "id DESC",
		}
		list, pageNum, total := duang.Pdo_getall("role", []string{}, w)
		c.JSON(http.StatusOK, duang.H{
			"code":    200,
			"msg":     "ok",
			"data":    list,
			"pageNum": pageNum,
			"total":   total,
		})
		return
	})
}

// RoleEdit 权限修改/添加
func (m Memu) RoleEdit() {
	m.group.POST("/role/edit", func(c *duang.Context) {
		data := app.GetRequest(c).JsonData
		var result interface{}
		fmt.Println(data)
		//if data["id"] == nil {
		//	result = duang.Pdo_insert("role", data)
		//} else {
		//	id := data["id"]
		//	delete(data, "id")
		//	result = duang.Pdo_update("role", data, map[string]interface{}{"id": id})
		//}
		if result == 0 {
			c.JSON(http.StatusOK, duang.H{
				"code": 400,
				"msg":  "error",
			})
			return
		} else {
			c.JSON(http.StatusOK, duang.H{
				"code": 200,
				"msg":  "ok",
				"data": result,
			})
			return
		}
	})
}

// RoleDel 权限删除
func (m Memu) RoleDel() {
	m.group.POST("/role/del", func(c *duang.Context) {
		data := app.GetRequest(c).JsonData
		ids := strings.Split(data["ids"].(string), ",")
		result := duang.Pdo_delete("role", map[string]interface{}{"id": ids})
		if result == 0 {
			c.JSON(http.StatusOK, duang.H{
				"code": 400,
				"msg":  "error",
			})
			return
		} else {
			c.JSON(http.StatusOK, duang.H{
				"code": 200,
				"msg":  "ok",
				"data": result,
			})
			return
		}
	})
}
func (m Memu) AdminData() {
	m.group.POST("/adminData", func(c *duang.Context) {
		list, _, _ := duang.Pdo_getall("admin", []string{}, map[string]interface{}{})
		c.JSON(http.StatusOK, duang.H{
			"code": 200,
			"msg":  "ok",
			"data": list,
		})
	})
}
