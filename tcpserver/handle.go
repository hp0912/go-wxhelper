package tcpserver

import (
	"bytes"
	"go-wechat/handler"
	"io"
	"log"
	"net"
)

// process
// @description: 处理连接
// @param conn
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
