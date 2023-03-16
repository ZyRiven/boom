/**
 * 1.引入了 github.com/gorilla/websocket 包，先创建 http 服务器，然后升级为 websocket 服务器；
 * 2.启动服务时，马上起一个协程作为处理中心，监听注册、注销、消息 3 个总的 channel（是的，通过channel通信）；
 * 3.当用户打开聊天页面，则建立连接，并新建一个 Client 实例，将实例指针推给处理中心的注册 channel 后，自己再起两个协程监听读、写操作；
 * 4.用户点击发送消息时，自己起的读协程将消息推给处理中心的广播 channel，循环将消息推给每个 client 的写 channel，对应 client  的写通道最将消息主动推给每个用户；
 * 5.这个思路比较复杂，但主要的抽象出来一个  Client 充当 客户端 和 websocket处理中心 的中介，每个 Client 都运行两个 goroutine ：
 *  5.1 读，循环监听是否该用户是否要发言，有则推给 websocket  的广播渠道 推给每个客户端
 *  5.2 写，循环读取 websocket 是否有新消息，有则主动推送消息给客户端
 *
 * @company: Co.预见（天津）智能科技有限公司
 * @Author:  ZhaoYi
 * @Date:    2023/1/31 18:08
 */

package duang

import (
	"encoding/json"
	"fmt"
	"github.com/gorilla/websocket"
	"gohello/config"
	"golang.org/x/sync/errgroup"
	"log"
	"net/http"
)

var (
	// 升级成 WebSocket 协议
	upgrader = websocket.Upgrader{
		// 允许CORS跨域请求
		CheckOrigin: func(r *http.Request) bool {
			return true
		},
	}
	conn *websocket.Conn
	err  error
	g    errgroup.Group
)

type (
	// CenterHandler 处理中心，关联着每个 Client 的注册、注销、广播通道，相当于每个用户的中心通讯的中介。
	CenterHandler struct {
		// 广播通道，有数据则循环每个用户广播出去
		broadcast chan []byte
		// 注册通道，有用户进来 则推到用户集合map中
		register chan *Client
		// 注销通道，有用户关闭连接 则将该用户剔出集合map中
		unregister chan *Client
		// 用户集合，每个用户本身也在跑两个协程，监听用户的读、写的状态
		clients map[*Client]string
	}
	// Client 抽象出来的 Client，里面有这个 websocket 连接的 读 和 写 操作
	Client struct {
		handler *CenterHandler
		conn    *websocket.Conn
		// 每个用户自己的循环跑起来的状态监控
		send chan []byte
	}
)

// 处理中心的一个接口，监控状态
func (ch *CenterHandler) monitoring() {
	for {
		select {
		// 注册，新用户连接过来会推进注册通道，这里接收推进来的用户指针
		case client := <-ch.register:
			err := Init("2023-02-20", 1)
			if err != nil {
				return
			}
			id := GenID()
			mess := map[string]string{
				"data": id,
				"type": "id",
			}
			msg, _ := json.Marshal(mess)
			client.send <- msg
			ch.clients[client] = id
			// 注销，关闭连接或连接异常会将用户推出群聊
		case client := <-ch.unregister:
			delete(ch.clients, client)
			// 消息，监听到有新消息到来
		case message := <-ch.broadcast:
			println("message：" + string(message))
			js := make(map[string]interface{})
			err := json.Unmarshal(message, &js)
			if err != nil {
				log.Println(err)
				break
			}
			// 推送给每个用户的通道，每个用户都有跑协程起了writePump的监听
			for client := range ch.clients {
				//if "chat" == js["type"] {
				//	if v == js["toId"] || v == js["selfId"] {
				//		client.send <- message
				//	}
				//} else {
				//	if v == js["toId"] {
				//		client.send <- message
				//	}
				//}
				//fmt.Println(v)
				//if js["groupId"] == "come" {
				client.send <- message
				//	break
				//} else if js["groupId"] == "go" {
				//	for k, vv := range ch.clients {
				//		if vv == js["selfId"] {
				//			k.send <- message
				//			break
				//		}
				//	}
				//}

			}
		}
	}
}

// 写，主动推送消息给客户端
func (c *Client) writePump() {
	defer func() {
		c.handler.unregister <- c
		c.conn.Close()
	}()
	for {
		// 广播推过来的新消息，马上通过websocket推给自己
		message, _ := <-c.send
		if err := c.conn.WriteMessage(websocket.TextMessage, message); err != nil {
			return
		}
	}
}

// 读，监听客户端是否有推送内容过来服务端
func (c *Client) readPump() {
	defer func() {
		c.handler.unregister <- c
		c.conn.Close()
	}()
	for {
		// 循环监听是否该用户是否要发言
		_, message, err := c.conn.ReadMessage()
		if err != nil {
			// 异常关闭的处理
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				log.Printf("error: %v", err)
			}
			break
		}
		// 要的话，推给广播中心，广播中心再推给每个用户
		c.handler.broadcast <- message
	}
}

func WebSocketMain() {
	// 应用一运行，就初始化 CenterHandler 处理中心对象
	handler := CenterHandler{
		broadcast:  make(chan []byte),
		register:   make(chan *Client),
		unregister: make(chan *Client),
		clients:    make(map[*Client]string),
	}
	// 起个协程跑起来，监听注册、注销、消息 3 个 channel
	go handler.monitoring()
	// websocket 请求，建立双工通讯连接
	http.HandleFunc("/ws", func(writer http.ResponseWriter, request *http.Request) {
		fmt.Println("Welcome")
		// 由 http 升级成为 websocket 服务
		if conn, err = upgrader.Upgrade(writer, request, nil); err != nil {
			log.Println(err)
			return
		}
		// 为每个连接创建一个 Client 实例，（实际上这里应该还有绑定用户真实信息的操作）
		client := &Client{&handler, conn, make(chan []byte, 256)}
		// 推给监控中心注册到用户集合中
		handler.register <- client
		// 每个 client 都挂起 2 个新的协程，监控读、写状态
		go client.writePump()
		go client.readPump()
	})
	server := &http.Server{
		Addr:    config.RetServer().WebSocketPort,
		Handler: nil,
	}
	g.Go(func() error {
		return server.ListenAndServe()
	})
}
