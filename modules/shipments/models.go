package shipments

import (
	"github.com/masadamsahid/golang-gin-goldship-api/modules/branches"
	"github.com/masadamsahid/golang-gin-goldship-api/modules/payments"
	"github.com/masadamsahid/golang-gin-goldship-api/modules/users"
)

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

	Sender    *users.User       `json:"sender,omitempty"`
	Payment   *payments.Payment `json:"payment"`
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

	Shipment *Shipment        `json:"shipment,omitempty"`
	Courier  *users.User      `json:"courier,omitempty"`
	Branch   *branches.Branch `json:"branch,omitempty"`
}
