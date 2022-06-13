package aform

import (
	"fmt"
	"github.com/go-playground/validator/v10"
	"strconv"
	"strings"
	"sync"
)

var doValidateInitOnce sync.Once
var validate *validator.Validate

func defaultValidate() *validator.Validate {
	doValidateInitOnce.Do(func() {
		validate = validator.New()
		_ = validate.RegisterValidation("boolean", func(fl validator.FieldLevel) bool {
			_, err := parseBool(fl.Field().String())
			return err == nil
		})
		setValidationTranslations(validate)
	})
	return validate
}

func validateValue(value string, tag string) []Error {
	err := defaultValidate().Var(value, tag)
	if err != nil {
		var all []Error
		for _, err := range err.(validator.ValidationErrors) {
			all = append(all, ErrorWrap(errorFromFieldError{fieldError: err}))
		}
		return all
	}
	return nil
}

func buildValidationRules(required bool, rules ...string) string {
	var parts []string
	if required {
		parts = append(parts, RequiredErrorCode)
	}
	for _, rule := range rules {
		parts = append(parts, rule)
	}
	return strings.Join(parts, ",")
}

func buildValidationChoicesRule(choices ...string) string {
	var parts []string
	for _, choice := range choices {
		parts = append(parts, fmt.Sprintf("'%s'", choice))
	}
	return ChoiceErrorCode + "=" + strings.Join(parts, " ")
}

func buildValidationMinRule(min uint) string {
	return MinLengthErrorCode + "=" + strconv.FormatUint(uint64(min), 10)
}

func buildValidationMaxRule(max uint) string {
	return MaxLengthErrorCode + "=" + strconv.FormatUint(uint64(max), 10)
}
