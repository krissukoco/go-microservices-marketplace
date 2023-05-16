package model

type Store struct {
	Id          string `gorm:"primaryKey"`
	Name        string
	Description string
	City        string
	Province    string
	OwnerId     string
	ImageFull   string
	ImageThumb  string
	IsPremium   bool
	IsOfficial  bool
	CreatedAt   int64 `gorm:"autoCreateTime:milli"`
	UpdatedAt   int64 `gorm:"autoUpdateTime:milli"`
}
