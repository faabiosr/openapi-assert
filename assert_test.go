package assert

import (
	"bytes"
	"github.com/stretchr/testify/suite"
	"net/http"
	"net/url"
	"testing"
)

type (
	AssertTestSuite struct {
		BaseTestSuite
	}
)

func (s *AssertTestSuite) TestRequestMediaTypeWithInvalidPath() {
	err := RequestMediaType("application/json", s.doc, "/pet", http.MethodPost)

	s.assert.Error(err)
	s.assert.Contains(err.Error(), ErrResourceURI)
}

func (s *AssertTestSuite) TestRequestMediaTypeWithInvalidType() {
	err := RequestMediaType("text/html", s.doc, "/api/food", http.MethodGet)

	s.assert.Error(err)
	s.assert.EqualError(err, failf(ErrMediaType, "text/html", "application/json").Error())
}

func (s *AssertTestSuite) TestRequestMediaTypeWithValidType() {
	err := RequestMediaType("application/json", s.doc, "/api/food", http.MethodGet)

	s.assert.Nil(err)
}

func (s *AssertTestSuite) TestResponseMediaTypeWithInvalidPath() {
	err := ResponseMediaType("application/json", s.doc, "/pet", http.MethodPost)

	s.assert.Error(err)
	s.assert.Contains(err.Error(), ErrResourceURI)
}

func (s *AssertTestSuite) TestResponseMediaTypeWithInvalidType() {
	err := ResponseMediaType("text/html", s.doc, "/api/food", http.MethodGet)

	s.assert.Error(err)
	s.assert.EqualError(err, failf(ErrMediaType, "text/html", "application/json").Error())
}

func (s *AssertTestSuite) TestResponseMediaTypeWithValidType() {
	err := ResponseMediaType("application/json", s.doc, "/api/food", http.MethodGet)

	s.assert.Nil(err)
}

func (s *AssertTestSuite) TestRequestHeadersWithInvalidPath() {
	headers := map[string][]string{}

	err := RequestHeaders(headers, s.doc, "/pet", http.MethodPost)

	s.assert.Error(err)
	s.assert.Contains(err.Error(), ErrResourceURI)
}

func (s *AssertTestSuite) TestRequestHeaderWithoutRequiredValues() {
	headers := map[string][]string{}

	err := RequestHeaders(headers, s.doc, "/api/pets/1", http.MethodPatch)

	s.assert.Error(err)
	s.assert.EqualError(err, failf(ErrRequestHeaders, "{}", "x-required-header is required").Error())
}

func (s *AssertTestSuite) TestResponseHeadersWithInvalidPath() {
	headers := map[string][]string{}

	err := ResponseHeaders(headers, s.doc, "/pet", http.MethodPost, http.StatusOK)

	s.assert.Error(err)
	s.assert.Contains(err.Error(), ErrResourceURI)
}

func (s *AssertTestSuite) TestResponseHeaderWithoutRequiredValues() {
	headers := map[string][]string{}

	err := ResponseHeaders(headers, s.doc, "/api/pets", http.MethodGet, http.StatusOK)

	s.assert.Error(err)
	s.assert.EqualError(err, failf(ErrResponseHeaders, "{}", "etag is required").Error())
}

func (s *AssertTestSuite) TestRequestQueryWithInvalidPath() {
	query := url.Values{}

	err := RequestQuery(query, s.doc, "/pet", http.MethodPost)

	s.assert.Error(err)
	s.assert.Contains(err.Error(), ErrResourceURI)
}

func (s *AssertTestSuite) TestRequestQueryWithoutRequiredValues() {
	query := url.Values{}

	err := RequestQuery(query, s.doc, "/api/pets", http.MethodGet)

	s.assert.Error(err)
	s.assert.EqualError(err, failf(ErrRequestQuery, "{}", "limit is required").Error())
}

func (s *AssertTestSuite) TestRequestBodyWithInvalidPath() {
	buf := bytes.NewBufferString("{}")
	err := RequestBody(buf, s.doc, "/pet", http.MethodPost)

	s.assert.Error(err)
	s.assert.Contains(err.Error(), ErrResourceURI)
}

func (s *AssertTestSuite) TestRequestBodyWithInvalidData() {
	buf := bytes.NewBufferString("")
	err := RequestBody(buf, s.doc, "/api/pets", http.MethodPost)

	s.assert.Error(err)
	s.assert.Contains(err.Error(), ErrValidation)
}

func (s *AssertTestSuite) TestRequestBodyWithoutRequiredValues() {
	buf := bytes.NewBufferString("{}")
	err := RequestBody(buf, s.doc, "/api/pets", http.MethodPost)

	s.assert.Error(err)
	s.assert.EqualError(err, failf(ErrRequestBody, "{}", "id is required, name is required, id is required, Must validate all the schemas (allOf)").Error())
}

func TestAssertTestSuite(t *testing.T) {
	suite.Run(t, new(AssertTestSuite))
}
