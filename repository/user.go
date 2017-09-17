package repository

import (
	"fmt"
	"log"
	"time"

	"github.com/google/uuid"
	_ "github.com/jinzhu/gorm/dialects/sqlite" // Force sqlite dialect to load
	"gopkg.in/hlandau/passlib.v1"
)

type User struct {
	ID        uint      `json:"-"     gorm:"primary_key"`
	UUID      string    `json:"uuid" gorm:"unique_index"`
	CreatedAt time.Time `json:"-"`
	UpdatedAt time.Time `json:"updated_at"`

	Email    string `json:"email" gorm:"unique_index"`
	Password string `json:"-"`
	Token    string `json:"-"`
}

func (u *User) SetPassword(password string) error {
	hash, err := passlib.Hash(password)
	if err != nil {
		return &UserError{
			What: "User",
			Type: "Unknown",
			Arg:  err.Error(),
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
			Arg:  u.Email,
		}
	}
	uuidToken, err := uuid.NewRandom()
	if err != nil {
		return "", err
	}
	u.Token = uuidToken.String()
	db.Save(&u)
	return u.Token, nil
}

func GetAllUsers() []User {
	var users []User

	db.Find(&users)

	return users
}

func (u *User) CheckPassword(password string) bool {
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

func CheckToken(token string) bool {
	if token == "" {
		// TODO(cad): Write a log about what happened here.
		return false
	}

	var user User
	db.Where(&User{Token: token}).First(&user)
	if db.NewRecord(&user) {
		return false
	}

	return true
}

func GetUserByUUID(uUID string) (User, error) {
	var user User
	if uUID == "" {
		return user, &UserError{What: "uUID", Type: "Empty", Arg: uUID}
	}

	db.Where(&User{UUID: uUID}).First(&user)
	if db.NewRecord(&user) {
		return user, UserError{
			What: "User",
			Type: "Not-Found",
			Arg:  uUID,
		}
	}
	return user, nil
}

func GetUserByEmail(email string) (User, error) {
	var user User
	if email == "" {
		return user, &UserError{What: "email", Type: "Empty", Arg: email}
	}

	db.Where(&User{Email: email}).First(&user)
	if db.NewRecord(&user) {
		return user, UserError{
			What: "User",
			Type: "Not-Found",
			Arg:  email,
		}
	}
	return user, nil
}

func GetUserByToken(token string) (User, error) {
	var user User
	if token == "" {
		// TODO(cad): Write a log about what happened here.
		return user, fmt.Errorf("token is expected to be not empty, but it's empty instead")
	}

	db.Where(&User{Token: token}).First(&user)
	if db.NewRecord(&user) {
		return user, fmt.Errorf("user is not found")
	}

	return user, nil
}

func CreateNewUser(email string, password string) (User, error) {
	var user User
	if email == "" || password == "" {
		return user, &UserError{What: "EmailOrPassword", Type: "Empty", Arg: ""}
	}

	uUID, err := uuid.NewRandom()
	user = User{
		Email: email,
		UUID:  uUID.String(),
	}
	if err != nil {
		return user, err
	}
	user.SetPassword(password)
	db.Create(&user)
	if db.NewRecord(&user) {
		return user, UserError{
			What: "User",
			Type: "Can-Not-Create",
			Arg:  email,
		}
	}
	return user, nil
}

func DeleteUserByUUID(uUID string) (User, error) {
	user, err := GetUserByUUID(uUID)

	if err != nil {
		return user, err
	}

	if db.NewRecord(user) {
		return user, UserError{
			What: "User",
			Type: "Can-Not-Delete",
			Arg:  uUID,
		}
	}

	db.Unscoped().Delete(&user)
	return user, nil
}

type UserError struct {
	What string
	Type string
	Arg  string
}

func (e UserError) Error() string {
	return fmt.Sprintf("%s: <%s> %s", e.Type, e.What, e.Arg)
}
