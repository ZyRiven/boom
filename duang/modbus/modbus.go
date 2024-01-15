/**
 *
 * @company: Co.预见（天津）智能科技有限公司
 * @Author:  ZhaoYi
 * @Date:    2023/2/24 10:34
 */

package modbus

import (
	"fmt"
	"github.com/goburrow/serial"
	"github.com/things-go/go-modbus"
	"boom/config"
	"log"
)

func ModbusRtuMain(slaveID byte, funCode int, addr uint16, quantity uint16, v ...uint16) interface{} {
	p := modbus.NewRTUClientProvider(modbus.WithEnableLogger(),
		modbus.WithSerialConfig(serial.Config{
			Address:  config.RetModbus().Address,
			BaudRate: config.RetModbus().BaudRate,
			DataBits: config.RetModbus().DataBits,
			StopBits: config.RetModbus().StopBits,
			Parity:   config.RetModbus().Parity,
			Timeout:  modbus.SerialDefaultTimeout,
		}))

	client := modbus.NewClient(p)
	err := client.Connect()
	if err != nil {
		fmt.Println("连接失败, ", err)
	}
	defer client.Close()
	//for {
	var results interface{}
	switch funCode {
	case 1:
		results, err = client.ReadCoils(slaveID, addr, quantity)
	case 3:
		results, err = client.ReadHoldingRegisters(slaveID, addr, quantity)
	case 6:
		err = client.WriteSingleRegister(slaveID, addr, quantity)
	case 10:
		err = client.WriteMultipleRegisters(slaveID, addr, quantity, v)
	default:
		return "功能码无效"
	}
	if err != nil {
		fmt.Println(err.Error())
	}
	return results
	//time.Sleep(time.Second * 2)
	//}
}

func ModbusTcpMain(slaveID byte, funCode int, addr uint16, quantity uint16, v ...uint16) interface{} {
	p := modbus.NewTCPClientProvider(config.RetModbus().TCPAddr, modbus.WithEnableLogger())
	client := modbus.NewClient(p)
	err := client.Connect()
	if err != nil {
		fmt.Println("connect failed, ", err)
	}
	defer client.Close()
	//for {
	var results interface{}
	switch funCode {
	case 1:
		results, err = client.ReadCoils(slaveID, addr, quantity)
	case 3:
		results, err = client.ReadHoldingRegisters(slaveID, addr, quantity)
	case 6:
		err = client.WriteSingleRegister(slaveID, addr, quantity)
	case 10:
		err = client.WriteMultipleRegisters(slaveID, addr, quantity, v)
	default:
		return "功能码无效"
	}
	if err != nil {
		fmt.Println(err.Error())
	}
	return results
	//time.Sleep(time.Second * 2)
	//}
}

func ModbusTCPServer() {
	srv := modbus.NewTCPServer()
	srv.LogMode(true)
	srv.AddNodes(
		modbus.NewNodeRegister(
			1,
			0, 10, 0, 10,
			0, 10, 0, 10),
		modbus.NewNodeRegister(
			2,
			0, 10, 0, 10,
			0, 10, 0, 10),
		modbus.NewNodeRegister(
			3,
			0, 10, 0, 10,
			0, 10, 0, 10))

	err := srv.ListenAndServe(":502")
	if err != nil {
		log.Panicln(err)
	}
}
