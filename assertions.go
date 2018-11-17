package assert

import (
	"net/http"
	"net/url"
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

// RequestHeaders asserts rquest headers againt a schema header list.
func (a *Assertions) RequestHeaders(header http.Header, path, method string) error {
	return RequestHeaders(header, a.doc, path, method)
}

// ResponseHeaders asserts response headers againt a schema header list.
func (a *Assertions) ResponseHeaders(header http.Header, path, method string, statusCode int) error {
	return ResponseHeaders(header, a.doc, path, method, statusCode)
}

// RequestQuery asserts request query againt a schema query list.
func (a *Assertions) RequestQuery(query url.Values, path, method string) error {
	return RequestQuery(query, a.doc, path, method)
}
