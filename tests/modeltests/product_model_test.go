package modeltests

import (
	"log"
	"testing"

	"github.com/bayuajik2403/golang-backend-crud-api/api/models"
	_ "github.com/jinzhu/gorm/dialects/mysql"
	"gopkg.in/go-playground/assert.v1"
)

func TestFindAllProducts(t *testing.T) {

	err := refreshAllTable()
	if err != nil {
		log.Fatalf("Error refreshing all table %v\n", err)
	}
	_, _, _, err = seedAllTable()
	if err != nil {
		log.Fatalf("Error seeding all table %v\n", err)
	}
	products, err := productInstance.FindAllProducts(server.DB)
	if err != nil {
		t.Errorf("this is the error getting the products: %v\n", err)
		return
	}
	assert.Equal(t, len(*products), 2)
}

func TestSaveProduct(t *testing.T) {

	err := refreshAllTable()
	if err != nil {
		log.Fatalf("Error all refreshing table %v\n", err)
	}

	user, err := seedOneUser()
	if err != nil {
		log.Fatalf("Cannot seed user %v\n", err)
	}

	newProduct := models.Product{
		ID:                 1,
		ProductName:        "Oreo Kelapa",
		ProductDescription: "Oreo Kelapa taste",
		AvailableStock:     100,
		Price:              15000,
		SellerID:           user.ID,
	}
	savedProduct, err := newProduct.SaveProduct(server.DB)
	if err != nil {
		t.Errorf("this is the error getting the product: %v\n", err)
		return
	}
	assert.Equal(t, newProduct.ID, savedProduct.ID)
	assert.Equal(t, newProduct.ProductName, savedProduct.ProductName)
	assert.Equal(t, newProduct.ProductDescription, savedProduct.ProductDescription)
	assert.Equal(t, newProduct.AvailableStock, savedProduct.AvailableStock)
	assert.Equal(t, newProduct.Price, savedProduct.Price)
}

func TestGetProductByID(t *testing.T) {

	err := refreshAllTable()
	if err != nil {
		log.Fatalf("Error refreshing all table: %v\n", err)
	}
	_, product, _, err := seedOneAllTable()
	if err != nil {
		log.Fatalf("Error Seeding User and product table")
	}
	foundProduct, err := productInstance.FindProductByID(server.DB, product.ID)
	if err != nil {
		t.Errorf("this is the error getting one user: %v\n", err)
		return
	}
	assert.Equal(t, foundProduct.ID, product.ID)
	assert.Equal(t, foundProduct.ProductName, product.ProductName)
	assert.Equal(t, foundProduct.ProductDescription, product.ProductDescription)
	assert.Equal(t, foundProduct.AvailableStock, product.AvailableStock)
	assert.Equal(t, foundProduct.Price, product.Price)
}

func TestGetProductByUser(t *testing.T) {

	err := refreshAllTable()
	if err != nil {
		log.Fatalf("Error refreshing all table: %v\n", err)
	}
	user, product, _, err := seedOneAllTable()
	if err != nil {
		log.Fatalf("Error Seeding User and product table")
	}
	foundProduct, err := productInstance.FindProductByID(server.DB, user.ID)
	if err != nil {
		t.Errorf("this is the error getting one user: %v\n", err)
		return
	}
	assert.Equal(t, foundProduct.ID, product.ID)
	assert.Equal(t, foundProduct.ProductName, product.ProductName)
	assert.Equal(t, foundProduct.ProductDescription, product.ProductDescription)
	assert.Equal(t, foundProduct.AvailableStock, product.AvailableStock)
	assert.Equal(t, foundProduct.Price, product.Price)
}

func TestUpdateAProduct(t *testing.T) {

	err := refreshAllTable()
	if err != nil {
		log.Fatalf("Error refreshing all table: %v\n", err)
	}
	user, product, _, err := seedOneAllTable()
	if err != nil {
		log.Fatalf("Error Seeding table")
	}
	productUpdate := models.Product{
		ID:                 1,
		ProductName:        "Oreo Kelapa Update",
		ProductDescription: "Oreo Kelapa Update taste",
		AvailableStock:     120,
		Price:              15500,
		SellerID:           user.ID,
	}
	updatedProduct, err := productUpdate.UpdateAProduct(server.DB)
	if err != nil {
		t.Errorf("this is the error updating the user: %v\n", err)
		return
	}
	assert.Equal(t, updatedProduct.ID, productUpdate.ID)
	assert.Equal(t, updatedProduct.ProductName, productUpdate.ProductName)
	assert.Equal(t, updatedProduct.ProductDescription, productUpdate.ProductDescription)
	assert.Equal(t, updatedProduct.AvailableStock, productUpdate.AvailableStock)
	assert.Equal(t, updatedProduct.Price, productUpdate.Price)
	assert.Equal(t, updatedProduct.SellerID, productUpdate.SellerID)
}

func TestDeleteAProduct(t *testing.T) {

	err := refreshAllTable()
	if err != nil {
		log.Fatalf("Error refreshing all table: %v\n", err)
	}
	_, product, _, err := seedOneAllTable()
	if err != nil {
		log.Fatalf("Error Seeding tables")
	}
	isDeleted, err := productInstance.DeleteAProduct(server.DB, product.ID, product.SellerID)
	if err != nil {
		t.Errorf("this is the error deleting the user: %v\n", err)
		return
	}
	//one shows that the record has been deleted or:
	// assert.Equal(t, int(isDeleted), 1)

	//Can be done this way too
	assert.Equal(t, isDeleted, int64(1))
}
