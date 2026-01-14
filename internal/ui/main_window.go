package ui

import (
	"fmt"
	"strconv"

    "fyne.io/fyne/v2"
    "fyne.io/fyne/v2/container"
    "fyne.io/fyne/v2/widget"

	"github.com/j-elliott3/crm/internal/data"
	"github.com/j-elliott3/crm/internal/domain"
)

type MainWindow struct {
	app 	fyne.App
	window 	fyne.Window
	repo 	data.DealRepository
	deals 	[]domain.Deal
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

	newDealBtn := widget.NewButton("New Deal", func() {
		mw.showNewDealForm()
	})

	content := container.NewBorder(
		container.NewVBox(newDealBtn, mw.status),
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
		mw.status.SetText("Error loading deals: " + err.Error())
		return
	}
	mw.deals = deals
	mw.list.Refresh()
}

func (mw *MainWindow) showNewDealForm() {
	dealNameEntry := widget.NewEntry()
	customerEntry := widget.NewEntry()
	contactEntry := widget.NewEntry()
	phoneEntry := widget.NewEntry()
	emailEntry := widget.NewEntry()
	valueEntry := widget.NewEntry()
	nextActionEntry := widget.NewEntry()


	// Dropdown for stage
	stageSelect := widget.NewSelect(
		[]string{
			string(domain.StageNewLead),
			string(domain.StageQualified),
			string(domain.StageSurvey),
			string(domain.StageQuoteSent),
			string(domain.StageWon),
			string(domain.StageLost),
		},
		nil,
	)
	stageSelect.SetSelected(string(domain.StageNewLead))

	form := &widget.Form{
		Items: []*widget.FormItem{
			{Text: "Deal Name", Widget: dealNameEntry},
            {Text: "Customer Name", Widget: customerEntry},
            {Text: "Contact Person", Widget: contactEntry},
            {Text: "Phone", Widget: phoneEntry},
            {Text: "Email", Widget: emailEntry},
            {Text: "Estimated Value", Widget: valueEntry},
            {Text: "Stage", Widget: stageSelect},
            {Text: "Next Action", Widget: nextActionEntry},
		},
		OnSubmit: func() {
			if err := mw.createDealFromForm(
				dealNameEntry.Text,
                customerEntry.Text,
                contactEntry.Text,
                phoneEntry.Text,
                emailEntry.Text,
                valueEntry.Text,
                stageSelect.Selected,
                nextActionEntry.Text,
			); err != nil {
				mw.status.SetText("Error: " + err.Error())
			} else {
				mw.status.SetText("Deal created")
				mw.refreshDeals()
			}
		},
		OnCancel: func() {},
	}

	dialog := widget.NewModalPopUp(
		container.NewVBox(
			widget.NewLabel("New Deal"),
			form,
		),
		mw.window.Canvas(),
	)

	// Close the dialog when form submits or cancels
    oldOnSubmit := form.OnSubmit
    form.OnSubmit = func() {
        oldOnSubmit()
        dialog.Hide()
    }
    oldOnCancel := form.OnCancel
    form.OnCancel = func() {
        oldOnCancel()
        dialog.Hide()
    }

    dialog.Show()
}


func (mw *MainWindow) createDealFromForm(
	dealName, customer, contact, phone, email, valueStr, stageStr, nextAction string,
) error {
	if dealName == "" || customer == "" {
		return fmt.Errorf("deal name and customer name are required")
	}

	var value float64
	if valueStr != "" {
		var err error
		value, err = strconv.ParseFloat(valueStr, 64)
		if err != nil {
			return fmt.Errorf("invalid estimated value: %w", err)
		}
	}

	d := &domain.Deal{
		DealName:       dealName,
        CustomerName:   customer,
        ContactPerson:  contact,
        Phone:          phone,
        Email:          email,
        EstimatedValue: value,
        Stage:          domain.Stage(stageStr),
        Source:         "manual",
        NextAction:     nextAction,
	}

	return mw.repo.Create(d)
}