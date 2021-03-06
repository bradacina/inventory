package services

import (
	"errors"
	"log"

	"github.com/asdine/storm"
	"github.com/bradacina/inventory/db"
)

const (
	DefaultInventoryName = "<Inventory name not set>"
	DefaultItemTitle     = "<Item name not set>"
	DefaultItemSKU       = "<Item SKU not set>"
	DefaultItemBarcode   = "<Item Barcode not set>"
)

var (
	ErrorOperationNotPermitted = errors.New("Operation is not permitted")
	ErrorNameNotLongEnough     = errors.New("inventory name must contain at least 3 characters")
)

type InventoryServicer interface {
	CreateWithName(name string, userID int) error
	Create(inventory *db.Inventory, userID int) error
	CreateByAdmin(inventory *db.Inventory) error
	Update(inventory *db.Inventory, userID int) error
	UpdateByAdmin(inventory *db.Inventory) error

	SoftDelete(id int, userID int) error
	SoftDeleteByAdmin(id int) error
	//UpdateInventoryList(inventory []Inventory, userID int) error

	GetByUserID(userID int) ([]db.Inventory, error)
	GetByID(id int, userID int) (*db.Inventory, error)
	GetByIDByAdmin(id int) (*db.Inventory, error)
	GetAll() ([]db.Inventory, error)
}

type inventoryService struct {
	inventoryRepo db.InventoryRepoer
	userRepo      db.UserRepoer
}

func NewInventoryService(
	inventoryRepo db.InventoryRepoer,
	userRepo db.UserRepoer) InventoryServicer {

	return &inventoryService{inventoryRepo, userRepo}
}

func NewInventoryServiceFromDB(database *storm.DB) InventoryServicer {
	return &inventoryService{
		inventoryRepo: db.NewInventoryRepo(database),
		userRepo:      db.NewUserRepo(database)}
}

func (is *inventoryService) GetByID(id int, userID int) (*db.Inventory, error) {

	inv, err := is.inventoryRepo.GetByID(id)
	if err != nil {
		return nil, err
	}

	if inv.UserID != userID {
		return nil, ErrorOperationNotPermitted
	}

	return inv, nil
}

func (is *inventoryService) GetByIDByAdmin(id int) (*db.Inventory, error) {

	inv, err := is.inventoryRepo.GetByID(id)
	if err != nil {
		return nil, err
	}

	return inv, nil
}

func (is *inventoryService) GetAll() ([]db.Inventory, error) {
	return is.inventoryRepo.GetAll()
}

func (is *inventoryService) GetByUserID(userID int) ([]db.Inventory, error) {
	inventories, err := is.inventoryRepo.GetByUserID(userID)
	if err != nil {
		return []db.Inventory{}, err
	}

	return inventories, nil
}

func (is *inventoryService) SoftDelete(id int, userID int) error {
	inv, err := is.inventoryRepo.GetByID(id)
	if err != nil {
		return err
	}

	if inv.UserID != userID {
		return ErrorOperationNotPermitted
	}

	inv.IsDeleted = true
	is.inventoryRepo.Upsert(inv)

	return nil
}

func (is *inventoryService) SoftDeleteByAdmin(id int) error {
	inv, err := is.inventoryRepo.GetByID(id)
	if err != nil {
		return err
	}

	inv.IsDeleted = true
	is.inventoryRepo.Upsert(inv)

	return nil
}

func (is *inventoryService) CreateWithName(name string, userID int) error {
	if len(name) <= 3 {
		return ErrorNameNotLongEnough
	}

	inv := db.Inventory{Name: name, UserID: userID}

	err := is.inventoryRepo.Upsert(&inv)
	if err != nil {
		log.Println(err)
	}

	return err
}

func (is *inventoryService) CreateByAdmin(inventory *db.Inventory) error {
	inventory.ID = 0

	_, err := is.userRepo.GetByID(inventory.UserID)
	if err != nil {
		return err
	}

	if len(inventory.Name) == 0 {
		inventory.Name = DefaultInventoryName
	}

	setDefaultValuesForItems(inventory.Items)

	err = is.inventoryRepo.Upsert(inventory)
	if err != nil {
		return err
	}

	return nil
}

func (is *inventoryService) Create(inventory *db.Inventory, userID int) error {
	inventory.ID = 0

	_, err := is.userRepo.GetByID(inventory.UserID)
	if err != nil {
		return err
	}

	inventory.UserID = userID
	if len(inventory.Name) == 0 {
		inventory.Name = DefaultInventoryName
	}

	setDefaultValuesForItems(inventory.Items)

	is.inventoryRepo.Upsert(inventory)

	return nil
}

func (is *inventoryService) Update(inventory *db.Inventory, userID int) error {
	if inventory.ID <= 0 {
		return ErrorNotFound
	}

	existingInv, err := is.inventoryRepo.GetByID(inventory.ID)
	if err != nil {
		return ErrorNotFound
	}

	if existingInv.UserID != inventory.UserID ||
		existingInv.UserID != userID {
		return ErrorOperationNotPermitted
	}

	is.inventoryRepo.Upsert(inventory)

	return nil
}

func (is *inventoryService) UpdateByAdmin(inventory *db.Inventory) error {
	if inventory.ID <= 0 {
		return ErrorNotFound
	}

	_, err := is.inventoryRepo.GetByID(inventory.ID)
	if err != nil {
		return ErrorNotFound
	}

	is.inventoryRepo.Upsert(inventory)

	return nil
}

func setDefaultValuesForItems(items []db.Item) {
	for i := range items {
		if len(items[i].Barcode) == 0 {
			items[i].Barcode = DefaultItemBarcode
		}

		if len(items[i].SKU) == 0 {
			items[i].SKU = DefaultItemSKU
		}

		if len(items[i].Title) == 0 {
			items[i].Title = DefaultItemTitle
		}
	}
}

/*
func (is *inventoryService) UpdateInventoryList(inventories []Inventory, userID int) error {
	var newInv []*Inventory
	var delInv []int
	var updateInv []*Inventory

	invMap := make(map[int]*Inventory)

	for k, v := range inventories {
		if v.ID == 0 {
			inventories[k].UserID = userID
			if len(inventories[k].Name) == 0 {
				inventories[k].Name = DefaultName
			}
			newInv = append(newInv, &inventories[k])
		} else if v.UserID != userID {
			return ErrorOperationNotPermitted
		} else {
			invMap[v.ID] = &inventories[k]
		}
	}

	existingInvs := is.GetByUserID(userID)

	for k, v := range existingInvs {
		if mapVal, ok := invMap[v.ID]; ok {
			updateInv = append(updateInv, mapVal)
		} else {
			delInv = append(delInv, existingInvs[k].ID)
		}
	}

	// create new inventories
	for k := range newInv {
		is.inventoryRepo.Upsert(newInv[k])
	}

	// delete old inventories
	for k := range delInv {
		is.inventoryRepo.Delete(delInv[k])
	}

	for k := range updateInv {
		is.inventoryRepo.Upsert(updateInv[k])
	}

	return nil
}
*/
