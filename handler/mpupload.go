package handler

import (
	"fmt"
	"github.com/garyburd/redigo/redis"
	"github.com/gin-gonic/gin"
	"math"
	"myyun/cache"
	"myyun/db"
	"myyun/util"
	"os"
	"os/exec"
	"path"
	"strconv"
	"strings"
	"time"
)

// MultipartUploadInfo : 初始化信息
type MultipartUploadInfo struct {
	FileHash   string
	FileSize   int
	UploadID   string
	ChunkSize  int
	ChunkCount int
}



// InitialMultipartUploadHandler : 初始化分块上传
func InitialMultipartUploadHandler(c *gin.Context){
	username := c.Request.FormValue("username")
	filehash := c.Request.FormValue("filehash")
	filesize, err := strconv.Atoi(c.Request.FormValue("filesize"))
	if err!=nil{
		c.Writer.Write(util.NewRespMsg(-1, "params invalid", nil).JSONBytes())
		return
	}

	conn := cache.RedisPool().Get()
	defer conn.Close()
	var chunkinfo MultipartUploadInfo
	chunkinfo=MultipartUploadInfo{
		FileHash:filehash,
		FileSize:filesize,
		UploadID:username+fmt.Sprintf("%x",time.Now().Nanosecond()),
		ChunkSize:5*1024*1024,
		ChunkCount:int(math.Ceil(float64(filesize)/(5*1024*1024))),
	}
	conn.Do("HSET","MP_"+chunkinfo.UploadID,"filesize",filesize)
	conn.Do("HSET","MP_"+chunkinfo.UploadID,"filehash",filehash)
	conn.Do("HSET","MP_"+chunkinfo.UploadID,"chunkcount",chunkinfo.ChunkCount)
	c.Writer.Write(util.NewRespMsg(0, "OK", chunkinfo).JSONBytes())
}




// UploadPartHandler : 上传文件分块
func UploadPartHandler(c *gin.Context) {
	//	username := r.Form.Get("username")
	uploadID := c.Request.FormValue("uploadid")
	chunkIndex := c.Request.FormValue("index")

	// 2. 获得redis连接池中的一个连接
	rConn := cache.RedisPool().Get()
	defer rConn.Close()

	// 3. 获得文件句柄，用于存储分块内容
	fpath := "/data/" + uploadID + "/" + chunkIndex
	os.MkdirAll(path.Dir(fpath), 0744)
	fd, err := os.Create(fpath)
	if err != nil {
		c.Writer.Write(util.NewRespMsg(-1, "Upload part failed", nil).JSONBytes())
		return
	}
	defer fd.Close()

	buf := make([]byte, 1024*1024)
	for {
		n, err := c.Request.Body.Read(buf)
		fd.Write(buf[:n])
		if err != nil {
			break
		}
	}

	// 4. 更新redis缓存状态
	rConn.Do("HSET", "MP_"+uploadID, "chkidx_"+chunkIndex, 1)

	// 5. 返回处理结果到客户端
	c.Writer.Write(util.NewRespMsg(0, "OK", nil).JSONBytes())

}

const dirPath  = "/data/"

// CompleteUploadHandler : 通知上传合并
func CompleteUploadHandler(c * gin.Context) {
	// 1. 解析请求参数
	upid := c.Request.FormValue("uploadid")
	username := c.Request.FormValue("username")
	filehash := c.Request.FormValue("filehash")
	filesize := c.Request.FormValue("filesize")
	filename := c.Request.FormValue("filename")

	// 2. 获得redis连接池中的一个连接
	rConn := cache.RedisPool().Get()
	defer rConn.Close()

	reply, err :=redis.Values(rConn.Do("HGETALL", upid))
	if err != nil {
		c.Writer.Write(util.NewRespMsg(-1, "complete upload failed", nil).JSONBytes())
		return
	}
	chunktotal:=0
	chunkok:=0
	for i:=0;i<len(reply);i+=2 {
		k:=string(reply[i].([]byte))
		v:=string(reply[i+1].([]byte))
		if strings.HasPrefix(k, "chkidx_") && v == "1"{
			chunkok++
	}
		if k=="chunkcount"{
			chunktotal,_=strconv.Atoi(v)
		}
	}
	if chunktotal != chunkok {
		c.Writer.Write(util.NewRespMsg(-2, "invalid request", nil).JSONBytes())
		return
	}
	//
	var cmd *exec.Cmd
	filepath := dirPath + upid+"/"
	filestore :=  "/home/lxd/go/src/myyun/temp/"+filename

	cmd = exec.Command("/home/lxd/go/src/myyun/bash/"+"chunk.sh", filepath, filestore)
	// cmd.Run()
	if _, err := cmd.Output(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	} else {
		fmt.Println(filestore, " has been merge complete")
	}

	fsize, _ := strconv.Atoi(filesize)
	db.SaveFileToMysql(filehash,filename,int64(fsize),"")
	db.OnUserFileUploadFinished(username,filename,filehash,int64(fsize))
	// 6. 响应处理结果
	c.Writer.Write(util.NewRespMsg(0, "OK", nil).JSONBytes())
}