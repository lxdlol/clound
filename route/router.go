package route

import (
	"github.com/gin-gonic/gin"
	"myyun/handler"
)

func Router()*gin.Engine{
	router := gin.Default()
	router.Static("/static/","./static")

	router.POST("/user/signin",handler.SignInHandlerPost)
	router.GET("/user/signin",handler.SignInHandler)
	router.POST("/user/signup",handler.RegisterHandler)
	router.GET("/user/signup",handler.RegisterHandlerPost)

	router.Use(handler.HTTPInterceptor())
	//http.HandleFunc("/user/info",handler.HTTPInterceptor(handler.UserInfoHandler))
	router.POST("/user/info",handler.UserInfoHandler)
	//http.HandleFunc("/file/upload",handler.UploadHandler)
	router.GET("/file/upload",handler.UploadHandler)
	router.POST("/file/upload",handler.UploadHandlerPsot)
	//http.HandleFunc("/file/upload/suc",handler.UploadSucHandler)
	router.POST("/file/upload/suc",handler.UploadSucHandler)
	//http.HandleFunc("/file/meta",handler.GetFileHandler)
	router.POST("/file/meta",handler.GetFileHandler)
	//http.HandleFunc("/file/download",handler.DownloadHandler)
	router.POST("/file/download",handler.DownloadHandler)
	//http.HandleFunc("/file/query",handler.FindFileHandler)
	router.POST("/file/query",handler.FindFileHandler)
	//http.HandleFunc("/file/update",handler.UpdateFileHandler)
	router.POST("/file/update",handler.UpdateFileHandler)
	//http.HandleFunc("file/delete",handler.DeleteFileHandler)
	router.POST("file/delete",handler.DeleteFileHandler)
	//http.HandleFunc("/file/downloadurl", handler.HTTPInterceptor(
	//	handler.DownloadURLHandler))
	router.POST("/file/downloadurl",handler.DownloadURLHandler)
	//// 分块上传接口
	//http.HandleFunc("/file/mpupload/init",
	//	handler.HTTPInterceptor(handler.InitialMultipartUploadHandler))
	router.POST("/file/mpupload/init",handler.InitialMultipartUploadHandler)
	//http.HandleFunc("/file/mpupload/uppart",
	//	handler.HTTPInterceptor(handler.UploadPartHandler))
	router.POST("/file/mpupload/uppart",handler.UploadPartHandler)
	//http.HandleFunc("/file/mpupload/complete",
	//	handler.HTTPInterceptor(handler.CompleteUploadHandler))
	router.POST("/file/mpupload/complete",handler.CompleteUploadHandler)
	return router
}
