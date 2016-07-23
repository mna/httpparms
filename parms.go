// Copyright 2016 Martin Angers. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Package httpparms provides helper functions and types to load the
// content of an HTTP request into a Go struct. It supports loading the
// query string parameters, the form-encoded body and the JSON-encoded
// body. If the struct implements the `Validator` interface, it also
// validates the values.
//
// It uses the github.com/gorilla/schema package to load form values
// into a struct.
//
// To use non-default struct field names for form values, use the `schema:"name"`
// struct tag as documented in the gorilla/schema package. To use non-
// default struct field names for JSON values, use the `json:"name"`
// struct tag as documented in the stdlib's encoding/json package.
//
package httpparms

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"net/url"

	"github.com/gorilla/schema"
)

// shared schema decoder, safe for concurrent use.
var formDecoder = schema.NewDecoder()

func init() {
	formDecoder.IgnoreUnknownKeys(true)
}

// Validator defines the method required for a type to validate itself.
type Validator interface {
	Validate() error
}

func schemaDecode(vals url.Values, dst interface{}) error {
	return formDecoder.Decode(dst, vals)
}

// ParseQueryForm parses the Form parameters of r into dst. The parameters
// may be provided in the query string or in the form-encoded body.
// The dst value must be a pointer to a struct that contains fields
// matching the form parameters, possibly using `schema` struct tags.
// If dst is a Validator, Validate is called and its error returned.
func ParseQueryForm(r *http.Request, dst interface{}) error {
	if err := r.ParseForm(); err != nil {
		return err
	}

	if err := schemaDecode(r.Form, dst); err != nil {
		return err
	}

	if val, ok := dst.(Validator); ok {
		return val.Validate()
	}
	return nil
}

// ParseJSON parses the body of the request as JSON and unmarshals it into
// dst. If dst is a Validator, Validate is called and its error returned.
// The body is parsed as JSON regardless of the content-type of the request.
func ParseJSON(r *http.Request, dst interface{}) error {
	b, err := ioutil.ReadAll(r.Body)
	if err != nil {
		return err
	}
	if len(b) > 0 {
		if err := json.Unmarshal(b, dst); err != nil {
			return err
		}
	}

	if val, ok := dst.(Validator); ok {
		return val.Validate()
	}
	return nil
}

// ParseQueryJSON parses the query values and the body of the request as JSON
// and stores the values in dst. If dst is a Validator, Validate is called and
// its error returned. The body is parsed as JSON regardless of the
// content-type of the request.
func ParseQueryJSON(r *http.Request, dst interface{}) error {
	vals := r.URL.Query()
	if err := schemaDecode(vals, dst); err != nil {
		return err
	}
	return ParseJSON(r, dst)
}

// ParseQuery parses the query values and stores the values in dst. If
// dst is a Validator, Validate is called and its error returned.
func ParseQuery(r *http.Request, dst interface{}) error {
	vals := r.URL.Query()
	if err := schemaDecode(vals, dst); err != nil {
		return err
	}
	if val, ok := dst.(Validator); ok {
		return val.Validate()
	}
	return nil
}
