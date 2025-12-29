package data

import (
	"database/sql"
	"time"

	"github.com/j-elliott3/projects/crm/internal/domain"
)

type DealRepository struct {
	db *sql.DB
}

func NewDealRepository(db *sql.DB) DealRepository {
	return DealRepository{db: db}
}

func (r DealRepository) Create(d *domain.Deal) error {
	now := time.now().UTC()
	d.CreatedAt = now
	d.UpdatedAt = now

	var nextDue *string
	if d.NextActionDue != nil {
		s := d.NextActionDue.UTC().Format(time)
		nextDue = &s
	}

	res, err := r.db.Exec(`
		INSERT INTO deals (
			deal_name, customer_name, contact_person, phone, email,
			estimated_value, stage, source, created_at, updated_at,
			next_action, next_action_due
		)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)
	`,
		d.DealName,
		d.CustomerName,
		d.ContactPerson,
		d.Phone,
		d.Email,
		d.EstimatedValue,
		string(d.Stage),
		d.Source,
		d.CreatedAt.Format(time.RFC3339),
		d.UpdatedAt.Format(time.RFC3339),
		d.NextAction,
		nextDue,
	)
	if err != nil {
		return err
	}

	id, err := res.LastInsertId()
	if err != nil {
		return err
	}
	d.ID = id
	return nil
}

func (r DealRepository) ListAll() ([]domain.Deal, error) {
	rows, err := r.db.Query(`
		SELECT *
		FROM deals
		ORDER BY created_at DESC
	`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var deals []domain.Deal
	for rows.Next() {
		var (
			d 				domain.Deal
			createdStr 		string
			updatedStr 		string
			stageStr 		string
			nextActionDueS 	*string
		)

		err := rows.Scan(
			&d.ID,
            &d.DealName,
            &d.CustomerName,
            &d.ContactPerson,
            &d.Phone,
            &d.Email,
            &d.EstimatedValue,
            &stageStr,
            &d.Source,
            &createdStr,
            &updatedStr,
            &d.NextAction,
            &nextActionDueS,
		)
		if err != nil {
			return nil, err
		}

		d.Stage = domain.Stage(stageStr)

		if t, err := time.Parse(time.RFC3339, createdStr); err == nil {
            d.CreatedAt = t
        }
        if t, err := time.Parse(time.RFC3339, updatedStr); err == nil {
            d.UpdatedAt = t
        }
        if nextActionDueS != nil {
            if t, err := time.Parse(time.RFC3339, *nextActionDueS); err == nil {
                d.NextActionDue = &t
            }
        }

        deals = append(deals, d)
	}

	return deals, rows.Err()
}