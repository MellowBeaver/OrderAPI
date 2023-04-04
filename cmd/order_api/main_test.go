package main

import (
	"bytes"
	"database/sql"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

// Mock DB connection for testing purposes
func MockDB() (*sql.DB, error) {
	db, err := sql.Open("postgres", "host=localhost port=5432 user=postgres password=system dbname=postgres sslmode=disable")
	if err != nil {
		return nil, err
	}

	return db, nil
}

func TestCreateOrderHandler(t *testing.T) {

	// Create a new gin router with the sortOrderHandler as the endpoint
	r := gin.New()
	r.POST("/create", createOrderHandler)

	// Set up the test request payload
	order := map[string]interface{}{
		"id":     "abcdef-123456",
		"status": "PENDING_INVOICE",
		"items": []map[string]interface{}{
			{
				"id":          "123456",
				"description": "a product description",
				"price":       12.40,
				"quantity":    1,
			},
		},
		"total":        12.40,
		"currencyUnit": "USD",
	}
	payload, err := json.Marshal(order)
	assert.NoError(t, err)

	// Set up the test request
	req, err := http.NewRequest("POST", "/orders", bytes.NewReader(payload))
	assert.NoError(t, err)

	// Set up the test response recorder
	w := httptest.NewRecorder()

	// Set up the test context and call the handler function
	c, _ := gin.CreateTestContext(w)
	c.Request = req
	createOrderHandler(c)

	// Verify the response status code and body
	assert.Equal(t, http.StatusOK, w.Code)
}

func TestReadOrderHandler(t *testing.T) {
	// Create a new HTTP request
	req, err := http.NewRequest("GET", "/read/abcdef-123456", nil)
	if err != nil {
		t.Fatal(err)
	}

	// Create a new recorder to record the HTTP response
	rr := httptest.NewRecorder()

	// Create a new router instance
	r := gin.New()
	r.GET("/read/:id", readOrderHandler)

	// Serve the HTTP request to the recorder
	r.ServeHTTP(rr, req)

	// Check the HTTP status code
	if status := rr.Code; status != http.StatusOK {
		t.Errorf("Unexpected status code, expected %v but got %v", http.StatusOK, status)
	}

	// Check the HTTP response body
	expected := `{"id":"abcdef-123456","status":"PENDING_INVOICE","items":{"id":123456,"description":"a product description","price":12.4,"quantity":1},"total":12.4,"currencyUnit":"USD"}`
	if rr.Body.String() != expected {
		t.Errorf("Unexpected response body, got %v but expected %v", rr.Body.String(), expected)
	}
}

func TestUpdateOrderHandler(t *testing.T) {
	// Set up test router and handler
	r := gin.Default()

	r.PUT("/orders/:id", updateOrderHandler)

	// Mock request payload
	payload := `{
        "id": "abcdef-123456",
        "status": "DONE"
    }`

	// Convert payload string to bytes for request
	reqBody := bytes.NewBuffer([]byte(payload))

	// Create a mock request to use for testing
	req, err := http.NewRequest("PUT", "/orders/abcdef-123456", reqBody)
	if err != nil {
		t.Fatalf("Failed to create request: %v", err)
	}

	// Create a mock response recorder
	rr := httptest.NewRecorder()

	// Get a mock database connection for testing
	db, err := MockDB()
	if err != nil {
		t.Fatalf("Failed to connect to mock database: %v", err)
	}

	defer db.Close()

	// Bind the database connection to the context
	r.Use(func(c *gin.Context) {
		c.Set("db", db)
		c.Next()
	})

	// Send the request to the server
	r.ServeHTTP(rr, req)

	// Check the response status code
	if status := rr.Code; status != http.StatusOK {
		t.Errorf("Unexpected response code: got %v, want %v", status, http.StatusOK)
	}

	// Check the response body
	expectedBody := `{"message":"Updated = abcdef-123456"}`
	if strings.TrimSpace(rr.Body.String()) != expectedBody {
		t.Errorf("Unexpected response body: got %v, want %v", rr.Body.String(), expectedBody)
	}

	// Check that the row was updated in the database
	var status string
	err = db.QueryRow("SELECT status FROM orders WHERE id='abcdef-123456'").Scan(&status)
	if err != nil {
		t.Fatalf("Failed to query database: %v", err)
	}

	if status != "COMPLETED" {
		t.Errorf("Row was not updated in database: got %v, want COMPLETED", status)
	}
}

func TestDeleteOrderHandler(t *testing.T) {
	// Setup test router
	r := gin.Default()
	r.DELETE("/orders/:id", deleteOrderHandler)

	// Create test request
	req, err := http.NewRequest("DELETE", "/orders/abcdef-123456", nil)
	if err != nil {
		t.Fatal(err)
	}

	// Create a response recorder to record the response
	rr := httptest.NewRecorder()

	// Create a mock database connection and inject it into the context
	db, err := sql.Open("postgres", psqlInfo)
	if err != nil {
		t.Fatalf("Failed to create database connection: %v", err)
	}
	defer db.Close()
	ctx := gin.Context{}
	ctx.Set("db", db)

	// Call the deleteOrderHandler function and record the response
	r.ServeHTTP(rr, req)

	// Assert that the response status code is 200 OK
	assert.Equal(t, http.StatusOK, rr.Code)

	// Assert that the response body is the expected message
	expected := "{\"message\":\"Deleted = abcdef-123456\"}"
	assert.Equal(t, expected, strings.TrimSpace(rr.Body.String()))

	// Assert that the order has been deleted from the database
	var count int
	err = db.QueryRow("SELECT COUNT(*) FROM orders WHERE id=$1", "abcdef-123456").Scan(&count)
	assert.NoError(t, err)
	assert.Equal(t, 0, count)

}
