package main

import (
	"log"

	"fyne.io/fyne/v2/app"

	"github.com/j-elliott3/crm/internal/data"
	"github.com/j-elliott3/crm/internal/ui"
)

func main() {
	db, err := data.OpenDB()
	if err != nil {
		log.Fatalf("open db: %v", err)
	}
	defer db.Close()

	if err := data.RunMigrations(db); err != nil {
		log.Fatalf("migrations: %v", err)
	}

	dealRepo := data.NewDealRepository(db)

	fyneApp := app.New()
	w := ui.NewMainWindow(fyneApp, dealRepo)
	w.ShowAndRun()
}