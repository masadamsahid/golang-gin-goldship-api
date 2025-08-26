package helpers

import (
	"regexp"
	"strings"

	"github.com/go-playground/validator/v10"
)

var (
	matchFirstCap = regexp.MustCompile("(.)([A-Z][a-z]+)")
	matchAllCap   = regexp.MustCompile("([a-z0-9])([A-Z])")
)

func toSnakeCase(str string) string {
	snake := matchFirstCap.ReplaceAllString(str, "${1}_${2}")
	snake = matchAllCap.ReplaceAllString(snake, "${1}_${2}")
	return strings.ToLower(snake)
}

func HandleValidationErrors(validationErrors validator.ValidationErrors) map[string]string {
	errs := make(map[string]string)
	for _, fieldErr := range validationErrors {
		jsonKey := toSnakeCase(fieldErr.Field())

		switch fieldErr.Tag() {
		case "required":
			errs[jsonKey] = fieldErr.Field() + " is required"
		case "min":
			errs[jsonKey] = fieldErr.Field() + " must be at least " + fieldErr.Param() + " characters long"
		case "alphanum":
			errs[jsonKey] = fieldErr.Field() + " must be alphanumeric"
		case "eqfield":
			errs[jsonKey] = fieldErr.Field() + " should be equal to " + fieldErr.Param()
		case "lte":
			errs[jsonKey] = fieldErr.Field() + " should be less than or equal to " + fieldErr.Param()
		case "gte":
			errs[jsonKey] = fieldErr.Field() + " should be greater than or equal to " + fieldErr.Param()
		case "lt":
			errs[jsonKey] = fieldErr.Field() + " should be less than " + fieldErr.Param()
		case "gt":
			errs[jsonKey] = fieldErr.Field() + " should be greater than " + fieldErr.Param()
		case "url":
			errs[jsonKey] = fieldErr.Field() + " must be a valid URL"
		default:
			errs[jsonKey] = "Validation failed for " + fieldErr.Field()
		}
	}

	return errs
}
