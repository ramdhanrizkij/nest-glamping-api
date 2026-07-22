package query

import (
	"strconv"

	"github.com/gofiber/fiber/v3"
	"gorm.io/gorm"
)

type Params struct {
	Page    int
	PerPage int
	Sort    string
	Order   string
}

func GetParams(c fiber.Ctx) Params {
	page, _ := strconv.Atoi(c.Query("page", "1"))
	perPage, _ := strconv.Atoi(c.Query("per_page", "10"))
	sort := c.Query("sort", "created_at")
	order := c.Query("order", "desc")

	if page < 1 {
		page = 1
	}
	if perPage < 1 {
		perPage = 10
	}
	if perPage > 100 {
		perPage = 100
	}
	if sort == "" {
		sort = "created_at"
	}
	if order != "asc" && order != "desc" {
		order = "desc"
	}

	return Params{Page: page, PerPage: perPage, Sort: sort, Order: order}
}

func (p Params) Offset() int {
	return (p.Page - 1) *p.PerPage
}

func ApplyPagination(db *gorm.DB, p Params) *gorm.DB {
	return db.Order(p.Sort + " " + p.Order).Offset(p.Offset()).Limit(p.PerPage)
}

func CountTotal(db *gorm.DB) int64 {
	var total int64
	db.Count(&total)
	return total
}

func TotalPages(total int64, perPage int) int {
	if perPage <= 0 {
		return 0
	}
	return int((total + int64(perPage) - 1) / int64(perPage))
}
