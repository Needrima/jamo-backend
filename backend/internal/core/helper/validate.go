package helper

import (
	"errors"
	"reflect"

	"github.com/go-playground/validator/v10"
)

func Validate(data interface{}) error {
	LogEvent("INFO", "Validating "+reflect.TypeOf(data).String()+" Data...")
	err := validator.New().Struct(data)
	if err != nil {
		var errMsgs string
		LogEvent("ERROR", "Error validating struct: "+err.Error())
		for _, errs := range err.(validator.ValidationErrors) {
			errMsgs += errs.Error() + "\n"
		}
		return errors.New(errMsgs)
	}
	LogEvent("INFO", reflect.TypeOf(data).String()+" Data Validated Successfully...")
	return nil
}
