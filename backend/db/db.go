package db

import (
	"database/sql"
	"fmt"
	"os"
	"pasteGo/backend/db/typesDB"

	_ "github.com/mattn/go-sqlite3"
)

type DBInstance struct {
	db *sql.DB
}

var instance *DBInstance

func CreateDirectory(path string) error {
	_, err := os.Stat(path)
	if os.IsNotExist(err) {
		err = os.Mkdir(path, 0660)
		if err != nil {
			fmt.Println(err)
			return err
		}
	}
	return nil
}

func createDBConnection() (*sql.DB, error) {
	err := CreateDirectory("data")
	if err != nil {
		return nil, err
	}

	db, err := sql.Open("sqlite3", "./data/database.db")
	db.SetMaxOpenConns(10)
	db.SetMaxIdleConns(5)
	if err != nil {
		return nil, err
	}
	return db, nil
}

func GetDBInstance() (*DBInstance, error) {
	if instance == nil {
		db, err := createDBConnection()
		if err != nil {
			return nil, err
		}
		instance = &DBInstance{db: db}
	}
	return instance, nil
}

func CloseDB() {
	if instance != nil {
		err := instance.db.Close()
		if err != nil {
			panic(err)
		}
		instance = nil
	}
}

func (instance *DBInstance) Init() error {
	_, err := instance.db.Exec("PRAGMA foreign_keys = ON;")
	if err != nil {
		return err
	}
	initSQL :=
		`CREATE TABLE IF NOT EXISTS users (
        id TEXT PRIMARY KEY,
		username TEXT NOT NULL UNIQUE,
		password TEXT NOT NULL
    );

	CREATE TABLE IF NOT EXISTS tokens (
        token TEXT PRIMARY KEY,
		user_id TEXT NOT NULL,
		FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE
    );

	CREATE TABLE IF NOT EXISTS pastes (
        id TEXT PRIMARY KEY,
		user_id TEXT NOT NULL,
		text TEXT NOT NULL,
		lifetime INTEGER NOT NULL,
		created INTEGER NOT NULL,
		updated INTEGER NOT NULL,
		password TEXT NOT NULL,
		public INTEGER NOT NULL DEFAULT 0,
		FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE
    );`

	_, err = instance.db.Exec(initSQL)
	if err != nil {
		return err
	}

	return nil
}

///USERS

func (instance *DBInstance) GetUserRecordById(id string) (typesDB.UserRecord, bool, error) {
	query := "SELECT username, password FROM users WHERE id = ?"
	record := typesDB.UserRecord{Id: id}
	err := instance.db.QueryRow(query, id).Scan(&record.Username, &record.Password)
	if err != nil {
		if err == sql.ErrNoRows {
			return typesDB.UserRecord{}, false, nil
		}
		return typesDB.UserRecord{}, false, err
	}
	return record, true, nil
}

func (instance *DBInstance) GetUserRecordByUsername(username string) (typesDB.UserRecord, bool, error) {
	query := "SELECT id, password FROM users WHERE username = ?"
	record := typesDB.UserRecord{Username: username}
	err := instance.db.QueryRow(query, username).Scan(&record.Id, &record.Password)
	if err != nil {
		if err == sql.ErrNoRows {
			return typesDB.UserRecord{}, false, nil
		}
		return typesDB.UserRecord{}, false, err
	}
	return record, true, nil
}

func (instance *DBInstance) AddUserRecord(record *typesDB.UserRecord) (bool, error) {
	exists, err := instance.checkRecordExistsUser(record.Username)
	if err != nil {
		return false, err
	}
	if exists {
		return false, nil
	}

	query := "INSERT INTO users (id, username, password) VALUES (?, ?, ?) ON CONFLICT(username) DO NOTHING"
	statement, err := instance.db.Prepare(query)
	if err != nil {
		return false, err
	}
	defer statement.Close()

	res, err := statement.Exec(record.Id, record.Username, record.Password)
	if err != nil {
		return false, err
	}
	rowsAffected, _ := res.RowsAffected()
	if rowsAffected > 0 {
		return true, nil
	}
	return false, nil
}

func (instance *DBInstance) checkRecordExistsUser(username string) (bool, error) {
	var exists bool
	query := "SELECT EXISTS(SELECT 1 FROM users WHERE username = ?)"
	err := instance.db.QueryRow(query, username).Scan(&exists)
	if err != nil {
		return false, err
	}
	return exists, nil
}

func (instance *DBInstance) EditUserRecord(record *typesDB.UserRecord) error {
	query := "UPDATE users SET username = ?, password = ? WHERE id = ?"
	statement, err := instance.db.Prepare(query)
	if err != nil {
		return err
	}
	defer statement.Close()

	_, err = statement.Exec(record.Username, record.Password, record.Id)
	return err
}

///PASTES

func (instance *DBInstance) GetPasteRecordById(pasteId string) (typesDB.PasteRecord, bool, error) {
	query := "SELECT user_id, text, created, updated, lifetime, password, public FROM pastes WHERE id = ?"
	record := typesDB.PasteRecord{Id: pasteId}
	err := instance.db.QueryRow(query, pasteId).Scan(&record.UserId, &record.Text, &record.Created, &record.Updated, &record.Lifetime, &record.Password, &record.Public)
	if err != nil {
		if err == sql.ErrNoRows {
			return typesDB.PasteRecord{}, false, nil
		}
		return typesDB.PasteRecord{}, false, err
	}
	return record, true, nil
}

func (instance *DBInstance) GetPasteRecordsByUserId(userId string) (*[]typesDB.PasteRecord, error) {
	query := "SELECT id, text, created, updated, lifetime, password, public FROM pastes WHERE user_id = ?"
	records := make([]typesDB.PasteRecord, 0, 10)
	rows, err := instance.db.Query(query, userId)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}

	for rows.Next() {
		var record typesDB.PasteRecord
		rows.Scan(&record.Id, &record.Text, &record.Created, &record.Updated, &record.Lifetime, &record.Password, &record.Public)
		record.UserId = userId
		records = append(records, record)
	}
	return &records, nil
}

func (instance *DBInstance) AddPasteRecord(record *typesDB.PasteRecord) (bool, error) {
	query := "INSERT INTO pastes (id, user_id, text, lifetime, created, updated, password, public) VALUES (?, ?, ?, ?, ?, ?, ?, ?)"
	statement, err := instance.db.Prepare(query)
	if err != nil {
		return false, err
	}
	defer statement.Close()

	res, err := statement.Exec(record.Id, record.UserId, record.Text, record.Lifetime, record.Created, record.Updated, record.Password, record.Public)
	if err != nil {
		return false, err
	}
	rowsAffected, _ := res.RowsAffected()
	if rowsAffected > 0 {
		return true, nil
	}
	return false, nil
}

func (instance *DBInstance) EditPasteRecord(record *typesDB.PasteRecord) error {
	query := "UPDATE pastes SET user_id = ?, text = ?, lifetime = ?, created = ?, updated = ?, password = ?, public = ? WHERE id = ?"
	statement, err := instance.db.Prepare(query)
	if err != nil {
		return err
	}
	defer statement.Close()

	_, err = statement.Exec(record.UserId, record.Text, record.Lifetime, record.Created, record.Updated, record.Password, record.Public, record.Id)
	return err
}

///TOKENS

func (instance *DBInstance) AddToken(record *typesDB.TokenRecord) (bool, error) {
	query := "INSERT INTO tokens (token, user_id) VALUES (?, ?)"
	statement, err := instance.db.Prepare(query)
	if err != nil {
		return false, err
	}
	defer statement.Close()

	res, err := statement.Exec(record.RefreshToken, record.UserId)
	if err != nil {
		return false, err
	}
	rowsAffected, _ := res.RowsAffected()
	if rowsAffected > 0 {
		return true, nil
	}
	return false, nil
}

func (instance *DBInstance) CheckIfExistToken(refreshToken string) (bool, error) {
	var exists bool
	query := "SELECT EXISTS(SELECT 1 FROM tokens WHERE token = ?)"
	err := instance.db.QueryRow(query, refreshToken).Scan(&exists)
	if err != nil {
		return false, err
	}
	return exists, nil
}

func (instance *DBInstance) GetTokenByUserId(userId string) (*[]typesDB.TokenRecord, error) {
	query := "SELECT user_id, token FROM tokens WHERE user_id = ?"
	records := make([]typesDB.TokenRecord, 0, 2)
	rows, err := instance.db.Query(query, userId)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}

	for rows.Next() {
		var record typesDB.TokenRecord
		rows.Scan(&record.UserId, &record.RefreshToken)
		records = append(records, record)
	}
	return &records, nil
}

func (instance *DBInstance) ChangeToken(record *typesDB.TokenRecord, oldToken string) error {
	query := "UPDATE tokens SET token = ? WHERE token = ?"
	statement, err := instance.db.Prepare(query)
	if err != nil {
		return err
	}
	defer statement.Close()

	_, err = statement.Exec(record.RefreshToken, oldToken)
	return err
}

func (instance *DBInstance) DeleteToken(token string) error {
	query := "DELETE FROM tokens WHERE token = ?"
	_, err := instance.db.Exec(query, token)
	return err
}

///USERS, PASTES

func (instance *DBInstance) DeleteRecord(id string, tableName string) error {
	query := fmt.Sprintf("DELETE FROM %s WHERE id = ?", tableName)
	_, err := instance.db.Exec(query, id)
	return err
}
