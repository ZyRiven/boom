/**
 *
 * @company: Co.预见（天津）智能科技有限公司
 * @Author:  ZhaoYi
 * @Date:    2023/2/27 11:54
 */

package Develop

import (
	"gohello/duang"
	"gohello/duang/modbus"
	"net/http"
)

func (m CModbus) Rtu() {
	m.group.GET("/rtu", func(c *duang.Context) {
		//a1 := []uint16{1, 2, 3, 4}
		a := modbus.ModbusRtuMain(1, 3, 0, 4)
		c.JSON(http.StatusBadRequest, duang.H{
			"code": 200,
			"data": a,
		})
		return
	})
}

func (m CModbus) Tcp() {
	m.group.GET("/tcp", func(c *duang.Context) {
		a := modbus.ModbusTcpMain(1, 6, 0, 3)
		c.JSON(http.StatusBadRequest, duang.H{
			"code": 200,
			"data": a,
		})
		return
	})
}
