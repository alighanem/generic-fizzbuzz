package zhttp

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"strings"
)

func ReadIntParam(r *http.Request, key string) (int, error) {
	raw, ok := r.URL.Query()[key]
	if !ok {
		return 0, fmt.Errorf("param %s not found", key)
	}

	if len(raw) == 0 || len(raw[0]) == 0 {
		return 0, fmt.Errorf("param %s empty", key)
	}

	value, err := strconv.Atoi(raw[0])
	if err != nil {
		return 0, fmt.Errorf("cannot parse %s: %s. err: %w", key, raw[0], err)
	}

	return value, nil
}

func ReadStringParam(r *http.Request, key string) (string, error) {
	raw, ok := r.URL.Query()[key]
	if !ok {
		return "", fmt.Errorf("param %s not found", key)
	}

	if len(raw) == 0 {
		return "", fmt.Errorf("param %s empty", key)
	}

	return raw[0], nil
}

type response struct {
	Data   interface{} `json:"data,omitempty"`
	Errors []Error     `json:"errors,omitempty"`
	Meta   interface{} `json:"meta,omitempty"`
}

// Error is an Error which can be written into the http response
type Error struct {
	Code   string `json:"code,omitempty"`
	Detail string `json:"detail,omitempty"`
}

func (e Error) Error() string {
	var msg []string

	if e.Code != "" {
		msg = append(msg, e.Code)
	}

	if e.Detail != "" {
		msg = append(msg, e.Detail)
	}

	return strings.Join(msg, ": ")
}

// Write writes a json response with status, errors, and data into a ResponseWriter.
func Write(w http.ResponseWriter, httpStatus int, data interface{}, errs ...error) error {
	return writeFull(w, httpStatus, data, errs...)
}

// WriteError is an helper to hide data param and Error struct creation in Write calls.
func WriteError(w http.ResponseWriter, httpStatus int, code string, detail string) error {
	return writeFull(w, httpStatus, nil, Error{
		Code:   code,
		Detail: detail,
	})
}

func writeFull(w http.ResponseWriter, httpStatus int, data interface{}, errs ...error) error {
	if data == nil && errs == nil {
		// If everything is nil, don't write anything to the ResponseWriter,
		// not even the empty wrapper "{}".
		// We also don't set a content type since there is no content.
		// This is best suited for the StatusNoContent (204) status.
		w.WriteHeader(httpStatus)
		return nil
	}

	var zerrs []Error
	for _, err := range errs {
		switch err := err.(type) {
		case Error: // default use case
			zerrs = append(zerrs, err)
		default:
			zerrs = append(zerrs, Error{Code: "internal_error", Detail: err.Error()})
		}

	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(httpStatus)
	raw, err := json.Marshal(response{
		Data:   data,
		Errors: zerrs,
	})
	if err != nil {
		return fmt.Errorf("marshaling payload: %w", err)
	}
	_, err = w.Write(raw)
	return err
}
