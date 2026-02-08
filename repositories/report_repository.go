package repositories

import (
	"database/sql"
	"kasir-api/models"
)

type ReportRepository struct {
	db *sql.DB
}

func NewReportRepository(db *sql.DB) *ReportRepository {
	return &ReportRepository{db: db}
}

func (repo *ReportRepository) GetReportToday() (*models.Report, error) {
	// query for total revenue and total transaksi
	var totalRevenue, totalTransaksi int
	err := repo.db.QueryRow(`
	       select coalesce(sum(total_amount),0) as total_revenue, count(id) as total_transaksi
	       from transactions
	       where date(created_at) = current_date;
       `).Scan(&totalRevenue, &totalTransaksi)
	if err != nil {
		return nil, err
	}

	// query for best selling product
	var nama string
	var qtyTerjual int
	err = repo.db.QueryRow(`
	       select p.name, coalesce(sum(td.quantity),0) as qty_terjual
	       from transaction_details td
	       join products p on td.product_id = p.id
	       join transactions t on td.transaction_id = t.id
	       where date(t.created_at) = current_date
	       group by p.name
	       order by qty_terjual desc
	       limit 1;
       `).Scan(&nama, &qtyTerjual)
	if err == sql.ErrNoRows {
		nama = ""
		qtyTerjual = 0
	} else if err != nil {
		return nil, err
	}

	report := &models.Report{
		TotalRevenue:   totalRevenue,
		TotalTransaksi: totalTransaksi,
		ProdukTerlaris: models.ProdukTerlaris{
			Nama:       nama,
			QtyTerjual: qtyTerjual,
		},
	}
	return report, nil
}

func (repo *ReportRepository) GetReportByDate(startDate string, endDate string) (*models.Report, error) {
	// query for total revenue and total transaksi
	var totalRevenue, totalTransaksi int
	err := repo.db.QueryRow(`
		select coalesce(sum(total_amount),0) as total_revenue, count(id) as total_transaksi
		from transactions
		where date(created_at) between $1 and $2;
	`, startDate, endDate).Scan(&totalRevenue, &totalTransaksi)
	if err != nil {
		return nil, err
	}

	// query for best selling product
	var nama string
	var qtyTerjual int
	err = repo.db.QueryRow(`
		select p.name, coalesce(sum(td.quantity),0) as qty_terjual
		from transaction_details td
		join products p on td.product_id = p.id
		join transactions t on td.transaction_id = t.id
		where date(t.created_at) between $1 and $2
		group by p.name
		order by qty_terjual desc
		limit 1;
	`, startDate, endDate).Scan(&nama, &qtyTerjual)
	if err == sql.ErrNoRows {
		nama = ""
		qtyTerjual = 0
	} else if err != nil {
		return nil, err
	}

	report := &models.Report{
		TotalRevenue:   totalRevenue,
		TotalTransaksi: totalTransaksi,
		ProdukTerlaris: models.ProdukTerlaris{
			Nama:       nama,
			QtyTerjual: qtyTerjual,
		},
	}
	return report, nil
}
