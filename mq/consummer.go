package mq

import "log"
var done chan bool
func Consummer (qName, cName string, callback func(msg []byte) bool){
	consume, err := channel.Consume(
		qName,
		cName,
		true,
		false,
		false,
		false,
		nil,
	)
	if err!=nil{
		log.Println(err.Error())
		return
	}

	done = make(chan bool)

	go func() {
		// 循环读取channel的数据
		for d := range consume {
			processErr := callback(d.Body)
			if processErr {
				// TODO: 将任务写入错误队列，待后续处理
			}
		}
	}()

	// 接收done的信号, 没有信息过来则会一直阻塞，避免该函数退出
	<-done

	// 关闭通道
	channel.Close()

}

// StopConsume : 停止监听队列
func StopConsume() {
	done <- true
}