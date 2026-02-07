package repositories

import (
	"database/sql"
	"errors"
	"kasir-api/models"
)

type ProductRepository struct {
	db *sql.DB
}

func NewProductRepository(db *sql.DB) *ProductRepository {
	return &ProductRepository{db: db}
}

func (repo *ProductRepository) GetAll(name string) ([]models.Product, error) {
	query := "SELECT p.id, p.name, p.price, p.stock, c.name AS category FROM products p LEFT JOIN categories c ON p.category_id = c.id"

	var args []interface{}
	if name != "" {
		query += " WHERE p.name ILIKE $1"
		args = append(args, "%"+name+"%")
	}

	query += " ORDER BY p.id"

	rows, err := repo.db.Query(query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	products := make([]models.Product, 0)
	for rows.Next() {
		var p models.Product

		err := rows.Scan(&p.ID, &p.Name, &p.Price, &p.Stock, &p.Category)
		if err != nil {
			return nil, err
		}

		products = append(products, p)
	}

	return products, nil
}

func (repo *ProductRepository) Create(input *models.ProductInput) (*models.Product, error) {
	var id int
	query := "INSERT INTO products (name, price, stock, category_id) VALUES ($1, $2, $3, $4) RETURNING id"
	err := repo.db.QueryRow(query, input.Name, input.Price, input.Stock, input.CategoryID).Scan(&id)
	if err != nil {
		return nil, err
	}
	return repo.GetByID(id)
}

func (repo *ProductRepository) GetByID(id int) (*models.Product, error) {
	query := "SELECT p.id, p.name, p.price, p.stock, c.name AS category FROM products p LEFT JOIN categories c ON p.category_id = c.id WHERE p.id = $1"

	var p models.Product
	err := repo.db.QueryRow(query, id).Scan(&p.ID, &p.Name, &p.Price, &p.Stock, &p.Category)
	if err == sql.ErrNoRows {
		return nil, errors.New("product not found")
	}
	if err != nil {
		return nil, err
	}

	return &p, nil
}

func (repo *ProductRepository) Update(id int, input *models.ProductInput) (*models.Product, error) {
	query := "UPDATE products SET name = $1, price = $2, stock = $3, category_id = $4 WHERE id = $5"
	result, err := repo.db.Exec(query, input.Name, input.Price, input.Stock, input.CategoryID, id)
	if err != nil {
		return nil, err
	}

	rows, err := result.RowsAffected()
	if err != nil {
		return nil, err
	}

	if rows == 0 {
		return nil, errors.New("product not found")
	}

	return repo.GetByID(id)
}

func (repo *ProductRepository) Delete(id int) error {
	query := "DELETE FROM products WHERE id = $1"
	result, err := repo.db.Exec(query, id)
	if err != nil {
		return err
	}

	rows, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rows == 0 {
		return errors.New("product not found")
	}

	return nil
}
