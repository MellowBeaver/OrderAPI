package router

import (
	"database/sql"
	"orderapp/pkg/handler"

	"github.com/gin-gonic/gin"
)

func GinRouter(r *gin.Engine, db *sql.DB) {

	//Defining the endpoint
	r.GET("/order/:id", func(c *gin.Context) {
		handler.ReadOrderHandler(c, db)
	})
	r.POST("/order", func(c *gin.Context) {
		handler.CreateOrderHandler(c, db)
	})
	r.POST("/ordersort", func(c *gin.Context) {
		handler.SortOrderHandler(c, db)
	})
	r.PUT("/order", func(c *gin.Context) {
		handler.UpdateOrderHandler(c, db)
	})
	r.DELETE("/order/:id", func(c *gin.Context) {
		handler.DeleteOrderHandler(c, db)
	})

}
