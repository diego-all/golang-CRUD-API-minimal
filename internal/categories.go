package models

import (
	"context"
	"time"
)

type Category struct {
	Id          int       `json:"id"`
	Name        string    `json:"name,omitempty"`
	Description string    `json:"description,omitempty"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

func (c *Category) Insert(category Category) (int, error) {
	ctx, cancel := context.WithTimeout(context.Background(), dbTimeout)
	defer cancel()

	var newID int

	stmt := `insert into categories (name, description, created_at, updated_at)
	values ($1, $2, $3, $4) returning id`

	err := db.QueryRowContext(ctx, stmt,
		category.Name,
		category.Description,
		time.Now(),
		time.Now(),
	).Scan(&newID)

	if err != nil {
		return 0, err
	}

	return newID, nil

}

func (c *Category) GetOneById(id int) (*Category, error) {
	ctx, cancel := context.WithTimeout(context.Background(), dbTimeout)
	defer cancel()

	query := `select id, name, description, created_at, updated_at from categories where id = $1`

	var category Category
	row := db.QueryRowContext(ctx, query, id)

	err := row.Scan(
		&category.Id,
		&category.Name,
		&category.Description,
		&category.CreatedAt,
		&category.UpdatedAt,
	)

	if err != nil {
		return nil, err
	}

	return &category, nil

}

func (c *Category) Update(category Category) (int, error) {
	ctx, cancel := context.WithTimeout(context.Background(), dbTimeout)
	defer cancel() // resource leaks

	stmt := `update categories set
		name = $1,
		description = $2,
		updated_at = $3
		where id = $4`

	_, err := db.ExecContext(ctx, stmt,
		category.Name,
		category.Description,
		time.Now(),
		category.Id,
	)

	if err != nil {
		return 0, err
	}

	return 0, nil

}

func (c *Category) GetAll() ([]*Category, error) {
	ctx, cancel := context.WithTimeout(context.Background(), dbTimeout)
	defer cancel()

	query := `select id, name, description, created_at, updated_at from categories order by name`

	rows, err := db.QueryContext(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var categories []*Category

	for rows.Next() {
		var category Category
		err := rows.Scan(
			&category.Id,
			&category.Name,
			&category.Description,
			&category.CreatedAt,
			&category.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}

		categories = append(categories, &category)
	}

	return categories, nil
}

func (c *Category) DeleteByID(id int) error {

	ctx, cancel := context.WithTimeout(context.Background(), dbTimeout)
	defer cancel()

	stmt := `delete from categories where id = $1`

	_, err := db.ExecContext(ctx, stmt, id)
	if err != nil {
		return err
	}
	return nil
}
