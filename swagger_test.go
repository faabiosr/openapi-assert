package assert

import (
	"github.com/stretchr/testify/suite"
	"net/http"
	"testing"
)

type (
	SwaggerTestSuite struct {
		BaseTestSuite
	}
)

func (s *SwaggerTestSuite) TestLoadFromUriWithEmptyParam() {
	doc, err := LoadFromURI("")

	s.assert.Nil(doc)
	s.assert.Error(err)
	s.assert.Contains(err.Error(), ErrSwaggerLoad)
}

func (s *SwaggerTestSuite) TestLoadFromUriWithInvalidFile() {
	doc, err := LoadFromURI("./fixtures/invalid-doc.json")

	s.assert.Nil(doc)
	s.assert.Error(err)
	s.assert.Contains(err.Error(), ErrSwaggerExpand)
}

func (s *SwaggerTestSuite) TestLoadFromUri() {
	doc, err := LoadFromURI(s.filePath)

	s.assert.IsType(&Swagger{}, doc)
	s.assert.Implements(new(Document), doc)
	s.assert.NoError(err)
}

func (s *SwaggerTestSuite) TestFindPathWithBrokenDocument() {
	doc, _ := LoadFromURI("./fixtures/invalid-path.json")

	path, err := doc.findPath("/api/food/a")

	s.assert.Empty(path)
	s.assert.Error(err)
	s.assert.Contains(err.Error(), ErrResourceURI)
}

func (s *SwaggerTestSuite) TestRequestMediaTypesWithInvalidPath() {
	types, err := s.doc.RequestMediaTypes("/some", http.MethodPost)

	s.assert.Len(types, 0)
	s.assert.Error(err)
}

func (s *SwaggerTestSuite) TestRequestMediaTypes() {
	types, err := s.doc.RequestMediaTypes("/api/pets", http.MethodGet)

	s.assert.Len(types, 4)
	s.assert.Nil(err)
}

func (s *SwaggerTestSuite) TestRequestMediaTypesReturnsDefaultTypes() {
	types, err := s.doc.RequestMediaTypes("/api/pets/1", http.MethodDelete)

	s.assert.Len(types, 1)
	s.assert.Nil(err)
}

func (s *SwaggerTestSuite) TestResponseMediaTypesWithInvalidPath() {
	types, err := s.doc.ResponseMediaTypes("/some", http.MethodPost)

	s.assert.Len(types, 0)
	s.assert.Error(err)
}

func (s *SwaggerTestSuite) TestResponseMediaTypes() {
	types, err := s.doc.ResponseMediaTypes("/api/pets/1", http.MethodPatch)

	s.assert.Len(types, 2)
	s.assert.Nil(err)
}

func (s *SwaggerTestSuite) TestResponseMediaTypesReturnsDefaultTypes() {
	types, err := s.doc.ResponseMediaTypes("/api/pets/1", http.MethodDelete)

	s.assert.Len(types, 1)
	s.assert.Nil(err)
}

func (s *SwaggerTestSuite) TestRequestHeadersWithInvalidPath() {
	headers, err := s.doc.RequestHeaders("/some", http.MethodPost)

	s.assert.Len(headers, 0)
	s.assert.Error(err)
}

func (s *SwaggerTestSuite) TestRequestHeaders() {
	headers, err := s.doc.RequestHeaders("/api/pets/1", http.MethodPatch)

	s.assert.Len(headers, 3)
	s.assert.Contains(headers, "x-required-header")
	s.assert.NoError(err)
}

func (s *SwaggerTestSuite) TestRequestHeadersRetrievesNoHeaders() {
	headers, err := s.doc.RequestHeaders("/api/food", http.MethodGet)

	s.assert.Len(headers, 0)
	s.assert.NoError(err)
}

func (s *SwaggerTestSuite) TestResponseHeadersWithInvalidPath() {
	headers, err := s.doc.ResponseHeaders("/some", http.MethodPost, http.StatusOK)

	s.assert.Len(headers, 0)
	s.assert.Error(err)
}

func (s *SwaggerTestSuite) TestResponseHeaders() {
	headers, err := s.doc.ResponseHeaders("/api/pets", http.MethodGet, http.StatusOK)

	s.assert.Len(headers, 2)
	s.assert.Contains(headers, "etag")
	s.assert.NoError(err)
}

func (s *SwaggerTestSuite) TestResponseHeadersDefault() {
	headers, err := s.doc.ResponseHeaders("/api/pets", http.MethodGet, http.StatusBadRequest)

	s.assert.Len(headers, 0)
	s.assert.NoError(err)
}

func (s *SwaggerTestSuite) TestRequestQueryWithInvalidPath() {
	query, err := s.doc.RequestQuery("/some", http.MethodPost)

	s.assert.Len(query, 0)
	s.assert.Error(err)
}

func (s *SwaggerTestSuite) TestRequestQuery() {
	query, err := s.doc.RequestQuery("/api/pets", http.MethodGet)

	s.assert.Len(query, 3)
	s.assert.Contains(query, "limit")
	s.assert.Contains(query, "tags")
	s.assert.NoError(err)
}

func (s *SwaggerTestSuite) TestRequestQueryRetrievesNoQuery() {
	query, err := s.doc.RequestQuery("/api/food", http.MethodGet)

	s.assert.Len(query, 0)
	s.assert.NoError(err)
}

func (s *SwaggerTestSuite) TestRequestBodyWithInvalidPath() {
	body, err := s.doc.RequestBody("/some", http.MethodPost)

	s.assert.Nil(body)
	s.assert.Error(err)
}

func (s *SwaggerTestSuite) TestRequestBodyWhenBodyNotExists() {
	body, err := s.doc.RequestBody("/api/pets", http.MethodGet)

	s.assert.Nil(body)
	s.assert.Error(err)
	s.assert.Contains(err.Error(), ErrBodyNotFound)
}

func (s *SwaggerTestSuite) TestRequestBody() {
	body, err := s.doc.RequestBody("/api/pets", http.MethodPost)

	s.assert.NotEmpty(body)
	s.assert.NoError(err)
}

func TestSwaggerTestSuite(t *testing.T) {
	suite.Run(t, new(SwaggerTestSuite))
}
