package hw09structvalidator

import (
	"errors"
	"regexp"
)

var (
	ErrIntMin       = errors.New("value is under min")
	ErrIntMax       = errors.New("value is above max")
	ErrIntIn        = errors.New("value is out of range")
	ErrIntWrongType = errors.New("value must be type of int")
	ErrStrRegexp    = errors.New("value must satisfy regex")
	ErrStrLen       = errors.New("invalid value length")
	ErrStrIn        = errors.New("value is out of range")
	ErrStrWrongType = errors.New("value must be type of string")
)

type Rule interface {
	Validate(fieldValue interface{}) error
}

type StringLenRule struct {
	len int
}

func (r StringLenRule) Validate(fieldValue interface{}) error {
	switch value := fieldValue.(type) {
	case string:
		if len(value) == r.len {
			return nil
		}

		return ErrStrLen
	default:
		return ErrStrWrongType
	}
}

type StringRegexpRule struct {
	regexp *regexp.Regexp
}

func (r StringRegexpRule) Validate(fieldValue interface{}) error {
	switch value := fieldValue.(type) {
	case string:
		if r.regexp.MatchString(value) {
			return nil
		}

		return ErrStrRegexp
	default:
		return ErrStrWrongType
	}
}

type StringInRule struct {
	in []string
}

func (r StringInRule) Validate(fieldValue interface{}) error {
	switch value := fieldValue.(type) {
	case string:
		err := ErrStrIn

		for _, str := range r.in {
			if str == value {
				err = nil

				break
			}
		}

		return err
	default:
		return ErrStrWrongType
	}
}

type IntInRule struct {
	in []int64
}

func (r IntInRule) Validate(fieldValue interface{}) error {
	switch value := fieldValue.(type) {
	case int64:
		err := ErrIntIn

		for _, number := range r.in {
			if number == value {
				err = nil
				break
			}
		}

		return err
	default:
		return ErrIntWrongType
	}
}

type IntMinRule struct {
	min int64
}

func (r IntMinRule) Validate(fieldValue interface{}) error {
	switch value := fieldValue.(type) {
	case int64:
		if value < r.min {
			return ErrIntMin
		}

		return nil
	default:
		return ErrIntWrongType
	}
}

type IntMaxRule struct {
	max int64
}

func (r IntMaxRule) Validate(fieldValue interface{}) error {
	switch value := fieldValue.(type) {
	case int64:
		if value > r.max {
			return ErrIntMax
		}

		return nil
	default:
		return ErrIntWrongType
	}
}
