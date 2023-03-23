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
	RegisterController(admin.NewMemu(r.Group("/admin")))
	// 用户接口
	RegisterController(Develop.New(r.Group("/user")))
	//modbus接口
	RegisterController(Develop.NewCModbus(r.Group("/modbus")))
}

// RegisterController 反射体
func RegisterController(controller interface{}) {
	val := reflect.ValueOf(controller)
	numOfMethod := val.NumMethod()
	for i := 0; i < numOfMethod; i++ {
		val.Method(i).Call(nil)
	}
}
