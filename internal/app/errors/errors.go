package errors

import "errors"

const StatusOK = "OK"
const StatusError = "ERROR"

var ErrNotFound = errors.New("not found")
var ErrNoPerm = errors.New("you not have permission")
var ErrNotInserted = errors.New("not inserted")
