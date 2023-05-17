package model

import (
	"errors"

	"gorm.io/gorm"
)

var (
	ErrNoRatings = errors.New("no ratings for product")
)

type reviewSummary struct {
	Total   int64
	Average float64
}

type ProductReviewMedia struct {
	Id       uint `gorm:"primaryKey"`
	ReviewId uint
	Url      string
}

type ProductReview struct {
	Id        uint `gorm:"primaryKey"`
	ProductId string
	UserId    string
	Comment   string
	Rating    int
	CreatedAt int64                 `gorm:"autoCreateTime:milli"`
	Media     []*ProductReviewMedia `gorm:"foreignKey:ReviewId"`
}

func GetProductReviews(db *gorm.DB, productId string) ([]*ProductReview, error) {
	var reviews []*ProductReview
	tx := db.Where(&ProductReview{ProductId: productId}).Preload("Media").Find(&reviews)
	return reviews, tx.Error
}

func GetProductReviewSummary(db *gorm.DB, productId string) (*reviewSummary, error) {
	var summary reviewSummary
	tx := db.Model(&ProductReview{}).
		Where(&ProductReview{ProductId: productId}).
		Select("COUNT(*) as total, AVG(rating) as average").
		Scan(&summary)
	if tx.Error != nil {
		return nil, tx.Error
	}
	if summary.Total == 0 {
		return nil, ErrNoRatings
	}
	return &summary, nil
}
