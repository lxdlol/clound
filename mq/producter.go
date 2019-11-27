package mq

import (
	"github.com/streadway/amqp"
	"log"
	"myyun/conf"
)

var channel *amqp.Channel
var conn *amqp.Connection

var notifyClose chan *amqp.Error

//init
func init(){

	if initChannel() {
		channel.NotifyClose(notifyClose)
	}
	// 断线自动重连
	go func() {
		for {
			select {
			case msg := <-notifyClose:
				conn = nil
				channel = nil
				log.Printf("onNotifyChannelClosed: %+v\n", msg)
				initChannel()
			}
		}
	}()
}

func initChannel()bool{
	if channel!=nil{
		return true
	}
	conn, e := amqp.Dial(conf.RabbitURL)
	if e!=nil{
		log.Println(e.Error())
		return false
	}
	channel, e = conn.Channel()
	if e!=nil{
		log.Println(e.Error())
		return false
	}
	return true
}

func Publish(exchange,key string,data []byte)bool{
	err := channel.Publish(
		exchange,
		key,
		false,
		false,
		amqp.Publishing{
			ContentType: "text/plain",
			Body:        data,
		},
	)
	if err!=nil{
		log.Println(err.Error())
		return false
	}
	return true
}



