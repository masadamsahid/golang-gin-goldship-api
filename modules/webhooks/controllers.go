package webhooks

import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/masadamsahid/golang-gin-goldship-api/db"
	"github.com/masadamsahid/golang-gin-goldship-api/helpers/models"
	xenditService "github.com/masadamsahid/golang-gin-goldship-api/helpers/xendit-service"
	"github.com/xendit/xendit-go/v7/invoice"
)

func XenditInvoiceNotification(ctx *gin.Context) {
	callBackToken := ctx.GetHeader("X-CALLBACK-TOKEN")
	if callBackToken != xenditService.XENDIT_WEBHOOK_VERIFICATION_TOKEN {
		ctx.JSON(http.StatusUnauthorized, gin.H{
			"message": "Invalid verify token",
		})
		return
	}

	var body XenditInvoiceNotificationDto
	if err := ctx.ShouldBindJSON(&body); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}

	var payment models.Payment
	sqlStatement := `SELECT id, shipment_id, status FROM payments WHERE invoice_id = $1 LIMIT 1`
	err := db.DB.QueryRow(sqlStatement, body.ID).Scan(&payment.ID, &payment.ShipmentID, &payment.Status)
	if err != nil {
		log.Printf("Error getting payment: %v\n", err)
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"message": "Failed to get payment",
		})
		return
	}

	tx, txErr := db.DB.BeginTx(ctx, nil)
	if txErr != nil {
		log.Printf("Error beginning transaction: %v\n", txErr)
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"message": "Failed to begin transaction",
		})
		return
	}
	defer db.CloseTx(tx, txErr)

	var updatedPayment models.Payment
	switch body.Status {
	case string(invoice.INVOICESTATUS_PAID):
		fallthrough
	case string(invoice.INVOICESTATUS_SETTLED):
		sqlUpdatePaymentStatus := `UPDATE payments SET status = $2, paid_at = $3 WHERE invoice_id = $1 RETURNING id, shipment_id, status`
		txErr = tx.QueryRow(sqlUpdatePaymentStatus, body.ID, body.Status, body.PaidAt).Scan(
			&updatedPayment.ID,
			&updatedPayment.ShipmentID,
			&updatedPayment.Status,
		)
		if txErr != nil {
			log.Printf("Error updating payment status: %v\n", txErr)
			ctx.JSON(http.StatusInternalServerError, gin.H{
				"message": "Failed to update payment status",
			})
			return
		}

		desc := "Payment is paid. A courier will be dispatched to pick up"

		sqlInsertHistory := `
		INSERT INTO shipment_histories (shipment_id, status, "desc")
		VALUES ($1, $2, $3)
		RETURNING id
		`

		_, txErr = tx.Exec(sqlInsertHistory, updatedPayment.ShipmentID, models.StatusReadyToPickup, desc)
		if txErr != nil {
			log.Printf("Error inserting shipment history: %v\n", txErr)
			ctx.JSON(http.StatusInternalServerError, gin.H{
				"message": "Failed to insert shipment history",
			})
			return
		}

		sqlUpdateShipmentStatus := `UPDATE shipments SET status = $2 WHERE id = $1`
		_, txErr = tx.Exec(sqlUpdateShipmentStatus, updatedPayment.ShipmentID, models.StatusReadyToPickup)
		if txErr != nil {
			log.Printf("Error updating shipment status: %v\n", txErr)
			ctx.JSON(http.StatusInternalServerError, gin.H{
				"message": "Failed to update shipment status",
			})
			return
		}
		tx.Commit()
	case string(invoice.INVOICESTATUS_EXPIRED):
		sqlUpdatePaymentStatus := `UPDATE payments SET status = $2 WHERE invoice_id = $1 RETURNING id, shipment_id, status`
		txErr = tx.QueryRow(sqlUpdatePaymentStatus, body.ID, body.Status).Scan(
			&updatedPayment.ID,
			&updatedPayment.ShipmentID,
			&updatedPayment.Status,
		)
		if txErr != nil {
			log.Printf("Error updating payment status: %v\n", txErr)
			ctx.JSON(http.StatusInternalServerError, gin.H{
				"message": "Failed to update payment status",
			})
			return
		}

		desc := "Payment is expired. Please create a new shipment"

		sqlInsertHistory := `
		INSERT INTO shipment_histories (shipment_id, status, "desc")
		VALUES ($1, $2, $3)
		RETURNING id
		`

		_, txErr = tx.Exec(sqlInsertHistory, updatedPayment.ShipmentID, models.StatusCancelled, desc)
		if txErr != nil {
			log.Printf("Error inserting shipment history: %v\n", txErr)
			ctx.JSON(http.StatusInternalServerError, gin.H{
				"message": "Failed to insert shipment history",
			})
			return
		}

		sqlUpdateShipmentStatus := `UPDATE shipments SET status = $2 WHERE id = $1`
		_, txErr = tx.Exec(sqlUpdateShipmentStatus, updatedPayment.ShipmentID, models.StatusCancelled)
		if txErr != nil {
			log.Printf("Error updating shipment status: %v\n", txErr)
			ctx.JSON(http.StatusInternalServerError, gin.H{
				"message": "Failed to update shipment status",
			})
			return
		}
		tx.Commit()
	default:
		log.Printf("Unhandled invoice status: %s\n", body.Status)
		ctx.JSON(http.StatusBadRequest, gin.H{
			"message": "Unhandled invoice status",
		})
		tx.Rollback()
	}

	ctx.JSON(http.StatusOK, gin.H{
		"message": "Notification received",
	})
}
