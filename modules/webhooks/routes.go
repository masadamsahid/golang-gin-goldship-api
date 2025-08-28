package webhooks

import "github.com/gin-gonic/gin"

func Routes(rg *gin.RouterGroup) {
	rg.POST("/xendit-payment-received", XenditInvoiceNotification)
}

// {
//   "data": {
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
//   },
//   "message": "Notification received"
// }
