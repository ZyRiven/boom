/*
* 配置
*
* @company: Co.预见（天津）智能科技有限公司
* @Author:  ZhaoYi
* @Date:    2023/1/31 18:08
 */

package config

type (
	Database struct {
		Driver string // 数据库引擎
		Prefix string // 表前缀
		Dns    string
	}

	Server struct {
		Address       string // 项目端口
		WebSocketPort string // websocket端口
	}

	// Modbus ModbusRtu相关配置
	Modbus struct {
		Address  string // 端口
		BaudRate int    // 波特率
		DataBits int    // 数据位
		StopBits int    // 停止位
		Parity   string // 奇偶校验:N -无，E -偶数，O -奇数(默认值)E)(使用无奇偶校验需要2个停止位。)
		TCPAddr  string //tcp地址+端口
	}
)

func RetDatabase() Database {
	return Database{
		Driver: "mysql",
		Prefix: "boom_",
		Dns:    "root:123456@tcp(127.0.0.1:3306)/boom?charset=utf8mb4&parseTime=True&loc=Local",
	}
}

func RetServer() Server {
	return Server{
		Address:       ":9999",
		WebSocketPort: ":8888",
	}
}

func RetModbus() Modbus {
	return Modbus{
		Address:  "COM3",
		BaudRate: 115200,
		DataBits: 8,
		StopBits: 1,
		Parity:   "N",
		TCPAddr:  "127.0.0.1:502",
	}
}
