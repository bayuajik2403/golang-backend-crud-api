package models

import (
	"errors"
	"time"

	"github.com/jinzhu/gorm"
)

type Transaction struct {
	ID         uint64    `gorm:"primary_key;auto_increment" json:"id"`
	ProductT   Product   `json:"productt"`
	ProductID  uint32    `gorm:"not null" json:"product_id"`
	Qty        uint32    `gorm:"not null" json:"qty"`
	TotalPrice uint32    `gorm:"not null" json:"total_price"`
	Buyer      User      `json:"buyer"`
	BuyerID    uint32    `gorm:"not null" json:"buyer_id"`
	CreatedAt  time.Time `gorsm:"default:CURRENT_TIMESTAMP" json:"created_at"`
	UpdatedAt  time.Time `gorm:"default:CURRENT_TIMESTAMP" json:"updated_at"`
}

func (t *Transaction) Prepare() {
	t.ID = 0
	// t.Product = html.EscapeString(strings.TrimSpace(t.Product))
	// t.Price = html.EscapeString(strings.TrimSpace(t.Price))
	t.ProductT = Product{}
	t.Buyer = User{}
	t.CreatedAt = time.Now()
	t.UpdatedAt = time.Now()
}

func (t *Transaction) Validate() error {

	if t.ProductID < 1 {
		return errors.New("Required ProductID")
	}
	if t.Qty < 1 {
		return errors.New("Required Qty")
	}
	if (t.ProductT.AvailableStock - t.Qty) < 0 {
		return errors.New("No available stock for current product")
	}
	if t.TotalPrice < 0 {
		return errors.New("Required Total Price")
	}
	if t.BuyerID < 1 {
		return errors.New("Required Buyer ID")
	}
	return nil
}

func (t *Transaction) SaveTransaction(db *gorm.DB) (*Transaction, error) {
	var err error
	//var err2 error
	err = db.Debug().Model(&Transaction{}).Create(&t).Error
	if err != nil {
		return &Transaction{}, err
	}

	Products := Product{}
	err = db.Debug().Model(&Product{}).Where("id = ?", t.ProductID).Take(&Products).Error
	if err != nil {
		return &Transaction{}, err
	}

	if t.ID != 0 {
		err = db.Debug().Model(&User{}).Where("id = ?", t.BuyerID).Take(&t.Buyer).Error
		if err != nil {
			return &Transaction{}, err
		}
		var err = db.Debug().Model(&Product{}).Where("id = ?", t.ProductID).Updates(Product{ProductName: t.ProductT.ProductName, ProductDescription: t.ProductT.ProductDescription, AvailableStock: (Products.AvailableStock - t.Qty), Price: t.ProductT.Price, UpdatedAt: time.Now()}).Error
		if err != nil {
			return &Transaction{}, err
		}

		err = db.Debug().Model(&Product{}).Where("id = ?", t.ProductID).Take(&t.ProductT).Error
		if err != nil {
			return &Transaction{}, err
		}

		err = db.Debug().Model(&User{}).Where("id = ?", t.ProductT.SellerID).Take(&t.ProductT.Seller).Error
		if err != nil {
			return &Transaction{}, err
		}
	}
	return t, nil
}

func (t *Transaction) FindAllTransactions(db *gorm.DB) (*[]Transaction, error) {
	var err error
	Transactions := []Transaction{}
	err = db.Debug().Model(&Transaction{}).Limit(100).Find(&Transactions).Error
	if err != nil {
		return &[]Transaction{}, err
	}
	if len(Transactions) > 0 {
		for i, _ := range Transactions {
			err := db.Debug().Model(&User{}).Where("id = ?", Transactions[i].BuyerID).Take(&Transactions[i].Buyer).Error
			if err != nil {
				return &[]Transaction{}, err
			}
		}
		for i, _ := range Transactions {
			err := db.Debug().Model(&Product{}).Where("id = ?", Transactions[i].ProductID).Take(&Transactions[i].ProductT).Error
			if err != nil {
				return &[]Transaction{}, err
			}
		}

		for i, _ := range Transactions {
			err := db.Debug().Model(&User{}).Where("id = ?", Transactions[i].ProductT.SellerID).Take(&Transactions[i].ProductT.Seller).Error
			if err != nil {
				return &[]Transaction{}, err
			}
		}
	}
	return &Transactions, nil
}

func (t *Transaction) FindTransactionByID(db *gorm.DB, pid uint32) (*Transaction, error) {
	var err error
	err = db.Debug().Model(&Transaction{}).Where("id = ?", pid).Take(&t).Error
	if err != nil {
		return &Transaction{}, err
	}
	if t.ID != 0 {
		err = db.Debug().Model(&User{}).Where("id = ?", t.BuyerID).Take(&t.Buyer).Error
		if err != nil {
			return &Transaction{}, err
		}

		err = db.Debug().Model(&Product{}).Where("id = ?", t.ProductID).Take(&t.ProductT).Error
		if err != nil {
			return &Transaction{}, err
		}

		err = db.Debug().Model(&User{}).Where("id = ?", t.ProductT.SellerID).Take(&t.ProductT.Seller).Error
		if err != nil {
			return &Transaction{}, err
		}

	}
	return t, nil
}

func (t *Transaction) FindTransactionByUser(db *gorm.DB, pid uint32) (*[]Transaction, error) {
	var err error
	Transactions := []Transaction{}
	err = db.Debug().Model(&Transaction{}).Where("buyer_id = ?", pid).Limit(100).Find(&Transactions).Error
	if err != nil {
		return &[]Transaction{}, err
	}
	if len(Transactions) > 0 {
		for i, _ := range Transactions {
			err := db.Debug().Model(&User{}).Where("id = ?", Transactions[i].BuyerID).Take(&Transactions[i].Buyer).Error
			if err != nil {
				return &[]Transaction{}, err
			}
		}
		for i, _ := range Transactions {
			err := db.Debug().Model(&Product{}).Where("id = ?", Transactions[i].ProductID).Take(&Transactions[i].ProductT).Error
			if err != nil {
				return &[]Transaction{}, err
			}
		}
		for i, _ := range Transactions {
			err := db.Debug().Model(&User{}).Where("id = ?", Transactions[i].ProductT.SellerID).Take(&Transactions[i].ProductT.Seller).Error
			if err != nil {
				return &[]Transaction{}, err
			}
		}
	}
	return &Transactions, nil
}

func (t *Transaction) UpdateATransaction(db *gorm.DB) (*Transaction, error) {

	var err error
	//var err2 error
	// var n Transaction
	// if err := db.Where("ID = ?", t.ID).Scan(&n).Error; err != nil {
	// 	return nil, err
	// }

	// Products := Product{}
	// err = db.Debug().Model(&Product{}).Where("id = ?", t.ProductID).Take(&Products).Error
	// if err != nil {
	// 	return &Transaction{}, err
	// }

	err = db.Debug().Model(&Transaction{}).Where("id = ?", t.ID).Updates(Transaction{ProductID: t.ProductID, Qty: t.Qty, TotalPrice: t.TotalPrice, UpdatedAt: time.Now()}).Error
	if err != nil {
		return &Transaction{}, err
	}

	// err = db.Debug().Model(&Product{}).Where("id = ?", t.ProductID).Updates(Product{ProductName: t.ProductT.ProductName, ProductDescription: t.ProductT.ProductDescription, AvailableStock: (t.ProductT.AvailableStock + (n.Qty - t.Qty)), Price: t.ProductT.Price, UpdatedAt: time.Now()}).Error
	// if err != nil {
	// 	return &Transaction{}, err
	// }

	if t.ID != 0 {
		err = db.Debug().Model(&User{}).Where("id = ?", t.BuyerID).Take(&t.Buyer).Error
		if err != nil {
			return &Transaction{}, err
		}

		err = db.Debug().Model(&Product{}).Where("id = ?", t.ProductID).Take(&t.ProductT).Error
		if err != nil {
			return &Transaction{}, err
		}

		err = db.Debug().Model(&User{}).Where("id = ?", t.ProductT.SellerID).Take(&t.ProductT.Seller).Error
		if err != nil {
			return &Transaction{}, err
		}
	}
	return t, nil
}

func (t *Transaction) DeleteATransaction(db *gorm.DB, pid uint64, uid uint32) (int64, error) {

	db = db.Debug().Model(&Transaction{}).Where("id = ? and Buyer_id = ?", pid, uid).Take(&Transaction{}).Delete(&Transaction{})

	if db.Error != nil {
		if gorm.IsRecordNotFoundError(db.Error) {
			return 0, errors.New("Transaction not found")
		}
		return 0, db.Error
	}
	return db.RowsAffected, nil
}
