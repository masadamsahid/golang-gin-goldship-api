package webhooks

// {
//     "id": "579c8d61f23fa4ca35e52da4",
//     "external_id": "invoice_123124123",
//     "user_id": "5781d19b2e2385880609791c",
//     "is_high": true,
//     "payment_method": "BANK_TRANSFER",
//     "status": "PAID",
//     "merchant_name": "Xendit",
//     "amount": 50000,
//     "paid_amount": 50000,
//     "bank_code": "PERMATA",
//     "paid_at": "2016-10-12T08:15:03.404Z",
//     "payer_email": "wildan@xendit.co",
//     "description": "This is a description",
//     "adjusted_received_amount": 47500,
//     "fees_paid_amount": 0,
//     "updated": "2016-10-10T08:15:03.404Z",
//     "created": "2016-10-10T08:15:03.404Z",
//     "currency": "IDR",
//     "payment_channel": "PERMATA",
//     "payment_destination": "888888888888"
// }

type XenditInvoiceNotificationDto struct {
	ID                     string  `json:"id"`
	ExternalID             string  `json:"external_id"`
	UserID                 string  `json:"user_id"`
	IsHigh                 bool    `json:"is_high"`
	PaymentMethod          string  `json:"payment_method"`
	Status                 string  `json:"status"`
	MerchantName           string  `json:"merchant_name"`
	Amount                 float64 `json:"amount"`
	PaidAmount             float64 `json:"paid_amount"`
	BankCode               string  `json:"bank_code"`
	PaidAt                 string  `json:"paid_at"`
	PayerEmail             string  `json:"payer_email"`
	Description            string  `json:"description"`
	AdjustedReceivedAmount float64 `json:"adjusted_received_amount"`
	FeesPaidAmount         float64 `json:"fees_paid_amount"`
	Updated                string  `json:"updated"`
	Created                string  `json:"created"`
	Currency               string  `json:"currency"`
	PaymentChannel         string  `json:"payment_channel"`
	PaymentDestination     string  `json:"payment_destination"`
}
