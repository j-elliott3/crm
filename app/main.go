package main

import (
	"log"
	"github.com/j-elliott3/projects/crm/internal/data"
)

func main() {
	db, err := data.openDB()
	if err != nil {
		log.Fatalf("open db: %v", err)
	}
	defer db.Close()

	if err := data.RunMigration(db); err != nil {
		log.Fatalf("migrations: %v", err)
	}
}