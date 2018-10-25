package assert

import (
	"net/http"
)

type (
	// Assertions packs all assert methods into one structure
	Assertions struct {
		doc Document
	}
)

// New returns the Assertions instance
func New(doc Document) *Assertions {
	return &Assertions{doc}
}

// RequestMediaType asserts request media type against a list.
func (a *Assertions) RequestMediaType(mediaType, path, method string) error {
	return RequestMediaType(mediaType, a.doc, path, method)
}

// ResponseMediaType asserts response media type against a list.
func (a *Assertions) ResponseMediaType(mediaType, path, method string) error {
	return ResponseMediaType(mediaType, a.doc, path, method)
}

// RequestHeaders asserts rquest headers againts a schema header list.
func (a *Assertions) RequestHeaders(header http.Header, path, method string) error {
	return RequestHeaders(header, a.doc, path, method)
}
