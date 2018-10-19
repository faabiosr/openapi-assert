package assert

import (
	"fmt"
	"strings"
)

var (
	// FailMessage define the base error message
	FailMessage = "failed asserting that %s"

	// ErrMediaType shows the media type error format
	ErrMediaType = `'%s' is an allowed media type (%s)`
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
