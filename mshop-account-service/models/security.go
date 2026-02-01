package models

type SetSecurityQuestionsRequest struct {
	Answer1              string `json:"answer1" binding:"required"`
	Answer2              string `json:"answer2" binding:"required"`
	Answer3              string `json:"answer3" binding:"required"`
	RecoveryCodeLocation string `json:"recovery_code_location" binding:"required"`
}