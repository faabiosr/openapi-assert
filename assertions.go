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

// Assertions packs all assert methods into one structure.
type Assertions struct {
	doc Document
}

// New returns the Assertions instance.
func New(doc Document) *Assertions {
	return &Assertions{doc}
}

// RequestMediaType asserts request media type against a list.
func (a *Assertions) RequestMediaType(mediaType, path, method string) error {
	types, err := a.doc.RequestMediaTypes(path, method)
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
func (a *Assertions) ResponseMediaType(mediaType, path, method string) error {
	types, err := a.doc.ResponseMediaTypes(path, method)
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
func (a *Assertions) RequestHeaders(header http.Header, path, method string) error {
	schema, err := a.doc.RequestHeaders(path, method)
	if err != nil {
		return err
	}

	headers := map[string]string{}

	for k, v := range header {
		headers[k] = strings.Join(v, ", ")
	}

	result, err := a.validate(schema, headers)
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
func (a *Assertions) ResponseHeaders(header http.Header, path, method string, statusCode int) error {
	schema, err := a.doc.ResponseHeaders(path, method, statusCode)
	if err != nil {
		return err
	}

	headers := map[string]string{}

	for k, v := range header {
		headers[k] = strings.Join(v, ", ")
	}

	result, err := a.validate(schema, headers)
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

// RequestQuery asserts request query againt a schema query list.
func (a *Assertions) RequestQuery(query url.Values, path, method string) error {
	schema, err := a.doc.RequestQuery(path, method)
	if err != nil {
		return err
	}

	result, err := a.validate(schema, query)
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
func (a *Assertions) RequestBody(body io.Reader, path, method string) error {
	schema, err := a.doc.RequestBody(path, method)
	if err != nil {
		return err
	}

	data, err := ioutil.ReadAll(body)
	if err != nil {
		return err
	}

	result, err := a.validate(schema, data)
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
func (a *Assertions) ResponseBody(body io.Reader, path, method string, statusCode int) error {
	schema, err := a.doc.ResponseBody(path, method, statusCode)
	if err != nil {
		return err
	}

	data, err := ioutil.ReadAll(body)
	if err != nil {
		return err
	}

	result, err := a.validate(schema, data)
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
func (a *Assertions) Request(req *http.Request) error {
	path := req.URL.String()
	method := req.Method

	if err := a.RequestHeaders(req.Header, path, method); err != nil {
		return err
	}

	if err := a.RequestMediaType(req.Header.Get("content-type"), path, method); err != nil && req.Body != nil {
		return err
	}

	if err := a.RequestQuery(req.URL.Query(), path, method); err != nil {
		return err
	}

	buf := bytes.NewBuffer(make([]byte, 0))
	reader := io.TeeReader(req.Body, buf)
	req.Body = ioutil.NopCloser(buf)

	err := a.RequestBody(reader, path, method)
	if err != nil && err == ErrBodyNotFound {
		return nil
	}

	return err
}

// Response asserts http response against a schema.
func (a *Assertions) Response(res *http.Response) error {
	path := res.Request.URL.Path
	method := res.Request.Method
	statusCode := res.StatusCode

	if err := a.ResponseHeaders(res.Header, path, method, statusCode); err != nil {
		return err
	}

	if err := a.ResponseMediaType(res.Header.Get("content-type"), path, method); err != nil && res.Body != nil {
		return err
	}

	buf := bytes.NewBuffer(make([]byte, 0))
	reader := io.TeeReader(res.Body, buf)
	res.Body = ioutil.NopCloser(buf)

	return a.ResponseBody(reader, path, method, statusCode)
}

func (a *Assertions) validate(schema, data interface{}) (*gojsonschema.Result, error) {
	loader := gojsonschema.NewGoLoader(data)

	if b, ok := data.([]byte); ok {
		loader = gojsonschema.NewBytesLoader(b)
	}

	return gojsonschema.Validate(
		gojsonschema.NewGoLoader(schema),
		loader,
	)
}

func failf(format string, a ...interface{}) error {
	return fmt.Errorf("failed asserting that %s", fmt.Sprintf(format, a...))
}
