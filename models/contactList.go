package models

import (
	"database/sql"
	"os"
)

type ContactsList struct {
	Name      string `json:"name"`
	Hash      string `json:"hash"`
	PublicKey string `json:"public_key"`
}

const CreateContactsListTable = `
CREATE TABLE IF NOT EXISTS contacts_list (
	id INTEGER PRIMARY KEY AUTOINCREMENT,
	name TEXT NOT NULL,
	hash TEXT NOT NULL,
	public_key TEXT NOT NULL
);`

func InitializeDB(dbPath string, createScript string) (*sql.DB, error) {
	_, err := os.Stat(dbPath)
	if os.IsNotExist(err) {
		db, err := sql.Open("sqlite3", dbPath)
		if err != nil {
			return nil, err
		}
		_, err = db.Exec(createScript)
		if err != nil {
			return nil, err
		}
		return db, nil
	} else if err != nil {
		return nil, err
	}
	db, err := sql.Open("sqlite3", dbPath)
	if err != nil {
		return nil, err
	}
	if err = db.Ping(); err != nil {
		return nil, err
	}
	return db, nil
}

type ContactsListDb struct {
	db *sql.DB
}

func NewContactsListRepo(db *sql.DB) *ContactsListDb {
	return &ContactsListDb{db: db}
}

func (c *ContactsListDb) insert(contactsList ContactsList) bool {
	_, err := c.db.Exec("INSERT INTO contacts_list (name, hash, public_key) VALUES (?, ?, ?)", contactsList.Name, contactsList.Hash, contactsList.PublicKey)
	return err == nil
}

func (c *ContactsListDb) getAll() []ContactsList {
	rows, err := c.db.Query("SELECT * FROM contacts_list")
	if err != nil {
		return nil
	}
	defer rows.Close()

	var contactsList []ContactsList
	for rows.Next() {
		var contact ContactsList
		err := rows.Scan(&contact.Name, &contact.Hash, &contact.PublicKey)
		if err != nil {
			return nil
		}
		contactsList = append(contactsList, contact)
	}
	return contactsList
}

func (c *ContactsListDb) deleteAll() bool {
	_, err := c.db.Exec("DELETE FROM contacts_list")
	return err == nil
}

func (c *ContactsListDb) delete(name string) bool {
	_, err := c.db.Exec("DELETE FROM contacts_list WHERE name = ?", name)
	return err == nil
}

func (c *ContactsListDb) updateName(id int, name string) bool {
	_, err := c.db.Exec("UPDATE contacts_list SET name = ? WHERE id = ?", name, id)
	return err == nil
}
