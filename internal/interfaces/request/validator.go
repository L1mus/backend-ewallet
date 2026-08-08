package request

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"

	apperror "github.com/L1mus/backend-ewallet/internal/AppError"
	"github.com/go-playground/validator/v10"
)

var validate = validator.New()

func DecodeAndValidate(r *http.Request, dst any) error {
	if err := json.NewDecoder(r.Body).Decode(dst); err != nil {
		return apperror.NewAppErrorValidate("Invalid request body format (make sure JSON is correct)", nil)
	}

	if err := validate.Struct(dst); err != nil {
		var verrs validator.ValidationErrors
		if errors.As(err, &verrs) {
			return apperror.NewAppErrorValidate("Validation failed", toFieldErrors(verrs))
		}
	}
}

func toFieldErrors(verrs validator.ValidationErrors) map[string]string {
	fieldErrors := make(map[string]string, len(verrs))
	for _, fe := range verrs {
		fieldErrors[fe.Field()] = message(fe)
	}
	return fieldErrors
}

func message(fe validator.FieldError) string {
	switch fe.Tag() {
	case "required":
		return fmt.Sprintf("%s required", fe.Field())
	case "email":
		return fmt.Sprintf("%s must be a valid email", fe.Field())
	case "min":
		return fmt.Sprintf("%s minimum %s characters", fe.Field(), fe.Param())
	case "max":
		return fmt.Sprintf("%s maximum %s characters", fe.Field(), fe.Param())
	case "gt":
		return fmt.Sprintf("%s must be greater than %s", fe.Field(), fe.Param())
	case "uuid7":
		return fmt.Sprintf("%s must be a valid UUID", fe.Field())
	case "nefield":
		return fmt.Sprintf("%s cannot be equal to %s", fe.Field(), fe.Param())
	default:
		return fmt.Sprintf("%s invalid (rule: %s)", fe.Field(), fe.Tag())
	}
}
