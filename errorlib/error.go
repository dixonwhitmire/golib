// Package errorlib provides uniform formatting for error types.
package errorlib

import "fmt"

// WrapError returns a wrapped error formatted as: [functionMethodName]:[message] error:[error].
func WrapError(functionMethodName, message string, err error) error {
	return fmt.Errorf("%s:%s error:%w", functionMethodName, message, err)
}

// CreateError returns an error formatted as:
// [functionMethodName]:[message]
func CreateError(functionMethodName, message string) error {
	return fmt.Errorf("%s:%s", functionMethodName, message)
}
