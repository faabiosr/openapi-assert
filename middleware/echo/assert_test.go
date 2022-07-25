package echo

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	ec "github.com/labstack/echo/v4"
	"gitlab.com/flimzy/testy"

	oapi "github.com/faabiosr/openapi-assert"
)

func TestMiddlewareWithConfig(t *testing.T) {
	type tt struct {
		cfg AssertConfig
		err string
	}

	doc, _ := oapi.LoadFromURI("../../fixtures/docs.json")

	tests := testy.NewTable()

	tests.Add("with config", tt{
		cfg: AssertConfig{Document: doc},
		err: `code=400, message=failed asserting that '{"Content-Type":"application/json"}' is a valid request header (x-required-header is required), internal=failed asserting that '{"Content-Type":"application/json"}' is a valid request header (x-required-header is required)`,
	})

	tests.Add("with skipper", tt{
		cfg: AssertConfig{
			Document: doc,
			Skipper: func(ctx ec.Context) bool {
				return true
			},
		},
	})

	tests.Add("without document", tt{
		cfg: AssertConfig{},
		err: "echo: assert middleware requires an openapi-assert document",
	})

	tests.Run(t, func(t *testing.T, tt tt) {
		defer func() {
			r := recover()
			if r != nil && r != tt.err {
				t.Errorf("want %v, got %v", tt.err, r)
			}
		}()

		req := httptest.NewRequest(http.MethodPatch, "/api/pets/1", nil)
		req.Header.Add("Content-Type", "application/json")

		rec := httptest.NewRecorder()
		ctx := ec.New().NewContext(req, rec)

		err := AssertWithConfig(tt.cfg)(func(ctx ec.Context) error {
			return ctx.String(http.StatusOK, "test")
		})(ctx)

		testy.Error(t, tt.err, err)
	})
}

func TestMiddleware(t *testing.T) {
	reader := strings.NewReader(`{"id": 1, "name": "doggo"}`)
	req := httptest.NewRequest(http.MethodPost, "/api/pets", reader)
	req.Header.Add("Content-Type", "application/json")
	rec := httptest.NewRecorder()

	c := ec.New().NewContext(req, rec)
	doc, _ := oapi.LoadFromURI("../../fixtures/docs.json")

	err := Assert(doc)(func(ctx ec.Context) error {
		return ctx.String(http.StatusOK, "test")
	})(c)
	if err != nil {
		t.Error(err)
	}
}
