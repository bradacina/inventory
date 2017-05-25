package services

import (
	"os"
	"testing"

	"github.com/asdine/storm"
	"github.com/bradacina/inventory/db"
)

const (
	userID        = 1
	anotherUserID = 2
	invID         = 1
	anotherInvID  = 2
)

func TestInventoryService(t *testing.T) {
	dbfile := "testInventoryService.db"
	database, err := storm.Open(dbfile)
	if err != nil {
		t.Fatal(err)
	}

	defer func() {
		database.Close()
		os.Remove(dbfile)
	}()

	is := NewInventoryServiceFromDB(database)

	ur := db.NewUserRepo(database)

	ir := db.NewInventoryRepo(database)
	_ = NewInventoryService(ir, ur)

	var items []db.Item

	items = append(items, db.Item{Barcode: "123", Quantity: 12, SKU: "123", Title: "coca"})
	items = append(items, db.Item{Barcode: "456", Quantity: 1, SKU: "456", Title: "bmw"})

	inv := db.Inventory{UserID: userID, Items: items}

	ur.Upsert(&db.User{Email: "test@email.com"})

	err = is.Create(&inv, userID)
	if err != nil {
		t.Error("Could not create an inventory")
	}

	if inv.UserID != userID {
		t.Error("Inventory.UserID should be set for a newly created inventory")
	}

	if inv.Name != DefaultInventoryName {
		t.Error("Inventory.Name should be set to default if it was empty")
	}

	if inv.ID != 1 {
		t.Error("Inventory.ID was not set properly")
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

	_, err = is.GetByID(invID, anotherUserID)
	if err == nil {
		t.Error("Should not be able to retrieve other users inventory")
	}

	_, err = is.GetByUserID(userID)
	if err != nil {
		t.Error("Could not retrieve existing inventory by userID")
	}

	_, err = is.GetByUserID(anotherUserID)
	if err == nil {
		t.Error("Retrieved inventoris that don't exist")
	}

	err = is.Update(&db.Inventory{}, userID)
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

	inv4, err := is.GetByID(invID, userID)
	if err != nil {
		t.Error("Could not retrieve updated inventory, Error:", err)
	}

	if inv4.Name != "warehouse3" {
		t.Error("Inventory was not updated")
	}

	err = is.CreateWithName("xxx", userID)
	if err == nil {
		t.Error("Was able to create an inventory with name <= 3 characters")
	}

	err = is.CreateWithName("xxxx", userID)
	if err != nil {
		t.Error("Could not create inventory with name, Error:", err)
	}

	err = is.SoftDelete(inv.ID, userID)
	if err != nil {
		t.Error("Could not soft delete an inventory, Error:", err)
	}

	inv5, err := is.GetByID(invID, userID)
	if err != nil {
		t.Error("Could not retrieve a soft deleted inventory, Error:", err)
	}

	if !inv5.IsDeleted {
		t.Error("Inventory was not soft deleted")
	}
}
