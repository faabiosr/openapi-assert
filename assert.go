package assert

import (
	"encoding/json"
	"fmt"
	"github.com/pkg/errors"
	"github.com/xeipuuv/gojsonschema"
	"net/http"
	"strings"
)

var (
	// FailMessage define the base error message.
	FailMessage = "failed asserting that %s"

	// ErrMediaType shows the media type error format.
	ErrMediaType = `'%s' is an allowed media type (%s)`

	// ErrRequestHeaders shows the request headers error format.
	ErrRequestHeaders = `'%s' is a valid request header (%s)`

	// ErrValidation wrap the json schema validation errors
	ErrValidation = "unable to load the validation"

	// ErrJson wrap the json marshall errors
	ErrJson = "unable to marshal"
)

func fail(msg string) error {
	return fmt.Errorf(FailMessage, msg)
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

	return fail(fmt.Sprintf(ErrMediaType, mediaType, strings.Join(types, ", ")))
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

	return fail(fmt.Sprintf(ErrMediaType, mediaType, strings.Join(types, ", ")))
}

// RequestHeaders asserts rquest headers againts a schema header list.
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

	return fail(
		fmt.Sprintf(ErrRequestHeaders, string(data), strings.Join(errorMessages, ", ")),
	)
}
