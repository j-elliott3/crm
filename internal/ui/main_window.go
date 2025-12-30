package ui

import (
	"fmt"
    "time"

    "fyne.io/fyne/v2"
    "fyne.io/fyne/v2/container"
    "fyne.io/fyne/v2/widget"

	"github.com/j-elliott3/crm/internal/data"
	"github.com/j-elliott3/crm/internal/domain"
)

type MainWindow struct {
	app 	fyne.app
	window 	fyne.window
	repo 	data.DealRepository
	deals 	[]domain.deal
	list 	*widget.List
	status 	*widget.Label
}

func NewMainWindow(app fyne.App, repo data.DealRepository) fyne.Window {
	mw := &MainWindow{
		app: app,
		repo: repo,
		status: widget.NewLabel(""),
	}

	mw.window = app.NewWindow("CRM Pipeline")
	mw.window.Resize(fyne.NewSize(800,600))

	mw.setupUI()
	mw.refreshDeals()

	return mw.window
}

func (mw *MainWindow) setupUI() {
	// List widget: displays deals slice
	mw.list = widget.NewList(
		func() int {
			return len(mw.deals)
		},
		func() fyne.CanvasObject {
			return widget.NewLabel("deal")
		},
		func(i widget.ListItemID, o fyne.CanvasObject) {
			if i < 0 || i >= len(mw.deals) {
				return
			}
			d := mw.deals[i]
			text := fmt.Sprintf("#%d | %s | %s | $%.2f",
                d.ID, d.DealName, d.Stage, d.EstimatedValue)
			o.(*widget.Label).SetText(text)
		},
	)

	newDummyBtn := widget.NewButton("Add Dummy Deal", func() {
		if err := mw.addDummyDeal(); err != nil {
			mw.status.SetText("Error: " + err)
		} else {
			mw.status.SetText("Dummy deal added")
			mw.refreshDeals()
		}
	})

	content := container.NewBorder(
		container.NewVBox(newDummyBtn, mw.status),
		nil,
		nil,
		nil,
		mw.list,
	)

	mw.window.SetContent(content)
}

func (mw *MainWindow) refreshDeals() {
	deals, err := mw.repo.ListAll()
	if err != nil {
		mw.status.SetText("Error loading deal")
		return
	}
	mw.deals = deals
	mw.list.Refresh()
}

func (mw *MainWindow) addDummyDeal() error {
	d := &domain.Deal{
		DealName:      "Test Deal",
        CustomerName:  "Acme Corp",
        ContactPerson: "Jane Doe",
        Phone:         "555-1234",
        Email:         "jane@example.com",
        EstimatedValue: 10000,
        Stage:         domain.StageNewLead,
        Source:        "manual",
        NextAction:    "Call customer",
        // NextActionDue: can leave nil or set:
        NextActionDue: func() *time.Time {
            t := time.Now().Add(24 * time.Hour).UTC()
            return &t
        }(),
	}
	return mw.repo.Create(d)
}