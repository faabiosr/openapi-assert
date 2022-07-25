package assert

// Headers is a list of headers in json schema format.
type Headers map[string]interface{}

// Query is a list of query parameters in json schema format.
type Query map[string]interface{}

// Param is a document parameter in json schema format.
type Param struct {
	Type        string
	Description string
	In          string
}

// Body is a document schema in json format.
type Body interface{}

// Required is a list of required parameters.
type Required []string

// Document that defines the contract for reading OpenAPI documents.
type Document interface {
	// RequestMediaTypes retrives a list of request media types allowed.
	RequestMediaTypes(path, method string) ([]string, error)

	// ResponseMediaTypes retrives a list of response media types allowed.
	ResponseMediaTypes(path, method string) ([]string, error)

	// RequestHeaders retrieves a list of request headers.
	RequestHeaders(path, method string) (Headers, error)

	// ResponseHeaders retrieves a list of response headers.
	ResponseHeaders(path, method string, statusCode int) (Headers, error)

	// RequestQuery retrieves a list of request query.
	RequestQuery(path, method string) (Query, error)

	// RequestBody retrieves the request body.
	RequestBody(path, method string) (Body, error)

	// ResponseBody retrieves the response body.
	ResponseBody(path, method string, statusCode int) (Body, error)
}
