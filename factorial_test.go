package main

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/kinbiko/jsonassert"

	"github.com/julienschmidt/httprouter"
)

func TestServer(t *testing.T) {
	testCases := []struct {
		name string
		body []byte
		want string
	}{
		{
			name: "float_a",
			body: []byte(`{"a": 5.5, "b": 4}`),
			want: `{"error": "Incorrect message"}`,
		},
		{
			name: "string_a",
			body: []byte(`{"a": "somth", "b": 4}`),
			want: `{"error": "Incorrect message"}`,
		},
		{
			name: "object_a",
			body: []byte(`{"a": {}, "b": 4}`),
			want: `{"error": "Incorrect message"}`,
		},
		{
			name: "array_a",
			body: []byte(`{"a": [], "b": 4}`),
			want: `{"error": "Incorrect message"}`,
		},
		{
			name: "null_a",
			body: []byte(`{"a": null, "b": 4}`),
			want: `{"error": "Incorrect message"}`,
		},
		{
			name: "float_b",
			body: []byte(`{"a": 4, "b": 5.5}`),
			want: `{"error": "Incorrect message"}`,
		},
		{
			name: "string_b",
			body: []byte(`{"a": 4, "b": "somth"}`),
			want: `{"error": "Incorrect message"}`,
		},
		{
			name: "object_b",
			body: []byte(`{"a": 4, "b": {}}`),
			want: `{"error": "Incorrect message"}`,
		},
		{
			name: "array_b",
			body: []byte(`{"a": 4, "b": []}`),
			want: `{"error": "Incorrect message"}`,
		},
		{
			name: "null_b",
			body: []byte(`{"a": 4, "b": null}`),
			want: `{"error": "Incorrect message"}`,
		},
		{
			name: "minus-five_minus-four",
			body: []byte(`{"a": -5, "b": -4}`),
			want: `{"error": "Incorrect message"}`,
		},
		{
			name: "minus-one_one",
			body: []byte(`{"a": -1, "b": 1}`),
			want: `{"error": "Incorrect message"}`,
		},
		{
			name: "one_minus-one",
			body: []byte(`{"a": 1, "b": -1}`),
			want: `{"error": "Incorrect message"}`,
		},
		{
			name: "zero",
			body: []byte(`{"a": 0, "b": 0}`),
			want: `{"error": "Incorrect message"}`,
		},
		{
			name: "zero_one",
			body: []byte(`{"a": 0, "b": 1}`),
			want: `{"error": "Incorrect message"}`,
		},
		{
			name: "one_zero",
			body: []byte(`{"a": 1, "b": 0}`),
			want: `{"error": "Incorrect message"}`,
		},
		{
			name: "five_four",
			body: []byte(`{"a": 5, "b": 4}`),
			want: `{"a!":120, "b!":24}`,
		},
		{
			name: "twenty_twenty",
			body: []byte(`{"a": 20, "b": 20}`),
			want: `{"a!":2432902008176640000, "b!":2432902008176640000}`,
		},
		{
			name: "twenty-one_twenty",
			body: []byte(`{"a": 21, "b": 20}`),
			want: `{"error": "Incorrect message"}`,
		},
		{
			name: "twenty_twenty-one",
			body: []byte(`{"a": 20, "b": 21}`),
			want: `{"error": "Incorrect message"}`,
		},
		{
			name: "twenty-two_twenty-two",
			body: []byte(`{"a": 22, "b": 22}`),
			want: `{"error": "Incorrect message"}`,
		},
	}

	router := httprouter.New()
	router.GET("/calculate", calculate)

	for _, tc := range testCases {
		t.Run(tc.name, func(*testing.T) {
			req, err := http.NewRequest("GET", "/calculate", bytes.NewBuffer(tc.body))
			if err != nil {
				t.Errorf("Error on request: %v", err)
			}
			rr := httptest.NewRecorder()

			router.ServeHTTP(rr, req)
			jsonassert.New(t).Assertf(rr.Body.String(), tc.want)
		})
	}
}
