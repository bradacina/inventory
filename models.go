package main

type User struct {
	ID        int    `storm:"id,increment"`
	Email     string `storm:"unique"`
	Password  []byte
	IsAdmin   bool
	IsDeleted bool
}

type Item struct {
	Quantity int
	SKU      string
	Title    string
	Barcode  string
}

type Inventory struct {
	ID        int    `storm:"id,increment"`
	UserID    int    `storm:"index"`
	Name      string `storm:"index"`
	Items     []Item
	IsDeleted bool
}
