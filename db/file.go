package db

import (
	"database/sql"
	"fmt"
)

//SAVE FILE TO MYSQL
func SaveFileToMysql(hash string,name string,size int64,path string)bool{
	db:=NewDB()
	stmt, e := db.Prepare("insert ignore into tbl_file (`file_sha1`,`file_name`,`file_size`," +
		"`file_addr`,`status`) values (?,?,?,?,1)")
	if e!=nil{
		fmt.Println("sql prepare fail e:"+e.Error())
		return false
	}
	defer stmt.Close()
	exec,e := stmt.Exec(hash, name, size, path)
	if e!=nil{
		fmt.Println("insert into file fail e:",e.Error())
		return false
	}
	i, _ := exec.RowsAffected()
	if i<=0{
		fmt.Println("this file was uploaded")
		return true
	}
	return true
}

// TableFile : 文件表结构体
type TableFile struct {
	FileHash string
	FileName sql.NullString
	FileSize sql.NullInt64
	FileAddr sql.NullString
}

// GetFileMeta : 从mysql获取文件元信息
func GetFileMeta(filehash string)(*TableFile,error){
	db:=NewDB()
	tableFile:=TableFile{}
	stmt, e := db.Prepare("select file_sha1,file_addr,file_name,file_size from tbl_file " +
		"where file_sha1=? and status=1 limit 1 ")
	if e!=nil{
		return &tableFile,e
	}
	defer stmt.Close()
	e = stmt.QueryRow(filehash).Scan(tableFile.FileHash, tableFile.FileAddr, tableFile.FileName, tableFile.FileSize)
	if e!=nil{
		return &tableFile,e
	}
	return &tableFile,nil
}


// GetFileMetaList : 从mysql批量获取文件元信息
func GetFileMetaList(limit int)([]TableFile,error){
	db:=NewDB()
	stmt,e:=db.Prepare("select file_sha1,file_addr,file_name,file_size from tbl_file " +
		"where status=1 limit ?")
	if e!=nil{
		return nil,e
	}
	defer stmt.Close()
	rows, e := stmt.Query(limit)
	if e!=nil{
		return nil,e
	}
	var tablefiles []TableFile
	cloumns, _ := rows.Columns()
	values := make([]sql.RawBytes, len(cloumns))
	for i:=0;i<len(values)&&rows.Next();i++{
		var tablefile TableFile
		e := rows.Scan(tablefile.FileHash, tablefile.FileAddr, tablefile.FileName, tablefile.FileSize)
		if e != nil {
			fmt.Println(e.Error())
			break
		}
		tablefiles = append(tablefiles, tablefile)
	}
	return tablefiles,nil
}

// UpdateFileLocation : 更新文件的存储地址(如文件被转移了)
func UpdateFileLocation (hash string,addr string)bool{
	db:=NewDB()
	stmt, e := db.Prepare("update tbl_file set file_addr = ? where file_hash = ?")
	if e!=nil{
		return false
	}
	defer stmt.Close()
	exec, e := stmt.Exec(addr, hash)
	if e != nil {
		fmt.Println(e.Error())
		return false
	}
	if rf, err := exec.RowsAffected(); nil == err {
		if rf <= 0 {
			fmt.Printf("更新文件location失败, filehash:%s", hash)
		}
		return true
	}
	return false
}

