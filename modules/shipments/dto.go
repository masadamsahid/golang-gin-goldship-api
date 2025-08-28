package shipments

// CREATE TABLE IF NOT EXISTS shipments (
//   id SERIAL PRIMARY KEY,
//   tracking_number VARCHAR(255) UNIQUE,
//   sender_id INT NOT NULL,
//   sender_name VARCHAR(255) NOT NULL,
//   sender_phone VARCHAR(20) NOT NULL,
//   sender_address TEXT NOT NULL,
//   recipient_name VARCHAR(255) NOT NULL,
//   recipient_address TEXT NOT NULL,
//   recipient_phone VARCHAR(20) NOT NULL,
//   item_name VARCHAR(255) NOT NULL,
//   item_weight DECIMAL(10, 2) NOT NULL,
//   distance DECIMAL(10, 2) NOT NULL,
//   status shipment_status_enum NOT NULL DEFAULT 'PENDING_PAYMENT',
//   created_at TIMESTAMP DEFAULT NOW() NOT NULL,
//   updated_at TIMESTAMP DEFAULT NULL,
//   CONSTRAINT fk_shipments_user FOREIGN KEY (sender_id) REFERENCES users(id)
// );

type CreateShipmentDto struct {
	SenderName       string  `json:"sender_name" binding:"required"`
	SenderPhone      string  `json:"sender_phone" binding:"required"`
	SenderAddress    string  `json:"sender_address" binding:"required"`
	RecipientName    string  `json:"recipient_name" binding:"required"`
	RecipientAddress string  `json:"recipient_address" binding:"required"`
	RecipientPhone   string  `json:"recipient_phone" binding:"required"`
	ItemName         string  `json:"item_name" binding:"required"`
	ItemWeight       float64 `json:"item_weight" binding:"required"`
	Distance         float64 `json:"distance" binding:"required"`
}

type TransitShipmentDto struct {
	BranchID float64 `json:"branch_id" binding:"required"`
}
