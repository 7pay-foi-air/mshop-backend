package models

type CreateReportRequest struct {
	StartDate string `json:"start_date" binding:"required"` // YYYY-MM-DD
	EndDate   string `json:"end_date" binding:"required"`   // YYYY-MM-DD
	Email     string `json:"email" binding:"required,email"`
}
