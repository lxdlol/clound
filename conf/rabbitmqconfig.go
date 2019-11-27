package conf


const (
	// AsyncTransferEnable : 是否开启文件异步转移(默认同步)
	AsyncTransferEnable = false
	// RabbitURL : rabbitmq服务的入口url
	RabbitURL = "amqp://admin:admint@127.0.0.1:5672/"
	// TransExchangeName : 用于文件transfer的交换机
	TransExchangeName = "oss.trans"
	// TransOSSQueueName : oss转移队列名
	TransOSSQueueName = "oss.trans.oss"
	// TransOSSErrQueueName : oss转移失败后写入另一个队列的队列名
	TransOSSErrQueueName = "uploadserver.trans.oss.err"
	// TransOSSRoutingKey : routingkey
	TransOSSRoutingKey = "oss"

)
