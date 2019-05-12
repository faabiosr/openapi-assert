package assert

type (
	// Headers is a list of headers in json schema format.
	Headers map[string]interface{}

	// Query is a list of query parameters in json schema format.
	Query map[string]interface{}

	// Param is a document parameter in json schema format.
	Param struct {
		Type        string
		Description string
		In          string
	}

	// Required is a list of required parameters.
	Required []string

	// Document that defines the contract for reading OpenAPI documents.
	Document interface {
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
	}
)
