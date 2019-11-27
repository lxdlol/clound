package db

import (
	"fmt"
	"time"
)

// UserFile : 用户文件表结构体
type UserFile struct {
	UserName    string
	FileHash    string
	FileName    string
	FileSize    int64
	UploadAt    string
	LastUpdated string
}



// OnUserFileUploadFinished : 更新用户文件表
func OnUserFileUploadFinished(username,filename ,hash string,size int64)bool{
	db:=NewDB()
	stmt, e := db.Prepare("insert into tbl_user_file (`user_name`,`file_sha1`,`file_name`," +
		"`file_size`,`upload_at`) values (?,?,?,?,?)")
	if e!=nil{
		return false
	}
	defer stmt.Close()
	_, e = stmt.Exec(username, hash, filename, size, time.Now())
	if e!=nil{
		fmt.Println(e)
		return false
	}
	return true
}


// QueryUserFileMetas : 批量获取用户文件信息
func QueryUserFileMetas(username string,limit int)([]UserFile,error){
	db:=NewDB()
	stmt, e := db.Prepare("select file_sha1,file_name,file_size,upload_at," +
		"last_update from tbl_user_file where user_name=? limit ? ")
	if e!=nil{
		return nil,e
	}
	defer stmt.Close()
	rows, e := stmt.Query(username, limit)
	if e!=nil{
		return nil,e
	}
	var userfiles []UserFile
	for rows.Next(){
		var userfile UserFile
		e := rows.Scan(&userfile.FileHash, &userfile.FileName, &userfile.FileSize, &userfile.UploadAt,&userfile.LastUpdated)
		if e!=nil{
			break
		}
		userfiles=append(userfiles,userfile)
	}
	fmt.Println(userfiles)
	return userfiles,nil
}
//delete userfile
func DeleteUserFile(username string,hash string)bool{
	db:=NewDB()
	stmt, e := db.Prepare("delete from tbl_user_file where user_name = ? and file_sha1 = ?")
	if e!=nil{
		return false
	}
	defer stmt.Close()
	exec, e := stmt.Exec(username, hash)
	if e!=nil{
		return false
	}
	i, e := exec.RowsAffected()
	if e!=nil||i<=0{
		return false
	}
	return true

}