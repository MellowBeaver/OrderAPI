package handler

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"orderapp/pkg/constant"
	"orderapp/pkg/model"

	"github.com/gin-gonic/gin"
)

var psqlInfo string = fmt.Sprintf("host=%s port=%d user=%s "+"password=%s dbname=%s sslmode=disable", constant.Host, constant.Port, constant.User, constant.Password, constant.Dbname)

//This function fetches data based on id value present in orders table

func ReadOrderHandler(c *gin.Context) {

	//Create postgres connection
	db, err := sql.Open("postgres", psqlInfo)
	fmt.Print("CONNECT")
	if err != nil {
		log.Println(err)
	}

	//Close postgres connection before exiting function
	defer db.Close()

	var invoice model.Order
	var row *sql.Row

	//Fetching parameter value passed as id
	id := c.Params.ByName("id")

	row = db.QueryRow("select o.id, o.status, oi.id, oi.description, oi.price, oi.quantity, o.total, o.currency_unit from orders o, orderitems oi where o.item_id = oi.id and o.id=$1", id)

	//Check error in row
	err = row.Scan(&invoice.Id, &invoice.Status, &invoice.Items.Id, &invoice.Items.Desc, &invoice.Items.Price, &invoice.Items.Qty, &invoice.Total, &invoice.CurrencyUnit)

	if err != nil {
		log.Println(err)
		c.JSON(http.StatusInternalServerError, constant.ErrorRowNotExist)
		return
	}

	//Output order row data
	c.JSON(http.StatusOK, invoice)

}

//This function creates a new order in orders table and orderitems table

func CreateOrderHandler(c *gin.Context) {

	// Create postgres connection
	db, err := sql.Open("postgres", psqlInfo)
	if err != nil {
		log.Println(err)
	}

	// Close postgres connection before exiting function
	defer db.Close()

	// Getting data from POST request body
	decoder := json.NewDecoder(c.Request.Body)

	var one model.Order

	err = decoder.Decode(&one)
	if err != nil {
		log.Println(err)
		return
	}

	// Check if there is data in orderitems table with specified id
	var count int
	err = db.QueryRow("SELECT count(*) FROM orderitems WHERE id=$1", one.Items.Id).Scan(&count)
	if err != nil {
		log.Println(err)
		c.JSON(http.StatusInternalServerError, constant.ErrorCheckOrderItems)
		return
	}
	fmt.Println(count)

	// Add row to the table orderitems containing foreign key if there is no data
	if count == 0 {
		_, err = db.Exec("INSERT INTO orderitems(id, description, price, quantity) VALUES ($1, $2, $3, $4)", one.Items.Id, one.Items.Desc, one.Items.Price, one.Items.Qty)
		if err != nil {
			log.Println(err)
			c.JSON(http.StatusInternalServerError, constant.ErrorAddToOrderItems)
			return
		}
	}

	// Add row to the table orders
	_, err = db.Exec("INSERT INTO orders(id, status, item_id, total, currency_unit) VALUES ($1, $2, $3, $4, $5)", one.Id, one.Status, one.Items.Id, one.Total, one.CurrencyUnit)
	if err != nil {
		log.Println(err)
		c.JSON(http.StatusInternalServerError, constant.ErrorAddToOrders)
		return
	}

	// Output row ID
	c.JSON(http.StatusOK, fmt.Sprintf("Added = %v", one.Id))
}

//This function fetches orders from the orders table and sorts them based on the column name provided and ascending/descending order

func SortOrderHandler(c *gin.Context) {

	//Create postgres connection
	db, err := sql.Open("postgres", psqlInfo)
	fmt.Print("CONNECT")
	if err != nil {
		log.Println(err)
	}

	//Close postgres connection before exiting function
	defer db.Close()

	// Getting data from POST request body
	decoder := json.NewDecoder(c.Request.Body)

	type SortOrder struct {
		FilterColumn  string
		FilterValue   interface{}
		SortingColumn string
		OrderBy       string
	}

	var one SortOrder
	var query string

	//Getting condition, column name and order from JSON body
	err = decoder.Decode(&one)
	if err != nil {
		log.Println(err)
		c.JSON(http.StatusInternalServerError, constant.ErrorRowNotExist)
		return
	}

	if one.OrderBy == "" {
		one.OrderBy = "ASC"
	}
	if one.SortingColumn == "" {
		one.SortingColumn = "o.id"
	}

	//Check if a condition is provided to filter the order table
	if one.FilterColumn != "" {
		//query with filter and sort function
		query = fmt.Sprintf("SELECT o.id, o.status, o.item_id, oi.description, oi.price, oi.quantity, o.total, o.currency_unit FROM orders o INNER JOIN orderitems oi on o.item_id=oi.id WHERE %s = '%v' ORDER BY %s %s", one.FilterColumn, one.FilterValue, one.SortingColumn, one.OrderBy)
	} else {
		//query with sort function
		query = fmt.Sprintf("SELECT o.id, o.status, o.item_id, oi.description, oi.price, oi.quantity, o.total, o.currency_unit FROM orders o INNER JOIN orderitems oi on o.item_id=oi.id ORDER BY %s %s", one.SortingColumn, one.OrderBy)
	}

	fmt.Println(query)
	rows, err := db.Query(query)
	if err != nil {
		log.Println(err)
		c.JSON(http.StatusInternalServerError, constant.ErrorRowNotExist)
		return
	}
	var orders []model.Order

	for rows.Next() {
		var order model.Order

		//Check error in row
		err := rows.Scan(&order.Id, &order.Status, &order.Items.Id, &order.Items.Desc, &order.Items.Price, &order.Items.Qty, &order.Total, &order.CurrencyUnit)

		if err != nil {
			log.Println(err)
			c.JSON(http.StatusInternalServerError, constant.ErrorRowNotExist)
			return
		}

		orders = append(orders, order)
	}

	//Output sorted rows
	c.JSON(http.StatusOK, orders)

}

//This function updates the status of the order based on id value present in orders table

func UpdateOrderHandler(c *gin.Context) {

	//Create postgres connection
	db, err := sql.Open("postgres", psqlInfo)
	if err != nil {
		log.Println(err)
	}

	//Close postgres connection before exiting function
	defer db.Close()

	// Getting data from PUT request body
	decoder := json.NewDecoder(c.Request.Body)

	type UpdateStatus struct {
		Id     string
		Status string
	}

	var one UpdateStatus

	err = decoder.Decode(&one)
	if err != nil {
		log.Println(err)
		return
	}

	//Updating existing row status
	_, err = db.Exec("Update orders set status = $1 where id = $2", one.Status, one.Id)

	if err != nil {
		log.Println(err)
		c.JSON(http.StatusInternalServerError, constant.ErrorUpdating)
		return
	}

	//Output row ID
	c.JSON(http.StatusOK, gin.H{"message": fmt.Sprintf("Updated = %v", one.Id)})

}

//This function deletes an order based on id value present in orders table

func DeleteOrderHandler(c *gin.Context) {

	//Create postgres connection
	db, err := sql.Open("postgres", psqlInfo)
	if err != nil {
		log.Println(err)
	}

	//Close postgres connection before exiting function
	defer db.Close()

	//Fetching parameter value passed as id
	deleteId := c.Params.ByName("id")

	//Delete row
	_, err = db.Exec("Delete from orders where id = $1", deleteId)
	if err != nil {
		log.Println(err)
		c.JSON(http.StatusInternalServerError, constant.ErrorDeleting)
		return
	}

	//Output row ID
	c.JSON(http.StatusOK, gin.H{"message": fmt.Sprintf("Deleted = %v", deleteId)})
}
