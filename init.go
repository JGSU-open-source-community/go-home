package main

import (
	"database/sql"

	"github.com/astaxie/beego/orm"
	_ "github.com/go-sql-driver/mysql"
)

func init() {
	borm.RegisterDriver("mysql", orm.DRMySQL)
	orm.RegisterDataBase("default", "mysql", "root:digitalx168@tcp(127.0.0.1:3306)/12306charset=utf8mb4")
}
