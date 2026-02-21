package errorx

import "errors"

var (
	ErrParseNotInt = errors.New("error of convert from string to int")
)
