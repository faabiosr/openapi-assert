package assert

import (
	"github.com/stretchr/testify/suite"
	"net/http"
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

func TestAssertionsTestSuite(t *testing.T) {
	suite.Run(t, new(AssertionsTestSuite))
}
