package model

import (
	"fmt"
	"strings"

	"github.com/krissukoco/go-microservices-marketplace/internal/utils"
	"gorm.io/gorm"
)

var (
	ErrProductNotFound = fmt.Errorf("product not found")
)

// ProductVariant represents data in db
type ProductVariant struct {
	Id        uint `gorm:"primaryKey"`
	ProductId string
	Name      string
	Value     string
}

type ProductVariantGroup struct {
	// Name: Type or Title of variant, e.g. Color, Size, etc.
	Name string
	// Values: Array of actual values, e.g. [Red, Blue, Green] or [S, M, L]
	Values []string
}

type Product struct {
	Id            string `gorm:"primaryKey"`
	StoreId       string
	Title         string
	Slug          string
	Description   string
	Category      string
	Tags          string                 // comma separated
	CreatedAt     int64                  `gorm:"autoCreateTime:milli"`
	UpdatedAt     int64                  `gorm:"autoUpdateTime:milli"`
	Variants      []*ProductVariantGroup `gorm:"-"`
	AverageRating float64                `gorm:"-"`
	TotalRating   int64                  `gorm:"-"`
	TotalReview   int64                  `gorm:"-"`
	TotalSold     int64                  `gorm:"-"`
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
			group = &ProductVariantGroup{Name: variant.Name, Values: make([]string, 0)}
			groupMap[variant.Name] = group
		}
		group.Values = append(group.Values, variant.Value)
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
	return nil
}

func (p *Product) Save(db *gorm.DB) error {
	if p.Id == "" {
		uid := utils.NewAlphanumericID(10)
		p.Id = "p_" + uid
	}
	tx := db.Save(p)
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

func GetAllProductsBySearch(db *gorm.DB, search string, page int, size int) ([]*Product, error) {
	var products []*Product
	offset := (page - 1) * size
	tx := db.Where("title LIKE ?", "%"+search+"%").
		Or("tags LIKE ?", "%"+search+"%").
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
