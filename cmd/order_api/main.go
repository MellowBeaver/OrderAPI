package main

//Using postgresSQL and Go lang APIs for basic CRUD operations.

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"

	"github.com/gin-gonic/gin"
	_ "github.com/go-sql-driver/mysql"
	_ "github.com/lib/pq"
)

// Update the postgres Database
const (
	user     = "postgres"
	host     = "localhost"
	dbname   = "postgres"
	password = "system"
	port     = 5432
)

var psqlInfo string = fmt.Sprintf("host=%s port=%d user=%s "+"password=%s dbname=%s sslmode=disable", host, port, user, password, dbname)

func main() {

	r := gin.Default()

	//Defining the endpoint
	r.GET("/order/:id", readOrderHandler)
	r.POST("/order", createOrderHandler)
	r.POST("/ordersort", sortOrderHandler)
	r.PUT("/order", updateOrderHandler)
	r.DELETE("/order/:id", deleteOrderHandler)

	//Starting GIN server at port 8000
	r.Run(":8000")
}

//Creating Order Item struct

type OrderItem struct {
	Id    int     `json:"id"`
	Desc  string  `json:"description"`
	Price float32 `json:"price"`
	Qty   int     `json:"quantity"`
}

//Creating Order struct

type Order struct {
	Id           string    `json:"id"`
	Status       string    `json:"status"`
	Items        OrderItem `json:"items"`
	Total        float32   `json:"total"`
	CurrencyUnit string    `json:"currencyUnit"`
}

//This function fetches data based on id value present in orders table

func readOrderHandler(c *gin.Context) {

	//Create postgres connection
	db, err := sql.Open("postgres", psqlInfo)
	fmt.Print("CONNECT")
	if err != nil {
		log.Println(err)
	}

	//Close postgres connection before exiting function
	defer db.Close()

	var invoice Order
	var row *sql.Row

	//Fetching parameter value passed as id
	id := c.Params.ByName("id")

	row = db.QueryRow("select o.id, o.status, oi.id, oi.description, oi.price, oi.quantity, o.total, o.currency_unit from orders o, orderitems oi where o.item_id = oi.id and o.id=$1", id)

	//Check error in row
	err = row.Scan(&invoice.Id, &invoice.Status, &invoice.Items.Id, &invoice.Items.Desc, &invoice.Items.Price, &invoice.Items.Qty, &invoice.Total, &invoice.CurrencyUnit)

	if err != nil {
		log.Println(err)
		c.JSON(500, "Row does not exist")
		return
	}

	//Output order row data
	c.JSON(200, invoice)

}

//This function creates a new order in orders table and orderitems table

func createOrderHandler(c *gin.Context) {

	// Create postgres connection
	db, err := sql.Open("postgres", psqlInfo)
	if err != nil {
		log.Println(err)
	}

	// Close postgres connection before exiting function
	defer db.Close()

	// Getting data from POST request body
	decoder := json.NewDecoder(c.Request.Body)

	var one Order

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
		c.JSON(500, "Error checking orderitems table!")
		return
	}
	fmt.Println(count)

	// Add row to the table orderitems containing foreign key if there is no data
	if count == 0 {
		_, err = db.Exec("INSERT INTO orderitems(id, description, price, quantity) VALUES ($1, $2, $3, $4)", one.Items.Id, one.Items.Desc, one.Items.Price, one.Items.Qty)
		if err != nil {
			log.Println(err)
			c.JSON(500, "Error adding to orderitems table!")
			return
		}
	}

	// Add row to the table orders
	_, err = db.Exec("INSERT INTO orders(id, status, item_id, total, currency_unit) VALUES ($1, $2, $3, $4, $5)", one.Id, one.Status, one.Items.Id, one.Total, one.CurrencyUnit)
	if err != nil {
		log.Println(err)
		c.JSON(500, "Error adding to orders table!")
		return
	}

	// Output row ID
	c.JSON(200, fmt.Sprintf("Added = %v", one.Id))
}

//This function fetches orders from the orders table and sorts them based on the column name provided and ascending/descending order

func sortOrderHandler(c *gin.Context) {

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
		c.JSON(500, "Rows do not exist!!!!")
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
		c.JSON(500, "Rows do not exist!")
		return
	}
	var orders []Order

	for rows.Next() {
		var order Order

		//Check error in row
		err := rows.Scan(&order.Id, &order.Status, &order.Items.Id, &order.Items.Desc, &order.Items.Price, &order.Items.Qty, &order.Total, &order.CurrencyUnit)

		if err != nil {
			log.Println(err)
			c.JSON(500, "Rows do not exist!!")
			return
		}

		orders = append(orders, order)
	}

	//Output sorted rows
	c.JSON(200, orders)

}

//This function updates the status of the order based on id value present in orders table

func updateOrderHandler(c *gin.Context) {

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
		c.JSON(500, gin.H{"error": "Error in updation in table!"})
		return
	}

	//Output row ID
	c.JSON(200, gin.H{"message": fmt.Sprintf("Updated = %v", one.Id)})

}

//This function deletes an order based on id value present in orders table

func deleteOrderHandler(c *gin.Context) {

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
		c.JSON(500, gin.H{"message": "Error in deletion from table!"})
		return
	}

	//Output row ID
	c.JSON(200, gin.H{"message": fmt.Sprintf("Deleted = %v", deleteId)})
}
