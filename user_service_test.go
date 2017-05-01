package main

import "testing"
import "github.com/asdine/storm"
import "os"
import "strconv"

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

	_, err = us.GetByID(3)
	if err == nil {
		t.Error("Found a user that should not be in database by ID")
	}

	_, err = us.GetByEmail("someEmail@email.com")
	if err == nil {
		t.Error("Found a user that should not be in database by email")
	}
}

func TestMultiThreading(t *testing.T) {
	db, err := storm.Open("test.db")
	if err != nil {
		t.Fatal(err)
	}

	defer func() {
		db.Close()
		os.Remove("test.db")
	}()

	us := NewUserServiceFromDB(db)

	numRoutines := 5000

	register := func(id int) {
		us.RegisterUser("email"+strconv.Itoa(id), "password")
	}

	getByEmail := func(id int) {
		email := "email" + strconv.Itoa(id)
		_, err := us.GetByEmail(email)
		if err != nil {
			t.Error("Could not find user with email ", email)
		}
	}

	getByID := func(id int) {
		_, err := us.GetByID(id)
		if err != nil {
			t.Error("Could not find user with id ", id)
		}
	}

	update := func(id int) {
		defer func() {
			if recover() != nil {
				t.Log("Recover in id ", id)
			}
		}()

		user, err := us.GetByID(id)
		if err != nil {
			t.Error("Could not find user with id ", id)
		}

		modified := user.Email + "modified"
		user.Email = modified
		us.Update(user)

		user, err = us.GetByID(id)
		if err != nil {
			t.Error("Could not find user with id ", id)
		}

		if user.Email != modified {
			t.Error("Could not find the modified user with id ", id)
		}
	}

	for i := 0; i < numRoutines; i++ {
		go register(i)
	}

	for i := 0; i < numRoutines; i++ {
		go getByEmail(i)
	}

	for i := 0; i < numRoutines; i++ {
		go getByID(i + 1)
		go update(i)
	}

}
