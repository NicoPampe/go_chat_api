// I was struggling to find a quick error handling method in Go.
// There are TONS of options, but I didn't want to get hung up on one part.
// I based the error class off the example from swaggo
// https://github.com/swaggo/swag/blob/master/example/celler/httputil/error.go


package error

import "github.com/gin-gonic/gin"

// DefaultError
func DefaultError(ctx *gin.Context, status int, err error) {
	er := HTTPError{
		Code:    status,
		Message: err.Error(),
	}
	ctx.JSON(status, er)
}


// NewError
func NewError(ctx *gin.Context, status int, err error) {
	er := HTTPError{
		Code:    status,
		Message: err.Error(),
	}
	ctx.JSON(status, er)
}

// HTTPError
type HTTPError struct {
	Code    int    `json:"code" example:"400"`
	Message string `json:"message" example:"status bad request"`
}
