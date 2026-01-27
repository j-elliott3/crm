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
	app 		fyne.App
	window 		fyne.Window
	repo 		data.DealRepository
	deals 		[]domain.Deal
	list 		*widget.List
	status 		*widget.Label
	stageFilter *widget.Select

	// selectedIndex tracks which row in mw.deals is currently selected in the list.
    // -1 means “no selection”.
	selectedIndex 	int
}

func NewMainWindow(app fyne.App, repo data.DealRepository) fyne.Window {
	mw := &MainWindow{
		app: app,
		repo: repo,
		status: widget.NewLabel(""),
		selectedIndex: -1,
	}

	mw.window = app.NewWindow("CRM Pipeline")
	mw.window.Resize(fyne.NewSize(800,600))

	mw.setupUI()
	mw.refreshDeals()

	return mw.window
}

func (mw *MainWindow) setupUI() {
	mw.stageFilter = widget.NewSelect(
		[]string{
			"All",
			string(domain.StageNewLead),
			string(domain.StageQualified),
			string(domain.StageSurvey),
			string(domain.StageQuoteSent),
			string(domain.StageWon),
			string(domain.StageLost),
		},
		nil,
	)
	mw.stageFilter.SetSelected("All")
	mw.stageFilter.OnChanged = func(selected string) {
    	mw.applyStageFilter(selected)
	}
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
	// NOTE: OnSelected is a Fyne List callback.
	// It is invoked when the user selects a row in the list (via click).
	// We store the index so other parts of the UI (e.g. Edit button) know which Deal to edit.
	mw.list.OnSelected = func(id widget.ListItemID) {
    	mw.selectedIndex = int(id)
	}
	mw.list.OnUnselected = func(id widget.ListItemID) {
    // If the selected row is unselected, reset to “no selection”.
    if mw.selectedIndex == int(id) {
        mw.selectedIndex = -1
    	}
	}

	newDealBtn := widget.NewButton("New Deal", func() { mw.showNewDealFormForNew() })

	editDealBtn := widget.NewButton("Edit Selected", func() {
		mw.showEditDealForm()
	})

	topBar := container.NewHBox(
		newDealBtn,
		editDealBtn,
		widget.NewLabel("Stage:"),
		mw.stageFilter,
	)

	content := container.NewBorder(
		container.NewVBox(topBar, mw.status),
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

				selected := mw.stageFilter.Selected
				if selected == "" || selected == "All" {
					mw.refreshDeals()
				} else {
					mw.applyStageFilter(selected)
				}
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

func (mw *MainWindow) applyStageFilter(selected string) {
    if selected == "All" || selected == "" {
        mw.refreshDeals()
        return
    }

    deals, err := mw.repo.ListByStage(domain.Stage(selected))
    if err != nil {
        mw.status.SetText("Error loading deals: " + err.Error())
        return
    }
    mw.deals = deals
    mw.list.Refresh()
}

// showDealForm opens a modal dialog to create or edit a deal.
//
// If d is nil, it shows an empty form for a new deal.
// If d is non-nil, it pre-fills the form with d's fields and, on submit,
// updates that deal via the repository.
func (mw *MainWindow) showDealForm(d *domain.Deal) {
    // Build widgets for all fields.
    dealNameEntry := widget.NewEntry()
    customerEntry := widget.NewEntry()
    contactEntry := widget.NewEntry()
    phoneEntry := widget.NewEntry()
    emailEntry := widget.NewEntry()
    valueEntry := widget.NewEntry()
    nextActionEntry := widget.NewEntry()

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

    // If editing, pre-fill fields from d.
    if d != nil {
        dealNameEntry.SetText(d.DealName)
        customerEntry.SetText(d.CustomerName)
        contactEntry.SetText(d.ContactPerson)
        phoneEntry.SetText(d.Phone)
        emailEntry.SetText(d.Email)
        if d.EstimatedValue != 0 {
            valueEntry.SetText(fmt.Sprintf("%.2f", d.EstimatedValue))
        }
        nextActionEntry.SetText(d.NextAction)
        stageSelect.SetSelected(string(d.Stage))
    } else {
        // For a new deal, pick a sensible default.
        stageSelect.SetSelected(string(domain.StageNewLead))
    }

    // Decide if this is a "new" or "edit" form by whether d is nil.
    title := "New Deal"
    if d != nil {
        title = "Edit Deal"
    }

    // We create the form. OnSubmit will either Create or Update, depending on d.
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
    }

    // We'll fill OnSubmit and OnCancel below after building the dialog.
    dialog := widget.NewModalPopUp(
        container.NewVBox(
            widget.NewLabel(title),
            form,
        ),
        mw.window.Canvas(),
    )

	form.OnSubmit = func() {
        // Extract values from fields and call the appropriate repo method.
        if err := mw.saveDealFromForm(
            d, // may be nil (new) or non-nil (edit)
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
            if d == nil {
                mw.status.SetText("Deal created")
            } else {
                mw.status.SetText("Deal updated")
            }

            // Reapply the current stage filter after change.
            selected := mw.stageFilter.Selected
            if selected == "" || selected == "All" {
                mw.refreshDeals()
            } else {
                mw.applyStageFilter(selected)
            }
            dialog.Hide()
        }
    }

    form.OnCancel = func() {
        dialog.Hide()
    }

    dialog.Show()
}

// showNewDealFormForNew opens a blank form for creating a new deal.
func (mw *MainWindow) showNewDealFormForNew() {
    mw.showDealForm(nil)
}

// showEditDealForm opens a pre-filled form for the currently selected deal.
func (mw *MainWindow) showEditDealForm() {
    if mw.selectedIndex < 0 || mw.selectedIndex >= len(mw.deals) {
        mw.status.SetText("No deal selected to edit")
        return
    }

    // Important: get a copy of the selected deal.
    d := mw.deals[mw.selectedIndex]
    mw.showDealForm(&d)
}

// saveDealFromForm validates and persists a deal based on form input.
//
// If base is nil, it constructs a new Deal and calls repo.Create.
// If base is non-nil, it modifies that Deal and calls repo.Update.
func (mw *MainWindow) saveDealFromForm(
    base *domain.Deal,
    dealName, customer, contact, phone, email,
    valueStr, stageStr, nextAction string,
) error {
    if dealName == "" || customer == "" {
        return fmt.Errorf("deal name and customer name are required")
    }

    var value float64
    if valueStr != "" {
        v, err := strconv.ParseFloat(valueStr, 64)
        if err != nil {
            return fmt.Errorf("invalid estimated value: %w", err)
        }
        value = v
    }

    if base == nil {
        // NEW deal
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

    // EDIT existing deal: modify the existing struct and call Update.
    base.DealName = dealName
    base.CustomerName = customer
    base.ContactPerson = contact
    base.Phone = phone
    base.Email = email
    base.EstimatedValue = value
    base.Stage = domain.Stage(stageStr)
    base.NextAction = nextAction

    return mw.repo.Update(base)
}