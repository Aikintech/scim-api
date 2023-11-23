package utils

type FormattedValidationErrs map[string][]string
type ValidationErrs map[string]map[string]string

// FormatValidationErrors formats validation errors
// returned by the gookit/validate package
// into a map of field names and error messages
func FormatValidationErrors(errs ValidationErrs) FormattedValidationErrs {
	formatted := make(FormattedValidationErrs)

	for field, err := range errs {
		for _, msg := range err {
			formatted[field] = append(formatted[field], msg)
		}
	}

	return formatted
}
