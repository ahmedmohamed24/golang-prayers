package requests

import (
	errsPackage "errors"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
)

type Validator interface {
	Validate(ctx *gin.Context) bool
}

func Validate(ctx *gin.Context, request Validator) bool {
	if err := ctx.ShouldBindJSON(&request); err != nil {
		if err.Error() == "EOF" {
			ctx.JSON(http.StatusUnprocessableEntity, gin.H{
				"errors": "Please provide a valid json body",
			})
			return false
		}
		var errors []string
		var validattionErrors validator.ValidationErrors
		if errsPackage.As(err, &validattionErrors) {
			for _, v := range err.(validator.ValidationErrors) {
				errors = append(errors, v.StructField()+" field "+v.Tag())

			}
		} else {
			errors = append(errors, err.Error())
		}
		ctx.JSON(http.StatusUnprocessableEntity, gin.H{
			"errors": errors,
		})
		return false
	}
	return true
}
