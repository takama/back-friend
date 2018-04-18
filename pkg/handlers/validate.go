package handlers

import (
	"bufio"
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/takama/backer"
	"github.com/takama/backer/datastore"

	"github.com/takama/bit"
)

const (
	couldNotRecognizeRequestData = "Could not recognize request data"
	incorrectPointsParameter     = "The Points parameter has incorrect type"
	incorrectParameter           = " parameter has incorrect type: "
)

func decodeRecord(record interface{}, c bit.Control) bool {
	if record != nil && c.Request().Body != nil {
		decoder := json.NewDecoder(bufio.NewReader(c.Request().Body))
		decoder.UseNumber()
		if err := decoder.Decode(&record); err == nil {
			return true
		}
	}

	c.Code(http.StatusBadRequest)
	c.Body(couldNotRecognizeRequestData)
	return false
}

func decodePoints(points interface{}, c bit.Control) (backer.Points, bool) {
	number, ok := points.(json.Number)
	if !ok {
		c.Code(http.StatusBadRequest)
		c.Body(incorrectPointsParameter)
		return 0, false
	}
	p, err := number.Float64()
	if err != nil {
		serviceError(err, c)
		return 0, false
	}

	return backer.Points(p), true
}

func decodeString(name string, c bit.Control) (string, bool) {
	str := c.Query(name)

	if !isValidString(str) {
		return "", false
	}

	return str, true
}

func decodeNumber(name string, c bit.Control) (uint64, bool) {
	id, err := strconv.ParseUint(c.Query(name), 10, 64)
	if err != nil {
		c.Code(http.StatusBadRequest)
		c.Body(name + incorrectParameter + err.Error())
		return 0, false
	}

	return id, true
}

// isValidString checks that string contains only values
// from alphanumeric characters and special symbols
func isValidString(str string) bool {
	isValid := true
	for _, b := range str {
		if !('0' <= b && b <= '9' || 'a' <= b && b <= 'z' ||
			'A' <= b && b <= 'Z' || b == '_' || b == '-' || b == '.') {
			isValid = false
			break
		}
	}

	return isValid
}

func badRequest(c bit.Control) {
	c.Code(http.StatusBadRequest)
	c.Body(http.StatusText(http.StatusBadRequest))
}

func notFound(c bit.Control) {
	c.Code(http.StatusNotFound)
	c.Body(datastore.ErrRecordNotFound.Error())
}

func notAcceptable(err error, c bit.Control) {
	c.Code(http.StatusNotAcceptable)
	c.Body(err.Error())
}

func serviceError(err error, c bit.Control) {
	c.Code(http.StatusInternalServerError)
	c.Body(err.Error())
}
