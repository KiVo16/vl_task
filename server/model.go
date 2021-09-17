package main

type User struct {
	ID      int
	Name    string
	Records []UserRecord `json:",omitempty"`
}

type Record struct {
	ID   int
	Name string
	Type string
}

type UserRecord struct {
	RecordID int
	UserID   int
	Record   Record `json:",omitempty"`
	User     User   `json:",omitempty"`
}