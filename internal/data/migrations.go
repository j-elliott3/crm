package data

import "database/sql"

func RunMigration(db *sql.DB) error {
	stmt := `
CREATE TABLE IF NOT EXISTS deals (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    deal_name TEXT NOT NULL,
    customer_name TEXT NOT NULL,
    contact_person TEXT,
    phone TEXT,
    email TEXT,
    estimated_value REAL,
    stage TEXT NOT NULL,
    source TEXT,
    created_at TEXT NOT NULL,
    updated_at TEXT NOT NULL,
    next_action TEXT,
    next_action_due TEXT
);
`
	_, err := db.exec(stmt)
	return err
}