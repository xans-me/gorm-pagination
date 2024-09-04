package transaction

import "time"

type Data struct {
	ID            int       `json:"id"`
	AccountNumber string    `json:"account_number"`
	TrxDate       time.Time `json:"trx_date"`
	TrxAmount     float64   `json:"trx_amount"`
	TrxType       string    `json:"trx_type"`
	CIF           string    `json:"cif"`
	CreateDate    time.Time `json:"create_date"`
}

func (*Data) TableName() string {
	return "db_name" // Pastikan nama tabel ini sesuai
}
