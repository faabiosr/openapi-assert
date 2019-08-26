package echo

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	oapi "github.com/faabiosr/openapi-assert"
	ec "github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

type (
	AssertTestSuite struct {
		suite.Suite
		assert *assert.Assertions

		doc oapi.Document

		srv *ec.Echo
	}
)

func (s *AssertTestSuite) SetupTest() {
	s.assert = assert.New(s.T())
	s.doc, _ = oapi.LoadFromURI("../../fixtures/docs.json")
	s.srv = ec.New()
}

func (s *AssertTestSuite) TestMiddlewareWithConfig() {
	req := httptest.NewRequest(ec.PATCH, "/api/pets/1", nil)
	req.Header.Add(ec.HeaderContentType, ec.MIMEApplicationJSON)

	rec := httptest.NewRecorder()

	c := s.srv.NewContext(req, rec)

	cfg := AssertConfig{Document: s.doc}

	err := AssertWithConfig(cfg)(func(ctx ec.Context) error {
		return ctx.String(http.StatusOK, "test")
	})(c)

	s.assert.Error(err)
}

func (s *AssertTestSuite) TestMiddleware() {
	req := httptest.NewRequest(
		ec.POST,
		"/api/pets",
		strings.NewReader(`{"id": 1, "name": "doggo"}`),
	)

	req.Header.Add(ec.HeaderContentType, ec.MIMEApplicationJSON)

	rec := httptest.NewRecorder()

	c := s.srv.NewContext(req, rec)

	err := Assert(s.doc)(func(ctx ec.Context) error {
		return ctx.String(http.StatusOK, "test")
	})(c)

	s.assert.NoError(err)
}

func (s *AssertTestSuite) TestMiddlewareWithSkipper() {
	req := httptest.NewRequest(ec.PATCH, "/api/pets/1", nil)
	req.Header.Add(ec.HeaderContentType, ec.MIMEApplicationJSON)

	rec := httptest.NewRecorder()

	c := s.srv.NewContext(req, rec)

	cfg := AssertConfig{
		Document: s.doc,
		Skipper: func(c ec.Context) bool {
			return true
		},
	}

	err := AssertWithConfig(cfg)(func(ctx ec.Context) error {
		return ctx.String(http.StatusOK, "test")
	})(c)

	s.assert.NoError(err)
}

func (s *AssertTestSuite) TestMiddlewareWithoutDocument() {
	req := httptest.NewRequest(ec.PATCH, "/api/pets/1", nil)
	rec := httptest.NewRecorder()

	c := s.srv.NewContext(req, rec)

	cfg := AssertConfig{
		Skipper: func(c ec.Context) bool {
			return true
		},
	}

	caller := func() {
		AssertWithConfig(cfg)(func(ctx ec.Context) error {
			return ctx.String(http.StatusOK, "test")
		})(c)
	}

	s.assert.Panics(caller)
}

func TestAssertTestSuite(t *testing.T) {
	suite.Run(t, new(AssertTestSuite))
}
