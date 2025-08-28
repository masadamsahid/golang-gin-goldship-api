package models

import "time"

// CREATE TYPE payment_status_enum AS ENUM (
//   'PENDING',
//   'PAID',
//   'EXPIRED',
//   'CANCELLED'
// );

// CREATE TABLE IF NOT EXISTS payments (
//   id SERIAL PRIMARY KEY,
//   shipment_id INT NOT NULL,
//   amount INT NOT NULL,
//   payment_date TIMESTAMP DEFAULT NOW() NOT NULL,
//   invoice_id VARCHAR(255) NOT NULL,
// 	external_id VARCHAR(255) NOT NULL,
//   invoice_url TEXT NOT NULL,
//   status payment_status_enum NOT NULL DEFAULT 'PENDING',
//   created_at TIMESTAMP DEFAULT NOW() NOT NULL,
//   updated_at TIMESTAMP DEFAULT NULL,
//   CONSTRAINT fk_payment_shipment FOREIGN KEY (shipment_id) REFERENCES shipments(id)
// );

type PaymentStatus string

const (
	PaymentStatusPending   PaymentStatus = "PENDING"
	PaymentStatusPaid      PaymentStatus = "PAID"
	PaymentStatusExpired   PaymentStatus = "EXPIRED"
	PaymentStatusCancelled PaymentStatus = "CANCELLED"
)

type Payment struct {
	ID         int           `json:"id"`
	ShipmentID int           `json:"shipment_id"`
	Amount     int           `json:"amount"`
	PaidAt     *time.Time    `json:"paid_at"`
	ExpiredAt  *time.Time    `json:"expired_at"`
	InvoiceID  string        `json:"invoice_id"`
	ExternalID string        `json:"external_id"`
	InvoiceURL string        `json:"invoice_url"`
	Status     PaymentStatus `json:"status"`
	CreatedAt  time.Time     `json:"created_at"`
	UpdatedAt  *time.Time    `json:"updated_at"`
}
