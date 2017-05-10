package main

import (
	"github.com/asdine/storm"
)

type stormUserRepo struct {
	db *storm.DB
}

type stormInventoryRepo struct {
	db *storm.DB
}

func newUserRepo(db *storm.DB) UserRepoer {
	return &stormUserRepo{db: db}
}

func newInventoryRepo(db *storm.DB) InventoryRepoer {
	return &stormInventoryRepo{db: db}
}

func (ur *stormUserRepo) GetByID(id int) (*User, error) {
	if ur.db == nil {
		panic("Database is not instantiated")
	}

	var user User
	err := ur.db.One("ID", id, &user)
	return &user, err
}

func (ur *stormUserRepo) GetByEmail(email string) (*User, error) {
	if ur.db == nil {
		panic("Database is not instantiated")
	}

	var user User
	err := ur.db.One("Email", email, &user)
	return &user, err
}

func (ur *stormUserRepo) GetAll() *[]User {
	if ur.db == nil {
		panic("Database is not instantiated")
	}

	var users []User
	ur.db.All(&users)
	return &users
}

func (ur *stormUserRepo) Upsert(user *User) {
	if ur.db == nil {
		panic("Database is not instantiated")
	}

	if user.ID > 0 {
		ur.db.Update(user)
	} else {
		ur.db.Save(user)
	}
}

func (ir *stormInventoryRepo) GetByID(id int) (*Inventory, error) {
	if ir.db == nil {
		panic("Database is not instantiated")
	}

	var inventory Inventory

	err := ir.db.One("ID", id, &inventory)
	return &inventory, err
}

func (ir *stormInventoryRepo) GetByUserID(userID int) (*[]Inventory, error) {
	if ir.db == nil {
		panic("Database is not instantiated")
	}

	var inv []Inventory
	err := ir.db.Find("UserID", userID, &inv)

	return &inv, err
}

func (ir *stormInventoryRepo) GetByName(name string) (*[]Inventory, error) {
	if ir.db == nil {
		panic("Database is not instantiated")
	}

	var inv []Inventory
	err := ir.db.Find("Name", name, &inv)

	return &inv, err
}

func (ir *stormInventoryRepo) GetAll() *[]Inventory {
	if ir.db == nil {
		panic("Database is not instantiated")
	}

	var invs []Inventory
	ir.db.All(&invs)
	return &invs
}

func (ir *stormInventoryRepo) Upsert(inventory *Inventory) {
	if ir.db == nil {
		panic("Database is not instantiated")
	}

	if inventory.ID > 0 {
		ir.db.Update(inventory)
	} else {
		ir.db.Save(inventory)
	}
}
