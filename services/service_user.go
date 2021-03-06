package services

import (
	"errors"

	"github.com/asdine/storm"
	"github.com/bradacina/inventory/db"
	"golang.org/x/crypto/bcrypt"
)

type UserServicer interface {
	RegisterUser(email, password string) error
	GetByEmail(email string) (*db.User, error)
	GetByID(id int) (*db.User, error)
	GetAll() ([]db.User, error)
	ValidateCredentials(email, password string) (*db.User, error)
	Update(user *db.User, userID int) error

	UpdateByAdmin(user *db.User) error
}

var (
	ErrorInvalidCredentials = errors.New("Invalid Credentials")
	ErrorEmailInUse         = errors.New("Email address is already")
	ErrorNotFound           = errors.New("Item not found")
)

type userService struct {
	userRepo db.UserRepoer
}

func NewUserService(userRepo db.UserRepoer) UserServicer {
	return &userService{userRepo}
}

func NewUserServiceFromDB(stormdb *storm.DB) UserServicer {
	return &userService{userRepo: db.NewUserRepo(stormdb)}
}

func (us *userService) GetAll() ([]db.User, error) {
	return us.userRepo.GetAll()
}

func (us *userService) RegisterUser(email, password string) error {
	repo := us.userRepo
	_, err := repo.GetByEmail(email)
	if err != nil {
		hashed, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
		if err != nil {
			return err
		}

		user := db.User{Email: email, Password: hashed}
		repo.Upsert(&user)

		if user.ID == 1 {
			user.IsAdmin = true
		}

		repo.Upsert(&user)
		return nil
	}

	return ErrorEmailInUse
}

func (us *userService) GetByEmail(email string) (*db.User, error) {
	repo := us.userRepo
	user, err := repo.GetByEmail(email)
	if err != nil {
		return nil, ErrorNotFound
	}

	return user, nil
}

func (us *userService) GetByID(id int) (*db.User, error) {
	repo := us.userRepo
	user, err := repo.GetByID(id)
	if err != nil {
		return nil, ErrorNotFound
	}

	return user, nil
}

func (us *userService) ValidateCredentials(email, password string) (*db.User, error) {
	user, err := us.GetByEmail(email)
	if err != nil {
		return nil, ErrorInvalidCredentials
	}

	err = bcrypt.CompareHashAndPassword(user.Password, []byte(password))
	if err != nil {
		return nil, ErrorInvalidCredentials
	}

	return user, nil
}

func (us *userService) UpdateByAdmin(user *db.User) error {
	if user == nil {
		return ErrorNotFound
	}

	if user.ID <= 0 {
		return ErrorNotFound
	}

	existing, err := us.userRepo.GetByID(user.ID)
	if err != nil {
		return err
	}

	if existing.Email != user.Email {
		_, err := us.userRepo.GetByEmail(user.Email)
		if err == nil {
			return ErrorEmailInUse
		}

		existing.Email = user.Email
	}

	existing.IsAdmin = user.IsAdmin
	existing.IsDeleted = user.IsDeleted
	err = us.userRepo.Upsert(existing)
	if err != nil {
		return err
	}

	return nil
}

// Updates a user record
// userID is the id of the user that's making the request
func (us *userService) Update(user *db.User, userID int) error {

	if userID != 0 && userID != user.ID {
		return ErrorNotFound

	}

	return us.UpdateByAdmin(user)
}
