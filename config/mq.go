package config

import "fmt"

// mq
// @description: MQ配置
type mq struct {
	RabbitMQ rabbitMq `json:"rabbitmq" yaml:"rabbitmq"` // RabbitMQ配置
}

// rabbitMq
// @description: RabbitMQ配置
type rabbitMq struct {
	Host     string `json:"host" yaml:"host"`         // 主机地址
	Port     int    `json:"port" yaml:"port"`         // 端口
	User     string `json:"user" yaml:"user"`         // 用户名
	Password string `json:"password" yaml:"password"` // 密码
	VHost    string `json:"vhost" yaml:"vhost"`       // 虚拟主机
}

// GetURL
// @description: 获取MQ连接地址
// @receiver r
// @return string
func (r rabbitMq) GetURL() string {
	port := r.Port
	if port == 0 {
		port = 5672
	}
	return fmt.Sprintf("amqp://%s:%s@%s:%d/%s", r.User, r.Password, r.Host, port, r.VHost)
}
