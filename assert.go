package assert

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"

	"github.com/xeipuuv/gojsonschema"
)

func failf(format string, a ...interface{}) error {
	return fmt.Errorf("failed asserting that %s", fmt.Sprintf(format, a...))
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

	ts := strings.Join(types, ", ")

	return failf(`'%s' is an allowed media type (%s)`, mediaType, ts)
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

	ts := strings.Join(types, ", ")

	return failf(`'%s' is an allowed media type (%s)`, mediaType, ts)
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
		return err
	}

	if result.Valid() {
		return nil
	}

	data, err := json.Marshal(headers)
	if err != nil {
		return err
	}

	errorMessages := []string{}

	for _, v := range result.Errors() {
		errorMessages = append(errorMessages, v.Description())
	}

	errs := strings.Join(errorMessages, ", ")

	return failf(`'%s' is a valid request header (%s)`, string(data), errs)
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
		return err
	}

	if result.Valid() {
		return nil
	}

	data, err := json.Marshal(headers)
	if err != nil {
		return err
	}

	errorMessages := []string{}

	for _, v := range result.Errors() {
		errorMessages = append(errorMessages, v.Description())
	}

	errs := strings.Join(errorMessages, ", ")

	return failf(`'%s' is a valid response header (%s)`, string(data), errs)
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
		return err
	}

	if result.Valid() {
		return nil
	}

	data, err := json.Marshal(query)
	if err != nil {
		return err
	}

	errorMessages := []string{}

	for _, v := range result.Errors() {
		errorMessages = append(errorMessages, v.Description())
	}

	errs := strings.Join(errorMessages, ", ")

	return failf(`'%s' is a valid request query (%s)`, string(data), errs)
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
		return err
	}

	if result.Valid() {
		return nil
	}

	errorMessages := []string{}

	for _, v := range result.Errors() {
		errorMessages = append(errorMessages, v.Description())
	}

	errs := strings.Join(errorMessages, ", ")

	return failf(`'%s' is a valid request body (%s)`, string(data), errs)
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
		return err
	}

	if result.Valid() {
		return nil
	}

	errorMessages := []string{}

	for _, v := range result.Errors() {
		errorMessages = append(errorMessages, v.Description())
	}

	errs := strings.Join(errorMessages, ", ")

	return failf(`'%s' is a valid response body (%s)`, string(data), errs)
}

// Request asserts http request against a schema.
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

	buf := bytes.NewBuffer(make([]byte, 0))
	reader := io.TeeReader(req.Body, buf)
	req.Body = ioutil.NopCloser(buf)

	err := RequestBody(reader, doc, path, method)
	if err != nil && err == ErrBodyNotFound {
		return nil
	}

	return err
}

// Response asserts http response against a schema.
func Response(res *http.Response, doc Document) error {
	path := res.Request.URL.Path
	method := res.Request.Method
	statusCode := res.StatusCode

	if err := ResponseHeaders(res.Header, doc, path, method, statusCode); err != nil {
		return err
	}

	if err := ResponseMediaType(res.Header.Get("content-type"), doc, path, method); err != nil && res.Body != nil {
		return err
	}

	buf := bytes.NewBuffer(make([]byte, 0))
	reader := io.TeeReader(res.Body, buf)
	res.Body = ioutil.NopCloser(buf)

	return ResponseBody(reader, doc, path, method, statusCode)
}
