package errorlib

import "fmt"

// WrapError returns a wrapped error formatted as:
// [functionMethodName]:[message] error:[error]
func WrapError(functionMethodName, message string, err error) error {
	return fmt.Errorf("%s:%s errorlib:%w", functionMethodName, message, err)
}
