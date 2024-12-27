package response

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// CommonResponse is a standardized response structure.
type CommonResponse struct {
	Success bool        `json:"success"`           // Indicates if the request was successful
	Message string      `json:"message,omitempty"` // Human-readable message
	Data    interface{} `json:"data,omitempty"`    // Successful response data
	Errors  interface{} `json:"errors,omitempty"`  // Validation errors or general errors
}

// fallbackMessages maps status codes to default messages
var fallbackMessages = map[int]string{
	http.StatusOK:                  "Request was successful.",
	http.StatusCreated:             "The record was created successfully.",
	http.StatusNoContent:           "No content available.",
	http.StatusBadRequest:          "The request was invalid.",
	http.StatusInternalServerError: "An internal server error occurred.",
	// Add any other status codes your API frequently uses
}

// Success sends a success response with an optional message and data.
// If the message is empty, it uses a default from fallbackMessages if available.
func Success(c *gin.Context, statusCode int, message string, data interface{}) {
	if message == "" {
		if msg, ok := fallbackMessages[statusCode]; ok {
			message = msg
		} else {
			message = "Operation was successful." // Very generic fallback
		}
	}

	c.JSON(statusCode, CommonResponse{
		Success: true,
		Message: message,
		Data:    data,
	})
}

// Error sends an error response with an optional message and an errors object.
// If the message is empty, it uses a default from fallbackMessages if available.
func Error(c *gin.Context, statusCode int, message string, err interface{}) {
	if message == "" {
		if msg, ok := fallbackMessages[statusCode]; ok {
			message = msg
		} else {
			message = "An error occurred." // Very generic fallback
		}
	}

	c.JSON(statusCode, CommonResponse{
		Success: false,
		Message: message,
		Errors:  err,
	})
}

// BadRequest is a shorthand for sending a 400 response.
func BadRequest(c *gin.Context, message string, err interface{}) {
	Error(c, http.StatusBadRequest, message, err)
}

// InternalServerError is a shorthand for sending a 500 response.
func InternalServerError(c *gin.Context, message string, err interface{}) {
	Error(c, http.StatusInternalServerError, message, err)
}

// Created is a convenience function for returning a 201 Created response.
// If the message is empty, defaults to "The record was created successfully."
func Created(c *gin.Context, message string, data interface{}) {
	Success(c, http.StatusCreated, message, data)
}

// Updated is a convenience function for returning a 200 response for an updated record.
func Updated(c *gin.Context, message string, data interface{}) {
	Success(c, http.StatusOK, message, data)
	// Or you might consider using 204 No Content if you don't return data:
	// Success(c, http.StatusNoContent, message, nil)
}

// Deleted is a convenience function for returning a 200 response for a deleted record.
func Deleted(c *gin.Context, message string) {
	// Typically you might not send data back on a delete
	Success(c, http.StatusOK, message, nil)
}

// NotFound is a convenience function for returning a 404 response
func NotFound(c *gin.Context, message string, err interface{}) {
	Error(c, http.StatusNotFound, message, err)
}
