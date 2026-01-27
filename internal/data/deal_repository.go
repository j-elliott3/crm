package data

import (
	"database/sql"
	"time"

	"github.com/j-elliott3/crm/internal/domain"
)

type DealRepository struct {
	db *sql.DB
}

func NewDealRepository(db *sql.DB) DealRepository {
	return DealRepository{db: db}
}

func (r DealRepository) Create(d *domain.Deal) error {
	now := time.Now().UTC()
	d.CreatedAt = now
	d.UpdatedAt = now

	var nextDue *string
	if d.NextActionDue != nil {
		s := d.NextActionDue.UTC().Format(time.RFC3339)
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
		SELECT id, deal_name, customer_name, contact_person, phone, email,
               estimated_value, stage, source, created_at, updated_at,
               next_action, next_action_due
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

		t, err := time.Parse(time.RFC3339, createdStr)
		if err != nil {
			return nil, err
		}
		d.CreatedAt = t

		t, err = time.Parse(time.RFC3339, updatedStr)
		if err != nil {
			return nil, err
		}
		d.UpdatedAt = t
        if nextActionDueS != nil {
            if t, err := time.Parse(time.RFC3339, *nextActionDueS); err == nil {
                d.NextActionDue = &t
            }
        }

        deals = append(deals, d)
	}

	return deals, rows.Err()
}

func (r DealRepository) GetByID(id int64) (domain.Deal, error) {
	var (
        d              domain.Deal
        createdStr     string
        updatedStr     string
        stageStr       string
        nextActionDueS *string
    )
	
	row := r.db.QueryRow(`
		SELECT id, deal_name, customer_name, contact_person, phone, email,
               estimated_value, stage, source, created_at, updated_at,
               next_action, next_action_due
		FROM deals
		WHERE id = ?
	`, id)

	err := row.Scan(
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
		return domain.Deal{}, err
	}

	d.Stage = domain.Stage(stageStr)

	t, err := time.Parse(time.RFC3339, createdStr)
	if err != nil {
		return domain.Deal{}, err
	}
	d.CreatedAt = t

	t, err = time.Parse(time.RFC3339, updatedStr)
	if err != nil {
		return domain.Deal{}, err
	}
	d.UpdatedAt = t
	
	if nextActionDueS != nil {
		if t, err := time.Parse(time.RFC3339, *nextActionDueS); err == nil {
			d.NextActionDue = &t
		}
	}

	return d, nil
}

func (r DealRepository) ListByStage(stage domain.Stage) ([]domain.Deal, error) {
	rows, err := r.db.Query(`
		SELECT id, deal_name, customer_name, contact_person, phone, email,
               estimated_value, stage, source, created_at, updated_at,
               next_action, next_action_due
		FROM deals
		WHERE stage = ?
		ORDER BY created_at DESC
	`, string(stage))
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

		t, err := time.Parse(time.RFC3339, createdStr)
		if err != nil {
			return nil, err
		}
		d.CreatedAt = t

		t, err = time.Parse(time.RFC3339, updatedStr)
		if err != nil {
			return nil, err
		}
		d.UpdatedAt = t
        if nextActionDueS != nil {
            if t, err := time.Parse(time.RFC3339, *nextActionDueS); err == nil {
                d.NextActionDue = &t
            }
        }

        deals = append(deals, d)
	}

	return deals, rows.Err()
}

func (r DealRepository) Update[T any](field, value T, id int64) error {
    d.UpdatedAt = time.Now().UTC()

    var nextDue *string
    if d.NextActionDue != nil {
        s := d.NextActionDue.UTC().Format(time.RFC3339)
        nextDue = &s
    }

    _, err := r.db.Exec(`
        UPDATE deals
        SET
            deal_name       = ?,
            customer_name   = ?,
            contact_person  = ?,
            phone           = ?,
            email           = ?,
            estimated_value = ?,
            stage           = ?,
            source          = ?,
            created_at      = ?,  -- optional; you could leave this alone
            updated_at      = ?,
            next_action     = ?,
            next_action_due = ?
        WHERE id = ?
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
        d.ID,
    )
    return err
}