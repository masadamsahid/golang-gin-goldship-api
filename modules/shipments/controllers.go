package shipments

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/masadamsahid/golang-gin-goldship-api/db"
	"github.com/masadamsahid/golang-gin-goldship-api/helpers"
	xenditService "github.com/masadamsahid/golang-gin-goldship-api/helpers/xendit-service"
	"github.com/masadamsahid/golang-gin-goldship-api/modules/payments"
	"github.com/xendit/xendit-go/v7/invoice"
)

func CreateNewShipment(ctx *gin.Context) {
	u, ok := ctx.Get("user")
	if !ok {
		log.Println("Failed get user from context")
		ctx.JSON(http.StatusUnauthorized, gin.H{
			"message": "Unauthorized",
		})
		return
	}

	user, ok := u.(helpers.AuthPayload)
	if !ok {
		log.Println("Failed convert user from context to AuthPayload")
		ctx.JSON(http.StatusUnauthorized, gin.H{
			"message": "Unauthorized",
		})
		return
	}

	var body CreateShipmentDto
	err := ctx.ShouldBind(&body)
	if err != nil {
		validationErrors, ok := err.(validator.ValidationErrors)
		if !ok {
			log.Println(err)
			ctx.JSON(http.StatusInternalServerError, gin.H{
				"message": "Internal server error",
			})
			return
		}

		log.Printf("%+v\n", validationErrors)
		errs := helpers.HandleValidationErrors(validationErrors)

		ctx.JSON(http.StatusBadRequest, gin.H{"message": "Validation failed", "errors": errs})
		return
	}

	trackingNumber := helpers.GenerateTrackingNumber()

	sqlCreateShipment := `
		INSERT INTO shipments (
			tracking_number,
			sender_id,
			sender_name,
			sender_phone,
			sender_address,
			recipient_name,
			recipient_address,
			recipient_phone,
			item_name,
			item_weight,
			distance,
			"status"
		) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12)
		RETURNING
			id,
			tracking_number,
			sender_id,
			sender_name,
			sender_phone,
			sender_address,
			recipient_name,
			recipient_address,
			recipient_phone,
			item_name,
			item_weight,
			distance,
			"status",
			created_at,
			updated_at
	`

	tx, txErr := db.DB.BeginTx(ctx, nil)
	if txErr != nil {
		log.Fatalf("Error beginning transaction: %v", txErr)
	}
	defer db.CloseTx(tx, txErr)

	var newShipment Shipment
	err = tx.QueryRow(
		sqlCreateShipment,
		trackingNumber,
		user.ID,
		body.SenderName,
		body.SenderPhone,
		body.SenderAddress,
		body.RecipientName,
		body.RecipientAddress,
		body.RecipientPhone,
		body.ItemName,
		body.ItemWeight,
		body.Distance,
		StatusPendingPayment,
	).Scan(
		&newShipment.ID,
		&newShipment.TrackingNumber,
		&newShipment.SenderID,
		&newShipment.SenderName,
		&newShipment.SenderPhone,
		&newShipment.SenderAddress,
		&newShipment.RecipientName,
		&newShipment.RecipientAddress,
		&newShipment.RecipientPhone,
		&newShipment.ItemName,
		&newShipment.ItemWeight,
		&newShipment.Distance,
		&newShipment.Status,
		&newShipment.CreatedAt,
		&newShipment.UpdatedAt,
	)
	if err != nil {
		log.Println("Failed creating shipment", err)
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"message": "Internal server error",
		})
		return
	}

	inv, resp, xenditErr := xenditService.Client.InvoiceApi.CreateInvoice(context.Background()).CreateInvoiceRequest(
		*invoice.NewCreateInvoiceRequest(
			"INV-"+newShipment.TrackingNumber,
			float64(20000),
		),
	).Execute()

	if xenditErr != nil {
		log.Println("EHEEEY", xenditErr.FullError())
		fmt.Fprintf(os.Stderr, "Error when calling `InvoiceApi.CreateInvoice``: %v\n", xenditErr.Error())

		b, _ := json.Marshal(xenditErr.FullError())
		fmt.Fprintf(os.Stderr, "Full Error Struct: %v\n", string(b))

		fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", resp)
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"message": "Failed to create invoice",
		})
		return
	}

	var payment payments.Payment
	sqlCreatePayment := `
	INSERT INTO payments (
		shipment_id,
		amount,
		invoice_id,
		external_id,
		invoice_url
	) VALUES ($1, $2, $3, $4, $5)
	RETURNING id, shipment_id, amount, payment_date, invoice_id, external_id, invoice_url, "status", created_at, updated_at
	`
	err = tx.QueryRow(sqlCreatePayment, newShipment.ID, inv.Amount, inv.Id, inv.ExternalId, inv.InvoiceUrl).Scan(
		&payment.ID,
		&payment.ShipmentID,
		&payment.Amount,
		&payment.PaymentDate,
		&payment.InvoiceID,
		&payment.ExternalID,
		&payment.InvoiceURL,
		&payment.Status,
		&payment.CreatedAt,
		&payment.UpdatedAt,
	)
	if err != nil {
		log.Println("Failed creating payment", err)
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"message": "Internal server error",
		})
		return
	}

	var initialHistory ShipmentHistory
	sqlInitHistory := `
	INSERT INTO shipment_histories (shipment_id, status, "desc")
	VALUES ($1, $2, $3)
	RETURNING  id, shipment_id, status, "desc", courier_id, branch_id, timestamp
	`

	desc := fmt.Sprintf("%s has requested a shipment. Shipment currently is %s", user.Username, newShipment.Status)

	err = tx.QueryRow(sqlInitHistory, newShipment.ID, newShipment.Status, desc).Scan(
		&initialHistory.ID,
		&initialHistory.ShipmentID,
		&initialHistory.Status,
		&initialHistory.Desc,
		&initialHistory.CourierID,
		&initialHistory.BranchID,
		&initialHistory.Timestamp,
	)
	if err != nil {
		log.Println("Failed initializing first history", err)
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"message": "Internal server error",
		})
		return
	}

	err = tx.Commit()
	if err != nil {
		log.Println(err)
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"message": "Internal server error",
		})
		return
	}

	newShipment.Payment = &payment
	newShipment.Histories = append(newShipment.Histories, initialHistory)

	ctx.JSON(http.StatusCreated, gin.H{
		"message": "Shipment created successfully",
		"data":    newShipment,
	})
}
