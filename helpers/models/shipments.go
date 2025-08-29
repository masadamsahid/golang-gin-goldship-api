package models

import "time"

const (
	StatusPendingPayment = "PENDING_PAYMENT"
	StatusReadyToPickup  = "READY_TO_PICKUP"
	StatusPickedUp       = "PICKED_UP"
	StatusInTransit      = "IN_TRANSIT"
	StatusDelivered      = "DELIVERED"
	StatusCancelled      = "CANCELLED"
)

type Shipment struct {
	ID               int     `json:"id"`
	TrackingNumber   string  `json:"tracking_number"`
	SenderID         int     `json:"sender_id"`
	SenderName       string  `json:"sender_name"`
	SenderPhone      string  `json:"sender_phone"`
	SenderAddress    string  `json:"sender_address"`
	RecipientName    string  `json:"recipient_name"`
	RecipientAddress string  `json:"recipient_address"`
	RecipientPhone   string  `json:"recipient_phone"`
	ItemName         string  `json:"item_name"`
	ItemWeight       float64 `json:"item_weight"`
	Distance         float64 `json:"distance"`
	Status           string  `json:"status"`
	CreatedAt        string  `json:"created_at"`
	UpdatedAt        *string `json:"updated_at"` // Use pointer for nullable timestamp

	Sender    *User             `json:"sender,omitempty"`
	Payment   *Payment          `json:"payment"`
	Histories []ShipmentHistory `json:"histories"`
}

type ShipmentHistory struct {
	ID         int    `json:"id"`
	ShipmentID int    `json:"shipment_id"`
	Status     string `json:"status"`
	Desc       string `json:"desc"`
	CourierID  *int   `json:"courier_id"` // Use pointer for nullable foreign key
	BranchID   *int   `json:"branch_id"`  // Use pointer for nullable foreign key
	Timestamp  string `json:"timestamp"`

	Shipment *Shipment `json:"shipment,omitempty"`

	Courier *ShCourier `json:"courier,omitempty"`
	Branch  *ShBranch  `json:"branch,omitempty"`
}

type ShBranch struct {
	ID        *uint      `json:"id,omitempty"`
	Name      *string    `json:"name,omitempty"`
	Phone     *string    `json:"phone,omitempty"`
	Address   *string    `json:"address,omitempty"`
	CreatedAt *time.Time `json:"created_at,omitempty"`
	UpdatedAt *time.Time `json:"updated_at,omitempty"`
}

type ShCourier struct {
	ID        *uint      `json:"id,omitempty"`
	Username  *string    `json:"username,omitempty"`
	Email     *string    `json:"email,omitempty"`
	Role      *string    `json:"string,omitempty"`
	CreatedAt *time.Time `json:"created_at,omitempty"`
	UpdatedAt *time.Time `json:"updated_at,omitempty"`
}
