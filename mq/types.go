package mq

// Producer 是生产者抽象，用于向指定Topic发送/生产消息,可以被多个协程并发访问
type Producer interface {
	Publish(message string) error

	// Close 用于释放资源，多次调用返回的error与第一次调用返回的error相同
	Close()
}

// Consumer 是消费者的抽象，用于从指定Topic接收/消费消息,可以被多个协程并发访问
type Consumer interface {
	Consume() error

	// Close 用于释放资源，多次调用返回的error与第一次调用返回的error相同
	Close()
}
