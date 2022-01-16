package validator

import (
	"fmt"
	"math"
	errorcode "wager/error_code"

	go_validate "github.com/go-playground/validator/v10"
	"github.com/sirupsen/logrus"
)

var validate *go_validate.Validate

func init() {
	validate = go_validate.New()
	err := validate.RegisterValidation("monetary-format", validateMonetaryFormat)
	if err != nil {
		logrus.WithError(err).Fatal("failed to register monetary validator")
	}
}

func Validate(v interface{}) error {
	return validate.Struct(v)
}

func validateMonetaryFormat(field go_validate.FieldLevel) bool {
	eps := 1e-8
	value := field.Field().Float()
	return value*1e2-math.Floor(value*1e2) < eps
}

func ErrorMsg(err error) errorcode.ErrorResponse {
	result := []string{}
	for _, e := range err.(go_validate.ValidationErrors) {
		result = append(result, fieldErrorMsg(e))
	}
	return errorcode.ErrorResponse{Error: result}
}

func fieldErrorMsg(fieldError go_validate.FieldError) string {
	switch fieldError.Tag() {
	case "gt":
		return fmt.Sprintf("%v must be larger than %s", fieldError.Field(), fieldError.Param())
	case "gte":
		return fmt.Sprintf("%v must be larger than or equal %s", fieldError.Field(), fieldError.Param())
	case "lte":
		return fmt.Sprintf("%v must be less than or equal %s", fieldError.Field(), fieldError.Param())
	case "monetary-format":
		return fmt.Sprintf("%v must be in monetary format with maximum 2 decimal places", fieldError.Field())
	default:
		return fieldError.Error()
	}
}
