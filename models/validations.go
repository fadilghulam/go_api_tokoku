package model

import (
	"fmt"

	"github.com/go-playground/validator/v10"
)

func FormatValidationError(err error) string {
	if errs, ok := err.(validator.ValidationErrors); ok {
		var errMsg string
		for _, e := range errs {
			switch e.Tag() {
			case "required":
				errMsg = fmt.Sprintf("%s field is required. ", e.Field())
			case "email":
				errMsg = fmt.Sprintf("%s field must be a valid email address. ", e.Field())
			case "min":
				errMsg = fmt.Sprintf("%s field must be at least %s characters long. ", e.Field(), e.Param())
			case "max":
				errMsg = fmt.Sprintf("%s field must be at most %s characters long. ", e.Field(), e.Param())
			default:
				errMsg = fmt.Sprintf("%s field is invalid. ", e.Field())
			}
		}
		return errMsg
	}
	return err.Error()
}

func ValidateHrPerson(hrPerson *HrPerson) error {
	validate := validator.New()

	// Use the default validation
	err := validate.Struct(hrPerson)

	// fmt.Println(err)
	if err != nil {
		return err
	}
	return nil
}
