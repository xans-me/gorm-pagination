package pfm

import "time"

type BrimoPFM struct {
	ID            int       `json:"id"`
	AccountNumber string    `json:"account_number"`
	TrxDate       time.Time `json:"trx_date"`
	TrxAmount     float64   `json:"trx_amount"` // Pastikan ini adalah float64
	TrxType       string    `json:"trx_type"`
	CIF           string    `json:"cif"`
	CreateDate    time.Time `json:"create_date"`
}

func (*BrimoPFM) TableName() string {
	return "brimo_pfm_xan" // Pastikan nama tabel ini sesuai
}
