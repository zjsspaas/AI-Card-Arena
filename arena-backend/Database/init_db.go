package main

import (
	"database/sql"
	"fmt"
	"log"

	_ "github.com/go-sql-driver/mysql"
)

func main() {
	dsn := "root:power@tcp(127.0.0.1:3306)/?charset=utf8mb4&parseTime=True&loc=Local"

	db, err := sql.Open("mysql", dsn)
	if err != nil {
		log.Fatal("数据库连接失败:", err)
	}
	defer db.Close()

	// 创建数据库 ddz
	sqlStmt := `
	CREATE DATABASE IF NOT EXISTS ddz
	DEFAULT CHARACTER SET utf8mb4
	COLLATE utf8mb4_general_ci;
	`
	_, err = db.Exec(sqlStmt)
	if err != nil {
		log.Fatal("创建数据库失败:", err)
	}
	fmt.Println("数据库 ddz 创建成功 ✅")
}
