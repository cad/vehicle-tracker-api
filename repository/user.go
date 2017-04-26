package repository

import (
	"fmt"
	"log"
	"time"
	_ "github.com/jinzhu/gorm/dialects/sqlite"
	"gopkg.in/hlandau/passlib.v1"
	"github.com/google/uuid"
)

type User struct {
	ID        uint       `json:"-"     gorm:"primary_key"`
	CreatedAt time.Time  `json:"-"`
	UpdatedAt time.Time  `json:"updated_at"`
	DeletedAt *time.Time `json:"-"`

	Email     string     `json:"email" gorm:"not null;unique_index"`
	Password  string     `json:"-"`
	Token     string     `json:"-"`
}

func (u *User) SetPassword(password string) error {
	hash, err := passlib.Hash(password)
	if err != nil {
		return &UserError{
			What: "User",
			Type: "Unknown",
			Arg: err.Error(),
		}
	}
	u.Password = hash
	return nil
}

func (u *User) RenewToken() (string, error) {
	if db.NewRecord(u) {
		return "", &UserError{
			What: "User",
			Type: "User-Not-Persisted",
			Arg: u.Email,
		}
	}
	uuid_token, err := uuid.NewRandom()
	if err != nil {
		return "", err
	}
	u.Token = uuid_token.String()
	db.Save(&u)
	return u.Token, nil
}


func GetAllUsers () []User {
	var users []User

	db.Find(&users)

	return users
}

func (u *User) CheckPassword(password string) (bool) {
	_, err := passlib.Verify(password, u.Password)
	if err != nil {
		log.Println(err.Error())
		return false
	}
	if db.NewRecord(&u) {
		log.Println("User.ID is 0 (not persisted record)")
		return false
	}
	return true
}

func CheckToken(token string) (bool) {
	var user User
	db.Where(&User{Token: token}).First(&user)
	if db.NewRecord(&user) {
		return false
	}

	return true
}

func GetUserByEmail(email string) (User, error) {
	var user User

	db.Where(&User{Email: email}).First(&user)
	if db.NewRecord(&user) {
		return user, UserError{
			What: "User",
			Type: "Not-Found",
			Arg: email,
		}
	}
	return user, nil
}

func CreateNewUser(email string, password string) (User, error) {
	user := User{
		Email: email,
	}
	user.SetPassword(password)
	db.Create(&user)
	if db.NewRecord(&user)  {
		return user, UserError{
			What: "User",
			Type: "Can-Not-Create",
			Arg: email,
		}
	}
	return user, nil
}

func DeleteUserByEmail(email string) (User, error) {
	user, err := GetUserByEmail(email)

	if err != nil {
		return user, err
	}

	if db.NewRecord(user) {
		return user, UserError{
			What: "User",
			Type: "Can-Not-Delete",
			Arg: email,
		}
	}

	db.Delete(&user)
	return user, nil
}

type UserError struct {
	What string
	Type string
	Arg string
}


func (e UserError) Error() string {
	return fmt.Sprintf("%s: <%s> %s", e.Type, e.What, e.Arg)
}
