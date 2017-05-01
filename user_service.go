package main

import (
	"errors"

	"github.com/asdine/storm"
	"golang.org/x/crypto/bcrypt"
)

type UserServicer interface {
	RegisterUser(email, password string) error
	GetByEmail(email string) (*User, error)
	GetByID(id int) (*User, error)
	ValidateCredentials(email, password string) (bool, error)
	Update(*User)
}

var (
	ErrorInvalidCredentials = errors.New("Invalid Credentials")
	ErrorEmailInUse         = errors.New("Email address is already")
	ErrorNotFound           = errors.New("Item not found")
)

type userService struct {
	userRepo UserRepoer
}

func NewUserService(userRepo UserRepoer) UserServicer {
	return &userService{userRepo}
}

func NewUserServiceFromDB(db *storm.DB) UserServicer {
	return &userService{userRepo: newUserRepo(db)}
}

func (us *userService) RegisterUser(email, password string) error {
	repo := us.userRepo
	_, err := repo.GetByEmail(email)
	if err != nil {
		hashed, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
		if err != nil {
			return err
		}

		repo.Upsert(&User{Email: email, Password: hashed})
		return nil
	}

	return ErrorEmailInUse
}

func (us *userService) GetByEmail(email string) (*User, error) {
	repo := us.userRepo
	user, err := repo.GetByEmail(email)
	if err != nil {
		return nil, ErrorNotFound
	}

	return user, nil
}

func (us *userService) GetByID(id int) (*User, error) {
	repo := us.userRepo
	user, err := repo.GetByID(id)
	if err != nil {
		return nil, ErrorNotFound
	}

	return user, nil
}

func (us *userService) ValidateCredentials(email, password string) (bool, error) {
	user, err := us.GetByEmail(email)
	if err != nil {
		return false, ErrorInvalidCredentials
	}

	err = bcrypt.CompareHashAndPassword(user.Password, []byte(password))
	if err != nil {
		return false, ErrorInvalidCredentials
	}

	return true, nil
}

func (us *userService) Update(user *User) {
	us.Update(user)
}
