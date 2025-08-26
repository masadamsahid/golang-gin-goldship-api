package branches

import common "github.com/masadamsahid/golang-gin-goldship-api/helpers/commons"

type Branch struct {
	Name    string `json:"name"`
	Phone   string `json:"phone"`
	Address string `json:"address"`
	common.BaseEntity
}
