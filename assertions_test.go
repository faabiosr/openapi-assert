package assert

import (
	"bytes"
	"io/ioutil"
	"net/http"
	"net/url"
	"testing"
)

func TestAssertions(t *testing.T) {
	doc, _ := LoadFromURI("./fixtures/docs.json")
	assertions := New(doc)

	t.Run("request media type", func(t *testing.T) {
		err := assertions.RequestMediaType("application/json", "/api/food", http.MethodGet)
		if err != nil {
			t.Error(err)
		}
	})

	t.Run("response media type", func(t *testing.T) {
		err := assertions.ResponseMediaType("application/json", "/api/food", http.MethodGet)
		if err != nil {
			t.Error(err)
		}
	})

	t.Run("request headers", func(t *testing.T) {
		headers := map[string][]string{
			"x-required-header": {"value"},
		}

		err := assertions.RequestHeaders(headers, "/api/pets/1", http.MethodPatch)
		if err != nil {
			t.Error(err)
		}
	})

	t.Run("request headers", func(t *testing.T) {
		headers := map[string][]string{
			"etag": {"value"},
		}

		err := assertions.ResponseHeaders(headers, "/api/pets", http.MethodGet, http.StatusOK)
		if err != nil {
			t.Error(err)
		}
	})

	t.Run("request query", func(t *testing.T) {
		query := url.Values{}
		query.Add("tags", "foo")
		query.Add("tags", "bar")
		query.Add("limit", "1")

		err := assertions.RequestQuery(query, "/api/pets", http.MethodGet)
		if err != nil {
			t.Error(err)
		}
	})

	t.Run("request body", func(t *testing.T) {
		buf := bytes.NewBufferString(`{"id": 1, "name": "doggo"}`)

		err := assertions.RequestBody(buf, "/api/pets", http.MethodPost)
		if err != nil {
			t.Error(err)
		}
	})

	t.Run("response body", func(t *testing.T) {
		buf := bytes.NewBufferString(`[{"id": 1, "name": "doggo"}]`)

		err := assertions.ResponseBody(buf, "/api/pets", http.MethodGet, http.StatusOK)
		if err != nil {
			t.Error(err)
		}
	})

	t.Run("request", func(t *testing.T) {
		buf := bytes.NewBufferString(`{"id": 1, "name": "doggo"}`)

		req, _ := http.NewRequest(http.MethodPost, "/api/pets", buf)
		req.Header.Add("Content-Type", "application/json")

		err := assertions.Request(req)
		if err != nil {
			t.Error(err)
		}
	})

	t.Run("response", func(t *testing.T) {
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

		err := assertions.Response(res)
		if err != nil {
			t.Error(err)
		}
	})
}
