package main

import (
	"errors"
	"fmt"
	"regexp"

	"github.com/go-playground/validator/v10"
)

var validate = newValidator()

var slugPattern = regexp.MustCompile(`^[a-z0-9]([a-z0-9-]*[a-z0-9])?$`)

func newValidator() *validator.Validate {
	v := validator.New()
	if err := v.RegisterValidation("slug", func(fl validator.FieldLevel) bool {
		return slugPattern.MatchString(fl.Field().String())
	}); err != nil {
		panic(err)
	}
	return v
}

// validateRequest checks a decoded request's struct tags and returns a 400 listing
// every failing field.
func validateRequest(data interface{}) error {
	err := validate.Struct(data)
	if err == nil {
		return nil
	}
	var verrs validator.ValidationErrors
	if !errors.As(err, &verrs) {
		return err
	}
	details := make([]string, 0, len(verrs))
	for _, fe := range verrs {
		details = append(details, fmt.Sprintf("field '%s' failed validation: %s", fe.Namespace(), fe.Tag()))
	}
	return NewValidationError(details)
}
