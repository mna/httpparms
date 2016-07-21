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

// ParameterError is an error that stores information about the specific
// parameter that caused the error.
type ParameterError struct {
	Parameter string
	Err       error
}

// Error returns the error message for the ParameterError.
func (e ParameterError) Error() string {
	return e.Err.Error()
}

// ParameterFromErr returns the parameter name that caused the error
// if err is a ParameterError, or an empty string otherwise.
func ParameterFromErr(err error) string {
	if pe, ok := err.(ParameterError); ok {
		return pe.Parameter
	}
	return ""
}

// Validator defines the method required for a type to validate itself.
type Validator interface {
	Validate() error
}

func schemaDecode(vals url.Values, dst interface{}) error {
	err := formDecoder.Decode(dst, vals)
	if err != nil {
		// try to grab information about a specific field
		switch err := err.(type) {
		case schema.ConversionError:
			return ParameterError{Parameter: err.Key, Err: err}
		case schema.MultiError:
			for _, e := range err {
				if ce, ok := e.(schema.ConversionError); ok {
					return ParameterError{Parameter: ce.Key, Err: err}
				}
			}
		}
		return err
	}
	return nil
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
