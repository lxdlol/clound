package oss

import (
	"fmt"
	"github.com/aliyun/aliyun-oss-go-sdk/oss"
)


//<yourEndpoint>", "<yourAccessKeyId>", "<yourAccessKeySecret>
const (
	Endpoint ="oss-cn-shenzhen.aliyuncs.com"
	AccessKeyId="LTAI4FoQt1eeouxVAWKv9Wby"
	AccessKeySecret="r4heOaw8XvMPYm00yTICLW9bipzuz3"
	bucketname ="myyuntest"
)

var ossclient *oss.Client


//// Client : 创建oss client对象
func OssClient()*oss.Client{
	if ossclient!=nil{
		return ossclient
	}
	var err error
	ossclient, err = oss.New(Endpoint, AccessKeyId, AccessKeySecret)
	if err != nil {
		fmt.Println(err.Error())
		return nil
	}
	return ossclient
}




// Bucket : 获取bucket存储空间
func OssBucket()*oss.Bucket{
	cli:=OssClient()
	if cli!=nil{
		bucket, e := cli.Bucket("bucketname")
		if e!=nil{
			fmt.Println(e.Error())
			return nil
		}
		return bucket
	}
	return nil
}


// DownloadURL : 临时授权下载url
func DownloadURL(name string)string{
	bucket := OssBucket()
	if bucket!=nil{
		url, err := bucket.SignURL(name, oss.HTTPGet, 3600)
		if err != nil {
			fmt.Println(err.Error())
			return ""
		}
		return url
	}
	return ""
}



// BuildLifecycleRule : 针对指定bucket设置生命周期规则