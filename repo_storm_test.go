package main

import (
	"errors"
	"testing"

	"os"

	"github.com/asdine/storm"
)

var ErrorRepositoryDidntPanic = errors.New("Repository didn't panic")

func catchPanic(x func()) error {
	defer func() {
		recover()
	}()

	x()
	return ErrorRepositoryDidntPanic
}

func TestDbNil(t *testing.T) {
	ur := newUserRepo((*storm.DB)(nil))
	ir := newInventoryRepo((*storm.DB)(nil))

	err := catchPanic(func() { ur.GetAll() })
	if err != nil {
		t.Error(err)
	}

	err = catchPanic(func() { ur.GetByEmail("1") })
	if err != nil {
		t.Error(err)
	}

	err = catchPanic(func() { ur.GetByID(1) })
	if err != nil {
		t.Error(err)
	}

	err = catchPanic(func() { ur.Upsert(&User{}) })
	if err != nil {
		t.Error(err)
	}

	err = catchPanic(func() { ir.GetAll() })
	if err != nil {
		t.Error(err)
	}

	err = catchPanic(func() { ir.GetByID(1) })
	if err != nil {
		t.Error(err)
	}

	err = catchPanic(func() { ir.GetByName("1") })
	if err != nil {
		t.Error(err)
	}

	err = catchPanic(func() { ir.GetByUserID(1) })
	if err != nil {
		t.Error(err)
	}

	err = catchPanic(func() { ir.Upsert(&Inventory{}) })
	if err != nil {
		t.Error(err)
	}
}

func TestUserRepo(t *testing.T) {
	dbfile := "testUserRepo.db"
	db, err := storm.Open(dbfile)
	if err != nil {
		t.Fatal(err)
	}
	defer func() {
		db.Close()
		os.Remove(dbfile)
	}()

	user := User{Email: "test@email.com", Password: []byte("123")}

	ur := newUserRepo(db)

	err = ur.Upsert(&user)
	if err != nil {
		t.Error("Could not insert a new user")
	}

	if user.ID == 0 {
		t.Errorf("Could not insert a new user record")
	}

	user.Email = "test2@email.com"

	err = ur.Upsert(&user)
	if err != nil {
		t.Error("Could not update a user")
	}

	if user.ID != 1 {
		t.Error("Could not update a user record. Expected userID to be", 1, "but was", user.ID)
	}

	user.ID = 14

	err = ur.Upsert(&user)
	if err == nil {
		t.Error("Expected an error when updating a non existing user")
	}

	_, err = ur.GetByEmail("test2@email.com")
	if err != nil {
		t.Error("Could not retrieve a user by email address, Error:", err)
	}

	_, err = ur.GetByID(1)
	if err != nil {
		t.Error("Could not retrieve a user by ID, Error:", err)
	}

	err = ur.Upsert(&User{})
	if err != nil {
		t.Error("Could not insert a new user")
	}

	allUsers := ur.GetAll()

	if len(allUsers) != 2 {
		t.Error("Did not retrieve ALL users")
	}
}

func TestInventoryRepo(t *testing.T) {
	dbfile := "testInventoryRepo.db"
	db, err := storm.Open(dbfile)
	if err != nil {
		t.Fatal(err)
	}
	defer func() {
		db.Close()
		os.Remove(dbfile)
	}()

	user := User{Email: "test@email.com", Password: []byte("123")}

	ur := newUserRepo(db)

	err = ur.Upsert(&user)
	if err != nil {
		t.Error("Could not insert a new user")
	}

	items := []Item{
		Item{Barcode: "123", Quantity: 1, SKU: "123", Title: "Item 1"},
		Item{Barcode: "456", Quantity: 2, SKU: "456", Title: "Item 2"},
	}

	inv := Inventory{Name: "warehouse1", Items: items, UserID: user.ID}

	ir := newInventoryRepo(db)

	err = ir.Upsert(&inv)
	if err != nil {
		t.Error("Could not insert a new inventory")
	}

	if inv.ID == 0 {
		t.Error("Could not insert a new inventory record")
	}

	inv.Name = "warehouse2"

	err = ir.Upsert(&inv)
	if err != nil {
		t.Error("Could not update an existing inventory")
	}

	if inv.ID != 1 {
		t.Error("Could not update an inventory record")
	}

	inv.ID = 23

	err = ir.Upsert(&inv)
	if err == nil {
		t.Error("Expected to not be able to update an non existing inventory")
	}

	inv.ID = 1

	_, err = ir.GetByID(1)
	if err != nil {
		t.Error("Could not retrieve inventory by ID, Error:", err)
	}

	_, err = ir.GetByName("warehouse2")
	if err != nil {
		t.Error("Could not retrieve inventory by Name, Error:", err)
	}

	_, err = ir.GetByUserID(1)
	if err != nil {
		t.Error("Could not retrieve inventory by UserID, Error:", err)
	}

	_, err = ir.GetByUserID(23)
	if err == nil {
		t.Error("Searching inventory by random user id should error")
	}

	inv.Items[0].Barcode = "987"

	err = ir.Upsert(&inv)
	if err != nil {
		t.Error("Could not update an existing inventory")
	}

	invPt, err := ir.GetByID(1)
	if err != nil {
		t.Error("Could not retrieve inventory with ID = 1, Error:", err)
	}

	if invPt.Items[0].Barcode != "987" {
		t.Error("Inventory Items did not update correctly")
	}

	err = ir.Upsert(&Inventory{})
	if err != nil {
		t.Error("Could not insert a new inventory")
	}

	allInv := ir.GetAll()
	if len(allInv) != 2 {
		t.Error("Did not retrieve all inventories")
	}

	err = ir.Delete(1)
	if err != nil {
		t.Error("Could not delete existing inventory", err)
	}

	_, err = ir.GetByID(1)
	if err == nil {
		t.Error("Retrieved non existing inventory")
	}
}
