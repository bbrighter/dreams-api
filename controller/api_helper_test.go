package controller

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestStringToUint(t *testing.T) {
	t.Parallel()
	var testString string
	var err error
	var ui uint

	testString = "1"
	ui, err = stringToUint(testString)
	assert.NoError(t, err)
	assert.EqualValues(t, 1, ui)

	testString = "12345678"
	ui, err = stringToUint(testString)
	assert.NoError(t, err)
	assert.EqualValues(t, 12345678, ui)

	testString = "a"
	_, err = stringToUint(testString)
	assert.Error(t, err)

	testString = "-1"
	_, err = stringToUint(testString)
	assert.Error(t, err)
}
