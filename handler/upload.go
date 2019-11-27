package handler

import (
	"encoding/json"
	"github.com/gin-gonic/gin"
	"myyun/mq"
	"fmt"
	"io"
	"io/ioutil"
	"myyun/conf"
	"myyun/db"
	"myyun/meta"
	"myyun/util"
	"net/http"
	"myyun/storage/oss"
	"os"
	"strconv"
	"strings"
)



func UploadHandler(c *gin.Context){
	//back home html
	c.Redirect(http.StatusFound,"/static/view/index.html")
}
func UploadHandlerPsot(c *gin.Context){
	file, header, e := c.Request.FormFile("file")
	if e!=nil{
		fmt.Printf("failed to get data ,err:%s/n",e.Error())
		return
	}
	defer file.Close()
	filemeta:=meta.FileMeta{}
	filemeta.FileName=header.Filename
	filemeta.FilePath="./temp/" + header.Filename
	create, e := os.Create("./temp/" + header.Filename)
	if e!=nil{
		fmt.Printf("failed to creat new file ,err:%s/n",e.Error())
		return
	}
	defer create.Close()
	filemeta.FileSize, e = io.Copy(create, file)
	if e!=nil{
		fmt.Printf("copy failed ,err:%s/n",e.Error())
		return
	}
	//filemeta.UploadTime=time.Now().Format("2006-01-02 15:04:05")
	create.Seek(0,0)
	filemeta.Filehash= util.FileSha1(create)
	//meta.UpdateFilemeta(filemeta)
	// update  file to mysql
	var osspath string
	osspath="oss/data/"+filemeta.Filehash
	if conf.AsyncTransferEnable{
		bucket := oss.OssBucket()
		err := bucket.PutObject(osspath, create)
		if err!=nil{
			fmt.Println(err.Error())
			c.JSON(http.StatusOK,gin.H{
				"code":-1,
				"msg":"Upload failed!",
			})
			return
		}
	}else{
		var msg mq.Msg
		msg=mq.Msg{
			filemeta.Filehash,
			filemeta.FilePath,
			osspath,
			conf.StoreOSS,
		}
		bytes, _ := json.Marshal(msg)
		ok := mq.Publish(conf.TransExchangeName, conf.TransOSSRoutingKey, bytes)
		if !ok{
			//todo
		}
	}
	filemeta.FilePath=osspath
	meta.UpdateFilemetaDb(filemeta)
	username := c.Request.FormValue("username")
	ok := db.OnUserFileUploadFinished(username, filemeta.FileName, filemeta.Filehash, filemeta.FileSize)
	if ok{
		c.Redirect(http.StatusFound,"/static/view/home.html")
	}else{
		c.JSON(http.StatusOK,gin.H{
			"code":-1,
			"msg":"Upload failed!",
		})
	}
}

func UploadSucHandler(c *gin.Context)  {
	c.JSON(http.StatusOK,gin.H{
		"code":0,
		"msg":"upload ok",
	})
}

//
func GetFileHandler(c *gin.Context) {
	fileMeta,e := meta.GetFileByHash(c.Request.Form["filehash"][0])
	bytes, e := json.Marshal(fileMeta)
	if e!=nil{
		c.JSON(http.StatusInternalServerError,gin.H{})
		return
	}
	c.Writer.Write(bytes)
}
//
func FindFileHandler(c *gin.Context){
	i, _ := strconv.Atoi(c.Request.FormValue("limit"))
	username := c.Request.FormValue("username")
	files, err := db.QueryUserFileMetas(username, i)
	if err!=nil{
		c.Writer.WriteHeader(http.StatusInternalServerError)
		return
	}
	bytes, e := json.Marshal(files)
	if e!=nil{
		c.Writer.WriteHeader(http.StatusInternalServerError)
		return
	}
	c.Writer.Write(bytes)
}

//
func DownloadHandler (c *gin.Context){
	filemate,e := meta.GetFileByHash(c.Request.Form["filehash"][0])
	file, e := os.Open(filemate.FilePath)
	if e!=nil{
		c.Writer.WriteHeader(http.StatusInternalServerError)
		return
	}
	defer file.Close()
	bytes, e := ioutil.ReadAll(file)
	if e!=nil{
		c.Writer.WriteHeader(http.StatusInternalServerError)
		return
	}
	c.Writer.Header().Set("Content-Type","application/octet-stream")
	c.Writer.Header().Set("Content-Disposition","attachment;filename=\""+filemate.FileName+"\"")
	c.Writer.Write(bytes)

}

func UpdateFileHandler (c *gin.Context){
	op := c.Request.FormValue("op")
	hash := c.Request.FormValue("filehash")
	newName:=c.Request.FormValue("name")
	if op!="0"{
		c.Writer.WriteHeader(http.StatusForbidden)
		return
	}
	fileMeta,err := meta.GetFileByHash(hash)
	if err!=nil{
		c.Writer.WriteHeader(http.StatusInternalServerError)
		return
	}
	fileMeta.FileName=newName
	bytes, e := json.Marshal(fileMeta)
	if e!=nil{
		c.Writer.WriteHeader(http.StatusInternalServerError)
		return
	}
	c.Writer.Write(bytes)
}


func DeleteFileHandler(c *gin.Context){

	hash:=c.Request.FormValue("filehash")
	name := c.Request.FormValue("username")
	fileMeta,err := meta.GetFileByHash(hash)
	if err!=nil{
		c.Writer.WriteHeader(http.StatusInternalServerError)
		return
	}
	os.Remove(fileMeta.FilePath)
	meta.Deletefile(hash,name)
	c.Writer.WriteHeader(http.StatusOK)
}


// DownloadURLHandler : 生成文件的下载地址
func DownloadURLHandler(c *gin.Context) {
	filehash :=c.Request.FormValue("filehash")
	// 从文件表查找记录
	row, _ := db.GetFileMeta(filehash)

	// TODO: 判断文件存在OSS，还是Ceph，还是在本地
	if strings.HasPrefix(row.FileAddr.String, "/tmp") {
		username := c.Request.FormValue("username")
		token := c.Request.FormValue("token")
		tmpUrl := fmt.Sprintf("http://%s/file/download?filehash=%s&username=%s&token=%s",
			c.Request.Host, filehash, username, token)
		c.Writer.Write([]byte(tmpUrl))
	} else if strings.HasPrefix(row.FileAddr.String, "/ceph") {
		// TODO: ceph下载url
	} else if strings.HasPrefix(row.FileAddr.String, "oss/") {
		// oss下载url
		signedURL := oss.DownloadURL(row.FileAddr.String)
		c.Writer.Write([]byte(signedURL))
	}
}

