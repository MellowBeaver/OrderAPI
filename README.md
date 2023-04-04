# OrderAPI

This is a simple API to demonstrate CRUD operations along with sort and filter functionality on a database.

# Usecases

This API has 6 usecases:
1. Create
2. Read
3. Update
4. Sort - takes column name from orders table and ASC/DESC.
5. Filter - takes condition, column name from orders table and ASC/DESC.
6. Delete

## Create

- URL: `localhost:8000/order`
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

- URL: `localhost:8000/order/[id]`
- Here we pass the id from Order table.
- Example:
  - URL: `localhost:8000/order/abcdef-123456`
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

- URL: localhost:8000/order
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
  "Updated = abcdef-123456 "
  ```

## Sort

- URL: `localhost:8000/ordersort`
- Here, we pass column name (id/status/item_id/total/currency_unit from the orders table) and the order type (ASC/DESC).
- Example:
  - Input: JSON
  ```
  {
    "columnname":"total",
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
- URL: `localhost:8000/ordersort`
- Here, we pass column name (id/status/item_id/total/currency_unit from the orders table) and the order type (ASC/DESC).
- Example:
  - Input: JSON
  ```
  {
    "condition": "status='DONE'",
    "columnname":"total",
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
    }
  ]
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
