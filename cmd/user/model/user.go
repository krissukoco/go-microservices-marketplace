package model

import (
	"errors"
	"time"

	"github.com/golang-jwt/jwt"
	"github.com/google/uuid"
	"gorm.io/gorm"

	"golang.org/x/crypto/bcrypt"
)

var (
	ErrUserNotFound         = errors.New("user not found")
	ErrTokenInvalid         = errors.New("token invalid")
	ErrTokenExpired         = errors.New("token expired")
	ErrInvalidSigningMethod = errors.New("invalid signing method")
)

type User struct {
	ID        string `gorm:"primaryKey"`
	Email     string
	Password  string
	FirstName string
	LastName  string
	ImageUrl  string
	City      string
	Province  string
	CreatedAt int64 `gorm:"autoCreateTime:milli"`
	UpdatedAt int64 `gorm:"autoUpdateTime:milli"`
}

func (u *User) GetFullName() string {
	return u.FirstName + " " + u.LastName
}

func (u *User) FindByID(db *gorm.DB, id string) error {
	tx := db.Where(&User{ID: id}).First(&u)
	return tx.Error
}

func (u *User) FindByEmail(db *gorm.DB, email string) error {
	tx := db.Where(&User{Email: email}).First(&u)
	return tx.Error
}

// UPDATE user if exists
// INSERT user if not exists
func (u *User) Save(db *gorm.DB) error {
	err := u.FindByID(db, u.ID)
	if err != nil && err != gorm.ErrRecordNotFound {
		return err
	}
	if err == nil {
		tx := db.Save(&u)
		return tx.Error
	}
	tx := db.Create(&u)
	return tx.Error
}

func (u *User) ComparePassword(password string) error {
	return bcrypt.CompareHashAndPassword([]byte(u.Password), []byte(password))
}

func (u *User) ChangePassword(password string) error {
	hashedPwd, err := bcrypt.GenerateFromPassword([]byte(password), 7)
	if err != nil {
		return err
	}
	u.Password = string(hashedPwd)
	return nil
}

func (u *User) GenerateJWT(secret string) (string, error) {
	n := time.Now()
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.StandardClaims{
		Subject:   u.ID,
		Audience:  "github.com/krissukoco/go-microservices-marketplace",
		Issuer:    "github.com/krissukoco",
		IssuedAt:  n.Unix(),
		NotBefore: n.Unix(),
		ExpiresAt: n.Add(time.Hour * 24).Unix(),
		Id:        uuid.NewString(),
	})
	tokenStr, err := token.SignedString([]byte(secret))
	if err != nil {
		return "", err
	}
	return tokenStr, nil
}

func (u *User) FromJWT(db *gorm.DB, jwtString string, jwtSecret string) error {
	token, err := jwt.Parse(jwtString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, ErrInvalidSigningMethod
		}
		return []byte(jwtSecret), nil
	})
	if err != nil {
		return err
	}
	// log.Println("Token claims: ", token.Claims)
	// log.Println("Token valid: ", token.Valid)
	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		// log.Println("claims: ", claims)
		expIntf, ok := claims["exp"]
		if !ok {
			return ErrTokenInvalid
		}
		expFl, ok := expIntf.(float64)
		if !ok {
			return ErrTokenInvalid
		}
		exp := int64(expFl)
		if time.Now().Unix() > exp { // Note: Token `eat` is always in second
			return ErrTokenExpired
		}
		userId, ok := claims["sub"].(string)
		if !ok {
			return ErrTokenInvalid
		}
		if err := u.FindByID(db, userId); err != nil {
			return err
		}

		return nil
	}
	// claims, ok := token.Claims.(jwt.StandardClaims)
	// log.Println("Claims OK? ", ok)
	// if ok && token.Valid {
	// 	if time.Now().Unix() > claims.ExpiresAt { // Note: Token `exp` is always in second
	// 		return ErrTokenExpired
	// 	}
	// 	userId := claims.Subject
	// 	if err := u.FindByID(db, userId); err != nil {
	// 		return err
	// 	}

	// 	return nil
	// }
	return ErrTokenInvalid
}
