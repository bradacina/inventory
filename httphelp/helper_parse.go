package httphelp

import (
	"errors"
	"net/http"
	"strconv"
)

var (
	ErrorCannotFindQueryStringParam = errors.New("Could not find query string param")
)

func ParseQueryStringParam(name string, r *http.Request) (string, error) {
	r.ParseForm()

	var value string
	if v, ok := r.Form[name]; ok && len(v) > 0 {
		value = v[0]
	} else {
		return "", ErrorCannotFindQueryStringParam
	}

	return value, nil
}

func ParseIDFromQueryString(r *http.Request) (int, error) {
	sval, err := ParseQueryStringParam("id", r)
	if err != nil {
		return 0, err
	}

	val, err := strconv.Atoi(sval)
	if err != nil {
		return 0, err
	}

	return val, nil
}
