# OrderAPI

This is a simple API to demonstrate CRUD operations alongwith sort functionality on a database.

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
