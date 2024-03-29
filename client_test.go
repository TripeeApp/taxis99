package taxis99

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"net/url"
	"reflect"
	"testing"
)

func newMockServer(hc *http.Client, handler http.HandlerFunc) (*Client, *httptest.Server) {
	s := httptest.NewServer(handler)

	c := NewClient(hc)
	u, _ := url.Parse(s.URL + "/")
	c.BaseURL = u

	return c, s
}

func TestClientRequestMethod(t *testing.T) {
	testCases := []struct {
		method string
	}{
		{http.MethodGet},
		{http.MethodPost},
	}

	for _, tc := range testCases {
		t.Run(tc.method, func(t *testing.T) {
			var got string
			handler := func(w http.ResponseWriter, r *http.Request) {
				got = r.Method
			}

			client, svr := newMockServer(nil, handler)
			defer svr.Close()

			err := client.Request(context.Background(), tc.method, "", nil, nil)
			if err != nil {
				t.Fatalf("Got error calling Request: %s; want it to be nil.", err.Error())
			}

			if got != tc.method {
				t.Errorf("Got Request.Method %s; want %s.", got, tc.method)
			}
		})
	}
}

func TestClientRequestPath(t *testing.T) {
	testCases := []struct {
		path string
		want string
	}{
		{"employees", "/employees"},
		{"rides", "/rides"},
	}

	for _, tc := range testCases {
		t.Run(tc.path, func(t *testing.T) {
			var got string
			handler := func(w http.ResponseWriter, r *http.Request) {
				got = r.URL.Path
			}

			client, srv := newMockServer(nil, handler)
			defer srv.Close()

			err := client.Request(context.Background(), http.MethodGet, tc.path, nil, nil)
			if err != nil {
				t.Fatalf("Got error calling Request: %s; want it to be nil.", err.Error())
			}

			if got != tc.want {
				t.Errorf("Got Request.URL.Path '%s'; want '%s'.", got, tc.want)
			}

		})
	}
}

func TestClientRequestHeader(t *testing.T) {
	testCases := []struct {
		body       interface{}
		wantHeader string
	}{
		{
			struct {
				Name string `json:"name"`
			}{"Test"},
			"application/json",
		},
		{
			nil,
			"",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.wantHeader, func(t *testing.T) {
			var got string
			handler := func(w http.ResponseWriter, r *http.Request) {
				got = r.Header.Get("Content-Type")
			}

			client, srv := newMockServer(nil, handler)
			defer srv.Close()

			err := client.Request(context.Background(), http.MethodGet, "", tc.body, nil)
			if err != nil {
				t.Fatalf("Got error calling Request: %s; want it to be nil.", err.Error())
			}

			if got != tc.wantHeader {
				t.Errorf("Got Content-Type Header: '%s'; want '%s'.", got, tc.wantHeader)
			}
		})
	}

}

func TestClientRequestBody(t *testing.T) {
	testCases := []struct {
		body     interface{}
		wantBody []byte
	}{
		{
			struct {
				Name string `json:"name"`
			}{"Test"},
			[]byte(`{"name":"Test"}`),
		},
	}

	for _, tc := range testCases {
		t.Run(string(tc.wantBody), func(t *testing.T) {
			var got []byte
			handler := func(w http.ResponseWriter, r *http.Request) {
				path, err := ioutil.ReadAll(r.Body)
				if err != nil {
					t.Fatalf("Got error while Request.Body: %s; want it to be nil.", err.Error())
				}
				got = path
			}

			client, srv := newMockServer(nil, handler)
			defer srv.Close()

			err := client.Request(context.Background(), http.MethodGet, "", tc.body, nil)
			if err != nil {
				t.Fatalf("Got error calling Request: %s; want it to be nil.", err.Error())
			}

			if !bytes.Contains(got, tc.wantBody) {
				t.Errorf("got body: %s, want %s.", got, tc.wantBody)
			}
		})
	}

}

func TestClientRequestResponseBody(t *testing.T) {
	type employee struct {
		Name string `json:"name"`
	}

	response := []byte(`{"name":"test"}`)

	t.Run(string(response), func(t *testing.T) {
		var emp employee
		handler := func(w http.ResponseWriter, r *http.Request) {
			w.Write(response)
		}

		client, srv := newMockServer(nil, handler)
		defer srv.Close()

		err := client.Request(context.Background(), http.MethodGet, "", nil, &emp)
		if err != nil {
			t.Fatalf("Got error calling Request: %s; want it to be nil.", err.Error())
		}

		if got, _ := json.Marshal(emp); !bytes.Contains(got, response) {
			t.Errorf("Got response '%s'; want '%s'.", got, response)
		}
	})

	t.Run("empty", func(t *testing.T) {
		handler := func(w http.ResponseWriter, r *http.Request) {}

		client, srv := newMockServer(nil, handler)
		defer srv.Close()

		err := client.Request(context.Background(), http.MethodGet, "", nil, nil)
		if err != nil {
			t.Fatalf("Got error calling Request: %s; want it to be nil.", err.Error())
		}
	})

}

type testRoundTripperFn func(*http.Request) (*http.Response, error)

func (fn testRoundTripperFn) RoundTrip(r *http.Request) (*http.Response, error) {
	return fn(r)
}

func TestClientRequestPathError(t *testing.T) {
	path := ":"
	c := NewClient(nil)

	err := c.Request(context.Background(), http.MethodGet, path, nil, nil)
	if err == nil {
		t.Error("Got error nil; want it not to be nil.")
	}
}

func TestClientRequestError(t *testing.T) {
	t.Run("NewRequest", func(t *testing.T) {
		c := NewClient(nil)
		err := c.Request(context.Background(), "ö", "", nil, nil)
		if err == nil {
			t.Error("Got error nil; want it not to be nil.")
		}
	})

	t.Run("Do", func(t *testing.T) {
		tripper := testRoundTripperFn(func(r *http.Request) (*http.Response, error) {
			defer r.Body.Close()
			return nil, errors.New("Testing error.")
		})
		hc := &http.Client{
			Transport: tripper,
		}

		c := NewClient(hc)

		err := c.Request(context.Background(), http.MethodGet, "", nil, nil)
		if err == nil {
			t.Error("Got error nil; want it not to be nil.")
		}
	})
}

func TestClienRequestAPIError(t *testing.T) {
	testCases := []struct {
		handler http.HandlerFunc
	}{
		{
			func(w http.ResponseWriter, r *http.Request) {
				w.Write([]byte("nvalidJSON"))
			},
		},
		{
			func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(http.StatusUnprocessableEntity)
				w.Write([]byte(`{"errors":[{"field":"employee.phone","message":"error.invalidPhoneNumber"}]}`))
			},
		},
	}

	for _, tc := range testCases {
		client, srv := newMockServer(nil, tc.handler)
		defer srv.Close()

		var out struct{}
		err := client.Request(context.Background(), http.MethodGet, "", nil, &out)
		if err == nil {
			t.Error("Got error nil; want it not to be nil.")
		}

		var apiErr *APIError
		if !errors.As(err, &apiErr) {
			t.Error("Got error type not to be of APIError type; want it to be.")
		}
	}
}

func TestClientRequestContextCancelled(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	cancel()

	c := NewClient(nil)

	err := c.Request(ctx, http.MethodGet, "", nil, nil)
	if err == nil {
		t.Error("Got error nil; want it not to be nil.")
	}
}

type mockRequester func(ctx context.Context, method, path string, body, output interface{}) error

func (fn mockRequester) Request(ctx context.Context, method, path string, body, output interface{}) error {
	return fn(ctx, method, path, body, output)
}

func newMockRequesterClient(rmck mockRequester) *Client {
	c := NewClient(nil)
	c.common.client = rmck

	return c
}

func testPath(t *testing.T, want string, run func(*Client) error) {
	t.Run("Path", func(t *testing.T) {
		var got string
		request := func(ctx context.Context, method, path string, body, output interface{}) error {
			got = path
			return nil
		}

		c := newMockRequesterClient(mockRequester(request))

		err := run(c)
		if err != nil {
			t.Fatalf("Got unexpected error '%s'; want nil.", err.Error())
		}

		if got != want {
			t.Errorf("Got path '%s'; want '%s'.", got, want)
		}
	})
}

func testMethod(t *testing.T, want string, run func(*Client) error) {
	t.Run("Method", func(t *testing.T) {
		var got string
		request := func(ctx context.Context, method, path string, body, output interface{}) error {
			got = method
			return nil
		}

		c := newMockRequesterClient(mockRequester(request))

		err := run(c)
		if err != nil {
			t.Fatalf("Got unexpected error '%s'; want nil.", err.Error())
		}

		if got != want {
			t.Errorf("Got method '%s'; want '%s'.", got, want)
		}
	})
}

func testQuery(t *testing.T, f []Filter, allowedFields map[string]struct{}, run func(*Client, Filter) error) {
	t.Run("Query", func(t *testing.T) {
		for _, filter := range f {
			var got url.Values
			request := func(ctx context.Context, method, path string, body, output interface{}) error {
				u, _ := url.Parse(path)
				got = u.Query()
				return nil
			}

			c := newMockRequesterClient(mockRequester(request))

			err := run(c, filter)
			if err != nil {
				t.Fatalf("Got unexpected error '%s'; want nil.", err.Error())
			}

			if want := filter.values(allowedFields); !reflect.DeepEqual(got, want) {
				t.Errorf("Got query values '%+v'; want '%+v'.", got, want)
			}
		}
	})
}

func testResponseBody(t *testing.T, responses [][]byte, run func(*Client) (interface{}, error)) {
	t.Run("ResponseBody", func(t *testing.T) {
		for _, response := range responses {
			request := func(ctx context.Context, method, path string, body, output interface{}) error {
				buf := bytes.NewBuffer(response)
				json.NewDecoder(buf).Decode(output)
				return nil
			}

			c := newMockRequesterClient(mockRequester(request))

			res, err := run(c)
			if err != nil {
				t.Fatalf("Got unexpected error '%s'; want nil.", err.Error())
			}

			want, _ := json.Marshal(res)
			if !bytes.Contains(response, want) {
				t.Errorf("Got response %s; want substring %s.", response, want)
			}
		}
	})
}

func testRequestBody(t *testing.T, tests []func(*Client) ([]byte, error)) {
	t.Run("RequestBody", func(t *testing.T) {
		for _, run := range tests {
			var got []byte
			request := func(ctx context.Context, method, path string, body, output interface{}) error {
				got, _ = json.Marshal(body)
				return nil
			}

			c := newMockRequesterClient(mockRequester(request))

			want, err := run(c)
			if err != nil {
				t.Fatalf("Got unexpected error '%s'; want nil.", err.Error())
			}

			if !bytes.Contains(got, want) {
				t.Errorf("Got request body %s; want %s.", got, want)
			}
		}
	})
}

func testError(t *testing.T, run func(*Client) error) {
	request := func(ctx context.Context, method, path string, body, output interface{}) error {
		return errors.New("Error!")
	}

	c := newMockRequesterClient(mockRequester(request))

	if err := run(c); err == nil {
		t.Error("Got error nil; want it not nil.")
	}
}
