package model

import (
	"fmt"
	"math/rand"
	"strings"
	"time"

	"gorm.io/gorm"
)

var (
	AlphaNumerics      = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	ErrProductNotFound = fmt.Errorf("product not found")
)

// ProductVariant represents data in db
type ProductVariant struct {
	Id        uint `gorm:"primaryKey"`
	ProductId string
	Name      string
	Value     string
	Price     int64
	Stock     int64
	CreatedAt int64 `gorm:"autoCreateTime:milli"`
	UpdatedAt int64 `gorm:"autoUpdateTime:milli"`
}

func (va *ProductVariant) Save(db *gorm.DB) error {
	tx := db.Save(va)
	return tx.Error
}

type VariantItem struct {
	Value string
	Price int64
	Stock int64
}

type ProductVariantGroup struct {
	// Name: Type or Title of variant, e.g. Color, Size, etc.
	Name string
	// Values: Array of actual values, e.g. [Red, Blue, Green] or [S, M, L]
	Values []*VariantItem
}

type Product struct {
	Id          string `gorm:"primaryKey"`
	StoreId     string
	Title       string
	Slug        string
	Description string
	// Category: should be written with hierarcy as 'Shoes->Sneakers'
	Category string
	// Tags: comma separated string, e.g. 'shoes,sneakers,running'
	Tags string
	// Price: in IDR; if product has variant, then this is the lowest price
	Price int64
	// Stock: if product has no variant, then this is the actual stock
	Stock         int64
	TotalViews    int64
	CreatedAt     int64                  `gorm:"autoCreateTime:milli"`
	UpdatedAt     int64                  `gorm:"autoUpdateTime:milli"`
	Variants      []*ProductVariantGroup `gorm:"-"`
	AverageRating float64                `gorm:"-"`
	TotalReview   int64                  `gorm:"-"`
	TotalSold     int64                  `gorm:"-"`
}

func NewProductId(count ...int) string {
	c := 10
	if len(count) > 0 {
		c = count[0]
	}
	id := "p_"
	for i := 0; i < c; i++ {
		rand.Seed(time.Now().UnixNano())
		id += string(AlphaNumerics[rand.Intn(len(AlphaNumerics))])
	}
	return id
}

func (p *Product) GetSlug() string {
	titleSplit := strings.Split(p.Title, " ")
	for i, word := range titleSplit {
		titleSplit[i] = strings.ToLower(word)
	}
	return fmt.Sprintf("%s-%s", p.Id, p.Title)
}

func (p *Product) FillVariants(db *gorm.DB) error {
	// Get variants from db
	var variants []*ProductVariant
	tx := db.Where(&ProductVariant{ProductId: p.Id}).Find(&variants)
	if tx.Error != nil {
		return tx.Error
	}
	// Group variants
	p.Variants = make([]*ProductVariantGroup, 0)
	var groupMap = make(map[string]*ProductVariantGroup)
	for _, variant := range variants {
		group, ok := groupMap[variant.Name]
		if !ok {
			group = &ProductVariantGroup{Name: variant.Name, Values: make([]*VariantItem, 0)}
			groupMap[variant.Name] = group
		}
		group.Values = append(group.Values, &VariantItem{
			Value: variant.Value,
			Price: variant.Price,
			Stock: variant.Stock,
		})
	}
	for _, group := range groupMap {
		p.Variants = append(p.Variants, group)
	}
	return nil
}

func (p *Product) GetById(db *gorm.DB, id string) error {
	tx := db.Where(&Product{Id: id}).Take(&p)
	if tx.Error != nil {
		if tx.Error == gorm.ErrRecordNotFound {
			return ErrProductNotFound
		}
		return tx.Error
	}
	// Fill variants
	err := p.FillVariants(db)
	if err != nil {
		return err
	}
	// Get reviews summary
	reviews, _ := GetProductReviewSummary(db, id)
	if reviews != nil {
		p.AverageRating = reviews.Average
		p.TotalReview = reviews.Total
	}
	return nil
}

func (p *Product) GetBySlug(db *gorm.DB, slug string) error {
	tx := db.Where(&Product{Slug: slug}).Take(&p)
	if tx.Error != nil {
		if tx.Error == gorm.ErrRecordNotFound {
			return ErrProductNotFound
		}
		return tx.Error
	}
	// Fill variants
	err := p.FillVariants(db)
	if err != nil {
		return err
	}
	// Get reviews summary
	reviews, _ := GetProductReviewSummary(db, p.Id)
	if reviews != nil {
		p.AverageRating = reviews.Average
		p.TotalReview = reviews.Total
	}
	return nil
}

func (p *Product) Save(db *gorm.DB) error {
	if p.Id == "" {
		p.Id = NewProductId()
	}
	p.Slug = p.GetSlug()
	tx := db.Save(p)
	return tx.Error
}

func (p *Product) Delete(db *gorm.DB) error {
	tx := db.Delete(&Product{Id: p.Id, StoreId: p.StoreId})
	return tx.Error
}

func GetAllProductsByStore(db *gorm.DB, storeId string) ([]*Product, error) {
	var products []*Product
	tx := db.Where(&Product{StoreId: storeId}).Find(&products)
	if tx.Error != nil {
		return nil, tx.Error
	}
	for _, product := range products {
		err := product.FillVariants(db)
		if err != nil {
			return nil, err
		}
	}
	return products, nil
}

func GetAllProductsBySearch(db *gorm.DB, search, category string, page, size int) ([]*Product, error) {
	var products []*Product
	offset := (page - 1) * size
	tx := db
	if search != "" {
		tx = tx.Where("title LIKE ?", "%"+search+"%")
	}
	if category != "" {
		tx = tx.Where("category LIKE ?", "%"+category+"%")
	}
	tx = tx.Order("created_at DESC"). // TODO: order by popularity
						Offset(offset).
						Limit(size).
						Find(&products)
	if tx.Error != nil {
		return nil, tx.Error
	}
	for _, product := range products {
		err := product.FillVariants(db)
		if err != nil {
			return nil, err
		}
	}
	return products, nil
}
