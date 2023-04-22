package model

type Service struct {
	ID          string   `json:"id" bson:"id"`
	Name        string   `json:"name" bson:"name"`
	Description string   `json:"description" bson:"description"`
	Category    string   `json:"category" bson:"category"`
	Price       int64    `json:"price" bson:"price"`
	Duration    int64    `json:"duration" bson:"duration"` // Hours
	Images      []string `json:"images" bson:"images"`
	CreatedAt   int64    `json:"created_at" bson:"created_at"`
	UpdatedAt   int64    `json:"updated_at" bson:"updated_at"`
}

func (s *Service) New() error {
	return nil
}
