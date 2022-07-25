package assert

import (
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"strconv"
	"strings"

	"github.com/go-openapi/jsonpointer"
	"github.com/go-openapi/loads"
	"github.com/go-openapi/spec"
	"github.com/yosida95/uritemplate/v3"
)

const (
	// ErrSwaggerLoad returns an error when load swagger document.
	ErrSwaggerLoad = err("unable to load the document by uri")

	// ErrBodyNotFound returns an error when body does not exists.
	ErrBodyNotFound = err("body does not exists")
)

// swagger stores the loaded swagger spec.
type swagger struct {
	spec *spec.Swagger
}

var _ Document = &swagger{}

// LoadFromURI loads and expands swagger document by uri.
func LoadFromURI(uri string) (Document, error) {
	doc, err := loads.Spec(uri)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", ErrSwaggerLoad, err)
	}

	return load(doc)
}

// LoadFromReader loads and expand swagger document from io.Reader.
func LoadFromReader(r io.Reader) (Document, error) {
	data, err := ioutil.ReadAll(r)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", ErrSwaggerLoad, err)
	}

	doc, err := loads.Analyzed(data, "")
	if err != nil {
		return nil, fmt.Errorf("%s: %w", ErrSwaggerLoad, err)
	}

	return load(doc)
}

func load(doc *loads.Document) (Document, error) {
	doc, err := doc.Expanded()
	if err != nil {
		return nil, fmt.Errorf("unable to expand the document: %w", err)
	}

	return &swagger{doc.Spec()}, nil
}

// findPath searches for an uri in document and returns the path.
func (s *swagger) findPath(uri string) (string, error) {
	for path := range s.spec.Paths.Paths {
		tmpl, err := uritemplate.New(s.spec.BasePath + path)
		if err != nil {
			return "", fmt.Errorf("resource uri does not match: %w", err)
		}

		if tmpl.Regexp().MatchString(uri) {
			return strings.ReplaceAll(path, "/", "~1"), nil
		}
	}

	return "", errors.New("resource uri does not match")
}

// findNode searches a node using segments in the schema.
func (s *swagger) findNode(segments ...string) (interface{}, error) {
	segments = append([]string{""}, segments...)

	pointer, err := jsonpointer.New(strings.Join(segments, "/"))
	if err != nil {
		return nil, err
	}

	data, _, err := pointer.Get(s.spec)
	if err != nil {
		return nil, fmt.Errorf("node does not exists: %w", err)
	}

	return data, nil
}

func (s *swagger) mediaTypes(path, method, segment string) ([]string, error) {
	path, err := s.findPath(path)
	if err != nil {
		return []string{}, err
	}

	method = strings.ToLower(method)

	data, err := s.findNode("paths", path, method, segment)
	if err != nil {
		return []string{}, err
	}

	var types []string

	if data != nil {
		types = data.([]string)
	}

	if len(types) > 0 {
		return types, nil
	}

	data, err = s.findNode(segment)
	if err != nil {
		return []string{}, err
	}

	if data != nil {
		types = data.([]string)
	}

	return types, nil
}

// RequestMediaTypes retrives a list of request media types allowed.
func (s *swagger) RequestMediaTypes(path, method string) ([]string, error) {
	return s.mediaTypes(path, method, "consumes")
}

// ResponseMediaTypes retrives a list of response media types allowed.
func (s *swagger) ResponseMediaTypes(path, method string) ([]string, error) {
	return s.mediaTypes(path, method, "produces")
}

func (s *swagger) requestParameters(path, method string) ([]spec.Parameter, error) {
	var params []spec.Parameter

	path, err := s.findPath(path)
	if err != nil {
		return params, err
	}

	data, _ := s.findNode("paths", path, "parameters")
	if data != nil {
		params = append(params, data.([]spec.Parameter)...)
	}

	data, _ = s.findNode("paths", path, strings.ToLower(method), "parameters")
	if data != nil {
		params = append(params, data.([]spec.Parameter)...)
	}

	return params, nil
}

func (s *swagger) response(path, method string, statusCode int) (spec.Response, error) {
	var res spec.Response

	path, err := s.findPath(path)
	if err != nil {
		return res, err
	}

	method = strings.ToLower(method)

	data, err := s.findNode("paths", path, method, "responses", strconv.Itoa(statusCode))
	if data != nil && err == nil {
		return data.(spec.Response), nil
	}

	data, err = s.findNode("paths", path, method, "responses", "default")
	if err != nil {
		return res, err
	}

	return *data.(*spec.Response), err
}

// RequestHeaders retrieves a list of request headers.
func (s *swagger) RequestHeaders(path, method string) (Headers, error) {
	headers := Headers{}

	params, err := s.requestParameters(path, method)
	if err != nil {
		return headers, err
	}

	required := Required{}

	for _, param := range params {
		if param.In != "header" {
			continue
		}

		name := strings.ToLower(param.Name)
		headers[name] = &Param{
			param.Type,
			param.Description,
			param.In,
		}

		if param.Required {
			required = append(required, name)
		}
	}

	if len(required) > 0 {
		headers["required"] = required
	}

	return headers, nil
}

// ResponseHeaders retrieves a list of response headers.
func (s *swagger) ResponseHeaders(path, method string, statusCode int) (Headers, error) {
	headers := Headers{}

	res, err := s.response(path, method, statusCode)
	if err != nil {
		return headers, err
	}

	required := []string{}

	for k, schema := range res.Headers {
		name := strings.ToLower(k)
		headers[name] = schema

		required = append(required, name)
	}

	if len(required) > 0 {
		headers["required"] = required
	}

	return headers, nil
}

// RequestQuery retrieves a list of request query.
func (s *swagger) RequestQuery(path, method string) (Query, error) {
	query := Query{}

	params, err := s.requestParameters(path, method)
	if err != nil {
		return query, err
	}

	required := Required{}

	for _, param := range params {
		if param.In != "query" {
			continue
		}

		name := strings.ToLower(param.Name)
		query[name] = &Param{
			param.Type,
			param.Description,
			param.In,
		}

		if param.Required {
			required = append(required, name)
		}
	}

	if len(required) > 0 {
		query["required"] = required
	}

	return query, nil
}

// RequestBody retrieves the request body.
func (s *swagger) RequestBody(path, method string) (Body, error) {
	params, err := s.requestParameters(path, method)
	if err != nil {
		return nil, err
	}

	for _, param := range params {
		if param.In == "body" {
			return Body(param.Schema), nil
		}
	}

	return nil, ErrBodyNotFound
}

// ResponseBody retrieves the response body.
func (s *swagger) ResponseBody(path, method string, statusCode int) (Body, error) {
	res, err := s.response(path, method, statusCode)
	if err != nil {
		return nil, err
	}

	if res.Schema != nil {
		return Body(res.Schema), nil
	}

	return nil, ErrBodyNotFound
}
