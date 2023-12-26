package models

import (
	"database/sql"
	"os"

	_ "github.com/mattn/go-sqlite3"
)

const (
	dbPath = "client.db"
)

var db *sql.DB

func InitializeDB() {
	var err error
	if _, err := os.Stat(dbPath); os.IsNotExist(err) {
		file, err := os.Create(dbPath)
		if err != nil {
			panic(err)
		}
		file.Close()
	}

	db, err = sql.Open("sqlite3", dbPath)
	if err != nil {
		panic(err)
	}
	_, err = db.Exec(CreateContactsListTable)
	if err != nil {
		panic(err)
	}
}

func CloseDB() {
	db.Close()
}

func ContactsInsert(contactsList ContactsList) bool {
	_, err := db.Exec("INSERT INTO contacts_list (name, uid) VALUES (?, ?)", contactsList.Name, contactsList.Uid)
	return err == nil
}

func ContactsGetAll() []ContactsList {
	rows, err := db.Query("SELECT name,uid FROM contacts_list")
	if err != nil {
		return nil
	}
	defer rows.Close()

	var contactsList []ContactsList
	for rows.Next() {
		var contact ContactsList
		err := rows.Scan(&contact.Name, &contact.Uid)
		if err != nil {
			return nil
		}
		contactsList = append(contactsList, contact)
	}
	return contactsList
}

func ContactsDeleteAll() bool {
	_, err := db.Exec("DELETE FROM contacts_list")
	return err == nil
}

func ContactsDelete(name string) bool {
	_, err := db.Exec("DELETE FROM contacts_list WHERE name = ?", name)
	return err == nil
}

func ContactsUpdateName(id int, name string) bool {
	_, err := db.Exec("UPDATE contacts_list SET name = ? WHERE id = ?", name, id)
	return err == nil
}
