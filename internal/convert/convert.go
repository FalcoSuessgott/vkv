package convert

import (
	"fmt"
	"strconv"
)

//nolint
var errConversionError = func(v interface{}) error {
	return fmt.Errorf("cannot convert value %v (type %T) to integer", v, v)
}

// ToInteger converts given value to integer.
func ToInteger(v interface{}) (int, error) {
	i, err := strconv.Atoi(v.(string))
	if err != nil {
		return i, errConversionError(v)
	}

	return i, nil
}
