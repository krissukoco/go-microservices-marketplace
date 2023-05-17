package model

import (
	"errors"

	"gorm.io/gorm"
)

var (
	ErrStoreNotFound = errors.New("store not found")
)

type Store struct {
	Id          string `gorm:"primaryKey"`
	Name        string
	Description string
	City        string
	Province    string
	UserId      string
	ImageFull   string
	ImageThumb  string
	IsPremium   bool
	IsOfficial  bool
	CreatedAt   int64 `gorm:"autoCreateTime:milli"`
	UpdatedAt   int64 `gorm:"autoUpdateTime:milli"`
}

func (s *Store) FindById(db *gorm.DB, id string) error {
	tx := db.Where(&Store{Id: id}).Take(&s)
	return tx.Error
}

func (s *Store) FindByUserId(db *gorm.DB, userId string) error {
	tx := db.Where(&Store{UserId: userId}).Take(&s)
	return tx.Error
}
