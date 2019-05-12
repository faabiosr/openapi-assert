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
	AssertionsTestSuite struct {
		BaseTestSuite

		assertions *Assertions
	}
)

func (s *AssertionsTestSuite) SetupTest() {
	s.BaseTestSuite.SetupTest()
	s.assertions = New(s.doc)
}

func (s *AssertionsTestSuite) TestRequestMediaType() {
	err := s.assertions.RequestMediaType("application/json", "/api/food", http.MethodGet)

	s.assert.Nil(err)
}

func (s *AssertionsTestSuite) TestResponseMediaType() {
	err := s.assertions.ResponseMediaType("application/json", "/api/food", http.MethodGet)

	s.assert.Nil(err)
}

func (s *AssertionsTestSuite) TestRequestHeaders() {
	headers := map[string][]string{
		"x-required-header": {"value"},
	}

	err := s.assertions.RequestHeaders(headers, "/api/pets/1", http.MethodPatch)

	s.assert.Nil(err)
}

func (s *AssertionsTestSuite) TestResponseHeaders() {
	headers := map[string][]string{
		"etag": {"value"},
	}

	err := s.assertions.ResponseHeaders(headers, "/api/pets", http.MethodGet, http.StatusOK)

	s.assert.Nil(err)
}

func (s *AssertionsTestSuite) TestRequestQuery() {
	query := url.Values{}
	query.Add("tags", "foo")
	query.Add("tags", "bar")
	query.Add("limit", "1")

	err := s.assertions.RequestQuery(query, "/api/pets", http.MethodGet)

	s.assert.Nil(err)
}

func (s *AssertionsTestSuite) TestRequestBody() {
	buf := bytes.NewBufferString(`{"id": 1, "name": "doggo"}`)

	err := s.assertions.RequestBody(buf, "/api/pets", http.MethodPost)

	s.assert.Nil(err)
}

func (s *AssertionsTestSuite) TestResponseBody() {
	buf := bytes.NewBufferString(`[{"id": 1, "name": "doggo"}]`)

	err := s.assertions.ResponseBody(buf, "/api/pets", http.MethodGet, http.StatusOK)

	s.assert.Nil(err)
}

func (s *AssertionsTestSuite) TestRequest() {
	buf := bytes.NewBufferString(`{"id": 1, "name": "doggo"}`)

	req, _ := http.NewRequest(http.MethodPost, "/api/pets", buf)
	req.Header.Add("Content-Type", "application/json")

	s.assert.Nil(s.assertions.Request(req))
}

func (s *AssertionsTestSuite) TestResponse() {
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

	s.assert.Nil(s.assertions.Response(res))
}

func TestAssertionsTestSuite(t *testing.T) {
	suite.Run(t, new(AssertionsTestSuite))
}
