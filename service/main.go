package main

import (
	"myyun/route"
)

func main() {

	// 静态资源处理
	// http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(assets.AssetFS())))
	//http.Handle("/static/",
	//	http.StripPrefix("/static/", http.FileServer(http.Dir("./static"))))
	//
	//http.HandleFunc("/file/upload",handler.UploadHandler)
	//http.HandleFunc("/file/upload/suc",handler.UploadSucHandler)
	//http.HandleFunc("/file/meta",handler.GetFileHandler)
	//http.HandleFunc("/file/download",handler.DownloadHandler)
	//http.HandleFunc("/file/query",handler.FindFileHandler)
	//http.HandleFunc("/file/update",handler.UpdateFileHandler)
	//http.HandleFunc("file/delete",handler.DeleteFileHandler)
	//
	//http.HandleFunc("/file/downloadurl", handler.HTTPInterceptor(
	//	handler.DownloadURLHandler))
	//
	//// 分块上传接口
	//http.HandleFunc("/file/mpupload/init",
	//	handler.HTTPInterceptor(handler.InitialMultipartUploadHandler))
	//http.HandleFunc("/file/mpupload/uppart",
	//	handler.HTTPInterceptor(handler.UploadPartHandler))
	//http.HandleFunc("/file/mpupload/complete",
	//	handler.HTTPInterceptor(handler.CompleteUploadHandler))
	////singin
	//http.HandleFunc("/user/signin",handler.SignInHandler)
	//http.HandleFunc("/user/signup",handler.RegisterHandler)
	//http.HandleFunc("/user/info",handler.HTTPInterceptor(handler.UserInfoHandler))
	router := route.Router()
	router.Run(":8080")

}
