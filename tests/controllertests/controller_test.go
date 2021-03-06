package controllertests

import (
	"fmt"
	"log"
	"os"
	"testing"

	"github.com/bayuajik2403/golang-backend-crud-api/api/controllers"
	"github.com/bayuajik2403/golang-backend-crud-api/api/models"
	"github.com/jinzhu/gorm"
	"github.com/joho/godotenv"
)

var server = controllers.Server{}
var userInstance = models.User{}
var productInstance = models.Product{}
var transactionInstance = models.Transaction{}

func TestMain(m *testing.M) {
	err := godotenv.Load(os.ExpandEnv("../../.env"))
	if err != nil {
		log.Fatalf("Error getting env %v\n", err)
	}
	Database()

	os.Exit(m.Run())

}

func Database() {

	var err error

	TestDbDriver := os.Getenv("TestDbDriver")

	DBURL := fmt.Sprintf("host=%s port=%s user=%s dbname=%s sslmode=disable password=%s", os.Getenv("TestDbHost"), os.Getenv("TestDbPort"), os.Getenv("TestDbUser"), os.Getenv("TestDbName"), os.Getenv("TestDbPassword"))
	server.DB, err = gorm.Open(TestDbDriver, DBURL)
	if err != nil {
		fmt.Printf("Cannot connect to %s database\n", TestDbDriver)
		log.Fatal("This is the error:", err)
	} else {
		fmt.Printf("We are connected to the %s database\n", TestDbDriver)
	}

}

func refreshUserTable() error {
	/*
		err := server.DB.DropTableIfExists(&models.User{}).Error
		if err != nil {
			return err
		}
		err = server.DB.AutoMigrate(&models.User{}).Error
		if err != nil {
			return err
		}
	*/
	// err := server.DB.DropTableIfExists(&models.User{}).Error
	// if err != nil {
	// 	return err
	// }
	// err = server.DB.AutoMigrate(&models.User{}).Error
	// if err != nil {
	// 	return err
	// }

	// log.Printf("Successfully refreshed table(s)")
	// return nil

	//server.DB.Exec("SET foreign_key_checks=0")
	err := server.DB.Debug().DropTableIfExists(&models.User{}).Error
	if err != nil {
		return err
	}
	//server.DB.Exec("SET foreign_key_checks=1")
	err = server.DB.Debug().AutoMigrate(&models.User{}).Error
	if err != nil {
		return err
	}
	log.Printf("Successfully refreshed table")
	log.Printf("refreshUserTable routine OK !!!")
	return nil
}

func seedOneUser() (models.User, error) {

	err := refreshUserTable()
	if err != nil {
		log.Fatal(err)
	}

	user := models.User{
		Nickname: "Pet",
		Email:    "pet@gmail.com",
		Password: "password",
	}

	err2 := server.DB.Debug().Model(&models.User{}).Create(&user).Error
	if err2 != nil {
		log.Fatalf("cannot seed users table: %v", err)
	}

	log.Printf("seedOneUser routine OK !!!")
	return user, nil
}

func seedUsers() ([]models.User, error) {

	var err error
	if err != nil {
		return nil, err
	}
	users := []models.User{
		models.User{
			Nickname: "Steven victor",
			Email:    "steven@gmail.com",
			Password: "password",
		},
		models.User{
			Nickname: "Kenny Morris",
			Email:    "kenny@gmail.com",
			Password: "password",
		},
	}
	for i, _ := range users {
		err := server.DB.Model(&models.User{}).Create(&users[i]).Error
		if err != nil {
			return []models.User{}, err
		}
	}
	return users, nil
}

func refreshAllTable() error {

	err := server.DB.DropTableIfExists(&models.Product{}, &models.Transaction{}, &models.User{}).Error
	if err != nil {
		return err
	}
	err = server.DB.AutoMigrate(&models.User{}, &models.Product{}, &models.Transaction{}).Error
	if err != nil {
		return err
	}
	log.Printf("Successfully refreshed tables")
	return nil
}

func seedOneAllTable() (models.User, models.Product, models.Transaction, error) {

	err := refreshAllTable()
	if err != nil {
		return models.User{}, models.Product{}, models.Transaction{}, err
	}
	user := models.User{
		Nickname: "Faker Phil",
		Email:    "faker@gmail.com",
		Password: "password",
	}
	err = server.DB.Debug().Model(&models.User{}).Create(&user).Error
	if err != nil {
		return models.User{}, models.Product{}, models.Transaction{}, err
	}
	product := models.Product{
		ProductName:        "Oreo Kelapa",
		ProductDescription: "Oreo Kelapa taste",
		AvailableStock:     100,
		Price:              15000,
		SellerID:           user.ID,
	}
	err = server.DB.Debug().Model(&models.Product{}).Create(&product).Error
	if err != nil {
		return models.User{}, models.Product{}, models.Transaction{}, err
	}
	transaction := models.Transaction{
		ProductID:  product.ID,
		Qty:        3,
		TotalPrice: 45000,
		BuyerID:    user.ID,
	}
	err = server.DB.Debug().Model(&models.Transaction{}).Create(&transaction).Error
	if err != nil {
		return models.User{}, models.Product{}, models.Transaction{}, err
	}

	log.Printf("seedOneAllTable routine OK !!!")
	return user, product, transaction, nil
}

func seedAllTable() ([]models.User, []models.Product, []models.Transaction, error) {

	var err error

	if err != nil {
		return []models.User{}, []models.Product{}, []models.Transaction{}, err
	}
	var users = []models.User{
		models.User{
			Nickname: "Yongwoo Grid",
			Email:    "yongwoo@gmail.com",
			Password: "password",
		},
		models.User{
			Nickname: "Yura Ann",
			Email:    "yura@gmail.com",
			Password: "password",
		},
	}
	var products = []models.Product{
		models.Product{
			ProductName:        "Oreo Kelapa",
			ProductDescription: "Oreo Kelapa taste",
			AvailableStock:     100,
			Price:              15000,
		},
		models.Product{
			ProductName:        "Oreo Kelapa Muda",
			ProductDescription: "Oreo Kelapa Muda taste",
			AvailableStock:     100,
			Price:              15000,
		},
	}
	var transactions = []models.Transaction{
		models.Transaction{
			ProductID:  1,
			Qty:        3,
			TotalPrice: 45000,
		},
		models.Transaction{
			ProductID:  2,
			Qty:        3,
			TotalPrice: 45000,
		},
	}

	for i := range users {
		err = server.DB.Debug().Model(&models.User{}).Create(&users[i]).Error
		if err != nil {
			log.Fatalf("cannot seed users table: %v", err)
		}
		products[i].SellerID = users[i].ID
		transactions[i].BuyerID = users[i].ID
		err = server.DB.Debug().Model(&models.Product{}).Create(&products[i]).Error
		if err != nil {
			log.Fatalf("cannot seed posts table: %v", err)
		}
		err = server.DB.Debug().Model(&models.Transaction{}).Create(&transactions[i]).Error
		if err != nil {
			log.Fatalf("cannot seed posts table: %v", err)
		}
	}
	log.Printf("seedAllTable routine OK !!!")
	return users, products, transactions, nil
}
