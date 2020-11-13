package controllertests

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/http/httptest"
	"strconv"
	"testing"

	"github.com/bayuajik2403/golang-backend-crud-api/api/models"
	"github.com/gorilla/mux"
	"gopkg.in/go-playground/assert.v1"
)

func TestCreateTransaction(t *testing.T) {

	err := refreshAllTable()
	if err != nil {
		log.Fatal(err)
	}
	user, product, _, err := seedOneAllTable()
	if err != nil {
		log.Fatalf("Cannot seed user %v\n", err)
	}
	token, err := server.SignIn(user.Email, "password") //Note the password in the database is already hashed, we want unhashed
	if err != nil {
		log.Fatalf("cannot login: %v\n", err)
	}
	tokenString := fmt.Sprintf("Bearer %v", token)

	samples := []struct {
		inputJSON    string
		statusCode   int
		product_id   uint64
		qty          uint32
		total_price  uint32
		buyer_id     uint64
		tokenGiven   string
		errorMessage string
	}{
		{
			inputJSON:    `{"product_id":1,"qty":3,"total_price":45000,"buyer_id":1}`,
			statusCode:   201,
			tokenGiven:   tokenString,
			product_id:   product.ID,
			qty:          3,
			total_price:  45000,
			buyer_id:     user.ID,
			errorMessage: "",
		},
		{
			// When no token is passed
			inputJSON:    `{"product_id":1,"qty":3,"total_price":45000,"buyer_id":1}`,
			statusCode:   401,
			tokenGiven:   "",
			errorMessage: "Unauthorized",
		},
		{
			// When incorrect token is passed
			inputJSON:    `{"product_id":1,"qty":3,"total_price":45000,"buyer_id":1}`,
			statusCode:   401,
			tokenGiven:   "This is an incorrect token",
			errorMessage: "Unauthorized",
		},
		// {
		// 	inputJSON:    `{"product_id":0,"qty":3,"total_price":45000,"buyer_id":1}`,
		// 	statusCode:   422,
		// 	tokenGiven:   tokenString,
		// 	errorMessage: "Required Product Id",
		// },
		{
			inputJSON:    `{"product_id":1,"qty":0,"total_price":45000,"buyer_id":1}`,
			statusCode:   422,
			tokenGiven:   tokenString,
			errorMessage: "Required Qty",
		},
		{
			inputJSON:    `{"product_id":1,"qty":3,"total_price":45000,"buyer_id":0}`,
			statusCode:   422,
			tokenGiven:   tokenString,
			errorMessage: "Required Buyer ID",
		},
		{
			// When user 2 uses user 1 token
			inputJSON:    `{"product_id":1,"qty":3,"total_price":45000,"buyer_id":2}`,
			statusCode:   401,
			tokenGiven:   tokenString,
			errorMessage: "Unauthorized",
		},
	}
	for _, v := range samples {

		req, err := http.NewRequest("POST", "/transactions", bytes.NewBufferString(v.inputJSON))
		if err != nil {
			t.Errorf("this is the error: %v\n", err)
		}
		rr := httptest.NewRecorder()
		handler := http.HandlerFunc(server.CreateTransaction)

		req.Header.Set("Authorization", v.tokenGiven)
		handler.ServeHTTP(rr, req)

		responseMap := make(map[string]interface{})
		err = json.Unmarshal([]byte(rr.Body.String()), &responseMap)
		if err != nil {
			fmt.Printf("Cannot convert to json: %v", err)
		}
		assert.Equal(t, rr.Code, v.statusCode)
		if v.statusCode == 201 {
			assert.Equal(t, responseMap["product_id"], float64(v.product_id))
			assert.Equal(t, responseMap["qty"], float64(v.qty))
			assert.Equal(t, responseMap["total_price"], float64(v.total_price))
			assert.Equal(t, responseMap["buyer_id"], float64(v.buyer_id)) //just for both ids to have the same type
		}
		if v.statusCode == 401 || v.statusCode == 422 || v.statusCode == 500 && v.errorMessage != "" {
			assert.Equal(t, responseMap["error"], v.errorMessage)
		}
	}
}

func TestGetTransactions(t *testing.T) {

	err := refreshAllTable()
	if err != nil {
		log.Fatal(err)
	}
	_, _, _, err = seedAllTable()
	if err != nil {
		log.Fatal(err)
	}

	req, err := http.NewRequest("GET", "/transactions", nil)
	if err != nil {
		t.Errorf("this is the error: %v\n", err)
	}
	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(server.GetTransactions)
	handler.ServeHTTP(rr, req)

	var transactions []models.Transaction
	err = json.Unmarshal([]byte(rr.Body.String()), &transactions)

	assert.Equal(t, rr.Code, http.StatusOK)
	assert.Equal(t, len(transactions), 2)
}
func TestGetTransactionByID(t *testing.T) {

	err := refreshAllTable()
	if err != nil {
		log.Fatal(err)
	}
	_, _, transaction, err := seedOneAllTable()
	if err != nil {
		log.Fatal(err)
	}
	transactionSample := []struct {
		id           string
		statusCode   int
		product_id   uint64
		qty          uint32
		total_price  uint32
		buyer_id     uint64
		errorMessage string
	}{
		{
			id:          strconv.Itoa(int(transaction.ID)),
			statusCode:  200,
			product_id:  transaction.ProductID,
			qty:         transaction.Qty,
			total_price: transaction.TotalPrice,
			buyer_id:    transaction.BuyerID,
		},
		{
			id:         "unknwon",
			statusCode: 400,
		},
	}
	for _, v := range transactionSample {

		req, err := http.NewRequest("GET", "/transactions", nil)
		if err != nil {
			t.Errorf("this is the error: %v\n", err)
		}
		req = mux.SetURLVars(req, map[string]string{"id": v.id})

		rr := httptest.NewRecorder()
		handler := http.HandlerFunc(server.GetTransaction)
		handler.ServeHTTP(rr, req)

		responseMap := make(map[string]interface{})
		err = json.Unmarshal([]byte(rr.Body.String()), &responseMap)
		if err != nil {
			log.Fatalf("Cannot convert to json: %v", err)
		}
		assert.Equal(t, rr.Code, v.statusCode)

		if v.statusCode == 200 {
			assert.Equal(t, responseMap["product_id"], float64(transaction.ProductID))
			assert.Equal(t, responseMap["qty"], float64(transaction.Qty))
			assert.Equal(t, responseMap["total_price"], float64(transaction.TotalPrice))
			assert.Equal(t, responseMap["buyer_id"], float64(transaction.BuyerID)) //just for both ids to have the same type
		}
	}
}

func TestUpdateTransaction(t *testing.T) {

	var TransactionUserEmail, TransactionUserPassword string
	var AuthTransactionBuyerID uint64
	var AuthTransactionID uint64

	err := refreshAllTable()
	if err != nil {
		log.Fatal(err)
	}
	users, _, transactions, err := seedAllTable()
	if err != nil {
		log.Fatal(err)
	}
	// Get only the first user
	for _, user := range users {
		if user.ID == 2 {
			continue
		}
		TransactionUserEmail = user.Email
		TransactionUserPassword = "password" //Note the password in the database is already hashed, we want unhashed
	}
	//Login the user and get the authentication token
	token, err := server.SignIn(TransactionUserEmail, TransactionUserPassword)
	if err != nil {
		log.Fatalf("cannot login: %v\n", err)
	}
	tokenString := fmt.Sprintf("Bearer %v", token)

	// Get only the first transaction
	for _, transaction := range transactions {
		if transaction.ID == 2 {
			continue
		}
		AuthTransactionID = transaction.ID
		AuthTransactionBuyerID = transaction.BuyerID
	}
	// fmt.Printf("this is the auth transaction: %v\n", AuthTransactionID)

	samples := []struct {
		id           string
		updateJSON   string
		statusCode   int
		product_id   uint64
		qty          uint32
		total_price  uint32
		buyer_id     uint64
		tokenGiven   string
		errorMessage string
	}{
		{
			// Convert int64 to int first before converting to string
			id:           strconv.Itoa(int(AuthTransactionID)),
			updateJSON:   `{"product_id":1,"qty":3,"total_price":45000,"buyer_id":1}`,
			statusCode:   200,
			product_id:   1,
			qty:          3,
			total_price:  45000,
			buyer_id:     AuthTransactionBuyerID,
			tokenGiven:   tokenString,
			errorMessage: "",
		},
		{
			// When no token is provided
			id:           strconv.Itoa(int(AuthTransactionID)),
			updateJSON:   `{"product_id":1,"qty":3,"total_price":45000,"buyer_id":1}`,
			tokenGiven:   "",
			statusCode:   401,
			errorMessage: "Unauthorized",
		},
		{
			// When incorrect token is provided
			id:           strconv.Itoa(int(AuthTransactionID)),
			updateJSON:   `{"product_id":1,"qty":3,"total_price":45000,"buyer_id":1}`,
			tokenGiven:   "this is an incorrect token",
			statusCode:   401,
			errorMessage: "Unauthorized",
		},
		{
			id:           strconv.Itoa(int(AuthTransactionID)),
			updateJSON:   `{"product_id":0,"qty":3,"total_price":45000,"buyer_id":1}`,
			statusCode:   422,
			tokenGiven:   tokenString,
			errorMessage: "Required ProductID",
		},
		{
			id:           strconv.Itoa(int(AuthTransactionID)),
			updateJSON:   `{"product_id":1,"qty":0,"total_price":45000,"buyer_id":1}`,
			statusCode:   422,
			tokenGiven:   tokenString,
			errorMessage: "Required Qty",
		},
		{
			id:         "unknwon",
			statusCode: 400,
		},
		{
			id:           strconv.Itoa(int(AuthTransactionID)),
			updateJSON:   `{"product_id":1,"qty":3,"total_price":45000,"buyer_id":2}`,
			tokenGiven:   tokenString,
			statusCode:   401,
			errorMessage: "Unauthorized",
		},
	}

	for _, v := range samples {

		req, err := http.NewRequest("POST", "/transactions", bytes.NewBufferString(v.updateJSON))
		if err != nil {
			t.Errorf("this is the error: %v\n", err)
		}
		req = mux.SetURLVars(req, map[string]string{"id": v.id})
		rr := httptest.NewRecorder()
		handler := http.HandlerFunc(server.UpdateTransaction)

		req.Header.Set("Authorization", v.tokenGiven)

		handler.ServeHTTP(rr, req)

		responseMap := make(map[string]interface{})
		err = json.Unmarshal([]byte(rr.Body.String()), &responseMap)
		if err != nil {
			t.Errorf("Cannot convert to json: %v", err)
		}
		assert.Equal(t, rr.Code, v.statusCode)
		if v.statusCode == 200 {
			assert.Equal(t, responseMap["product_id"], float64(v.product_id))
			assert.Equal(t, responseMap["qty"], float64(v.qty))
			assert.Equal(t, responseMap["total_price"], float64(v.total_price))
			assert.Equal(t, responseMap["buyer_id"], float64(v.buyer_id)) //just for both ids to have the same type
		}
		if v.statusCode == 401 || v.statusCode == 422 || v.statusCode == 500 && v.errorMessage != "" {
			assert.Equal(t, responseMap["error"], v.errorMessage)
		}
	}
}

func TestDeleteTransaction(t *testing.T) {

	var TransactionUserEmail, TransactionUserPassword string
	var TransactionBuyerID uint64
	var AuthTransactionID uint64

	err := refreshAllTable()
	if err != nil {
		log.Fatal(err)
	}
	users, _, transactions, err := seedAllTable()
	if err != nil {
		log.Fatal(err)
	}
	//Let's get only the Second user
	for _, user := range users {
		if user.ID == 1 {
			continue
		}
		TransactionUserEmail = user.Email
		TransactionUserPassword = "password" //Note the password in the database is already hashed, we want unhashed
	}
	//Login the user and get the authentication token
	token, err := server.SignIn(TransactionUserEmail, TransactionUserPassword)
	if err != nil {
		log.Fatalf("cannot login: %v\n", err)
	}
	tokenString := fmt.Sprintf("Bearer %v", token)

	// Get only the second transaction
	for _, transaction := range transactions {
		if transaction.ID == 1 {
			continue
		}
		AuthTransactionID = transaction.ID
		TransactionBuyerID = transaction.BuyerID
	}
	transactionSample := []struct {
		id           string
		statusCode   int
		buyer_id     uint64
		tokenGiven   string
		errorMessage string
	}{
		{
			// Convert int64 to int first before converting to string
			id:           strconv.Itoa(int(AuthTransactionID)),
			buyer_id:     TransactionBuyerID,
			tokenGiven:   tokenString,
			statusCode:   200,
			errorMessage: "",
		},
		{
			// When empty token is passed
			id:           strconv.Itoa(int(AuthTransactionID)),
			buyer_id:     TransactionBuyerID,
			tokenGiven:   "",
			statusCode:   401,
			errorMessage: "Unauthorized",
		},
		{
			// When incorrect token is passed
			id:           strconv.Itoa(int(AuthTransactionID)),
			buyer_id:     TransactionBuyerID,
			tokenGiven:   "This is an incorrect token",
			statusCode:   401,
			errorMessage: "Unauthorized",
		},
		{
			id:         "unknwon",
			tokenGiven: tokenString,
			statusCode: 400,
		},
		{
			id:           strconv.Itoa(int(1)),
			buyer_id:     1,
			statusCode:   401,
			errorMessage: "Unauthorized",
		},
	}
	for _, v := range transactionSample {

		req, _ := http.NewRequest("GET", "/transactions", nil)
		req = mux.SetURLVars(req, map[string]string{"id": v.id})

		rr := httptest.NewRecorder()
		handler := http.HandlerFunc(server.DeleteTransaction)

		req.Header.Set("Authorization", v.tokenGiven)

		handler.ServeHTTP(rr, req)

		assert.Equal(t, rr.Code, v.statusCode)

		if v.statusCode == 401 && v.errorMessage != "" {

			responseMap := make(map[string]interface{})
			err = json.Unmarshal([]byte(rr.Body.String()), &responseMap)
			if err != nil {
				t.Errorf("Cannot convert to json: %v", err)
			}
			assert.Equal(t, responseMap["error"], v.errorMessage)
		}
	}
}
