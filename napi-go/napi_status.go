package napi

import "frm/napi-go/internal/napi"

// Process status to return error if StatusOK return nil on error
func mustValueErr[T any](input T, status napi.Status) (T, error) {
	return input, status.ToError()
}

// return error from status
func singleMustValueErr(status napi.Status) error {
	return status.ToError()
}

// Process status to return error if StatusOK return nil on error
func mustValueErr2[T any](input T, _ bool, status napi.Status) (T, error) {
	return input, status.ToError()
}
