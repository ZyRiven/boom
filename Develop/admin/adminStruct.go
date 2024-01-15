/**
 *
 * @company: Co.预见（天津）智能科技有限公司
 * @Author:  ZhaoYi
 * @Date:    2023/3/6 15:09
 */

package admin

import "boom/duang"

type Memu struct {
	group *duang.RouterGroup
}

func NewMemu(g *duang.RouterGroup) Memu {
	return Memu{group: g}
}
