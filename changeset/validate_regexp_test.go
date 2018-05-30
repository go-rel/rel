package changeset

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestValidateRegexp(t *testing.T) {
	exp := regexp.MustCompile(`foo.*`)

	tests := []interface{}{
		"seafood",
		1,
		2.0,
		false,
	}

	for _, tt := range tests {
		t.Run(fmt.Sprintf("%T", tt), func(t *testing.T) {
			ch := &Changeset{
				changes: map[string]interface{}{
					"field": tt,
				},
			}

			ValidateRegexp(ch, "field", exp)
			assert.Nil(t, ch.Errors())
		})
	}
}

func TestValidateRegexp_error(t *testing.T) {
	ch := &Changeset{
		changes: map[string]interface{}{
			"field": "seafood",
		},
	}

	ValidateRegexp(ch, "field", regexp.MustCompile(`boo.*`))
	assert.NotNil(t, ch.Errors())
	assert.Equal(t, "field's format is invalid", ch.Error().Error())
}

func TestValidateRegexp_missing(t *testing.T) {
	ch := &Changeset{}
	ValidateRegexp(ch, "field", regexp.MustCompile(`foo.*`))
	assert.Nil(t, ch.Errors())
}
