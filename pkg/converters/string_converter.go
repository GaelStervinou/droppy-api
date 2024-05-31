package converters

import (
	"fmt"
	"strconv"
)

func StringToUint(s string) (uint, error) {
	value, err := strconv.ParseUint(s, 10, 64) // base 10, bit size 64
	if err != nil {
		return 0, err
	}

	// Check if the parsed value fits into uint
	if value > uint64(^uint(0)) {
		return 0, fmt.Errorf("value out of range for uint")
	}

	return uint(value), nil
}