package assert

type (
	// Document that defines the contract for reading OpenAPI documents.
	Document interface {
		// FindPath searches for an uri in document and returns the path.
		FindPath(uri string) (string, error)

		// FindNode searches a node by using the segments in the document.
		FindNode(segments ...string) (interface{}, error)
	}
)
