package main

import "github.com/asdine/storm"
import "errors"
import "log"

const (
	DefaultName = "<Inventory name not set>"
)

var (
	ErrorOperationNotPermitted = errors.New("Operation is not permitted")
	ErrorNameNotLongEnough     = errors.New("inventory name must contain at least 3 characters")
)

type InventoryServicer interface {
	CreateWithName(name string, userID int) error
	Create(inventory *Inventory, userID int) error
	Update(inventory *Inventory, userID int) error

	SoftDelete(id int, userID int) error
	//UpdateInventoryList(inventory []Inventory, userID int) error

	GetByUserID(userID int) ([]Inventory, error)
	GetByID(id int, userID int) (*Inventory, error)
}

type inventoryService struct {
	inventoryRepo InventoryRepoer
}

func NewInventoryService(inventoryRepo InventoryRepoer) InventoryServicer {
	return &inventoryService{inventoryRepo}
}

func NewInventoryServiceFromDB(db *storm.DB) InventoryServicer {
	return &inventoryService{inventoryRepo: newInventoryRepo(db)}
}

func (is *inventoryService) GetByID(id int, userID int) (*Inventory, error) {

	inv, err := is.inventoryRepo.GetByID(id)
	if err != nil {
		return nil, err
	}

	if inv.UserID != userID {
		return nil, ErrorOperationNotPermitted
	}

	return inv, nil
}

func (is *inventoryService) GetByUserID(userID int) ([]Inventory, error) {
	inventories, err := is.inventoryRepo.GetByUserID(userID)
	if err != nil {
		return []Inventory{}, err
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

func (is *inventoryService) CreateWithName(name string, userID int) error {
	if len(name) <= 3 {
		return ErrorNameNotLongEnough
	}

	inv := Inventory{Name: name, UserID: userID}

	err := is.inventoryRepo.Upsert(&inv)
	if err != nil {
		log.Println(err)
	}

	return err
}

func (is *inventoryService) Create(inventory *Inventory, userID int) error {
	if inventory.ID != 0 {
		return ErrorOperationNotPermitted
	}

	inventory.UserID = userID
	if len(inventory.Name) == 0 {
		inventory.Name = DefaultName
	}

	is.inventoryRepo.Upsert(inventory)

	return nil
}

func (is *inventoryService) Update(inventory *Inventory, userID int) error {
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
