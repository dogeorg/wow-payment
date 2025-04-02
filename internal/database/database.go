package database

import (
	"database/sql"

	"github.com/dogeorg/wow-payment/internal/models"

	_ "github.com/mattn/go-sqlite3"
)

func InitDB(dbPath string) (*sql.DB, error) {
	db, err := sql.Open("sqlite3", dbPath)
	if err != nil {
		return nil, err
	}

	createTable := `
    CREATE TABLE IF NOT EXISTS shibes (
        id INTEGER PRIMARY KEY AUTOINCREMENT,
        name TEXT,
        email TEXT,
        country TEXT,
        address TEXT,
        postalCode TEXT,
        dogeAddress TEXT,
        size TEXT,
        bname TEXT,
        bemail TEXT,
        bcountry TEXT,
        baddress TEXT,
        bpostalCode TEXT,
        amount REAL,
        paytoDogeAddress TEXT,
        paid INTEGER DEFAULT 0,
        date TIMESTAMP DEFAULT CURRENT_TIMESTAMP
    )`
	_, err = db.Exec(createTable)
	return db, err
}

func InsertShibe(db *sql.DB, shibe models.RegistrationRequest) (int64, error) {
	stmt, err := db.Prepare(`
        INSERT INTO shibes (name, email, country, address, postalCode, dogeAddress, size, 
            bname, bemail, bcountry, baddress, bpostalCode, amount, paytoDogeAddress) 
        VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)`)
	if err != nil {
		return 0, err
	}
	defer stmt.Close()

	result, err := stmt.Exec(shibe.Name, shibe.Email, shibe.Country, shibe.Address,
		shibe.PostalCode, shibe.DogeAddress, shibe.Size, shibe.BName, shibe.BEmail,
		shibe.BCountry, shibe.BAddress, shibe.BPostalCode, shibe.Amount, shibe.PaytoDogeAddress)
	if err != nil {
		return 0, err
	}
	return result.LastInsertId()
}
