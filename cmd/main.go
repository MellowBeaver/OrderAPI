package main

//Using postgresSQL and Go lang APIs for basic CRUD operations.

import (
	"database/sql"
	"fmt"
	"log"
	"orderapp/pkg/constant"
	"orderapp/pkg/router"

	"github.com/gin-gonic/gin"
	_ "github.com/go-sql-driver/mysql"
	_ "github.com/lib/pq"
)

var psqlInfo string = fmt.Sprintf("host=%s port=%d user=%s "+"password=%s dbname=%s sslmode=disable", constant.Host, constant.Port, constant.User, constant.Password, constant.Dbname)
var db *sql.DB

func main() {

	r := gin.Default()
	var err error

	// Create postgres connection pool
	db, err = sql.Open("postgres", psqlInfo)
	if err != nil {
		log.Fatalf("Error creating connection pool: %s", err.Error())
	}

	// Set maximum number of idle connections in the pool.
	// This helps to reuse idle connections instead of creating new ones.
	db.SetMaxIdleConns(10)

	// Set maximum number of open connections in the pool.
	db.SetMaxOpenConns(100)

	// Validate the connection pool
	err = db.Ping()
	if err != nil {
		log.Fatalf("Error validating connection pool: %s", err.Error())
	}

	router.GinRouter(r, db)
	//Starting GIN server at port 8000
	r.Run(":8000")
}
