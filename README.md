# OrderAPI

This is a simple API to demonstrate CRUD operations along with sort and filter functionality on a database.

Test cases are covered for Create, Read, Update and Delete functions and not Sort/Filter function.

# Usecases

This API has 6 use cases:
1. Create
2. Read
3. Update
4. Sort - takes column name from orders and sorts by ASC/DESC.
5. Filter - takes condition, column name from orders and sorts by ASC/DESC.
6. Delete

## Create

- POST `/order`
- Example:
  - Input: JSON 
  ```
  {
      "id": "abcdef-123456",
      "status": "PENDING_INVOICE",
      "items": {
          "id": 123456,
          "description": "a product description",
          "price": 12.40,
          "quantity": 1
      },
      "total": 12.40,
      "currencyUnit": "USD"
  }
  ```
  - Output: `"Added = abcdef-123456"`

## Read 

- GET `/order/[id]`
- Here we pass the id from Order table.
- Example:
  - GET `/order/abcdef-123456`
  - Output:
  ```
  {
      "id": "abcdef-123456",
      "status": "PENDING_INVOICE",
      "items": {
          "id": 123456,
          "description": "a product description",
         "price": 12.4,
          "quantity": 1
      },
      "total": 12.4,
      "currencyUnit": "USD"
  }
  ```

## Update

- POST `/order`
- Example:
  - Input: JSON
  ```
  {
    "id":"abcdef-123456",
    "status":"DONE"
  }
  ```
  - Output:
  ```
  {
    "message": "Updated = abcdef-123456 "
  }
  ```

## Sort

- POST `/ordersort`
- Here, we pass column name (o.id/status/item_id/total/currency_unit from the orders table) and the order type (ASC/DESC).
- Example:
  - Input: JSON
  ```
  {
    "sortingcolumn": "total",
    "orderby":"ASC"
  }
  ```
  - Output:
  ```
  [
    {
        "id": "abcdef-123456",
        "status": "DONE",
        "items": {
            "id": 123456,
            "description": "a product description",
            "price": 12.4,
            "quantity": 1
        },
        "total": 12.4,
        "currencyUnit": "USD"
    },
    {
        "id": "ab-1",
        "status": "DONE",
        "items": {
            "id": 105,
            "description": "yellow",
            "price": 32,
            "quantity": 1
        },
        "total": 32,
        "currencyUnit": "HKD"
    },
    {
        "id": "abc-123",
        "status": "PENDING_INVOICE",
        "items": {
            "id": 125,
            "description": "z products ",
            "price": 32,
            "quantity": 3
        },
        "total": 96,
        "currencyUnit": "INR"
    }
  ]
  ```
## Filter
- POST `/ordersort`
- Here, we pass filter column name (o.id/status/item_id/total/currency_unit from the orders table) and filter value.
- Example:
  - Input: JSON
  ```
  {
    "filtercolumn":"currency_unit",
    "filtervalue":"INR"
  }
  ```
  - Output:
  ```
  [
    {
        "id": "abc-123",
        "status": "PENDING_INVOICE",
        "items": {
            "id": 125,
            "description": "z products ",
            "price": 32,
            "quantity": 3
        },
        "total": 96,
        "currencyUnit": "INR"
    }
  ]
  ```
## Sort and Filter
- POST `/ordersort`
- Here, we pass filter column name (o.id/status/item_id/total/currency_unit from the orders table) and filter value alone with sorting column name (o.id/status/item_id/total/currency_unit from the orders table) and orderby (ASC/DESC).
- Example:
  - Input: JSON
  ```
  {
    "filtercolumn":"status",
    "filtervalue":"DONE",
    "sortingcolumn": "total",
    "orderby":"DESC"
  }
  ```
  - Output:
  ```
  [
    {
        "id": "ab-1",
        "status": "DONE",
        "items": {
            "id": 105,
            "description": "yellow",
            "price": 32,
            "quantity": 1
        },
        "total": 32,
        "currencyUnit": "HKD"
    },
    {
        "id": "abcdef-123456",
        "status": "DONE",
        "items": {
            "id": 123456,
            "description": "a product description",
            "price": 12.4,
            "quantity": 1
        },
        "total": 12.4,
        "currencyUnit": "USD"
    }
  ]
  ```
## Delete
- DELETE `/order/[id]`
- Here we pass the id from Order table.
- Example:
  - DELETE `/order/ab-1`
  - Output:
  ```
  {
    "message": "Deleted = ab-1"
  }
  ```

# Steps to run the API

## Steps to run DB using Docker

- Install docker on system.
- Run `docker compose up`.
- Run `docker ps to check if container is running`.

## Steps to connect to DB container

- Run `docker exec -it [container name] psql -U [username] -W [dbname]`.
- This will open postgreSQL for you to run query to create databases and tables.

### CREATE DATABASE ORDERITEMS

```
create table orderitems (
  id int primary key, 
  description varchar(100), 
  price numeric, 
  quantity int
);
```

### CREATE DATABASE ORDERS

```
create table orders (
  id varchar(15) not null, 
  status varchar(25), 
  item_id int,
  total numeric, 
  currency_unit varchar(10),
  FOREIGN KEY (item_id) REFERENCES orderitems (id) ON DELETE CASCADE
);
```

## Steps to run the app

- First download dependencies, run `go get ./...`
- Finally run `go run main.go`

## Steps to run the test file

- Run `go test -run main_test.go`
