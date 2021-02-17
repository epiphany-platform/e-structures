package test

import (
	"bytes"
	"fmt"
	"strings"
)

type TestValidationErrors []TestValidationError

func (e TestValidationErrors) Error() string {
	buff := bytes.NewBufferString("")

	for _, te := range e {

		buff.WriteString(te.Error())
		buff.WriteString("\n")
	}

	return strings.TrimSpace(buff.String())
}

type TestValidationError struct {
	Key   string
	Field string
	Tag   string
}

func (e TestValidationError) Error() string {
	return fmt.Sprintf("Key: '%s' Error:Field validation for '%s' failed on the '%s' tag", e.Key, e.Field, e.Tag)
}
