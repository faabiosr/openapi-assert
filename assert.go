package assert

import (
	"encoding/json"
	"fmt"
	"github.com/pkg/errors"
	"github.com/xeipuuv/gojsonschema"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"
)

var (
	// FailMessage define the base error message.
	FailMessage = "failed asserting that %s"

	// ErrMediaType shows the media type error format.
	ErrMediaType = `'%s' is an allowed media type (%s)`

	// ErrRequestHeaders shows the request headers error format.
	ErrRequestHeaders = `'%s' is a valid request header (%s)`

	// ErrRequestQuery shows the request query error format.
	ErrRequestQuery = `'%s' is a valid request query (%s)`

	// ErrResponseHeaders shows the response headers error format.
	ErrResponseHeaders = `'%s' is a valid response header (%s)`

	// ErrValidation wrap the json schema validation errors
	ErrValidation = "unable to load the validation"

	// ErrJson wrap the json marshall errors
	ErrJson = "unable to marshal"

	// ErrRequestBody shows the request body error format.
	ErrRequestBody = `'%s' is a valid request body (%s)`

	// ErrResponseBody shows the response body error format.
	ErrResponseBody = `'%s' is a valid response body (%s)`
)

func failf(format string, a ...interface{}) error {
	return fmt.Errorf(FailMessage, fmt.Sprintf(format, a...))
}

// RequestMediaType asserts request media type against a list.
func RequestMediaType(mediaType string, doc Document, path, method string) error {
	types, err := doc.RequestMediaTypes(path, method)

	if err != nil {
		return err
	}

	for _, v := range types {
		if v == mediaType {
			return nil
		}
	}

	return failf(ErrMediaType, mediaType, strings.Join(types, ", "))
}

// ResponseMediaType asserts response media type against a list.
func ResponseMediaType(mediaType string, doc Document, path, method string) error {
	types, err := doc.ResponseMediaTypes(path, method)

	if err != nil {
		return err
	}

	for _, v := range types {
		if v == mediaType {
			return nil
		}
	}

	return failf(ErrMediaType, mediaType, strings.Join(types, ", "))
}

// RequestHeaders asserts rquest headers againt a schema header list.
func RequestHeaders(header http.Header, doc Document, path, method string) error {
	schema, err := doc.RequestHeaders(path, method)

	if err != nil {
		return err
	}

	headers := map[string]string{}

	for k, v := range header {
		headers[k] = strings.Join(v, ", ")
	}

	result, err := gojsonschema.Validate(
		gojsonschema.NewGoLoader(schema),
		gojsonschema.NewGoLoader(headers),
	)

	if err != nil {
		return errors.Wrap(err, ErrValidation)
	}

	if result.Valid() {
		return nil
	}

	data, err := json.Marshal(headers)

	if err != nil {
		return errors.Wrap(err, ErrJson)
	}

	errorMessages := []string{}

	for _, v := range result.Errors() {
		errorMessages = append(errorMessages, v.Description())
	}

	return failf(ErrRequestHeaders, string(data), strings.Join(errorMessages, ", "))
}

// ResponseHeaders asserts response headers againt a schema header list.
func ResponseHeaders(header http.Header, doc Document, path, method string, statusCode int) error {
	schema, err := doc.ResponseHeaders(path, method, statusCode)

	if err != nil {
		return err
	}

	headers := map[string]string{}

	for k, v := range header {
		headers[k] = strings.Join(v, ", ")
	}

	result, err := gojsonschema.Validate(
		gojsonschema.NewGoLoader(schema),
		gojsonschema.NewGoLoader(headers),
	)

	if err != nil {
		return errors.Wrap(err, ErrValidation)
	}

	if result.Valid() {
		return nil
	}

	data, err := json.Marshal(headers)

	if err != nil {
		return errors.Wrap(err, ErrJson)
	}

	errorMessages := []string{}

	for _, v := range result.Errors() {
		errorMessages = append(errorMessages, v.Description())
	}

	return failf(ErrResponseHeaders, string(data), strings.Join(errorMessages, ", "))
}

// RequestQuery asserts request query againt a schema.
func RequestQuery(query url.Values, doc Document, path, method string) error {
	schema, err := doc.RequestQuery(path, method)

	if err != nil {
		return err
	}

	result, err := gojsonschema.Validate(
		gojsonschema.NewGoLoader(schema),
		gojsonschema.NewGoLoader(query),
	)

	if err != nil {
		return errors.Wrap(err, ErrValidation)
	}

	if result.Valid() {
		return nil
	}

	data, err := json.Marshal(query)

	if err != nil {
		return errors.Wrap(err, ErrJson)
	}

	errorMessages := []string{}

	for _, v := range result.Errors() {
		errorMessages = append(errorMessages, v.Description())
	}

	return failf(ErrRequestQuery, string(data), strings.Join(errorMessages, ", "))
}

// RequestBody asserts request body against a schema.
func RequestBody(body io.Reader, doc Document, path, method string) error {
	schema, err := doc.RequestBody(path, method)

	if err != nil {
		return err
	}

	data, err := ioutil.ReadAll(body)

	if err != nil {
		return err
	}

	result, err := gojsonschema.Validate(
		gojsonschema.NewGoLoader(schema),
		gojsonschema.NewBytesLoader(data),
	)

	if err != nil {
		return errors.Wrap(err, ErrValidation)
	}

	if result.Valid() {
		return nil
	}

	errorMessages := []string{}

	for _, v := range result.Errors() {
		errorMessages = append(errorMessages, v.Description())
	}

	return failf(ErrRequestBody, string(data), strings.Join(errorMessages, ", "))
}

// ResponseBody asserts response body against a schema.
func ResponseBody(body io.Reader, doc Document, path, method string, statusCode int) error {
	schema, err := doc.ResponseBody(path, method, statusCode)

	if err != nil {
		return err
	}

	data, err := ioutil.ReadAll(body)

	if err != nil {
		return err
	}

	result, err := gojsonschema.Validate(
		gojsonschema.NewGoLoader(schema),
		gojsonschema.NewBytesLoader(data),
	)

	if err != nil {
		return errors.Wrap(err, ErrValidation)
	}

	if result.Valid() {
		return nil
	}

	errorMessages := []string{}

	for _, v := range result.Errors() {
		errorMessages = append(errorMessages, v.Description())
	}

	return failf(ErrResponseBody, string(data), strings.Join(errorMessages, ", "))
}

// Requery asserts http request against a schema.
func Request(req *http.Request, doc Document) error {

	path := req.URL.String()
	method := req.Method

	if err := RequestHeaders(req.Header, doc, path, method); err != nil {
		return err
	}

	if err := RequestMediaType(req.Header.Get("content-type"), doc, path, method); err != nil && req.Body != nil {
		return err
	}

	if err := RequestQuery(req.URL.Query(), doc, path, method); err != nil {
		return err
	}

	if err := RequestBody(req.Body, doc, path, method); err != nil {
		return err
	}

	return nil
}
