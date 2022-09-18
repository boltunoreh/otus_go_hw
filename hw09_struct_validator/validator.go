package hw09structvalidator

import (
	"errors"
	"fmt"
	"reflect"
	"regexp"
	"strconv"
	"strings"
)

type ValidationError struct {
	Field string
	Err   error
}

const validationTag = "validate"

type ValidationErrors []ValidationError

func (v ValidationErrors) Error() string {
	var errorString string
	for _, validationError := range v {
		errorString += validationError.Field
	}

	return errorString
}

func (v ValidationErrors) Unwrap() error {
	return v[0].Err
}

type ValidationRuleError struct {
	FieldName string
	Err       error
}

func (v ValidationRuleError) Error() string {
	return "field: " + v.FieldName + " has validation rule error: " + v.Err.Error()
}

type InvalidFieldError struct {
	FieldName string
	Type      string
}

func (v InvalidFieldError) Error() string {
	return "field: " + v.FieldName + " has invalid field type: " + v.Type
}

func Validate(v interface{}) error {
	typeOf := reflect.TypeOf(v)

	if typeOf.Kind() != reflect.Struct {
		return errors.New("value is not structure")
	}

	var validationErrors ValidationErrors
	for i := 0; i < typeOf.NumField(); i++ {
		field := typeOf.Field(i)
		fieldValue := reflect.ValueOf(v).Field(i)

		validationRules, err := getValidationRules(field, fieldValue)
		if err != nil {
			return err
		}
		if len(validationRules) == 0 {
			continue
		}

		err = ValidateField(fieldValue, validationRules)
		if err != nil {
			validationError := ValidationError{
				field.Name,
				err,
			}
			validationErrors = append(validationErrors, validationError)
		}
	}

	if len(validationErrors) > 0 {
		return validationErrors
	}

	return nil
}

func getValidationRules(field reflect.StructField, fieldValue reflect.Value) ([]Rule, error) {
	var rules []Rule
	var err error

	ruleStrings := strings.Split(field.Tag.Get(validationTag), "|")
	kind := field.Type.Kind()
	//exhaustive:ignore
	switch kind {
	case reflect.Int:
		rules, err = getIntRules(ruleStrings)
	case reflect.String:
		rules, err = getStringRules(ruleStrings)
	case reflect.Slice:
		if fieldValue.Len() > 0 {
			//exhaustive:ignore
			switch fieldValue.Index(1).Kind() {
			case reflect.Int:
				rules, err = getIntRules(ruleStrings)
			case reflect.String:
				rules, err = getStringRules(ruleStrings)
			default:
			}
		}
	default:
		return nil, InvalidFieldError{
			field.Name,
			kind.String(),
		}
	}

	if err != nil {
		return nil, ValidationRuleError{
			field.Name,
			err,
		}
	}

	return rules, nil
}

func getIntRules(ruleStrings []string) ([]Rule, error) {
	var rules []Rule

	for _, ruleString := range ruleStrings {
		switch {
		case strings.HasPrefix(ruleString, "min:"):
			valueString := strings.TrimPrefix(ruleString, "min:")
			valueInt, err := strconv.ParseInt(valueString, 10, 0)
			if err != nil {
				return nil, err
			}

			rules = append(rules, IntMinRule{
				min: valueInt,
			})
		case strings.HasPrefix(ruleString, "max:"):
			valueString := strings.TrimPrefix(ruleString, "max:")
			valueInt, err := strconv.ParseInt(valueString, 10, 0)
			if err != nil {
				return nil, err
			}

			rules = append(rules, IntMaxRule{
				max: valueInt,
			})
		case strings.HasPrefix(ruleString, "in:"):
			valueString := strings.TrimPrefix(ruleString, "in:")
			values := strings.Split(valueString, ",")

			var valuesInt []int64
			for _, value := range values {
				valueInt, err := strconv.ParseInt(value, 10, 0)
				if err != nil {
					return nil, err
				}

				valuesInt = append(valuesInt, valueInt)
			}

			rules = append(rules, IntInRule{
				in: valuesInt,
			})
		}
	}

	return rules, nil
}

func getStringRules(ruleStrings []string) ([]Rule, error) {
	var rules []Rule

	for _, ruleString := range ruleStrings {
		switch {
		case strings.HasPrefix(ruleString, "len:"):
			valueString := strings.TrimPrefix(ruleString, "len:")
			valueInt, err := strconv.Atoi(valueString)
			if err != nil {
				return nil, err
			}

			rules = append(rules, StringLenRule{
				len: valueInt,
			})
		case strings.HasPrefix(ruleString, "regexp:"):
			valueString := strings.TrimPrefix(ruleString, "regexp:")
			regExpr, err := regexp.Compile(valueString)
			if err != nil {
				return nil, err
			}

			rules = append(rules, StringRegexpRule{
				regexp: regExpr,
			})
		case strings.HasPrefix(ruleString, "in:"):
			valueString := strings.TrimPrefix(ruleString, "in:")
			values := strings.Split(valueString, ",")

			rules = append(rules, StringInRule{
				in: values,
			})
		}
	}

	return rules, nil
}

func ValidateField(value reflect.Value, validationRules []Rule) error {
	var fieldErrors []error

	for _, rule := range validationRules {
		var fieldError error

		//exhaustive:ignore
		switch value.Kind() {
		case reflect.Slice:
			fieldError = validateSlice(value, rule)
		case reflect.String:
			fieldError = rule.Validate(value.String())
		case reflect.Int:
			fieldError = rule.Validate(value.Int())
		default:
			return nil
		}

		if fieldError != nil {
			fieldErrors = append(fieldErrors, fieldError)
		}
	}

	if len(fieldErrors) == 0 {
		return nil
	}

	var errorsString string
	for _, fieldError := range fieldErrors {
		errorsString += fmt.Sprint(fieldError)
	}

	return errors.New(errorsString)
}

func validateSlice(value reflect.Value, rule Rule) error {
	var fieldErrors []error

	switch value.Interface().(type) {
	case []int:
		for _, element := range value.Interface().([]int) {
			fieldError := rule.Validate(element)
			if fieldError != nil {
				fieldErrors = append(fieldErrors, fieldError)
			}
		}
	case []string:
		for _, element := range value.Interface().([]string) {
			fieldError := rule.Validate(element)
			if fieldError != nil {
				fieldErrors = append(fieldErrors, fieldError)
			}
		}
	}

	if len(fieldErrors) > 0 {
		var errorsString string
		for _, fieldError := range fieldErrors {
			errorsString += fmt.Sprint(fieldError)
		}

		return errors.New(errorsString)
	}

	return nil
}
