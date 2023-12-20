package main

import (
	"database/sql"
	"fmt"
	"os"
	"time"

	_ "github.com/go-sql-driver/mysql"
	"github.com/joho/godotenv"
)

// 在連線的同時創建新的 SQL 資料進去
func CreateNewSqlUser(user User) {

	dataSourceName := loadDataSourceName()

	db, err := sql.Open("mysql", dataSourceName)

	db.SetConnMaxLifetime(time.Minute * 3)
	db.SetMaxOpenConns(10)
	db.SetMaxIdleConns(10)

	checkErr(err)

	tryCreateUsersTable(db)

	insertUserTable(db, user)

	showUsersData(db)

	db.Close()
}

func UpdateSqlUser(user User) {
	fmt.Println("Update sql user.")

	dataSourceName := loadDataSourceName()

	db, err := sql.Open("mysql", dataSourceName)

	db.SetConnMaxLifetime(time.Minute * 3)
	db.SetMaxOpenConns(10)
	db.SetMaxIdleConns(10)

	checkErr(err)

	tryUpdateUserTable(db, user)

	db.Close()

}

// 用來測試 SQL 能不能正常跑動的
func SqlTest() {
	fmt.Println("SqlTest")

	dataSourceName := loadDataSourceName()

	db, err := sql.Open("mysql", dataSourceName)

	db.SetConnMaxLifetime(time.Minute * 3)
	db.SetMaxOpenConns(10)
	db.SetMaxIdleConns(10)

	checkErr(err)

	showTable(db)

	checkErr(err)

	db.Close()
}

func loadDataSourceName() string {

	err := godotenv.Load()
	if err != nil {
		fmt.Println("Error loading .env file")
	}

	dbUser := os.Getenv("MYSQL_ROOT_USER")
	dbPassword := os.Getenv("MYSQL_ROOT_PASSWORD")
	dbHost := os.Getenv("MYSQL_HOST")
	dbPort := os.Getenv("MYSQL_PORT")
	dbName := os.Getenv("MYSQL_DATABASE")

	dataSourceName := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s", dbUser, dbPassword, dbHost, dbPort, dbName)

	return dataSourceName
}

func tryUpdateUserTable(db *sql.DB, user User) {
	stmt, err := db.Prepare("update users set level=? where id=?")
	checkErr(err)

	res, err := stmt.Exec(user.level, user.Addr)
	checkErr(err)

	affect, err := res.RowsAffected()
	checkErr(err)

	fmt.Println(affect)

	showUsersData(db)
}

func tryCreateUsersTable(db *sql.DB) {

	fmt.Println("Create table.")

	_, err := db.Exec("CREATE TABLE IF NOT EXISTS users (id VARCHAR(255) NOT NULL, level int NOT NULL, created DATE NOT NULL)")

	checkErr(err)

	showTable(db)
}

func showTable(db *sql.DB) {
	fmt.Println("Check table.")

	res, err := db.Query("SHOW TABLES")

	checkErr(err)

	var table string

	for res.Next() {
		res.Scan(&table)
		fmt.Println(table)
	}
}

func showUsersData(db *sql.DB) {
	fmt.Println("Show user data.")
	rows, err := db.Query("SELECT * FROM users")
	checkErr(err)

	for rows.Next() {
		var id string
		var level int
		var created string
		err = rows.Scan(&id, &level, &created)
		checkErr(err)
		fmt.Println(id)
		fmt.Println(level)
		fmt.Println(created)
	}
}

func insertUserTable(db *sql.DB, user User) {
	stmt, err := db.Prepare("INSERT users SET id=?,level=?,created=?")
	checkErr(err)

	timestamp := time.Now().Format("2006-01-02")
	fmt.Println(timestamp)

	res, err := stmt.Exec(user.Addr, user.level, timestamp)
	checkErr(err)

	id, err := res.LastInsertId()
	checkErr(err)

	fmt.Println(id)
}

// 不要 user table 時，把整個資料都刪除
func dropTable(db *sql.DB) {
	fmt.Println("Drop table.")
	stmt, err := db.Prepare("DROP TABLE users")
	checkErr(err)

	_, err = stmt.Exec()
	checkErr(err)
}

func checkErr(err error) {
	if err != nil {
		panic(err)
	}
}
