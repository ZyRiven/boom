/**
 *
 * @company: Co.预见（天津）智能科技有限公司
 * @Author:  ZhaoYi
 * @Date:    2023/3/3 11:53
 */

package Develop

import "gohello/duang"

type (
	Controller struct {
		group *duang.RouterGroup
	}
	// CModbus modbus结构体
	CModbus struct {
		group *duang.RouterGroup
	}
)

func New(g *duang.RouterGroup) Controller {
	return Controller{group: g}
}

func NewCModbus(g *duang.RouterGroup) CModbus {
	return CModbus{group: g}
}
