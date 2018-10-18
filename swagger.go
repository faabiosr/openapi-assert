package assert

import (
	"github.com/go-openapi/jsonpointer"
	"github.com/go-openapi/loads"
	"github.com/go-openapi/spec"
	"github.com/pkg/errors"
	"github.com/yosida95/uritemplate"
	"strings"
)

var (
	// ErrResourceURI returns an error when uri does not match.
	ErrResourceURI = "resource uri does not match"

	// ErrSwaggerLoad returns an error when load swagger document.
	ErrSwaggerLoad = "unable to load the document by uri"

	// ErrSwaggerExpand returns an error of expanding schema.
	ErrSwaggerExpand = "unable to expand the document"

	// ErrNodeNotFound returns an error when node does not exists.
	ErrNodeNotFound = "node does not exists"
)

type (
	// Swagger stores the loaded swagger spec.
	Swagger struct {
		spec *spec.Swagger
	}
)

// LoadFromURI loads and expands swagger document by uri.
func LoadFromURI(uri string) (*Swagger, error) {
	doc, err := loads.Spec(uri)

	if err != nil {
		return nil, errors.Wrap(err, ErrSwaggerLoad)
	}

	doc, err = doc.Expanded()

	if err != nil {
		return nil, errors.Wrap(err, ErrSwaggerExpand)
	}

	return &Swagger{doc.Spec()}, nil
}

// FindPath searches for an uri in document and returns the path.
func (s *Swagger) FindPath(uri string) (string, error) {
	for path := range s.spec.Paths.Paths {
		tmpl, err := uritemplate.New(s.spec.BasePath + path)

		if err != nil {
			return "", errors.Wrap(err, ErrResourceURI)
		}

		if tmpl.Regexp().MatchString(uri) {
			return strings.Replace(path, "/", "~1", -1), nil
		}
	}

	return "", errors.New(ErrResourceURI)
}

// FindNode searches a node using segments in the schema.
func (s *Swagger) FindNode(segments ...string) (interface{}, error) {
	segments = append([]string{""}, segments...)

	pointer, err := jsonpointer.New(strings.Join(segments, "/"))

	if err != nil {
		return nil, err
	}

	data, _, err := pointer.Get(s.spec)

	if err != nil {
		return nil, errors.Wrap(err, ErrNodeNotFound)
	}

	return data, nil
}
