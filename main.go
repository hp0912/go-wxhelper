package main

import (
	"bytes"
	"go-wechat/config"
	"go-wechat/handler"
	"go-wechat/initialization"
	"go-wechat/tasks"
	"go-wechat/utils"
	"io"
	"log"
	"net"
	"time"
)

func init() {
	initialization.InitConfig()
	tasks.InitTasks()
}

func process(conn net.Conn) {
	// 处理完关闭连接
	defer func() {
		log.Printf("处理完成: -> %s", conn.RemoteAddr())
		_ = conn.Close()
	}()

	var buf bytes.Buffer
	if _, err := io.Copy(&buf, conn); err != nil {
		log.Printf("[%s]返回数据失败，错误信息: %v", conn.RemoteAddr(), err)
	}
	log.Printf("[%s]数据长度: %d", conn.RemoteAddr(), buf.Len())
	go handler.Parse(conn.RemoteAddr(), buf.Bytes())
	// 将接受到的数据返回给客户端
	if _, err := conn.Write([]byte("200 OK")); err != nil {
		log.Printf("[%s]返回数据失败，错误信息: %v", conn.RemoteAddr(), err)
	}
}

func main() {
	// 如果启用了自动配置回调，就设置一下
	if config.Conf.Wechat.AutoSetCallback {
		utils.ClearCallback()
		time.Sleep(500 * time.Millisecond) // 休眠五百毫秒再设置
		utils.SetCallback(config.Conf.Wechat.Callback)
	}

	// 建立 tcp 服务
	listen, err := net.Listen("tcp", "0.0.0.0:19099")
	if err != nil {
		log.Printf("listen failed, err:%v", err)
		return
	}

	for {
		// 等待客户端建立连接
		conn, err := listen.Accept()
		if err != nil {
			log.Printf("accept failed, err:%v", err)
			continue
		}
		// 启动一个单独的 goroutine 去处理连接
		go process(conn)
	}
}
