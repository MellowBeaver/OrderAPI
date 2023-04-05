package constant

// Update the postgres Database
const (
	User     = "postgres"
	Host     = "localhost"
	Dbname   = "postgres"
	Password = "system"
	Port     = 5432

	// Error messages
	ErrorRowNotExist     = "Row does not exist"
	ErrorCheckOrderItems = "Error checking orderitems table!"
	ErrorAddToOrderItems = "Error adding to orderitems table!"
	ErrorAddToOrders     = "Error adding to orders table!"
	ErrorUpdating        = "Error in updating table"
	ErrorDeleting        = "Error in deleting table"
)
