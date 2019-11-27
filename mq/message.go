package mq

import "myyun/conf"

type Msg struct {
	FileHash      string
	CurLocation   string
	DestLocation  string
	DestStoreType conf.StoreType
}