package shipments

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/masadamsahid/golang-gin-goldship-api/db"
	"github.com/masadamsahid/golang-gin-goldship-api/helpers"
	"github.com/masadamsahid/golang-gin-goldship-api/helpers/models"
	xenditService "github.com/masadamsahid/golang-gin-goldship-api/helpers/xendit-service"
	"github.com/masadamsahid/golang-gin-goldship-api/modules/users/roles"
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
		// TODO: Handle response
		log.Fatalf("Error beginning transaction: %v", txErr)
	}
	defer db.CloseTx(tx, txErr)

	var newShipment models.Shipment
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
		models.StatusPendingPayment,
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

	createInvoiceReq := *invoice.NewCreateInvoiceRequest(
		"INV-"+newShipment.TrackingNumber,
		float64(20000), // TODO: change to dynamic rate using Google Map API or OpenCage. Currently let's use a flat price
	)

	createInvoiceReq.SetInvoiceDuration(30 * 60) // 30 mins

	inv, resp, xenditErr := xenditService.Client.InvoiceApi.CreateInvoice(context.Background()).CreateInvoiceRequest(
		createInvoiceReq,
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

	var payment models.Payment
	sqlCreatePayment := `
	INSERT INTO payments (
		shipment_id,
		amount,
		invoice_id,
		external_id,
		invoice_url
	) VALUES ($1, $2, $3, $4, $5)
	RETURNING id, shipment_id, amount, paid_at, expired_at, invoice_id, external_id, invoice_url, "status", created_at, updated_at
	`
	err = tx.QueryRow(sqlCreatePayment, newShipment.ID, inv.Amount, inv.Id, inv.ExternalId, inv.InvoiceUrl).Scan(
		&payment.ID,
		&payment.ShipmentID,
		&payment.Amount,
		&payment.PaidAt,
		&payment.ExpiredAt,
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

	var initialHistory models.ShipmentHistory
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

func GetShipmentsList(ctx *gin.Context) {
	page, pageSize, err := helpers.ParsePaginationFromQueryParams(ctx)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"message": err.Error(),
		})
		return
	}

	var shipments []models.Shipment

	sqlGetShipments := `
		SELECT
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
		FROM shipments
		LIMIT $1 OFFSET $2
	`

	rows, err := db.DB.Query(sqlGetShipments, pageSize, (page-1)*pageSize)
	if err != nil {
		log.Println("Failed to get shipments", err)
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"message": "Internal server error",
		})
		return
	}
	defer rows.Close()

	for rows.Next() {
		var shipment models.Shipment
		err := rows.Scan(
			&shipment.ID,
			&shipment.TrackingNumber,
			&shipment.SenderID,
			&shipment.SenderName,
			&shipment.SenderPhone,
			&shipment.SenderAddress,
			&shipment.RecipientName,
			&shipment.RecipientAddress,
			&shipment.RecipientPhone,
			&shipment.ItemName,
			&shipment.ItemWeight,
			&shipment.Distance,
			&shipment.Status,
			&shipment.CreatedAt,
			&shipment.UpdatedAt,
		)
		if err != nil {
			log.Println("Failed to scan shipment", err)
			ctx.JSON(http.StatusInternalServerError, gin.H{
				"message": "Internal server error",
			})
			return
		}
		shipments = append(shipments, shipment)
	}

	var totalShipments int
	err = db.DB.QueryRow("SELECT COUNT(id) FROM shipments").Scan(&totalShipments)
	if err != nil {
		log.Println("Failed to get total shipments", err)
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"message": "Internal server error",
		})
		return
	}

	if len(shipments) < 1 {
		ctx.JSON(http.StatusNotFound, gin.H{
			"message": "No shipments found",
			"data":    []models.Shipment{},
			"meta": gin.H{
				"total":     totalShipments,
				"page":      page,
				"page_size": pageSize,
			},
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"message": "Shipments retrieved successfully",
		"data":    shipments,
		"meta": gin.H{
			"total":     totalShipments,
			"page":      page,
			"page_size": pageSize,
		},
	})
}

func GetShipmentByID(ctx *gin.Context) {
	strId := ctx.Param("id")
	id, err := strconv.Atoi(strId)
	if err != nil {
		log.Println(strId)
		log.Println(err)
		ctx.JSON(http.StatusBadRequest, gin.H{
			"message": "Invalid shipment ID",
		})
		return
	}

	var s models.Shipment
	var p models.Payment
	sqlGetShipment := `
		SELECT
			s.id,
			s.tracking_number,
			s.sender_id,
			s.sender_name,
			s.sender_phone,
			s.sender_address,
			s.recipient_name,
			s.recipient_address,
			s.recipient_phone,
			s.item_name,
			s.item_weight,
			s.distance,
			s.status,
			s.created_at,
			s.updated_at,
			p.id,
			p.shipment_id,
			p.amount,
			p.paid_at,
			p.expired_at,
			p.invoice_id,
			p.external_id,
			p.invoice_url,
			p.status,
			p.created_at,
			p.updated_at
		FROM shipments s
		JOIN payments p ON p.shipment_id = s.id
		WHERE s.id = $1
	`

	err = db.DB.QueryRow(sqlGetShipment, id).Scan(
		&s.ID,
		&s.TrackingNumber,
		&s.SenderID,
		&s.SenderName,
		&s.SenderPhone,
		&s.SenderAddress,
		&s.RecipientName,
		&s.RecipientAddress,
		&s.RecipientPhone,
		&s.ItemName,
		&s.ItemWeight,
		&s.Distance,
		&s.Status,
		&s.CreatedAt,
		&s.UpdatedAt,
		&p.ID,
		&p.ShipmentID,
		&p.Amount,
		&p.PaidAt,
		&p.ExpiredAt,
		&p.InvoiceID,
		&p.ExternalID,
		&p.InvoiceURL,
		&p.Status,
		&p.CreatedAt,
		&p.UpdatedAt,
	)
	if err != nil {
		log.Println("Failed to get shipment by ID", err)
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"message": "Internal server error",
		})
		return
	}

	var histories []models.ShipmentHistory
	sqlGetHistories := `
		SELECT
			id,
			shipment_id,
			status,
			"desc",
			courier_id,
			branch_id,
			timestamp
		FROM shipment_histories
		WHERE shipment_id = $1
		ORDER BY timestamp ASC
	`

	rows, err := db.DB.Query(sqlGetHistories, id)
	if err != nil {
		log.Println("Failed to get shipment histories", err)
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"message": "Internal server error",
		})
		return
	}

	defer rows.Close()
	for rows.Next() {
		var history models.ShipmentHistory
		err := rows.Scan(
			&history.ID,
			&history.ShipmentID,
			&history.Status,
			&history.Desc,
			&history.CourierID,
			&history.BranchID,
			&history.Timestamp,
		)
		if err != nil {
			log.Println("Failed to scan shipment history", err)
			ctx.JSON(http.StatusInternalServerError, gin.H{
				"message": "Internal server error",
			})
			return
		}
		histories = append(histories, history)
	}

	s.Payment = &p
	s.Histories = histories

	ctx.JSON(http.StatusOK, gin.H{
		"message": "Shipment retrieved successfully",
		"data":    s,
	})
}

func CancelShipmentByID(ctx *gin.Context) {
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

	strId := ctx.Param("id")
	id, err := strconv.Atoi(strId)
	if err != nil {
		log.Println(strId)
		log.Println(err)
		ctx.JSON(http.StatusBadRequest, gin.H{
			"message": "Invalid shipment ID",
		})
		return
	}

	tx, txErr := db.DB.BeginTx(ctx, nil)
	if txErr != nil {
		log.Fatalf("Error beginning transaction: %v", txErr)
	}
	defer db.CloseTx(tx, txErr)

	var currentShipment models.Shipment
	err = tx.QueryRow(`SELECT id, sender_id, status FROM shipments WHERE id = $1 FOR UPDATE`, id).Scan(&currentShipment.ID, &currentShipment.SenderID, &currentShipment.Status)
	if err != nil {
		log.Println("Failed to get shipment for cancellation", err)
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"message": "Internal server error",
		})
		return
	}

	if user.Role != roles.RoleSuperAdmin && user.Role != roles.RoleAdmin && user.ID != uint(currentShipment.SenderID) {
		ctx.JSON(http.StatusForbidden, gin.H{
			"message": "You are not authorized to cancel this shipment",
		})
		return
	}

	if currentShipment.Status == models.StatusCancelled {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"message": "Shipment is already cancelled",
		})
		return
	}

	if currentShipment.Status != models.StatusPendingPayment {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"message": "Shipment can only be cancelled if it's in pending payment status",
		})
		return
	}

	_, err = tx.Exec(`UPDATE shipments SET status = $1 WHERE id = $2`, models.StatusCancelled, id)
	if err != nil {
		log.Println("Failed to update shipment status to cancelled", err)
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"message": "Internal server error",
		})
		return
	}

	var cancelHistory models.ShipmentHistory
	sqlInitHistory := `
	INSERT INTO shipment_histories (shipment_id, status, "desc")
	VALUES ($1, $2, $3)
	RETURNING  id, shipment_id, status, "desc", courier_id, branch_id, timestamp
	`

	desc := fmt.Sprintf("%s has cancelled the shipment. Shipment currently is %s", user.Username, models.StatusCancelled)

	err = tx.QueryRow(sqlInitHistory, id, models.StatusCancelled, desc).Scan(
		&cancelHistory.ID,
		&cancelHistory.ShipmentID,
		&cancelHistory.Status,
		&cancelHistory.Desc,
		&cancelHistory.CourierID,
		&cancelHistory.BranchID,
		&cancelHistory.Timestamp,
	)
	if err != nil {
		log.Println("Failed to insert cancellation history", err)
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

	ctx.JSON(http.StatusOK, gin.H{
		"message": "Shipment cancelled successfully",
	})
}

func PickupPackageByShipmentID(ctx *gin.Context) {
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

	strId := ctx.Param("id")
	id, err := strconv.Atoi(strId)
	if err != nil {
		log.Println(strId)
		log.Println(err)
		ctx.JSON(http.StatusBadRequest, gin.H{
			"message": "Invalid shipment ID",
		})
		return
	}

	tx, txErr := db.DB.BeginTx(ctx, nil)
	if txErr != nil {
		log.Fatalf("Error beginning transaction: %v", txErr)
	}
	defer db.CloseTx(tx, txErr)

	var currentShipment models.Shipment
	err = tx.QueryRow(`SELECT id, status FROM shipments WHERE id = $1 FOR UPDATE`, id).Scan(&currentShipment.ID, &currentShipment.Status)
	if err != nil {
		log.Println("Failed to get shipment for pickup", err)
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"message": "Internal server error",
		})
		return
	}

	if currentShipment.Status == models.StatusPickedUp ||
		currentShipment.Status == models.StatusInTransit ||
		currentShipment.Status == models.StatusDelivered {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"message": "Shipment is already picked up",
		})
		return
	}

	if currentShipment.Status != models.StatusReadyToPickup {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"message": "Shipment can only be picked up if it's in paid status",
		})
		return
	}

	_, err = tx.Exec(`UPDATE shipments SET status = $1 WHERE id = $2`, models.StatusPickedUp, id)
	if err != nil {
		log.Println("Failed to update shipment status to picked up", err)
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"message": "Internal server error",
		})
		return
	}

	var pickupHistory models.ShipmentHistory
	sqlInitHistory := `
	INSERT INTO shipment_histories (shipment_id, status, "desc", courier_id)
	VALUES ($1, $2, $3, $4)
	RETURNING  id, shipment_id, status, "desc", courier_id, branch_id, timestamp
	`

	desc := fmt.Sprintf("%s has picked up the package. Shipment currently is %s", user.Username, models.StatusPickedUp)

	err = tx.QueryRow(sqlInitHistory, id, models.StatusPickedUp, desc, user.ID).Scan(
		&pickupHistory.ID,
		&pickupHistory.ShipmentID,
		&pickupHistory.Status,
		&pickupHistory.Desc,
		&pickupHistory.CourierID,
		&pickupHistory.BranchID,
		&pickupHistory.Timestamp,
	)
	if err != nil {
		log.Println("Failed to insert picked up history", err)
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

	ctx.JSON(http.StatusOK, gin.H{
		"message": "Shipment picked up successfully",
	})
}

func TransitPackageByShipmentID(ctx *gin.Context) {
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

	strId := ctx.Param("id")
	id, err := strconv.Atoi(strId)
	if err != nil {
		log.Println(strId)
		log.Println(err)
		ctx.JSON(http.StatusBadRequest, gin.H{
			"message": "Invalid shipment ID",
		})
		return
	}

	var body TransitShipmentDto
	err = ctx.ShouldBind(&body)
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

	tx, txErr := db.DB.BeginTx(ctx, nil)
	if txErr != nil {
		log.Fatalf("Error beginning transaction: %v", txErr)
	}
	defer db.CloseTx(tx, txErr)

	var currentShipment models.Shipment
	err = tx.QueryRow(`SELECT id, status FROM shipments WHERE id = $1 FOR UPDATE`, id).Scan(&currentShipment.ID, &currentShipment.Status)
	if err != nil {
		log.Println("Failed to get shipment for transit", err)
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"message": "Internal server error",
		})
		return
	}

	var transitBranch models.Branch
	err = tx.QueryRow(`SELECT id, name, address FROM branches WHERE id = $1 LIMIT 1`, body.BranchID).Scan(&transitBranch.ID, &transitBranch.Name, &transitBranch.Address)
	if err != nil {
		log.Println("Failed to get branch for transit", err)
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"message": "Internal server error",
		})
		return
	}

	if currentShipment.Status == models.StatusDelivered {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"message": "Shipment is already in delivered",
		})
		return
	}

	if currentShipment.Status != models.StatusReadyToPickup && currentShipment.Status != models.StatusPickedUp && currentShipment.Status != models.StatusInTransit {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"message": "Shipment can only be transited if it's in picked up or in transit status",
		})
		return
	}

	_, err = tx.Exec(`UPDATE shipments SET status = $1 WHERE id = $2`, models.StatusInTransit, id)
	if err != nil {
		log.Println("Failed to update shipment status to in transit", err)
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"message": "Internal server error",
		})
		return
	}

	var transitHistory models.ShipmentHistory
	sqlInitHistory := `
	INSERT INTO shipment_histories (shipment_id, status, "desc", courier_id, branch_id)
	VALUES ($1, $2, $3, $4, $5)
	RETURNING  id, shipment_id, status, "desc", courier_id, branch_id, timestamp
	`

	desc := fmt.Sprintf(
		"%s has transited the package to branch %s [%d | %s]. Shipment currently is %s",
		user.Username, transitBranch.Name, transitBranch.ID, transitBranch.Address, models.StatusInTransit,
	)

	err = tx.QueryRow(sqlInitHistory, id, models.StatusInTransit, desc, user.ID, body.BranchID).Scan(
		&transitHistory.ID,
		&transitHistory.ShipmentID,
		&transitHistory.Status,
		&transitHistory.Desc,
		&transitHistory.CourierID,
		&transitHistory.BranchID,
		&transitHistory.Timestamp,
	)
	if err != nil {
		log.Println("Failed to insert in transit history", err)
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

	ctx.JSON(http.StatusOK, gin.H{
		"message": "Shipment transited successfully",
	})

}

func DeliverPackageByShipmentID(ctx *gin.Context) {
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

	strId := ctx.Param("id")
	id, err := strconv.Atoi(strId)
	if err != nil {
		log.Println(strId)
		log.Println(err)
		ctx.JSON(http.StatusBadRequest, gin.H{
			"message": "Invalid shipment ID",
		})
		return
	}

	tx, txErr := db.DB.BeginTx(ctx, nil)
	if txErr != nil {
		log.Fatalf("Error beginning transaction: %v", txErr)
	}
	defer db.CloseTx(tx, txErr)

	var currentShipment models.Shipment
	err = tx.QueryRow(`SELECT id, status FROM shipments WHERE id = $1 FOR UPDATE`, id).Scan(&currentShipment.ID, &currentShipment.Status)
	if err != nil {
		log.Println("Failed to get shipment for delivery", err)
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"message": "Internal server error",
		})
		return
	}

	if currentShipment.Status == models.StatusDelivered {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"message": "Shipment is already delivered",
		})
		return
	}

	if currentShipment.Status != models.StatusInTransit {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"message": "Shipment can only be delivered if it's in transit status",
		})
		return
	}

	_, err = tx.Exec(`UPDATE shipments SET status = $1 WHERE id = $2`, models.StatusDelivered, id)
	if err != nil {
		log.Println("Failed to update shipment status to delivered", err)
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"message": "Internal server error",
		})
		return
	}

	var deliveredHistory models.ShipmentHistory
	sqlInitHistory := `
	INSERT INTO shipment_histories (shipment_id, status, "desc", courier_id)
	VALUES ($1, $2, $3, $4)
	RETURNING  id, shipment_id, status, "desc", courier_id, branch_id, timestamp
	`

	desc := fmt.Sprintf("%s has delivered the package. Shipment currently is %s", user.Username, models.StatusDelivered)

	err = tx.QueryRow(sqlInitHistory, id, models.StatusDelivered, desc, user.ID).Scan(
		&deliveredHistory.ID,
		&deliveredHistory.ShipmentID,
		&deliveredHistory.Status,
		&deliveredHistory.Desc,
		&deliveredHistory.CourierID,
		&deliveredHistory.BranchID,
		&deliveredHistory.Timestamp,
	)
	if err != nil {
		log.Println("Failed to insert delivered history", err)
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

	ctx.JSON(http.StatusOK, gin.H{
		"message": "Shipment delivered successfully",
	})

}
