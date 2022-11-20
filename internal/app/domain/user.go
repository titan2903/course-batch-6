package domain

import (
	"errors"
	"os"
	"time"

	"github.com/golang-jwt/jwt/v4"
	"golang.org/x/crypto/bcrypt"
)

var signature = []byte(os.Getenv("JWT_SECRET"))

type User struct {
	ID        int       `json:"id"`
	Name      string    `json:"name"`
	Email     string    `json:"email"`
	Password  string    `json:"password"`
	NoHp      string    `json:"no_hp"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type UserRegister struct {
	Name     string `json:"name"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

type UserLogin struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

func NewUser(email, name, password string) (*User, error) {
	if name == "" {
		return nil, errors.New("name is required")
	}

	if email == "" {
		return nil, errors.New("email is required")
	}

	if password == "" {
		return nil, errors.New("password is required")
	}

	if len(password) < 7 {
		return nil, errors.New("password length must be 6 character or more")
	}

	hash, _ := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	return &User{
		Name:     name,
		Email:    email,
		Password: string(hash),
	}, nil
}

func (u User) GenerateJWT() (string, error) {
	claims := jwt.MapClaims{
		"user_id": u.ID,
		"exp":     time.Now().Add(24 * time.Hour).Unix(),
		"iss":     "edspert",
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	stringToken, err := token.SignedString(signature)
	return stringToken, err
}

func (u User) DecryptJWT(token string) (map[string]interface{}, error) {
	parsedToken, err := jwt.Parse(token, func(t *jwt.Token) (interface{}, error) {
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("invalid token")
		}
		return signature, nil
	})
	data := make(map[string]interface{})
	if err != nil {
		return data, err
	}

	if !parsedToken.Valid {
		return data, errors.New("invalid token")
	}

	return parsedToken.Claims.(jwt.MapClaims), nil
}

func (u User) CorrectPassword(password string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(u.Password), []byte(password))
	return err == nil
}
