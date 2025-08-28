package users

import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/masadamsahid/golang-gin-goldship-api/db"
	"github.com/masadamsahid/golang-gin-goldship-api/helpers"
	"github.com/masadamsahid/golang-gin-goldship-api/helpers/models"
)

func GetMyShipments(ctx *gin.Context) {
	var user helpers.AuthPayload
	err := helpers.ParseJWTUserFromCtx(ctx, &user)
	if err != nil {
		log.Println("Failed convert user from context to AuthPayload")
		ctx.JSON(http.StatusUnauthorized, gin.H{
			"message": "Unauthorized",
		})
		return
	}

	page, pageSize, err := helpers.ParsePaginationFromQueryParams(ctx)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"message": err.Error(),
		})
		return
	}

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
		status,
		created_at,
		updated_at
	FROM shipments
	WHERE sender_id = $1
	ORDER BY created_at DESC
	LIMIt $2
	OFFSET $3
	`

	var myShipments []models.Shipment
	rows, err := db.DB.Query(sqlGetShipments, user.ID, pageSize, pageSize*(page-1))
	if err != nil {
		log.Println("Failed to get shipments", err)
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"message": "Internal server error",
		})
		return
	}
	defer rows.Close()

	for rows.Next() {
		var s models.Shipment
		err := rows.Scan(
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
		)
		if err != nil {
			log.Println("Failed to scan shipment", err)
			ctx.JSON(http.StatusInternalServerError, gin.H{
				"message": "Internal server error",
			})
			return
		}

		myShipments = append(myShipments, s)
	}

	if len(myShipments) < 1 {
		ctx.JSON(http.StatusNotFound, gin.H{
			"message": "No shipments found",
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"message": "Shipments retrieved successfully",
		"data":    myShipments,
	})

}
