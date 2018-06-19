package resources

import (
	"os"

	"ProxyService/proxyservice/internal/queue"
)

var (
	// RabbitMQHost rabbit service host
	RabbitMQHost = /*"localhost"*/ "172.18.0.2"
	// RabbitMQPort rabbit service port
	RabbitMQPort = "5672"
	// RabbitMQUser rabbit service user
	RabbitMQUser = "guest"
	// RabbitMQPassword rabbit service pass
	RabbitMQPassword = "guest"
)

func init() {
	collectEnvs()
}

// NewQueueResource ...
func NewQueueResource() *queue.Queue {
	return queue.NewQueue(RabbitMQHost, RabbitMQPort, RabbitMQUser, RabbitMQPassword)
}

// ---- private functions ----

func collectEnvs() {
	if host := os.Getenv("rabbitmq_host"); host != "" {
		RabbitMQHost = host
	}

	if port := os.Getenv("rabbitmq_port"); port != "" {
		RabbitMQPassword = port
	}

	if user := os.Getenv("rabbitmq_user"); user != "" {
		RabbitMQUser = user
	}

	if pwd := os.Getenv("rabbitmq_password"); pwd != "" {
		RabbitMQPassword = pwd
	}
}
