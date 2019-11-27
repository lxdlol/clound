package main


import (
	"bufio"
	"encoding/json"
	"myyun/conf"
	"myyun/mq"
	"log"
	"os"
	"myyun/storage/oss"
	"myyun/db"
)

func ProcessTransfer(msg []byte) bool {
	log.Println(string(msg))

	pubData := mq.Msg{}
	err := json.Unmarshal(msg, &pubData)
	if err != nil {
		log.Println(err.Error())
		return false
	}

	fin, err := os.Open(pubData.CurLocation)
	if err != nil {
		log.Println(err.Error())
		return false
	}

	err = oss.OssBucket().PutObject(
		pubData.DestLocation,
		bufio.NewReader(fin))
	if err != nil {
		log.Println(err.Error())
		return false
	}

	_ = db.UpdateFileLocation(
		pubData.FileHash,
		pubData.DestLocation)
	return true
}

func main() {
	if !conf.AsyncTransferEnable {
		log.Println("异步转移文件功能目前被禁用，请检查相关配置")
		return
	}
	log.Println("文件转移服务启动中，开始监听转移任务队列...")
	mq.Consummer(
		conf.TransOSSQueueName,
		"transfer_oss",
		ProcessTransfer)
}
