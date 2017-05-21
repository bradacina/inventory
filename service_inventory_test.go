package main

import (
	"os"
	"testing"

	"github.com/asdine/storm"
)

const (
	userID        = 1
	anotherUserID = 2
	invID         = 1
	anotherInvID  = 2
)

func TestInventoryService(t *testing.T) {
	dbfile := "testInventoryService.db"
	db, err := storm.Open(dbfile)
	if err != nil {
		t.Fatal(err)
	}

	defer func() {
		db.Close()
		os.Remove(dbfile)
	}()

	is := NewInventoryServiceFromDB(db)

	ir := newInventoryRepo(db)
	_ = NewInventoryService(ir)

	var items []Item

	items = append(items, Item{Barcode: "123", Quantity: 12, SKU: "123", Title: "coca"})
	items = append(items, Item{Barcode: "456", Quantity: 1, SKU: "456", Title: "bmw"})

	inv := Inventory{Items: items}

	err = is.Create(&inv, userID)
	if err != nil {
		t.Error("Could not create an inventory")
	}

	if inv.UserID != userID {
		t.Error("Inventory.UserID should be set for a newly created inventory")
	}

	if inv.Name != DefaultName {
		t.Error("Inventory.Name should be set to default if it was empty")
	}

	if inv.ID != 1 {
		t.Error("Inventory.ID was not set properly")
	}

	err = is.Create(&inv, userID)
	if err == nil {
		t.Error("Should not be able to create an inventory that has ID")
	}

	inv2, err := is.GetByID(inv.ID, userID)
	if err != nil {
		t.Error("Could not retrieve inventory by ID")
	}

	if inv2 == nil {
		t.Error("The retrieved inventory was nil")
	}

	if inv2.UserID != userID {
		t.Error("Retrieved the wrong inventory")
	}

	inv3, err := is.GetByID(anotherInvID, userID)
	if err == nil {
		t.Error("Retrieved a non existing inventory")
	}

	if inv3 != nil {
		t.Error("Retrieved a non existing inventory")
	}

	_, err = is.GetByUserID(userID)
	if err != nil {
		t.Error("Could not retrieve existing inventory by userID")
	}

	_, err = is.GetByUserID(anotherUserID)
	if err == nil {
		t.Error("Retrieved inventoris that don't exist")
	}

	err = is.Update(&Inventory{}, userID)
	if err == nil {
		t.Error("Should not be able to update a non existing inventory")
	}

	inv.Name = "warehouse3"

	err = is.Update(&inv, anotherUserID)
	if err == nil {
		t.Error("Should not be able to update another user's inventory")
	}

	err = is.Update(&inv, userID)
	if err != nil {
		t.Error("Should be able to update inventory, Error:", err)
	}

	err = is.CreateWithName("xxx", userID)
	if err == nil {
		t.Error("Was able to create an inventory with name <= 3 characters")
	}

	err = is.CreateWithName("xxxx", userID)
	if err != nil {
		t.Error("Could not create inventory with name, Error:", err)
	}
}
