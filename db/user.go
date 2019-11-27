package db

import "fmt"

// User : 用户表model
type User struct {
	Username     string
	Email        string
	Phone        string
	SignupAt     string
	LastActiveAt string
	Status       int
}

//insert  user info
func InsertUserToMysql(name string,pw string) bool {
	db := NewDB()
	stmt, e := db.Prepare("insert ignore into tbl_user (`user_name`,`user_pwd`) values (?,?)")
	if e!=nil{
		fmt.Println("insert into userinfo prepare failed e:"+e.Error())
		return false
	}
	defer stmt.Close()
	result, e := stmt.Exec(name, pw)
	if e!=nil{
		fmt.Println("insert userinfo failed e:"+e.Error())
		return false
	}
	if affect,e:=result.RowsAffected();e==nil&&affect>=0{
		return true
	}else {
		return false
	}
	
}


//username and password is ok

func UserSignin(name string,pw string)bool  {
	db:=NewDB()
	stmt, e := db.Prepare("select * from tbl_user where user_name = ? limit 1")
	if e!=nil{
		fmt.Println("usersignin prepare failed e:"+e.Error())
		return false
	}
	defer stmt.Close()
	rows, e := stmt.Query(name)
	if e!=nil{
		fmt.Println("usersignin failed e:"+e.Error())
		return false
	}
	if rows==nil{
		fmt.Println("username not exist")
	}
	parseRows := ParseRows(rows)
	if string(parseRows[0]["user_pwd"].([]byte))==pw{
		return true
	}else{
		return false
	}

}


//update token
func UpdateToken (username string,token string)bool{
	db:=NewDB()
	stmt, e := db.Prepare("replace into tbl_user_token (`user_name`,`user_token`) values (?,?)")
	if e!=nil{
		fmt.Println("update token prepare failed e:"+e.Error())
		return false
	}
	defer stmt.Close()
	_, e = stmt.Exec(username, token)
	if e!=nil{
		fmt.Println("update token failed e:"+e.Error())
		return false
	}
	return true
}


//get user info
func GetUserInfo (username string)(User,  error){
	db:=NewDB()
	user:=User{}
	stmt, e := db.Prepare("select user_name,signup_at from tbl_user where user_name=? limit 1")
	if e!=nil{
		return user,e
	}
	defer stmt.Close()
	e = stmt.QueryRow(username).Scan(&user.Username, &user.SignupAt)
	if e!=nil{
		return user,e
	}
	return user,nil

}









