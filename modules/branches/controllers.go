package branches

import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/masadamsahid/golang-gin-goldship-api/db"
	"github.com/masadamsahid/golang-gin-goldship-api/helpers"
	"github.com/masadamsahid/golang-gin-goldship-api/helpers/models"
)

func HandleCreateBranch(ctx *gin.Context) {
	var body CreateBranchDto
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

	sqlCreateBranch := `INSERT INTO branches (name, phone, address) VALUES ($1, $2, $3) RETURNING id, name, phone, address, created_at, updated_at`
	var newBranch models.Branch
	err = db.DB.QueryRow(sqlCreateBranch, body.Name, body.Phone, body.Address).Scan(
		&newBranch.ID,
		&newBranch.Name,
		&newBranch.Phone,
		&newBranch.Address,
		&newBranch.CreatedAt,
		&newBranch.UpdatedAt,
	)
	if err != nil {
		status := http.StatusInternalServerError
		msg := "Failed creating branch"

		log.Println(err)
		ctx.JSON(status, gin.H{
			"message": msg,
		})
		return
	}

	ctx.JSON(http.StatusCreated, gin.H{
		"message": "Branch created successfully",
		"data":    newBranch,
	})
}

func HandleGetBranchesList(ctx *gin.Context) {
	page, pageSize, err := helpers.ParsePaginationFromQueryParams(ctx)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"message": err.Error(),
		})
		return
	}

	sqlGetBranches := `SELECT id, name, phone, address, created_at, updated_at FROM branches ORDER BY created_at ASC LIMIT $1 OFFSET $2`
	offset := (page - 1) * pageSize
	rows, err := db.DB.Query(sqlGetBranches, pageSize, offset)
	if err != nil {
		log.Println(err)
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"message": "Failed retrieving branches",
		})
		return
	}

	var branches []models.Branch
	defer rows.Close()
	for rows.Next() {
		var b models.Branch
		err := rows.Scan(
			&b.ID,
			&b.Name,
			&b.Phone,
			&b.Address,
			&b.CreatedAt,
			&b.UpdatedAt,
		)
		if err != nil {
			log.Println(err)
			ctx.JSON(http.StatusInternalServerError, gin.H{
				"message": "Failed retrieving branches",
			})
			return
		}
		branches = append(branches, b)
	}

	ctx.JSON(http.StatusOK, gin.H{
		"message": "Success retrieving branches",
		"data":    branches,
	})
}

func HandleGetBranchByID(ctx *gin.Context) {
	id := ctx.Param("id")

	sqlGetBranch := `SELECT id, name, phone, address, created_at, updated_at FROM branches WHERE id = $1`
	var branch models.Branch
	err := db.DB.QueryRow(sqlGetBranch, id).Scan(
		&branch.ID,
		&branch.Name,
		&branch.Phone,
		&branch.Address,
		&branch.CreatedAt,
		&branch.UpdatedAt,
	)
	if err != nil {
		log.Println(err)
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"message": "Failed retrieving branch",
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"message": "Success retrieving branch",
		"data":    branch,
	})
}

// TODO: Implement optional field updating
func HandleUpdateBranch(ctx *gin.Context) {
	id := ctx.Param("id")
	var body UpdateBranchDto
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

	sqlUpdateBranch := `UPDATE branches SET name = $2, phone = $3, address = $4, updated_at = NOW() WHERE id = $1 RETURNING id, name, phone, address, created_at, updated_at`
	var updatedBranch models.Branch
	err = db.DB.QueryRow(sqlUpdateBranch, id, body.Name, body.Phone, body.Address).Scan(
		&updatedBranch.ID,
		&updatedBranch.Name,
		&updatedBranch.Phone,
		&updatedBranch.Address,
		&updatedBranch.CreatedAt,
		&updatedBranch.UpdatedAt,
	)
	if err != nil {
		status := http.StatusInternalServerError
		msg := "Failed updating branch"

		log.Println(err)
		ctx.JSON(status, gin.H{
			"message": msg,
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"message": "Branch updated successfully",
		"data":    updatedBranch,
	})
}

func HandleDeleteBranch(ctx *gin.Context) {
	id := ctx.Param("id")

	sqlDeleteBranch := `DELETE FROM branches WHERE id = $1`
	_, err := db.DB.Exec(sqlDeleteBranch, id)
	if err != nil {
		log.Println(err)
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"message": "Failed deleting branch",
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"message": "Branch deleted successfully",
	})
}
