package handler

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"myyun/db"
	"myyun/util"
	"net/http"
	"time"
)
const (
	pwdsalt ="$&*.941027"
	tokensalt="*&^.6464"
)
//sign in

func SignInHandler(c *gin.Context)  {
	c.Redirect(http.StatusFound,"/static/view/signin.html")
	return
}
func SignInHandlerPost(c *gin.Context){
	name:=c.Request.FormValue("username")
	pw := c.Request.FormValue("password")
	//
	pwdsha1:=util.Sha1([]byte(pw+pwdsalt))
	signin := db.UserSignin(name, pwdsha1)
	if signin==false{
		c.JSON(http.StatusOK,
			gin.H{
				"code":-1,
				"msg":"username or password is mistake",
			})
		return
	}
	//
	token := Token(name)
	ok := db.UpdateToken(name, token)
	if ok==false{
		c.JSON(http.StatusOK,
			gin.H{
				"code":-1,
				"msg":"signin faild",
			})
		return
	}
	//
	// 3. 登录成功后重定向到首页
	//w.Write([]byte("http://" + r.Host + "/static/view/home.html"))
	resp := util.RespMsg{
		Code: 0,
		Msg:  "OK",
		Data: struct {
			Location string
			Username string
			Token    string
		}{
			Location: "http://" + c.Request.Host + "/static/view/home.html",
			Username: name,
			Token:    token,
		},
	}
	c.Data(http.StatusOK,"ContentType/JSON",resp.JSONBytes())
}





//register

func RegisterHandler(c *gin.Context)  {
		//http.Redirect(w,r,"/static/view/signup.html",http.StatusOK)
		c.Redirect(http.StatusFound,"/static/view/signup.html")
		return
}

func RegisterHandlerPost(c *gin.Context){

	name := c.Request.FormValue("username")
	pw :=c.Request.FormValue("password")
	//search from mysql
	if len(name) < 3 || len(pw) < 5 {
		c.JSON(http.StatusOK,gin.H{
			"code":-1,
			"msg":"Invalid parameter",
		})
		return
	}
	pwsha1:=util.Sha1([]byte(pw+pwdsalt))
	ok := db.InsertUserToMysql(name, pwsha1)
	if !ok{
		c.JSON(http.StatusOK,gin.H{
			"code":-1,
			"msg":"register fail",
		})
		return
	}else{
		c.JSON(http.StatusOK,gin.H{
			"code":0,
			"msg":"register suc",
		})
	}



}








//get usrinfo
func  UserInfoHandler(c *gin.Context){
	// 1. 解析请求参数
	username := c.Request.FormValue("username")
	//	token := r.Form.Get("token")

	// // 2. 验证token是否有效
	// isValidToken := IsTokenValid(token)
	// if !isValidToken {
	// 	w.WriteHeader(http.StatusForbidden)
	// 	return
	// }

	// 3. 查询用户信息
	user, err := db.GetUserInfo(username)
	if err != nil {
		fmt.Println(err.Error())
		c.JSON(http.StatusForbidden,gin.H{
			"code":-2,
			"msg":"not found",
		})
		return
	}

	// 4. 组装并且响应用户数据
	resp := util.RespMsg{
		Code: 0,
		Msg:  "OK",
		Data: user,
	}
	c.Data(http.StatusOK,"ContentType/json",resp.JSONBytes())
}









//token
func Token(name string) string{
	// 40位字符:md5(username+timestamp+token_salt)+timestamp[:8]
	md5 := util.MD5([]byte(name + fmt.Sprintf("%x", time.Now().Unix()) + tokensalt))
	return md5+time.Now().Format("20060102150405")

}
func IsTokenValid(token string) bool{
	if len(token)!=46{
		return false
	}
	//if token[33:]
	return true

}

