/**
 * 路由注册
 *
 * @company: Co.预见（天津）智能科技有限公司
 * @Author:  ZhaoYi
 * @Date:    2023/2/1 11:08
 */

package route

import (
	"gohello/Develop"
	"gohello/Develop/admin"
	"gohello/duang"
	"reflect"
)

func Init(r *duang.Engine) {
	// 后台接口
	a := r.Group("/admin")
	cAdmin := admin.NewMemu(a)
	RegisterController(cAdmin)
	// web接口
	web := r.Group("/web")
	cWeb := Develop.New(web)
	RegisterController(cWeb)
	//modbus接口
	modbusTest := r.Group("/modbus")
	cModbusTest := Develop.NewCModbus(modbusTest)
	RegisterController(cModbusTest)
}

// RegisterController 反射体
func RegisterController(controller interface{}) {
	val := reflect.ValueOf(controller)
	numOfMethod := val.NumMethod()
	for i := 0; i < numOfMethod; i++ {
		val.Method(i).Call(nil)
	}
}
