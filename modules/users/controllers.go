package users

import (
	"log"
	"net/http"
	"slices"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/masadamsahid/golang-gin-goldship-api/db"
	"github.com/masadamsahid/golang-gin-goldship-api/helpers"
	"github.com/masadamsahid/golang-gin-goldship-api/helpers/models"
	"github.com/masadamsahid/golang-gin-goldship-api/modules/users/roles"
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

func ChangeUserRole(ctx *gin.Context) {
	targetUsername := ctx.Param("username")

	var user helpers.AuthPayload
	err := helpers.ParseJWTUserFromCtx(ctx, &user)
	if err != nil {
		log.Println("Failed convert user from context to AuthPayload")
		ctx.JSON(http.StatusUnauthorized, gin.H{
			"message": "Unauthorized",
		})
		return
	}

	var body ChangeUserRoleDto
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

	if !slices.Contains(roles.ROLE_LIST, body.Role) {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"message": "Invalid role",
		})
		return
	}

	// Only SUPERADMIN who can grant a user role as ADMIN or SUPERADMIN
	if (body.Role == roles.RoleSuperAdmin || body.Role == roles.RoleAdmin) && user.Role != roles.RoleSuperAdmin {
		ctx.JSON(http.StatusForbidden, gin.H{
			"message": "You are not allowed to grant this role",
		})
		return
	}

	var targetUser models.User
	err = db.DB.QueryRow("SELECT id, username, email, role, created_at, updated_at FROM users WHERE username = $1 LIMIT 1", targetUsername).
		Scan(&targetUser.ID, &targetUser.Username, &targetUser.Email, &targetUser.Role, &targetUser.CreatedAt, &targetUser.UpdatedAt)
	if err != nil {
		log.Println("Failed to get target user", err)
		ctx.JSON(http.StatusNotFound, gin.H{
			"message": "User not found",
		})
		return
	}

	if targetUser.Role == body.Role {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"message": "User already has this role",
			"data":    targetUser,
		})
		return
	}

	if targetUser.Role == roles.RoleSuperAdmin {
		ctx.JSON(http.StatusForbidden, gin.H{
			"message": "You are not allowed to change this user",
		})
		return
	}

	// Only SUPERADMIN who can change ADMINs' role
	if targetUser.Role == roles.RoleAdmin && user.Role != roles.RoleSuperAdmin {
		ctx.JSON(http.StatusForbidden, gin.H{
			"message": "You are not allowed to change this user",
		})
		return
	}

	sqlUpdateRole := `UPDATE users SET role = $2, updated_at = NOW() WHERE id = $1 AND role != 'SUPERADMIN' RETURNING id, username, email, role, created_at, updated_at`
	err = db.DB.QueryRow(sqlUpdateRole, targetUser.ID, body.Role).Scan(&targetUser.ID, &targetUser.Username, &targetUser.Email, &targetUser.Role, &targetUser.CreatedAt, &targetUser.UpdatedAt)
	if err != nil {
		log.Println("Failed to update user role", err)
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"message": "Internal server error",
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"message": "User role updated successfully",
		"data":    targetUser,
	})

}
