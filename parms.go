// Copyright 2016 Martin Angers. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Package httpparms provides helper functions and types to load the
// content of an HTTP request into a Go struct. It supports loading the
// query string parameters, the form-encoded body and the JSON-encoded
// body. If the struct implements the `Validator` interface, it also
// validates the values.
//
// It supports various form decoders and JSON unmarshalers. Common
// such packages that satisfy the FormDecoder function are
//     - github.com/go-playground/form
//     - github.com/gorilla/schema
//
// Common packages that satisfy the JSONUmarshaler function are
//     - encoding/json in the standard library
//     - pquerna/ffjson/ffjson
//
package httpparms

import (
	"encoding/json"
	"errors"
	"io/ioutil"
	"net/http"
	"net/url"
)

// FormDecoderAdapter is an adapter for form decoder functions that
// take an url.Values type instead of a map[string][]string.
func FormDecoderAdapter(fn func(v interface{}, vals url.Values) error) FormDecoder {
	return func(v interface{}, vals map[string][]string) error {
		return fn(v, url.Values(vals))
	}
}

// FormDecoder is the function signature of a function
// that decodes values from vals into v.
type FormDecoder func(v interface{}, vals map[string][]string) error

// JSONUnmarshaler is the function signature of a function
// that unmarshals JSON from data to v.
type JSONUnmarshaler func(data []byte, v interface{}) error

// Parser decodes request parameters into a struct and validates
// the values if the struct implements Validator.
//
// If the FormDecoder and JSONUnmarshaler used by the Parser are
// safe for concurrent use, the the Parser is also safe for concurrent
// use.
type Parser struct {
	// Form is the FormDecoder function to use to decode form
	// values from a map. If it is nil, form decoding will fail
	// with an error.
	Form FormDecoder

	// JSON is the JSONUnmarshaler function to use to unmarshal
	// JSON from a slice of bytes. If it is nil, json.Unmarshal
	// from the standard library is used.
	JSON JSONUnmarshaler
}

// Validator defines the method required for a type to validate itself.
type Validator interface {
	Validate() error
}

func (p *Parser) schemaDecode(dst interface{}, vals url.Values) error {
	if p.Form == nil {
		return errors.New("httpparms: no form decoder")
	}
	return p.Form(dst, vals)
}

// ParseQueryForm parses the Form parameters of r into dst. The parameters
// may be provided in the query string or in the form-encoded body.
// The dst value must be a pointer to a struct that contains fields
// matching the form parameters, possibly using `schema` struct tags.
// If dst is a Validator, Validate is called and its error returned.
func (p *Parser) ParseQueryForm(r *http.Request, dst interface{}) error {
	if err := r.ParseForm(); err != nil {
		return err
	}

	if err := p.schemaDecode(dst, r.Form); err != nil {
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
func (p *Parser) ParseJSON(r *http.Request, dst interface{}) error {
	b, err := ioutil.ReadAll(r.Body)
	if err != nil {
		return err
	}
	if len(b) > 0 {
		fn := p.JSON
		if fn == nil {
			fn = json.Unmarshal
		}
		if err := fn(b, dst); err != nil {
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
func (p *Parser) ParseQueryJSON(r *http.Request, dst interface{}) error {
	vals := r.URL.Query()
	if err := p.schemaDecode(dst, vals); err != nil {
		return err
	}
	return p.ParseJSON(r, dst)
}

// ParseQuery parses the query values and stores the values in dst. If
// dst is a Validator, Validate is called and its error returned.
func (p *Parser) ParseQuery(r *http.Request, dst interface{}) error {
	vals := r.URL.Query()
	if err := p.schemaDecode(dst, vals); err != nil {
		return err
	}
	if val, ok := dst.(Validator); ok {
		return val.Validate()
	}
	return nil
}
