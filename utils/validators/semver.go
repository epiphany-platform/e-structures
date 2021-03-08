package validators

import (
	"fmt"
	"reflect"

	"github.com/Masterminds/semver"
	"github.com/go-playground/validator/v10"
)

func HasVersion(fl validator.FieldLevel) bool {
	field := fl.Field()

	switch field.Kind() {
	case reflect.String:
		f := field.String()
		r, err := checkVersionConstraint(f, fl.Param())
		if err != nil {
			panic(err)
		}
		return r
	}

	panic(fmt.Sprintf("Bad field type %T", field.Interface()))
}

func checkVersionConstraint(version, pattern string) (bool, error) {
	c, err := semver.NewConstraint(pattern)
	if err != nil {
		return false, err
	}
	v, err := semver.NewVersion(version)
	if err != nil {
		return false, err
	}
	if c.Check(v) {
		return true, nil
	}
	return false, nil
}
