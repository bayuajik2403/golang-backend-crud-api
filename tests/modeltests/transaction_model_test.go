package modeltests

import (
	"log"
	"testing"

	"github.com/bayuajik2403/golang-backend-crud-api/api/models"
	_ "github.com/jinzhu/gorm/dialects/mysql"
	"gopkg.in/go-playground/assert.v1"
)

func TestFindAllTransactions(t *testing.T) {

	err := refreshAllTable()
	if err != nil {
		log.Fatalf("Error refreshing all table %v\n", err)
	}
	_, _, _, err = seedAllTable()
	if err != nil {
		log.Fatalf("Error seeding all table %v\n", err)
	}
	transactions, err := transactionInstance.FindAllTransactions(server.DB)
	if err != nil {
		t.Errorf("this is the error getting the transactions: %v\n", err)
		return
	}
	assert.Equal(t, len(*transactions), 2)
}

func TestSaveTransaction(t *testing.T) {

	err := refreshAllTable()
	if err != nil {
		log.Fatalf("Error all refreshing table %v\n", err)
	}

	user, product, _, err := seedOneAllTable()
	if err != nil {
		log.Fatalf("Cannot seed All table %v\n", err)
	}

	newTransaction := models.Transaction{
		ID:         1,
		ProductID:  product.ID,
		Qty:        3,
		TotalPrice: 45000,
		BuyerID:    user.ID,
	}
	savedTransaction, err := newTransaction.SaveTransaction(server.DB)
	if err != nil {
		t.Errorf("this is the error getting the transaction: %v\n", err)
		return
	}
	assert.Equal(t, newTransaction.ID, savedTransaction.ID)
	assert.Equal(t, newTransaction.ProductID, savedTransaction.ProductID)
	assert.Equal(t, newTransaction.Qty, savedTransaction.Qty)
	assert.Equal(t, newTransaction.TotalPrice, savedTransaction.TotalPrice)
	assert.Equal(t, newTransaction.BuyerID, savedTransaction.BuyerID)
}

func TestGetTransactionByID(t *testing.T) {

	err := refreshAllTable()
	if err != nil {
		log.Fatalf("Error refreshing all table: %v\n", err)
	}
	_, _, transaction, err := seedOneAllTable()
	if err != nil {
		log.Fatalf("Error Seeding one all table")
	}
	foundTransaction, err := transactionInstance.FindTransactionByID(server.DB, transaction.ID)
	if err != nil {
		t.Errorf("this is the error getting one user: %v\n", err)
		return
	}
	assert.Equal(t, foundTransaction.ID, transaction.ID)
	assert.Equal(t, foundTransaction.ProductID, transaction.ProductID)
	assert.Equal(t, foundTransaction.Qty, transaction.Qty)
	assert.Equal(t, foundTransaction.TotalPrice, transaction.TotalPrice)
	assert.Equal(t, foundTransaction.BuyerID, transaction.BuyerID)
}

func TestGetTransactionByUser(t *testing.T) {

	err := refreshAllTable()
	if err != nil {
		log.Fatalf("Error refreshing all table: %v\n", err)
	}
	user, transaction, _, err := seedOneAllTable()
	if err != nil {
		log.Fatalf("Error seeding one all table")
	}
	foundTransaction, err := transactionInstance.FindTransactionByID(server.DB, user.ID)
	if err != nil {
		t.Errorf("this is the error getting one user: %v\n", err)
		return
	}
	assert.Equal(t, foundTransaction.ID, transaction.ID)
	assert.Equal(t, foundTransaction.ProductID, transaction.ProductID)
	assert.Equal(t, foundTransaction.Qty, transaction.Qty)
	assert.Equal(t, foundTransaction.TotalPrice, transaction.TotalPrice)
	assert.Equal(t, foundTransaction.BuyerID, transaction.BuyerID)
}

func TestUpdateATransaction(t *testing.T) {

	err := refreshAllTable()
	if err != nil {
		log.Fatalf("Error refreshing all table: %v\n", err)
	}
	user, product, transaction, err := seedOneAllTable()
	if err != nil {
		log.Fatalf("Error Seeding table")
	}
	transactionUpdate := models.Transaction{
		ID:         1,
		ProductID:  product.ID,
		Qty:        4,
		TotalPrice: 50000,
		BuyerID:    user.ID,
	}
	updatedTransaction, err := transactionUpdate.UpdateATransaction(server.DB)
	if err != nil {
		t.Errorf("this is the error updating the user: %v\n", err)
		return
	}
	assert.Equal(t, updatedTransaction.ID, transactionUpdate.ID)
	assert.Equal(t, updatedTransaction.ProductID, transactionUpdate.ProductID)
	assert.Equal(t, updatedTransaction.Qty, transactionUpdate.Qty)
	assert.Equal(t, updatedTransaction.TotalPrice, transactionUpdate.TotalPrice)
	assert.Equal(t, updatedTransaction.BuyerID, transactionUpdate.BuyerID)
}

func TestDeleteATransaction(t *testing.T) {

	err := refreshAllTable()
	if err != nil {
		log.Fatalf("Error refreshing all table: %v\n", err)
	}
	_, transaction, _, err := seedOneAllTable()
	if err != nil {
		log.Fatalf("Error Seeding tables")
	}
	isDeleted, err := transactionInstance.DeleteATransaction(server.DB, transaction.ID, transaction.BuyerID)
	if err != nil {
		t.Errorf("this is the error deleting the user: %v\n", err)
		return
	}
	//one shows that the record has been deleted or:
	// assert.Equal(t, int(isDeleted), 1)

	//Can be done this way too
	assert.Equal(t, isDeleted, int64(1))
}
