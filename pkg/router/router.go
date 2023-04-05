package router

import (
	"orderapp/pkg/handler"

	"github.com/gin-gonic/gin"
)

func GinRouter(r *gin.Engine) {

	//Defining the endpoint
	r.GET("/order/:id", handler.ReadOrderHandler)
	r.POST("/order", handler.CreateOrderHandler)
	r.POST("/ordersort", handler.SortOrderHandler)
	r.PUT("/order", handler.UpdateOrderHandler)
	r.DELETE("/order/:id", handler.DeleteOrderHandler)

}
