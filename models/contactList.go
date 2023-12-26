package models

type ContactsList struct {
	Name string `json:"name"`
	Uid  string `json:"uid"`
}

const CreateContactsListTable string = `
	CREATE TABLE IF NOT EXISTS contacts_list (
		id INTEGER PRIMARY KEY,
		name TEXT NOT NULL,
		uid TEXT NOT NULL
	);`
