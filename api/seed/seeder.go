package seed

import (
	"log"

	"github.com/bayuajik2403/go-crud-api/api/models"
	"github.com/jinzhu/gorm"
)

var users = []models.User{
	models.User{
		Nickname: "Bayu Aji",
		Email:    "bayuaji@gmail.com",
		Password: "password",
	},
	models.User{
		Nickname: "Bayu Kurniawan",
		Email:    "bayukurniawan@gmail.com",
		Password: "password",
	},
}

var products = []models.Product{
	models.Product{
		ProductName:        "Oreo Strawberry",
		ProductDescription: "Oreo Strawberry Taste",
		AvailableStock:     100,
		Price:              15000,
	},
	models.Product{
		ProductName:        "Oreo Vanila",
		ProductDescription: "Oreo Vanila Taste",
		AvailableStock:     100,
		Price:              14000,
	},
}

var transaction = []models.Transaction{
	models.Transaction{
		ProductID:  1,
		Qty:        3,
		TotalPrice: 45000,
		BuyerID:    2,
	},
	models.Transaction{
		ProductID:  2,
		Qty:        3,
		TotalPrice: 42000,
		BuyerID:    1,
	},
}

func Load(db *gorm.DB) {

	err := db.Debug().DropTableIfExists(&models.Product{}, &models.User{}, &models.Transaction{}).Error
	if err != nil {
		log.Fatalf("cannot drop table: %v", err)
	}
	err = db.Debug().AutoMigrate(&models.User{}, &models.Product{}, &models.Transaction{}).Error
	if err != nil {
		log.Fatalf("cannot migrate table: %v", err)
	}

	err = db.Debug().Model(&models.Product{}).AddForeignKey("seller_id", "users(id)", "cascade", "cascade").Error
	if err != nil {
		log.Fatalf("attaching foreign key error: %v", err)
	}

	err = db.Debug().Model(&models.Transaction{}).AddForeignKey("buyer_id", "users(id)", "cascade", "cascade").Error
	if err != nil {
		log.Fatalf("attaching foreign key error: %v", err)
	}

	err = db.Debug().Model(&models.Transaction{}).AddForeignKey("product_id", "products(id)", "cascade", "cascade").Error
	if err != nil {
		log.Fatalf("attaching foreign key error: %v", err)
	}

	for i, _ := range users {
		err = db.Debug().Model(&models.User{}).Create(&users[i]).Error
		if err != nil {
			log.Fatalf("cannot seed users table: %v", err)
		}
		products[i].SellerID = users[i].ID

		err = db.Debug().Model(&models.Product{}).Create(&products[i]).Error
		if err != nil {
			log.Fatalf("cannot seed posts table: %v", err)
		}

	}
	for i, _ := range transaction {
		err = db.Debug().Model(&models.Transaction{}).Create(&transaction[i]).Error
		if err != nil {
			log.Fatalf("cannot seed posts table: %v", err)
		}

	}
}
