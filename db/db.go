package db

import (
	"database/sql"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
)

var (
	db  *sql.DB
	err error
)

const (
	MaxCons int = 100
	MinCons int = 2
)

func init() {
	db, err = sql.Open("mysql", "root:root@tcp(127.0.0.1)/shici?charset=utf8&parseTime=true")

	// 处理连接错误
	if err != nil {
		panic(err)
	}

	//defer db.Close() //不能添加改行:数据库被关闭

	//设置最大和最小连接数
	db.SetMaxIdleConns(MaxCons)
	db.SetMaxIdleConns(MinCons)

	err = db.Ping()
	if err != nil {
		panic(err)
	}
}

func checkError(err error) bool {
	if err != nil {
		fmt.Println("错误信息：",err)
		return true
	}
	return false
}
