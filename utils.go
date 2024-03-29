package main

import "strconv"

// e.g. x = "0x32439"
func hexStringToInt(x string) (int64, error) {
	y, err := strconv.ParseInt(x[2:], 16, 64) // x[2:] to ignore 0x of hex
	if err != nil {
		return 0, err
	}
	return y, nil
}
