package meta

import (
	"myyun/db"

)

type FileMeta struct {
	Filehash string
	FileName string
	FileSize int64
	FilePath string
	UploadTime  string
}

var fileMap map[string]FileMeta

func init() {
	fileMap = make(map[string]FileMeta)
}


//new filemeta
func UpdateFilemetaDb(f FileMeta)bool{
	return db.SaveFileToMysql(f.Filehash, f.FileName, f.FileSize, f.FilePath)
}

//get file by hash
func GetFileByHash (h string)(*FileMeta,error){
	file, e := db.GetFileMeta(h)
	if e!=nil||file==nil{
		return nil,e
	}
	filemeta:=FileMeta{
		file.FileHash,
		file.FileName.String,
		file.FileSize.Int64,
		file.FileAddr.String,
		"",
	}
	return &filemeta,nil
}
//change file name

//func ChangeFileName(h string,name string )  {
//	meta := GetFileByHash(h)
//	meta.FileName=name
//}

// GetLastFileMetas : 获取批量的文件元信息列表
func GetLastFileMetas(count int) ([]FileMeta,error) {
	files, e := db.GetFileMetaList(count)
	if e!=nil||files==nil{
		return nil,e
	}
	filemeta:=make([]FileMeta,len(files))
	for _,v:=range files{
		var file FileMeta
		file=FileMeta{
			v.FileHash,
			v.FileName.String,
			v.FileSize.Int64,
			v.FileAddr.String,
			"",
		}
		filemeta=append(filemeta,file)
	}
	return filemeta,nil

}
func Deletefile(h string,name string){
	 db.DeleteUserFile(h, name)
}