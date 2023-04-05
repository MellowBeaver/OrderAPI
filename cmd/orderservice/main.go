package main

//Using postgresSQL and Go lang APIs for basic CRUD operations.

import (
	"orderapp/pkg/router"

	"github.com/gin-gonic/gin"
	_ "github.com/go-sql-driver/mysql"
	_ "github.com/lib/pq"
)

func main() {

	r := gin.Default()

	router.GinRouter(r)
	//Starting GIN server at port 8000
	r.Run(":8000")
}
