/**
 * 工具库
 *
 * @company: Co.预见（天津）智能科技有限公司
 * @Author:  ZhaoYi
 * @Date:    2023/2/1 12:00
 */

package app

import (
	"crypto/md5"
	"encoding/hex"
	"gohello/duang"
	"net/url"
)

type (
	Request struct {
		Header   map[string][]string
		JsonData map[string]interface{}
		FormData url.Values
		Query    map[string][]string
		Group    string
		Path     string
		Host     string
	}
	Tree []map[string]interface{}
)

// GetRequest 请求头参数
func GetRequest(c *duang.Context) Request {
	path := c.Path
	request := Request{
		FormData: c.FormValues(),
		JsonData: c.PostJson(),
		Header:   c.Req.Header,
		Query:    c.Req.URL.Query(),
		Path:     path,
		Host:     c.Req.Host,
	}
	return request
}

// Md5_1 md5加密
func Md5_1(s string) string {
	m := md5.New()
	m.Write([]byte(s))
	return hex.EncodeToString(m.Sum(nil))
}

// ListToTree 菜单树
func ListToTree(list *[]map[string]interface{}, pid uint32) []map[string]interface{} {
	var tree Tree
	for _, v := range *list {
		if pid == v["pid"].(uint32) {
			v["children"] = ListToTree(list, v["id"].(uint32))
			if v["children"] == nil {
				delete(v, "children")
			}
			tree = append(tree, v)
		}
	}
	return tree
}
