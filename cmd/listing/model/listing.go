package model

type Category struct {
	ID      int    `json:"id" bson:"id"`
	Name    string `json:"name" bson:"name"`
	SubID   string `json:"sub_id" bson:"sub_id"`
	SubName string `json:"sub_name" bson:"sub_name"`
}

type Listing struct {
	ID             int         `json:"id" bson:"id"`
	Title          string      `json:"title" bson:"title"`
	Description    string      `json:"description" bson:"description"`
	Category       *Category   `json:"category" bson:"category"`
	Price          int64       `json:"price" bson:"price"`
	Stock          int64       `json:"stock" bson:"stock"`
	Year           int         `json:"year" bson:"year"`
	Second         bool        `json:"second" bson:"second"`
	Images         []string    `json:"images" bson:"images"`
	Specifications interface{} `json:"specifications" bson:"specifications"`
	CreatedAt      int64       `json:"created_at" bson:"created_at"`
	UpdatedAt      int64       `json:"updated_at" bson:"updated_at"`
}
