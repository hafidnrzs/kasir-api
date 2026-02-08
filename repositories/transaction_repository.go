package repositories

import (
	"database/sql"
	"fmt"
	"kasir-api/models"
	"strings"
)

type TransactionRepository struct {
	db *sql.DB
}

func NewTransactionRepository(db *sql.DB) *TransactionRepository {
	return &TransactionRepository{db: db}
}

func (repo *TransactionRepository) CreateTransaction(items []models.CheckoutItem) (*models.Transaction, error) {
	tx, err := repo.db.Begin()
	if err != nil {
		return nil, err
	}
	defer tx.Rollback()

	// validasi item yang dicheckout tidak kosong
	if len(items) == 0 {
		return nil, fmt.Errorf("no items provided")
	}

	// validasi quantity item dan hitung jumlah item yang di-checkout
	qtyMap := make(map[int]int)
	for _, item := range items {
		if item.Quantity <= 0 {
			return nil, fmt.Errorf("quantity for product id %d must be greater than 0", item.ProductID)
		}
		qtyMap[item.ProductID] += item.Quantity
	}

	// siapkan string query ke dalam placeholders dan args-nya
	var (
		placeholders []string
		args         []interface{}
		idx          = 1
	)
	for id := range qtyMap {
		// buat placeholder sesuai format PostgreSQL: $1, $2, dst
		placeholders = append(placeholders, fmt.Sprintf("$%d", idx))
		args = append(args, id)
		idx++
	}

	// query produk sekaligus berdasarkan ProductID yang dibutuhkan
	// misal ... WHERE id IN ($1, $2, $3)
	// lalu args diisi dengan variable ProductID, misal []interface{1, 3, 4}
	query := fmt.Sprintf("SELECT id, name, price, stock FROM products WHERE id IN (%s)", strings.Join(placeholders, ","))
	rows, err := tx.Query(query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	// siapkan map "products" yang melakukan mapping hasil query (struct hasil scan). key = id
	products := map[int]struct {
		name  string
		price int
		stock int
	}{}
	var id, price, stock int
	var name string
	for rows.Next() {
		if err := rows.Scan(&id, &name, &price, &stock); err != nil {
			return nil, err
		}
		products[id] = struct {
			name  string
			price int
			stock int
		}{name, price, stock}
	}

	// cek produk ada dan stok cukup
	for id, qty := range qtyMap {
		p, ok := products[id]
		if !ok {
			return nil, fmt.Errorf("product id %d not found", id)
		}
		if p.stock < qty {
			return nil, fmt.Errorf("insufficient stock for product id %d", id)
		}
	}

	// inisialisasi subtotal -> jumlah total transaksi keseluruhan
	totalAmount := 0
	// inisialisasi transactionDetails -> nanti kita insert ke db
	details := make([]models.TransactionDetail, 0)

	// siapkan detail transaksi
	for _, item := range items {
		p := products[item.ProductID]
		subtotal := item.Quantity * p.price
		totalAmount += subtotal
		details = append(details, models.TransactionDetail{
			ProductID:   item.ProductID,
			ProductName: p.name,
			Quantity:    item.Quantity,
			Subtotal:    subtotal,
		})
	}

	// insert transaction
	var transactionID int
	err = tx.QueryRow("INSERT INTO transactions (total_amount) VALUES ($1) RETURNING id", totalAmount).Scan(&transactionID)
	if err != nil {
		return nil, err
	}

	// update stok produk setelah transaksi
	for id, qty := range qtyMap {
		res, err := tx.Exec("UPDATE products SET stock = stock - $1 WHERE id = $2 AND stock >= $1", qty, id)
		if err != nil {
			return nil, err
		}
		ra, err := res.RowsAffected()
		if err != nil {
			return nil, err
		}
		if ra == 0 {
			return nil, fmt.Errorf("insufficient stock for product id %d", id)
		}
	}

	// insert transaction details
	for i := range details {
		details[i].TransactionID = transactionID
		_, err := tx.Exec("INSERT INTO transaction_details (transaction_id, product_id, quantity, subtotal) VALUES ($1, $2, $3, $4)", transactionID, details[i].ProductID, details[i].Quantity, details[i].Subtotal)
		if err != nil {
			return nil, err
		}
	}

	if err := tx.Commit(); err != nil {
		return nil, err
	}

	return &models.Transaction{
		ID:          transactionID,
		TotalAmount: totalAmount,
		Details:     details,
	}, nil
}
