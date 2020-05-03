package assert

import (
	"bytes"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"
	"testing"

	"gitlab.com/flimzy/testy"
)

func TestAssertRequestMediaType(t *testing.T) {
	type tt struct {
		path      string
		method    string
		mediaType string
		err       string
	}

	tests := testy.NewTable()

	tests.Add("invalid path", tt{
		path:      "/pet",
		method:    http.MethodPost,
		mediaType: "application/json",
		err:       "resource uri does not match",
	})

	tests.Add("invalid type", tt{
		path:      "/api/food",
		method:    http.MethodGet,
		mediaType: "text/html",
		err:       "failed asserting that 'text/html' is an allowed media type (application/json)",
	})

	tests.Add("success", tt{
		path:      "/api/food",
		method:    http.MethodGet,
		mediaType: "application/json",
	})

	tests.Run(t, func(t *testing.T, tt tt) {
		doc, _ := LoadFromURI("./fixtures/docs.json")

		err := RequestMediaType(tt.mediaType, doc, tt.path, tt.method)
		testy.Error(t, tt.err, err)
	})
}

func TestAssertResponseMediaType(t *testing.T) {
	type tt struct {
		path      string
		method    string
		mediaType string
		err       string
	}

	tests := testy.NewTable()

	tests.Add("invalid path", tt{
		path:      "/pet",
		method:    http.MethodPost,
		mediaType: "application/json",
		err:       "resource uri does not match",
	})

	tests.Add("invalid type", tt{
		path:      "/api/food",
		method:    http.MethodGet,
		mediaType: "text/html",
		err:       "failed asserting that 'text/html' is an allowed media type (application/json)",
	})

	tests.Add("success", tt{
		path:      "/api/food",
		method:    http.MethodGet,
		mediaType: "application/json",
	})

	tests.Run(t, func(t *testing.T, tt tt) {
		doc, _ := LoadFromURI("./fixtures/docs.json")

		err := ResponseMediaType(tt.mediaType, doc, tt.path, tt.method)
		testy.Error(t, tt.err, err)
	})
}

func TestAssertRequestHeaders(t *testing.T) {
	type tt struct {
		path    string
		method  string
		headers map[string][]string
		err     string
	}

	tests := testy.NewTable()

	tests.Add("invalid path", tt{
		path:   "/pet",
		method: http.MethodPost,
		err:    "resource uri does not match",
	})

	tests.Add("required values", tt{
		path:   "/api/pets/1",
		method: http.MethodPatch,
		err:    "failed asserting that '{}' is a valid request header (x-required-header is required)",
	})

	tests.Run(t, func(t *testing.T, tt tt) {
		doc, _ := LoadFromURI("./fixtures/docs.json")

		err := RequestHeaders(tt.headers, doc, tt.path, tt.method)
		testy.Error(t, tt.err, err)
	})
}

func TestAssertResponseHeaders(t *testing.T) {
	type tt struct {
		path    string
		method  string
		headers map[string][]string
		status  int
		err     string
	}

	tests := testy.NewTable()

	tests.Add("invalid path", tt{
		path:   "/pet",
		method: http.MethodPost,
		status: http.StatusOK,
		err:    "resource uri does not match",
	})

	tests.Add("required values", tt{
		path:   "/api/pets",
		method: http.MethodGet,
		status: http.StatusOK,
		err:    "failed asserting that '{}' is a valid response header (etag is required)",
	})

	tests.Run(t, func(t *testing.T, tt tt) {
		doc, _ := LoadFromURI("./fixtures/docs.json")

		err := ResponseHeaders(tt.headers, doc, tt.path, tt.method, tt.status)
		testy.Error(t, tt.err, err)
	})
}

func TestAssertRequestQuery(t *testing.T) {
	type tt struct {
		path   string
		method string
		query  url.Values
		err    string
	}

	tests := testy.NewTable()

	tests.Add("invalid path", tt{
		path:   "/pet",
		method: http.MethodPost,
		err:    "resource uri does not match",
	})

	tests.Add("required values", tt{
		path:   "/api/pets",
		method: http.MethodGet,
		err:    "failed asserting that '{}' is a valid request query (limit is required)",
	})

	tests.Run(t, func(t *testing.T, tt tt) {
		doc, _ := LoadFromURI("./fixtures/docs.json")

		err := RequestQuery(url.Values{}, doc, tt.path, tt.method)
		testy.Error(t, tt.err, err)
	})
}

func TestAssertRequestBody(t *testing.T) {
	type tt struct {
		path   string
		method string
		body   string
		err    string
	}

	tests := testy.NewTable()

	tests.Add("invalid path", tt{
		path:   "/pet",
		method: http.MethodPost,
		body:   "{}",
		err:    "resource uri does not match",
	})

	tests.Add("invalid data", tt{
		path:   "/api/pets",
		method: http.MethodPost,
		err:    "unable to load the validation: EOF",
	})

	tests.Add("required values", tt{
		path:   "/api/pets",
		method: http.MethodPost,
		body:   "{}",
		err:    "failed asserting that '{}' is a valid request body (id is required, name is required, id is required, Must validate all the schemas (allOf))",
	})

	tests.Run(t, func(t *testing.T, tt tt) {
		doc, _ := LoadFromURI("./fixtures/docs.json")

		buf := bytes.NewBufferString(tt.body)
		err := RequestBody(buf, doc, tt.path, tt.method)
		testy.Error(t, tt.err, err)
	})
}

func TestAssertResponseBody(t *testing.T) {
	type tt struct {
		path   string
		method string
		status int
		body   string
		err    string
	}

	tests := testy.NewTable()

	tests.Add("invalid path", tt{
		path:   "/pet",
		method: http.MethodPost,
		status: http.StatusOK,
		body:   "{}",
		err:    "resource uri does not match",
	})

	tests.Add("invalid data", tt{
		path:   "/api/pets",
		method: http.MethodGet,
		status: http.StatusOK,
		err:    "unable to load the validation: EOF",
	})

	tests.Add("required values", tt{
		path:   "/api/pets",
		method: http.MethodGet,
		status: http.StatusOK,
		body:   "{}",
		err:    "failed asserting that '{}' is a valid response body (Invalid type. Expected: array, given: object)",
	})

	tests.Run(t, func(t *testing.T, tt tt) {
		doc, _ := LoadFromURI("./fixtures/docs.json")

		buf := bytes.NewBufferString(tt.body)
		err := ResponseBody(buf, doc, tt.path, tt.method, tt.status)
		testy.Error(t, tt.err, err)
	})
}

func TestAssertRequest(t *testing.T) {
	type tt struct {
		path      string
		method    string
		mediaType string
		body      io.Reader
		err       string
	}

	tests := testy.NewTable()

	tests.Add("without required headers", tt{
		path:   "/api/pets/1",
		method: http.MethodPatch,
		err:    `failed asserting that '{"Content-Type":""}' is a valid request header (x-required-header is required)`,
	})

	tests.Add("without media type", tt{
		path:      "/api/food",
		method:    http.MethodGet,
		mediaType: "text/html",
		body:      bytes.NewBufferString("{}"),
		err:       "failed asserting that 'text/html' is an allowed media type (application/json)",
	})

	tests.Add("without query", tt{
		path:   "/api/pets",
		method: http.MethodGet,
		err:    "failed asserting that '{}' is a valid request query (limit is required)",
	})

	tests.Add("without required body", tt{
		path:      "/api/pets",
		method:    http.MethodPost,
		mediaType: "application/json",
		body:      bytes.NewBufferString("{}"),
		err:       "failed asserting that '{}' is a valid request body (id is required, name is required, id is required, Must validate all the schemas (allOf))",
	})

	tests.Add("without body", tt{
		path:      "/api/food",
		method:    http.MethodGet,
		mediaType: "application/json",
		body:      bytes.NewBufferString("{}"),
		err:       "failed asserting that '' is an allowed media type (application/json)",
	})

	tests.Add("read body", tt{
		path:      "/api/pets",
		method:    http.MethodPost,
		mediaType: "application/json",
		body:      bytes.NewBufferString(`{"id": 1, "name": "doggo"}`),
	})

	tests.Run(t, func(t *testing.T, tt tt) {
		doc, _ := LoadFromURI("./fixtures/docs.json")

		req, _ := http.NewRequest(tt.method, tt.path, tt.body)
		req.Header.Add("Content-Type", tt.mediaType)

		err := Request(req, doc)
		if err != nil {
			testy.Error(t, tt.err, err)
		}

		if d := testy.DiffJSON(testy.Snapshot(t), req.Body); d != nil {
			t.Error(d)
		}
	})
}

func TestAssertResponse(t *testing.T) {
	type tt struct {
		path    string
		method  string
		status  int
		headers map[string][]string
		body    io.ReadCloser
		err     string
	}

	tests := testy.NewTable()

	tests.Add("without required headers", tt{
		path:   "/api/pets",
		method: http.MethodGet,
		status: http.StatusOK,
		err:    "failed asserting that '{}' is a valid response header (etag is required)",
	})

	tests.Add("without media type", tt{
		path:   "/api/food",
		method: http.MethodGet,
		status: http.StatusOK,
		headers: map[string][]string{
			"Content-Type": {"text/html"},
		},
		body: ioutil.NopCloser(bytes.NewBufferString("{}")),
		err:  "failed asserting that 'text/html' is an allowed media type (application/json)",
	})

	tests.Add("without required body", tt{
		path:   "/api/pets",
		method: http.MethodGet,
		status: http.StatusOK,
		headers: map[string][]string{
			"Content-Type": {"text/html"},
			"etag":         {"value"},
		},
		body: ioutil.NopCloser(bytes.NewBufferString("{}")),
		err:  "failed asserting that '{}' is a valid response body (Invalid type. Expected: array, given: object)",
	})

	tests.Add("read body", tt{
		path:   "/api/pets",
		method: http.MethodGet,
		status: http.StatusOK,
		headers: map[string][]string{
			"Content-Type": {"application/json"},
			"etag":         {"value"},
		},
		body: ioutil.NopCloser(bytes.NewBufferString(`[{"id": 1, "name": "doggo"}]`)),
	})

	tests.Run(t, func(t *testing.T, tt tt) {
		req, _ := http.NewRequest(tt.method, tt.path, nil)
		res := &http.Response{
			StatusCode: tt.status,
			Request:    req,
			Header:     tt.headers,
			Body:       tt.body,
		}

		doc, _ := LoadFromURI("./fixtures/docs.json")
		err := Response(res, doc)
		if err != nil {
			testy.Error(t, tt.err, err)
		}

		if d := testy.DiffJSON(testy.Snapshot(t), res.Body); d != nil {
			t.Error(d)
		}
	})
}
