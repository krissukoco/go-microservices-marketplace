package model

import "gorm.io/gorm"

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
	CreatedAt int64                 `gorm:"autoCreateTime:milli"`
	Media     []*ProductReviewMedia `gorm:"foreignKey:ReviewId"`
}

type ProductRating struct {
	Id        uint `gorm:"primaryKey"`
	ProductId string
	UserId    string
	Rating    int
	CreatedAt int64 `gorm:"autoCreateTime:milli"`
}

func GetProductReviews(db *gorm.DB, productId string) ([]*ProductReview, error) {
	var reviews []*ProductReview
	tx := db.Where(&ProductReview{ProductId: productId}).Preload("Media").Find(&reviews)
	return reviews, tx.Error
}

func GetProductRatings(db *gorm.DB, productId string) ([]*ProductRating, error) {
	var ratings []*ProductRating
	tx := db.Where(&ProductRating{ProductId: productId}).Find(&ratings)
	return ratings, tx.Error
}
