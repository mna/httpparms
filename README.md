# httpparms [![GoDoc](https://godoc.org/github.com/PuerkitoBio/httpparms?status.png)][godoc] [![Build Status](https://semaphoreci.com/api/v1/mna/httpparms/branches/master/badge.svg)](https://semaphoreci.com/mna/httpparms)

Package httpparms provides helper functions and mechanisms to load the content of an HTTP request into a Go struct. It supports loading the query string parameters, the form-encoded body and/or the JSON-encoded body. If the struct implements the `Validator` interface, it also validates the values.

See the [godoc][] for full documentation.

## License

The [BSD 3-clause][bsd] license, see LICENSE file.

[bsd]: http://opensource.org/licenses/BSD-3-Clause
[godoc]: http://godoc.org/github.com/PuerkitoBio/httpparms

