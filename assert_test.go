package assert

import (
	"bytes"
	"github.com/stretchr/testify/suite"
	"io/ioutil"
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

func (s *AssertTestSuite) TestResponseBodyWithInvalidPath() {
	buf := bytes.NewBufferString("{}")
	err := ResponseBody(buf, s.doc, "/pet", http.MethodPost, http.StatusOK)

	s.assert.Error(err)
	s.assert.Contains(err.Error(), ErrResourceURI)
}

func (s *AssertTestSuite) TestResponseBodyWithInvalidData() {
	buf := bytes.NewBufferString("")
	err := ResponseBody(buf, s.doc, "/api/pets", http.MethodGet, http.StatusOK)

	s.assert.Error(err)
	s.assert.Contains(err.Error(), ErrValidation)
}

func (s *AssertTestSuite) TestResponseBodyWithoutRequiredValues() {
	buf := bytes.NewBufferString("{}")
	err := ResponseBody(buf, s.doc, "/api/pets", http.MethodGet, http.StatusOK)

	s.assert.Error(err)
	s.assert.EqualError(err, failf(ErrResponseBody, "{}", "Invalid type. Expected: array, given: object").Error())
}

func (s *AssertTestSuite) TestRequestWithoutRequiredHeaders() {
	req, _ := http.NewRequest(http.MethodPatch, "/api/pets/1", nil)

	s.assert.Error(Request(req, s.doc))
}

func (s *AssertTestSuite) TestRequestWithoutRequiredMediaType() {
	buf := bytes.NewBufferString("{}")

	req, _ := http.NewRequest(http.MethodGet, "/api/food", buf)
	req.Header.Add("Content-Type", "text/html")

	s.assert.Error(Request(req, s.doc))
}

func (s *AssertTestSuite) TestRequestWithoutRequiredQuery() {
	req, _ := http.NewRequest(http.MethodGet, "/api/pets", nil)

	s.assert.Error(Request(req, s.doc))
}

func (s *AssertTestSuite) TestRequestWithoutRequiredBody() {
	req, _ := http.NewRequest(http.MethodPost, "/api/pets", bytes.NewBufferString("{}"))
	req.Header.Add("Content-Type", "application/json")

	s.assert.Error(Request(req, s.doc))
}

func (s *AssertTestSuite) TestRequestReadBodyAfterValidation() {
	buf := bytes.NewBufferString(`{"id": 1, "name": "doggo"}`)

	req, _ := http.NewRequest(http.MethodPost, "/api/pets", buf)
	req.Header.Add("Content-Type", "application/json")

	s.assert.Nil(Request(req, s.doc))

	body, err := ioutil.ReadAll(req.Body)

	s.assert.NotEmpty(body)
	s.assert.NoError(err)
}

func (s *AssertTestSuite) TestResponseWithoutRequiredHeaders() {
	req, _ := http.NewRequest(http.MethodGet, "/api/pets", nil)

	res := &http.Response{
		StatusCode: http.StatusOK,
		Request:    req,
	}

	s.assert.Error(Response(res, s.doc))
}

func (s *AssertTestSuite) TestResponseWithoutRequiredMediaType() {
	req, _ := http.NewRequest(http.MethodGet, "/api/food", nil)

	buf := bytes.NewBufferString("{}")

	res := &http.Response{
		Request: req,
		Header: map[string][]string{
			"Content-Type": {"text/html"},
		},
		Body: ioutil.NopCloser(buf),
	}

	s.assert.Error(Response(res, s.doc))
}

func (s *AssertTestSuite) TestResponseWithoutRequiredBody() {
	req, _ := http.NewRequest(http.MethodGet, "/api/pets", nil)

	buf := bytes.NewBufferString("{}")

	res := &http.Response{
		StatusCode: http.StatusOK,
		Request:    req,
		Header: map[string][]string{
			"Content-Type": {"application/json"},
			"etag":         {"value"},
		},
		Body: ioutil.NopCloser(buf),
	}

	s.assert.Error(Response(res, s.doc))
}

func (s *AssertTestSuite) TestResponseReadBodyAfterValidation() {
	req, _ := http.NewRequest(http.MethodGet, "/api/pets", nil)

	buf := bytes.NewBufferString(`[{"id": 1, "name": "doggo"}]`)

	res := &http.Response{
		StatusCode: http.StatusOK,
		Request:    req,
		Header: map[string][]string{
			"Content-Type": {"application/json"},
			"etag":         {"value"},
		},
		Body: ioutil.NopCloser(buf),
	}

	s.assert.NoError(Response(res, s.doc))

	body, err := ioutil.ReadAll(res.Body)

	s.assert.NotEmpty(body)
	s.assert.NoError(err)
}

func TestAssertTestSuite(t *testing.T) {
	suite.Run(t, new(AssertTestSuite))
}
