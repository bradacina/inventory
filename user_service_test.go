package main

import "testing"
import "github.com/asdine/storm"
import "os"

const (
	email = "test1@email.com"
)

func TestUserService(t *testing.T) {
	db, err := storm.Open("test.db")
	if err != nil {
		t.Fatal(err)
	}

	defer func() {
		db.Close()
		os.Remove("test.db")
	}()

	us := NewUserServiceFromDB(db)

	ur := newUserRepo(db)
	_ = NewUserService(ur)

	err = us.RegisterUser(email, "password1")
	if err != nil {
		t.Error("Cannot register a new user, Error:", err)
	}

	err = us.RegisterUser(email, "password2")
	if err == nil {
		t.Error("Email re-usage detected")
	}

	err = us.RegisterUser("test2@email.com", "password")
	if err != nil {
		t.Error("Cannot register a second user, Error:", err)
	}

	user, err := us.GetByEmail(email)
	if err != nil {
		t.Error("Could not find user by email, Error:", err)
	}

	if user.Email != email {
		t.Error("Retrieved the wrong user")
	}

	result, err := us.ValidateCredentials(email, "password1")
	if err != nil {
		t.Error("Cannot validate credentials")
	}

	if !result {
		t.Error("Wrong result when validating credentials")
	}

	_, err = us.ValidateCredentials("other@email.com", "password1")
	if err == nil {
		t.Error("Found a user that should not exist")
	}

	_, err = us.ValidateCredentials(email, "password")
	if err == nil {
		t.Error("Passwords should not match")
	}

	user2, err := us.GetByID(user.ID)
	if err != nil {
		t.Error("Could not retrieve user by ID, Error:", err)
	}

	if user2.Email != email {
		t.Error("Retrived the wrong user by ID")
	}
}
