package utils

import "github.com/gin-gonic/gin"

func SuccessResponse(message string, data interface{}) gin.H {
	resp := gin.H{
		"status":  true,
		"message": message,
	}
	if data != nil {
		resp["data"] = data
	}
	return resp
}

func ErrorResponse(message string, detail string) gin.H {
	resp := gin.H{
		"status":  false,
		"message": message,
	}
	if detail != "" {
		resp["error"] = detail
	}
	return resp
}