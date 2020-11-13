# golang-backend-crud-api

## How to run and test

if want to run in local, uncomment 'DB_HOST=127.0.0.1' and 'TestDbHost=127.0.0.1' in .env file first

### run local
**-go run main.go**

### test
**-open tests folder**\
**-go test -v ./...**

### run docker
**sometimes need to stop postgre server locally first**\
**-docker-compose up**\
**-then stop: docker-compose down**\
**-and run again in background: docker-compose up -d**

## Api Functions:

## Users routes
### /users Methods("POST")
**-create a new user, to create a new account**

### /users Methods("GET")
**-get list all user, to show all user account**

### /users/{id} Methods("GET")
**-get detailed user filter by id, to get user by id**

### /users/{id} Methods("PUT")
**-update user account, to change nickname, email, or user password**\
**-need token/login first**

### /users/{id} Methods("DELETE")
**-delete user account**\
**-need token/login first**

## Products routes
### /product Methods("POST")
**-create a new product with product name, description, price, and stock, and make sure the seller_id is authorized**\
**-need token/login first**

### /product Methods("GET")
**-get list of all save product**

### /product/{id} Methods("GET")
**-get a detailed product filtered id, to show product listed**

### /product/{id} Methods("PUT")
**-for seller to update product info/price**\
**-need token/login first**

### /product/{id} Methods("DELETE")
**-for seller to delete the product listed**\
**-need token/login first**

### /product/find/{id} Methods("GET")
**-to list all product sold by a user( filtered by id)**

## Transactions routes
### /transaction Methods("POST")
**-record a new transaction, one of parameter is qty, product stock will be reduced automatically by the amount of 'qty', another parameter is total_price, its up to frontend to post the amount of total_price (backend wont automatically count the total price)**\
**-need token/login first**

### /transaction Methods("GET")
**get all transaction listed**

### /transaction/{id} Methods("GET")
**show detailed transaction for user, as history transaction**

### /transaction/{id} Methods("PUT")
**update transaction, if there is error made in transaction, front end could update the transaction, the product available_stock is also recalculated based on amount of qty changed**\
**need token/login first**

### /transaction/{id} Methods("DELETE")
**delete a transaction, for user to delete history**\
**need token/login first**

### /transaction/find/{id} Methods("GET")
**show all transaction history of a user**\
**need token/login first**

### /login Methods("POST")
**login to get JWT token**

### detailed API example
https://www.getpostman.com/collections/ca30d151fdc9afc1581d \
or import 'go crud api.postman_collection.json' in postman