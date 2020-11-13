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

func TestCreateProduct(t *testing.T) {

	err := refreshAllTable()
	if err != nil {
		log.Fatal(err)
	}
	user, err := seedOneUser()
	if err != nil {
		log.Fatalf("Cannot seed user %v\n", err)
	}
	token, err := server.SignIn(user.Email, "password") //Note the password in the database is already hashed, we want unhashed
	if err != nil {
		log.Fatalf("cannot login: %v\n", err)
	}
	tokenString := fmt.Sprintf("Bearer %v", token)

	samples := []struct {
		inputJSON           string
		statusCode          int
		product_name        string
		product_description string
		available_stock     uint32
		price               uint32
		seller_id           uint64
		tokenGiven          string
		errorMessage        string
	}{
		{
			inputJSON:           `{"product_name":"Oreo Kelapa Muda","product_description":"Oreo Kelapa Muda taste","available_stock":100,"price":15000,"seller_id":1}`,
			statusCode:          201,
			tokenGiven:          tokenString,
			product_name:        "Oreo Kelapa Muda",
			product_description: "Oreo Kelapa Muda taste",
			available_stock:     100,
			price:               15000,
			seller_id:           user.ID,
			errorMessage:        "",
		},
		{
			// When no token is passed
			inputJSON:    `{"product_name":"When no token is passed","product_description":"Oreo Kelapa taste","available_stock":100,"price":15000,"seller_id":1}`,
			statusCode:   401,
			tokenGiven:   "",
			errorMessage: "Unauthorized",
		},
		{
			// When incorrect token is passed
			inputJSON:    `{"product_name":"When incorrect token is passed","product_description":"Oreo Kelapa taste","available_stock":100,"price":15000,"seller_id":1}`,
			statusCode:   401,
			tokenGiven:   "This is an incorrect token",
			errorMessage: "Unauthorized",
		},
		{
			inputJSON:    `{"product_name":"","product_description":"Oreo Kelapa taste","available_stock":100,"price":15000,"seller_id":1}`,
			statusCode:   422,
			tokenGiven:   tokenString,
			errorMessage: "Required Product Name",
		},
		{
			inputJSON:    `{"product_name":"Oreo Kelapa","product_description":"","available_stock":100,"price":15000,"seller_id":1}`,
			statusCode:   422,
			tokenGiven:   tokenString,
			errorMessage: "Required Product Description",
		},
		{
			inputJSON:    `{"product_name":"Oreo Kelapa","product_description":"Oreo Kelapa taste","available_stock":100,"price":15000,"seller_id":0}`,
			statusCode:   422,
			tokenGiven:   tokenString,
			errorMessage: "Required Seller ID",
		},
		{
			// When user 2 uses user 1 token
			inputJSON:    `{"product_name":"Oreo Kelapa","product_description":"Oreo Kelapa taste","available_stock":100,"price":15000,"seller_id":2}`,
			statusCode:   401,
			tokenGiven:   tokenString,
			errorMessage: "Unauthorized",
		},
	}
	for _, v := range samples {

		req, err := http.NewRequest("POST", "/products", bytes.NewBufferString(v.inputJSON))
		if err != nil {
			t.Errorf("this is the error: %v\n", err)
		}
		rr := httptest.NewRecorder()
		handler := http.HandlerFunc(server.CreateProduct)

		req.Header.Set("Authorization", v.tokenGiven)
		handler.ServeHTTP(rr, req)

		responseMap := make(map[string]interface{})
		err = json.Unmarshal([]byte(rr.Body.String()), &responseMap)
		if err != nil {
			fmt.Printf("Cannot convert to json: %v", err)
		}
		assert.Equal(t, rr.Code, v.statusCode)
		if v.statusCode == 201 {
			assert.Equal(t, responseMap["product_name"], v.product_name)
			assert.Equal(t, responseMap["product_description"], v.product_description)
			assert.Equal(t, responseMap["available_stock"], float64(v.available_stock))
			assert.Equal(t, responseMap["price"], float64(v.price))
			assert.Equal(t, responseMap["seller_id"], float64(v.seller_id)) //just for both ids to have the same type
		}
		if v.statusCode == 401 || v.statusCode == 422 || v.statusCode == 500 && v.errorMessage != "" {
			assert.Equal(t, responseMap["error"], v.errorMessage)
		}
	}
}

func TestGetProducts(t *testing.T) {

	err := refreshAllTable()
	if err != nil {
		log.Fatal(err)
	}
	_, _, _, err = seedAllTable()
	if err != nil {
		log.Fatal(err)
	}

	req, err := http.NewRequest("GET", "/products", nil)
	if err != nil {
		t.Errorf("this is the error: %v\n", err)
	}
	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(server.GetProducts)
	handler.ServeHTTP(rr, req)

	var products []models.Product
	err = json.Unmarshal([]byte(rr.Body.String()), &products)

	assert.Equal(t, rr.Code, http.StatusOK)
	assert.Equal(t, len(products), 2)
}
func TestGetProductByID(t *testing.T) {

	err := refreshAllTable()
	if err != nil {
		log.Fatal(err)
	}
	_, product, _, err := seedOneAllTable()
	if err != nil {
		log.Fatal(err)
	}
	productSample := []struct {
		id                  string
		statusCode          int
		product_name        string
		product_description string
		available_stock     uint32
		price               uint32
		seller_id           uint64
		errorMessage        string
	}{
		{
			id:                  strconv.Itoa(int(product.ID)),
			statusCode:          200,
			product_name:        product.ProductName,
			product_description: product.ProductDescription,
			available_stock:     product.AvailableStock,
			price:               product.Price,
			seller_id:           product.SellerID,
		},
		{
			id:         "unknwon",
			statusCode: 400,
		},
	}
	for _, v := range productSample {

		req, err := http.NewRequest("GET", "/products", nil)
		if err != nil {
			t.Errorf("this is the error: %v\n", err)
		}
		req = mux.SetURLVars(req, map[string]string{"id": v.id})

		rr := httptest.NewRecorder()
		handler := http.HandlerFunc(server.GetProduct)
		handler.ServeHTTP(rr, req)

		responseMap := make(map[string]interface{})
		err = json.Unmarshal([]byte(rr.Body.String()), &responseMap)
		if err != nil {
			log.Fatalf("Cannot convert to json: %v", err)
		}
		assert.Equal(t, rr.Code, v.statusCode)

		if v.statusCode == 200 {
			// assert.Equal(t, product.Title, responseMap["title"])
			// assert.Equal(t, product.Content, responseMap["content"])
			// assert.Equal(t, float64(product.AuthorID), responseMap["author_id"]) //the response author id is float64

			assert.Equal(t, responseMap["product_name"], product.ProductName)
			assert.Equal(t, responseMap["product_description"], product.ProductDescription)
			assert.Equal(t, responseMap["available_stock"], float64(product.AvailableStock))
			assert.Equal(t, responseMap["price"], float64(product.Price))
			assert.Equal(t, responseMap["seller_id"], float64(product.SellerID)) //just for both ids to have the same type

		}
	}
}

func TestUpdateProduct(t *testing.T) {

	var ProductUserEmail, ProductUserPassword string
	var AuthProductSellerID uint64
	var AuthProductID uint64

	err := refreshAllTable()
	if err != nil {
		log.Fatal(err)
	}
	users, products, _, err := seedAllTable()
	if err != nil {
		log.Fatal(err)
	}
	// Get only the first user
	for _, user := range users {
		if user.ID == 2 {
			continue
		}
		ProductUserEmail = user.Email
		ProductUserPassword = "password" //Note the password in the database is already hashed, we want unhashed
	}
	//Login the user and get the authentication token
	token, err := server.SignIn(ProductUserEmail, ProductUserPassword)
	if err != nil {
		log.Fatalf("cannot login: %v\n", err)
	}
	tokenString := fmt.Sprintf("Bearer %v", token)

	// Get only the first product
	for _, product := range products {
		if product.ID == 2 {
			continue
		}
		AuthProductID = product.ID
		AuthProductSellerID = product.SellerID
	}
	// fmt.Printf("this is the auth product: %v\n", AuthProductID)

	samples := []struct {
		id                  string
		updateJSON          string
		statusCode          int
		product_name        string
		product_description string
		available_stock     uint32
		price               uint32
		seller_id           uint64
		tokenGiven          string
		errorMessage        string
	}{
		{
			// Convert int64 to int first before converting to string
			id:                  strconv.Itoa(int(AuthProductID)),
			updateJSON:          `{"product_name":"update","product_description":"update","available_stock":100,"price":15000,"seller_id":1}`,
			statusCode:          200,
			product_name:        "update",
			product_description: "update",
			available_stock:     100,
			price:               15000,
			seller_id:           AuthProductSellerID,
			tokenGiven:          tokenString,
			errorMessage:        "",
		},
		{
			// When no token is provided
			id:           strconv.Itoa(int(AuthProductID)),
			updateJSON:   `{"product_name":"update","product_description":"update","available_stock":101,"price":16000,"seller_id":1}`,
			tokenGiven:   "",
			statusCode:   401,
			errorMessage: "Unauthorized",
		},
		{
			// When incorrect token is provided
			id:           strconv.Itoa(int(AuthProductID)),
			updateJSON:   `{"product_name":"update","product_description":"update","available_stock":100,"price":15000,"seller_id":1}`,
			tokenGiven:   "this is an incorrect token",
			statusCode:   401,
			errorMessage: "Unauthorized",
		},
		{
			id:           strconv.Itoa(int(AuthProductID)),
			updateJSON:   `{"product_name":"","product_description":"update","available_stock":100,"price":15000,"seller_id":1}`,
			statusCode:   422,
			tokenGiven:   tokenString,
			errorMessage: "Required Product Name",
		},
		{
			id:           strconv.Itoa(int(AuthProductID)),
			updateJSON:   `{"product_name":"update","product_description":"","available_stock":100,"price":15000,"seller_id":1}`,
			statusCode:   422,
			tokenGiven:   tokenString,
			errorMessage: "Required Product Description",
		},
		// {
		// 	id:           strconv.Itoa(int(AuthProductID)),
		// 	updateJSON:   `{"product_name":"update","product_description":"update","available_stock":100,"price":15000,"seller_id":0}`,
		// 	statusCode:   422,
		// 	tokenGiven:   tokenString,
		// 	errorMessage: "Required Seller ID",
		// },
		{
			id:         "unknwon",
			statusCode: 400,
		},
	}

	for _, v := range samples {

		req, err := http.NewRequest("POST", "/products", bytes.NewBufferString(v.updateJSON))
		if err != nil {
			t.Errorf("this is the error: %v\n", err)
		}
		req = mux.SetURLVars(req, map[string]string{"id": v.id})
		rr := httptest.NewRecorder()
		handler := http.HandlerFunc(server.UpdateProduct)

		req.Header.Set("Authorization", v.tokenGiven)

		handler.ServeHTTP(rr, req)

		responseMap := make(map[string]interface{})
		err = json.Unmarshal([]byte(rr.Body.String()), &responseMap)
		if err != nil {
			t.Errorf("Cannot convert to json: %v", err)
		}
		assert.Equal(t, rr.Code, v.statusCode)
		if v.statusCode == 200 {
			assert.Equal(t, responseMap["product_name"], v.product_name)
			assert.Equal(t, responseMap["product_description"], v.product_description)
			assert.Equal(t, responseMap["available_stock"], float64(v.available_stock))
			assert.Equal(t, responseMap["price"], float64(v.price))
			assert.Equal(t, responseMap["seller_id"], float64(v.seller_id)) //just for both ids to have the same type
		}
		if v.statusCode == 401 || v.statusCode == 422 || v.statusCode == 500 && v.errorMessage != "" {
			assert.Equal(t, responseMap["error"], v.errorMessage)
		}
	}
}

func TestDeleteProduct(t *testing.T) {

	var ProductUserEmail, ProductUserPassword string
	var ProductUserID uint64
	var AuthProductID uint64

	err := refreshAllTable()
	if err != nil {
		log.Fatal(err)
	}
	users, products, _, err := seedAllTable()
	if err != nil {
		log.Fatal(err)
	}
	//Let's get only the Second user
	for _, user := range users {
		if user.ID == 1 {
			continue
		}
		ProductUserEmail = user.Email
		ProductUserPassword = "password" //Note the password in the database is already hashed, we want unhashed
	}
	//Login the user and get the authentication token
	token, err := server.SignIn(ProductUserEmail, ProductUserPassword)
	if err != nil {
		log.Fatalf("cannot login: %v\n", err)
	}
	tokenString := fmt.Sprintf("Bearer %v", token)

	// Get only the second product
	for _, product := range products {
		if product.ID == 1 {
			continue
		}
		AuthProductID = product.ID
		ProductUserID = product.SellerID
	}
	productSample := []struct {
		id           string
		seller_id    uint64
		tokenGiven   string
		statusCode   int
		errorMessage string
	}{
		{
			// Convert int64 to int first before converting to string
			id:           strconv.Itoa(int(AuthProductID)),
			seller_id:    ProductUserID,
			tokenGiven:   tokenString,
			statusCode:   200,
			errorMessage: "",
		},
		{
			// When empty token is passed
			id:           strconv.Itoa(int(AuthProductID)),
			seller_id:    ProductUserID,
			tokenGiven:   "",
			statusCode:   401,
			errorMessage: "Unauthorized",
		},
		{
			// When incorrect token is passed
			id:           strconv.Itoa(int(AuthProductID)),
			seller_id:    ProductUserID,
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
			seller_id:    1,
			statusCode:   401,
			errorMessage: "Unauthorized",
		},
	}
	for _, v := range productSample {

		req, _ := http.NewRequest("GET", "/products", nil)
		req = mux.SetURLVars(req, map[string]string{"id": v.id})

		rr := httptest.NewRecorder()
		handler := http.HandlerFunc(server.DeleteProduct)

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
