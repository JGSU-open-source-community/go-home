package main

import (
	"database/sql"
	"fmt"
	"log"
	"os"

	"github.com/astaxie/beego/config" // refercen beego config parser
	"github.com/astaxie/beego/orm"
	_ "github.com/go-sql-driver/mysql"
)

var (
	dbType       string
	user         string
	password     string
	host         string
	port         string
	databaseName string
	dns          string
	conf, _      = config.NewConfig("ini", "app.conf")
)

func init() {
	dbType = conf.String("type")
	user = conf.String("user")
	password = conf.String("pass")
	host = conf.String("host")
	port = conf.String("port")
	databaseName = conf.String("databaseName")

	dns = fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8", user, password, host, port, databaseName)

	// auto create database
	createDb()

	orm.RegisterDriver("mysql", orm.DRMySQL)
	orm.RegisterDataBase("default", "mysql", dns)
}

func createDb() {
	isCreate, _ := conf.Bool("isCreate")

	// if database not created!
	if !isCreate {
		createDbSQL := fmt.Sprintf("CREATE DATABASE  if not exists `%s` DEFAULT CHARACTER SET = 'utf8' DEFAULT COLLATE 'utf8_general_ci'", databaseName)

		// connect mysql generate db object
		db, err := sql.Open(dbType, dns)

		defer db.Close()

		if err != nil {
			log.Fatal(err)
			os.Exit(2)
		}

		if _, err := db.Exec(createDbSQL); err != nil {
			log.Fatal(err)
			os.Exit(2)
		}

		if err := conf.Set("isCreate", "1"); err != nil {
			log.Fatal(err)
			os.Exit(2)
		}
	}
}
