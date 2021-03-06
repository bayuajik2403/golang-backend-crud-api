package models

import (
	"errors"
	"html"
	"strings"
	"time"

	"github.com/jinzhu/gorm"
)

type Product struct {
	ID                 uint64    `gorm:"primary_key;auto_increment" json:"id"`
	ProductName        string    `gorm:"size:255;not null" json:"product_name"`
	ProductDescription string    `gorm:"size:255;not null" json:"product_description"`
	AvailableStock     uint32    `gorm:"not null;" json:"available_stock"`
	Price              uint32    `gorm:"not null;" json:"price"`
	Seller             User      `json:"seller"`
	SellerID           uint64    `gorm:"not null" json:"seller_id"`
	CreatedAt          time.Time `gorm:"default:CURRENT_TIMESTAMP" json:"created_at"`
	UpdatedAt          time.Time `gorm:"default:CURRENT_TIMESTAMP" json:"updated_at"`
}

func (p *Product) Prepare() {
	p.ID = 0
	p.ProductName = html.EscapeString(strings.TrimSpace(p.ProductName))
	p.ProductDescription = html.EscapeString(strings.TrimSpace(p.ProductDescription))
	p.Seller = User{}
	p.CreatedAt = time.Now()
	p.UpdatedAt = time.Now()
}

func (p *Product) Validate() error {

	if p.ProductName == "" {
		return errors.New("Required Product Name")
	}
	if p.ProductDescription == "" {
		return errors.New("Required Product Description")
	}
	var n uint64 = 1
	if p.SellerID < n {
		return errors.New("Required Seller ID")
	}
	return nil
}

func (p *Product) SaveProduct(db *gorm.DB) (*Product, error) {
	var err error
	err = db.Debug().Model(&Product{}).Create(&p).Error
	if err != nil {
		return &Product{}, err
	}
	if p.ID != 0 {
		err = db.Debug().Model(&User{}).Where("id = ?", p.SellerID).Take(&p.Seller).Error
		if err != nil {
			return &Product{}, err
		}
	}
	return p, nil
}

func (p *Product) FindAllProducts(db *gorm.DB) (*[]Product, error) {
	var err error
	Products := []Product{}
	err = db.Debug().Model(&Product{}).Limit(100).Find(&Products).Error
	if err != nil {
		return &[]Product{}, err
	}
	if len(Products) > 0 {
		for i, _ := range Products {
			err := db.Debug().Model(&User{}).Where("id = ?", Products[i].SellerID).Take(&Products[i].Seller).Error
			if err != nil {
				return &[]Product{}, err
			}
		}
	}
	return &Products, nil
}

func (p *Product) FindProductByID(db *gorm.DB, pid uint64) (*Product, error) {
	var err error
	err = db.Debug().Model(&Product{}).Where("id = ?", pid).Take(&p).Error
	if err != nil {
		return &Product{}, err
	}
	if p.ID != 0 {
		err = db.Debug().Model(&User{}).Where("id = ?", p.SellerID).Take(&p.Seller).Error
		if err != nil {
			return &Product{}, err
		}
	}
	return p, nil
}

func (p *Product) FindProductByUser(db *gorm.DB, pid uint64) (*[]Product, error) {
	var err error
	Products := []Product{}
	err = db.Debug().Model(&Product{}).Where("seller_id = ?", pid).Limit(100).Find(&Products).Error
	if err != nil {
		return &[]Product{}, err
	}
	if len(Products) > 0 {
		for i, _ := range Products {
			err := db.Debug().Model(&User{}).Where("id = ?", Products[i].SellerID).Take(&Products[i].Seller).Error
			if err != nil {
				return &[]Product{}, err
			}
		}
	}
	return &Products, nil
}

func (p *Product) UpdateAProduct(db *gorm.DB) (*Product, error) {

	var err error

	err = db.Debug().Model(&Product{}).Where("id = ?", p.ID).Updates(Product{ProductName: p.ProductName, ProductDescription: p.ProductDescription, AvailableStock: p.AvailableStock, Price: p.Price, UpdatedAt: time.Now()}).Error
	if err != nil {
		return &Product{}, err
	}
	if p.ID != 0 {
		err = db.Debug().Model(&User{}).Where("id = ?", p.SellerID).Take(&p.Seller).Error
		if err != nil {
			return &Product{}, err
		}
	}
	return p, nil
}

func (p *Product) DeleteAProduct(db *gorm.DB, pid uint64, uid uint64) (int64, error) {

	db = db.Debug().Model(&Product{}).Where("id = ? and seller_id = ?", pid, uid).Take(&Product{}).Delete(&Product{})

	if db.Error != nil {
		if gorm.IsRecordNotFoundError(db.Error) {
			return 0, errors.New("Product not found")
		}
		return 0, db.Error
	}
	return db.RowsAffected, nil
}
