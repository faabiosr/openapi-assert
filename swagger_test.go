package assert

import (
	"io"
	"net/http"
	"os"
	"strings"
	"testing"

	"gitlab.com/flimzy/testy"
)

func TestLoadFromURI(t *testing.T) {
	type tt struct {
		uri string
		err string
	}

	tests := testy.NewTable()

	tests.Add("empty param", tt{
		uri: "",
		err: "unable to load the document by uri: open : no such file or directory",
	})

	tests.Add("invalid file", tt{
		uri: "./fixtures/invalid-doc.json",
		err: `unable to expand the document: object has no key "definitions"`,
	})

	tests.Add("success", tt{
		uri: "./fixtures/docs.json",
	})

	tests.Run(t, func(t *testing.T, tt tt) {
		_, err := LoadFromURI(tt.uri)
		testy.Error(t, tt.err, err)
	})
}

func TestLoadFromReader(t *testing.T) {
	type tt struct {
		reader io.Reader
		err    string
	}

	tests := testy.NewTable()

	tests.Add("invalid content", tt{
		reader: strings.NewReader("{"),
		err:    "unable to load the document by uri: unexpected end of JSON input",
	})

	tests.Add("invalid file", func() interface{} {
		f, _ := os.Open("./fixtures/invalid-doc.json")

		return tt{
			reader: f,
			err:    `unable to expand the document: object has no key "ErrorModel"`,
		}
	})

	tests.Add("success", func() interface{} {
		f, _ := os.Open("./fixtures/docs.json")

		return tt{
			reader: f,
		}
	})

	tests.Run(t, func(t *testing.T, tt tt) {
		_, err := LoadFromReader(tt.reader)
		testy.Error(t, tt.err, err)
	})
}

func TestFindPath(t *testing.T) {
	doc, _ := LoadFromURI("./fixtures/invalid-path.json")

	_, err := doc.findPath("/api/food/a")
	testy.Error(t, "resource uri does not match: uritemplate:11:invalid varname", err)
}

func TestRequestMediaTypes(t *testing.T) {
	type tt struct {
		path   string
		method string
		err    string
	}

	tests := testy.NewTable()

	tests.Add("invalid path", tt{
		path:   "/some",
		method: http.MethodPost,
		err:    "resource uri does not match",
	})

	tests.Add("default types", tt{
		path:   "/api/pets/1",
		method: http.MethodDelete,
	})

	tests.Add("success", tt{
		path:   "/api/pets",
		method: http.MethodGet,
	})

	tests.Run(t, func(t *testing.T, tt tt) {
		doc, _ := LoadFromURI("./fixtures/docs.json")

		got, err := doc.RequestMediaTypes(tt.path, tt.method)
		if err != nil {
			testy.Error(t, tt.err, err)
		}

		if d := testy.DiffInterface(testy.Snapshot(t), got); d != nil {
			t.Error(d)
		}
	})
}

func TestResponseMediaTypes(t *testing.T) {
	type tt struct {
		path   string
		method string
		err    string
	}

	tests := testy.NewTable()

	tests.Add("invalid path", tt{
		path:   "/some",
		method: http.MethodPost,
		err:    "resource uri does not match",
	})

	tests.Add("default types", tt{
		path:   "/api/pets/1",
		method: http.MethodDelete,
	})

	tests.Add("success", tt{
		path:   "/api/pets/1",
		method: http.MethodPatch,
	})

	tests.Run(t, func(t *testing.T, tt tt) {
		doc, _ := LoadFromURI("./fixtures/docs.json")

		got, err := doc.ResponseMediaTypes(tt.path, tt.method)
		if err != nil {
			testy.Error(t, tt.err, err)
		}

		if d := testy.DiffInterface(testy.Snapshot(t), got); d != nil {
			t.Error(d)
		}
	})
}

func TestRequestHeaders(t *testing.T) {
	type tt struct {
		path   string
		method string
		err    string
	}

	tests := testy.NewTable()

	tests.Add("invalid path", tt{
		path:   "/some",
		method: http.MethodPost,
		err:    "resource uri does not match",
	})

	tests.Add("no headers", tt{
		path:   "/api/food",
		method: http.MethodGet,
	})

	tests.Add("success", tt{
		path:   "/api/pets/1",
		method: http.MethodPatch,
	})

	tests.Run(t, func(t *testing.T, tt tt) {
		doc, _ := LoadFromURI("./fixtures/docs.json")

		got, err := doc.RequestHeaders(tt.path, tt.method)
		if err != nil {
			testy.Error(t, tt.err, err)
		}

		if d := testy.DiffInterface(testy.Snapshot(t), got); d != nil {
			t.Error(d)
		}
	})
}

func TestResponseHeaders(t *testing.T) {
	type tt struct {
		path   string
		method string
		status int
		err    string
	}

	tests := testy.NewTable()

	tests.Add("invalid path", tt{
		path:   "/some",
		method: http.MethodPost,
		status: http.StatusOK,
		err:    "resource uri does not match",
	})

	tests.Add("default", tt{
		path:   "/api/pets",
		method: http.MethodGet,
		status: http.StatusBadRequest,
	})

	tests.Add("success", tt{
		path:   "/api/pets",
		method: http.MethodGet,
		status: http.StatusOK,
	})

	tests.Run(t, func(t *testing.T, tt tt) {
		doc, _ := LoadFromURI("./fixtures/docs.json")

		got, err := doc.ResponseHeaders(tt.path, tt.method, tt.status)
		if err != nil {
			testy.Error(t, tt.err, err)
		}

		if d := testy.DiffInterface(testy.Snapshot(t), got); d != nil {
			t.Error(d)
		}
	})
}

func TestRequestQuery(t *testing.T) {
	type tt struct {
		path   string
		method string
		err    string
	}

	tests := testy.NewTable()

	tests.Add("invalid path", tt{
		path:   "/some",
		method: http.MethodPost,
		err:    "resource uri does not match",
	})

	tests.Add("no query", tt{
		path:   "/api/food",
		method: http.MethodGet,
	})

	tests.Add("success", tt{
		path:   "/api/pets",
		method: http.MethodGet,
	})

	tests.Run(t, func(t *testing.T, tt tt) {
		doc, _ := LoadFromURI("./fixtures/docs.json")

		got, err := doc.RequestQuery(tt.path, tt.method)
		if err != nil {
			testy.Error(t, tt.err, err)
		}

		if d := testy.DiffInterface(testy.Snapshot(t), got); d != nil {
			t.Error(d)
		}
	})
}

func TestRequestBody(t *testing.T) {
	type tt struct {
		path   string
		method string
		err    string
	}

	tests := testy.NewTable()

	tests.Add("invalid path", tt{
		path:   "/some",
		method: http.MethodPost,
		err:    "resource uri does not match",
	})

	tests.Add("not exists", tt{
		path:   "/api/pets",
		method: http.MethodGet,
		err:    "body does not exists",
	})

	tests.Add("success", tt{
		path:   "/api/pets",
		method: http.MethodPost,
	})

	tests.Run(t, func(t *testing.T, tt tt) {
		doc, _ := LoadFromURI("./fixtures/docs.json")

		got, err := doc.RequestBody(tt.path, tt.method)
		if err != nil {
			testy.Error(t, tt.err, err)
		}

		if d := testy.DiffInterface(testy.Snapshot(t), got); d != nil {
			t.Error(d)
		}
	})
}

func TestResponseBody(t *testing.T) {
	type tt struct {
		path   string
		method string
		status int
		err    string
	}

	tests := testy.NewTable()

	tests.Add("invalid path", tt{
		path:   "/some",
		method: http.MethodPost,
		status: http.StatusOK,
		err:    "resource uri does not match",
	})

	tests.Add("not exists", tt{
		path:   "/api/food",
		method: http.MethodGet,
		status: http.StatusNotModified,
		err:    "body does not exists",
	})

	tests.Add("success", tt{
		path:   "/api/pets",
		method: http.MethodPost,
		status: http.StatusOK,
	})

	tests.Run(t, func(t *testing.T, tt tt) {
		doc, _ := LoadFromURI("./fixtures/docs.json")

		got, err := doc.ResponseBody(tt.path, tt.method, tt.status)
		if err != nil {
			testy.Error(t, tt.err, err)
		}

		if d := testy.DiffInterface(testy.Snapshot(t), got); d != nil {
			t.Error(d)
		}
	})
}
