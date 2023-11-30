package tcpserver

import (
	"log"
	"net"
)

// Start
// @description: 启动服务
func Start() {
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
