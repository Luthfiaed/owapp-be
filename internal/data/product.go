package data

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"time"
)

type ProductModel struct {
	DB *sql.DB
}

type Product struct {
	ID          int64     `json:"id"`
	Name        string    `json:"name"`
	Description string    `json:"description"`
	Price       int64     `json:"price"`
	IsDraft     bool      `json:"isDraft"`
	Image       string    `json:"image"`
	Reviews     []*Review `json:"reviews"`
}

type Review struct {
	ID        int64     `json:"id"`
	Review    string    `json:"review"`
	CreatedAt time.Time `json:"created_at"`
	UserID    int64     `json:"user_id"`
	Username  string    `json:"username"`
	ProductID int64     `json:"productId,omitempty"`
}

func (m ProductModel) GetByName(name string) ([]*Product, error) {
	query := ""
	/*
		VULNERABILITY POINT: SQL INJECTION
		TO PREVENT SQL INJECTION USE PARAMETERIZED QUERY
		REFER TO GetById function
	*/
	if name == "" {
		query = "SELECT * FROM products"
	} else {
		query = fmt.Sprintf("SELECT * FROM products WHERE name = '%s'", name)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	rows, err := m.DB.QueryContext(ctx, query)
	if err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			return nil, ErrRecordNotFound
		default:
			return nil, err
		}
	}
	defer rows.Close()

	products := make([]*Product, 0)
	for rows.Next() {
		var product Product
		err := rows.Scan(
			&product.ID,
			&product.Name,
			&product.Description,
			&product.Price,
			&product.IsDraft,
			&product.Image,
		)
		if err != nil {
			return nil, err
		}
		products = append(products, &product)
	}

	return products, nil
}

func (m ProductModel) GetById(id int) (*Product, error) {
	query := `
		SELECT products.id, products.name, products.description, products.price, products.image, products.is_draft 
		FROM products WHERE products.id = ?
	`

	var product Product
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	err := m.DB.QueryRowContext(ctx, query, id).Scan(
		&product.ID,
		&product.Name,
		&product.Description,
		&product.Price,
		&product.Image,
		&product.IsDraft,
	)
	if err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			return nil, ErrRecordNotFound
		default:
			return nil, err
		}
	}

	query = `
		SELECT reviews.id, reviews.review, reviews.created_at, users.id, users.username
		FROM reviews JOIN users ON reviews.user_id = users.id
		WHERE reviews.product_id = ?
	`
	reviews := make([]*Review, 0)
	rows, err := m.DB.QueryContext(ctx, query, id)
	if err != nil && !errors.Is(err, sql.ErrNoRows) {
		return nil, err
	}
	for rows.Next() {
		var review Review
		err := rows.Scan(
			&review.ID,
			&review.Review,
			&review.CreatedAt,
			&review.UserID,
			&review.Username,
		)
		if err != nil {
			return nil, err
		}
		reviews = append(reviews, &review)
	}

	product.Reviews = reviews

	return &product, nil
}

func (m ProductModel) Insert(data Review) error {
	query := `SELECT id FROM users WHERE users.username = ?`

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	var userId int
	err := m.DB.QueryRowContext(ctx, query, data.Username).Scan(&userId)
	if err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			return ErrRecordNotFound
		default:
			return err
		}
	}

	query = `INSERT INTO reviews (review, user_id, product_id) VALUES(?, ?, ?)`

	_, err = m.DB.ExecContext(ctx, query, data.Review, userId, data.ProductID)
	if err != nil {
		return err
	}

	return nil
}

func (m ProductModel) Update(data Review) error {
	query := `SELECT id FROM users WHERE users.username = ?`

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	var userId int
	err := m.DB.QueryRowContext(ctx, query, data.Username).Scan(&userId)
	if err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			return ErrRecordNotFound
		default:
			return err
		}
	}

	query = `UPDATE reviews SET review = ? WHERE id = ?`

	_, err = m.DB.ExecContext(ctx, query, data.Review, data.ID)
	if err != nil {
		return err
	}

	return nil
}
