package utils

import "github.com/gin-gonic/gin"

func NewBaseSuccessApiResponse(successMessage string) gin.H {
	return gin.H{
		"status":  "success",
		"message": successMessage,
	}
}

func NewSuccessApiResponseWithData(successMessage string, data interface{}) gin.H {
	return gin.H{
		"status":  "success",
		"message": successMessage,
		"data":    data,
	}
}

func NewErrorApiResponse(errMessage string) gin.H {
	return gin.H{
		"status":  "failed",
		"message": errMessage,
	}
}
