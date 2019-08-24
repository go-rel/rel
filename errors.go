package grimoire

import (
	"fmt"
	"strconv"
)

type NoResultError struct{}

func (no NoResultError) Error() string {
	return "No result found"
}

type FieldInvalidError struct {
	Field string
}

func (fie FieldInvalidError) Error() string {
	return fie.Field + " is invalid"
}

type ValidateRequiredError struct {
	Field string
}

func (vre ValidateRequiredError) Error() string {
	return vre.Field + " is required"
}

type ValidateMaxError struct {
	Field string
	Max   int
}

func (vme ValidateMaxError) Error() string {
	return vme.Field + " must be less than " + strconv.Itoa(vme.Max)
}

type ValidateMinError struct {
	Field string
	Min   int
}

func (vme ValidateMinError) Error() string {
	return vme.Field + " must be more than " + strconv.Itoa(vme.Min)
}

type ValidateRangeError struct {
	Field string
	Min   int
	Max   int
}

func (vre ValidateRangeError) Error() string {
	return vre.Field + " must be between " + strconv.Itoa(vre.Min) + " and " + strconv.Itoa(vre.Max)
}

type ValidateInclusionError struct {
	Field  string
	Values []interface{}
}

func (vie ValidateInclusionError) Error() string {
	return vie.Field + " must be one of " + fmt.Sprintf("%v", vie.Values)
}

type ValidateExclusionError struct {
	Field  string
	Values []interface{}
}

func (vee ValidateExclusionError) Error() string {
	return vee.Field + " must not be any of " + fmt.Sprintf("%v", vee.Values)
}

type ValidateFormatError struct {
	Field string
}

func (vfe ValidateFormatError) Error() string {
	return vfe.Field + " format is invalid"
}

type ForeignKeyConstraintError struct {
	Field string
	Value interface{}
	Err   error
}

func (fkce ForeignKeyConstraintError) Error() string {
	return fmt.Sprintf("%v", fkce.Value) + " is not a valid " + fkce.Field
}

type UniqueConstraintError struct {
	Field string
	Value interface{}
	Err   error
}

func (uce UniqueConstraintError) Error() string {
	return "duplicate value " + fmt.Sprintf("%v", uce.Value) + " for " + uce.Field
}
