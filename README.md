# httpparms [![GoDoc](https://godoc.org/github.com/PuerkitoBio/httpparms?status.png)][godoc] [![Build Status](https://semaphoreci.com/api/v1/mna/httpparms/branches/master/badge.svg)](https://semaphoreci.com/mna/httpparms)

Package httpparms provides helper functions and types to load the content of an HTTP request into a Go struct. It supports loading the query string parameters, the form-encoded body and the JSON-encoded body. If the struct implements the `Validator` interface, it also validates the values.

See the [godoc][] for full documentation.

## Installation

```
$ go get github.com/PuerkitoBio/httpparms
```

Use `-u` to update, `-t` to install test dependencies.

## Example

```
type parmTest struct {
	S string
	I int    `schema:"-"`
	Q string `schema:"q" json:"q_value"`
}

func (pt *parmTest) Validate() error {
    if pt.S == "" {
        return errors.New("parameter `s` is required")
    }
	if pt.I > 2 {
		return errors.New("parameter `i` is too big")
	}
	return nil
}

var parser = &httpparms.Parser{
    // use github.com/gorilla/schema as form decoder
    Form: schema.NewDecoder().Decode,
}

func myHandler(w http.ResponseWriter, r *http.Request) {
    var pt parmTest
    if err := parser.ParseQueryJSON(r, &pt); err != nil {
        // Optionally get the parameter names in error with
        // parms := parser.ParametersFromErr(err)
        // and use this in the error message.
        http.Error(w, err.Error(), http.StatusBadRequest)
        return
    }
    // process the request with valid parameters...
}
```

## License

The [BSD 3-clause][bsd] license, see LICENSE file.

[bsd]: http://opensource.org/licenses/BSD-3-Clause
[godoc]: http://godoc.org/github.com/PuerkitoBio/httpparms

