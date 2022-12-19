package database

import (
	"database/sql"
	"fmt"
	"os"
	"strconv"

	_ "github.com/go-sql-driver/mysql"
	"github.com/joho/godotenv"
)

var (
	DbConn            *sql.DB
	HOST_DATABASE     string
	USER_DATABASE     string
	PASSWORD_DATABASE string
	PORT_DATABASE     int
)

func init() {
	var err error

	godotenv.Load()

	HOST_DATABASE = os.Getenv("HOST_DATABASE")
	USER_DATABASE = os.Getenv("USER_DATABASE")
	PASSWORD_DATABASE = os.Getenv("PASSWORD_DATABASE")
	PORT_DATABASE, _ = strconv.Atoi(os.Getenv("PORT_DATABASE"))
	DATABASE := "app"

	URI_CONN := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s", USER_DATABASE, PASSWORD_DATABASE, HOST_DATABASE, PORT_DATABASE, DATABASE)

	DbConn, err = sql.Open("mysql", URI_CONN)

	if err != nil {
		fmt.Println("Database connection fail")
	}
}
