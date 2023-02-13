package model

import (
	"errors"
	"log"
	"time"

	"github.com/golang-jwt/jwt"
	"github.com/google/uuid"
	"github.com/krissukoco/go-microservices-marketplace/cmd/user/config"
	"github.com/krissukoco/go-microservices-marketplace/cmd/user/database"

	"golang.org/x/crypto/bcrypt"
)

type User struct {
	ID        string `json:"id"`
	Email     string `json:"email"`
	Password  string `json:"-"`
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
}

func NewUser(email, password, firstName, lastName string) (*User, error) {
	hashedPwd, err := bcrypt.GenerateFromPassword([]byte(password), 7)
	if err != nil {
		return nil, err
	}
	user := &User{
		ID:        "user_" + uuid.NewString(),
		Email:     email,
		Password:  string(hashedPwd),
		FirstName: firstName,
		LastName:  lastName,
	}

	return user, nil
}

func (u *User) GetFullName() string {
	return u.FirstName + " " + u.LastName
}

func (u *User) FindByID(id string) error {
	exists := false
	rows, err := database.PG.Query("SELECT * FROM users WHERE id = $1 LIMIT 1", id)
	if err != nil {
		log.Println("ERROR on query: ", err)
		return err
	}
	defer rows.Close()

	for rows.Next() {
		err := rows.Scan(&u.ID, &u.Email, &u.Password, &u.FirstName, &u.LastName)
		if err != nil {
			return err
		}
		exists = true
	}
	if !exists {
		return errors.New("user not found")
	}

	return nil
}

func (u *User) FindByEmail(email string) error {
	rows, err := database.PG.Query("SELECT * FROM users WHERE email = $1 LIMIT 1", email)
	if err != nil {
		log.Println("ERROR on query: ", err)
		return err
	}
	defer rows.Close()

	for rows.Next() {
		err := rows.Scan(&u.ID, &u.Email, &u.Password, &u.FirstName, &u.LastName)
		if err != nil {
			return err
		}
	}
	if u.ID == "" {
		return errors.New("user not found")
	}

	return nil
}

// UPDATE user if exists
// INSERT user if not exists
func (u *User) Save() error {
	exists := false
	rows, err := database.PG.Query("SELECT * FROM users WHERE id = $1 LIMIT 1", u.ID)
	if err != nil {
		log.Println("ERROR on query: ", err)
		return err
	}
	defer rows.Close()
	for rows.Next() {
		exists = true
	}
	if exists {
		_, err = database.PG.Exec(
			"UPDATE users SET email = $1, password = $2, first_name = $3, last_name = $4 WHERE id = $5",
			u.Email, u.Password, u.FirstName, u.LastName, u.ID,
		)
		if err != nil {
			return err
		}
		return nil
	}

	_, err = database.PG.Exec(
		"INSERT INTO users (id, email, password, first_name, last_name) VALUES ($1, $2, $3, $4, $5)",
		u.ID, u.Email, u.Password, u.FirstName, u.LastName,
	)
	if err != nil {
		return err
	}

	return nil
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

func (u *User) GenerateJWT() (string, error) {
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
	tokenStr, err := token.SignedString([]byte(config.Cfg.JWTSecret))
	if err != nil {
		return "", err
	}
	return tokenStr, nil
}

func (u *User) FromJWT(jwtString string) error {
	token, err := jwt.Parse(jwtString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("unexpected signing method")
		}
		return []byte(config.Cfg.JWTSecret), nil
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
			return errors.New("token invalid exp not found")
		}
		expFl, ok := expIntf.(float64)
		if !ok {
			return errors.New("token invalid exp not float64")
		}
		exp := int64(expFl)
		if time.Now().Unix() > exp {
			return errors.New("token expired")
		}
		userId, ok := claims["sub"].(string)
		if !ok {
			return errors.New("token invalid sub not string")
		}
		if err := u.FindByID(userId); err != nil {
			return err
		}

		return nil
	}
	return errors.New("invalid token not ok or not valid")
}
