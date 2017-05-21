package main

type UserRepoer interface {
	GetByID(id int) (*User, error)
	GetByEmail(email string) (*User, error)
	GetAll() ([]User, error)

	Upsert(user *User) error
}

type InventoryRepoer interface {
	GetByID(id int) (*Inventory, error)
	GetByUserID(userID int) ([]Inventory, error)
	GetByName(name string) ([]Inventory, error)
	GetAll() ([]Inventory, error)

	Upsert(inventory *Inventory) error

	Delete(id int) error
}
