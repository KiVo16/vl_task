package main

type User struct {
	ID      int          `json:"id"`
	Name    string       `json:"name"`
	Records []UserRecord `json:",omitempty"`
}

type Record struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
	Type string `json:"type"`
}

type UserRecord struct {
	RecordID int    `gorm:"primaryKey"`
	UserID   int    `gorm:"primaryKey"`
	Record   Record `json:",omitempty"`
	User     User   `json:",omitempty"`
}
