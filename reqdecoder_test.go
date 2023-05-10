package reqdecoder

import (
	"bytes"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestRequestDecoder(t *testing.T) {
	type request struct {
		contentEncoding string
	}
	type response struct {
		status int
		body   string
	}
	tests := []struct {
		name  string
		dFunc decoderFunc
		req   request
		resp  response
	}{
		{
			name: "Empty Content-Encoding",
			req: request{
				contentEncoding: "",
			},
			resp: response{
				status: http.StatusOK,
				body:   "",
			},
		},
		{
			name: "Content-Encoding identity",
			req: request{
				contentEncoding: "identity",
			},
			resp: response{
				status: http.StatusOK,
				body:   "",
			},
		},
		{
			name: "Success decode",
			dFunc: func(reader io.ReadCloser) io.ReadCloser {
				return reader
			},
			req: request{
				contentEncoding: "good_enc",
			},
			resp: response{
				status: http.StatusOK,
				body:   "",
			},
		},
		{
			name: "Bad decode",
			dFunc: func(reader io.ReadCloser) io.ReadCloser {
				return nil
			},
			req: request{
				contentEncoding: "bad_enc",
			},
			resp: response{
				status: http.StatusInternalServerError,
				body:   "Can't decode body",
			},
		},
		{
			name: "Unknown enc",
			req: request{
				contentEncoding: "unknown",
			},
			resp: response{
				status: http.StatusUnsupportedMediaType,
				body:   "Unsupported Content-Encoding",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.dFunc != nil {
				AddDecoder(tt.req.contentEncoding, tt.dFunc)
			}
			svr := httptest.NewServer(RequestDecoder(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				assert.Equal(t, tt.req.contentEncoding, r.Header.Get("Content-Encoding"))
				w.WriteHeader(http.StatusOK)
				_, err := w.Write(nil)
				assert.NoError(t, err)
			})))
			defer svr.Close()
			req, err := http.NewRequest(http.MethodPost, svr.URL, bytes.NewBuffer([]byte("test")))
			assert.NoError(t, err)
			req.Header.Set("Content-Encoding", tt.req.contentEncoding)
			client := &http.Client{}
			resp, err := client.Do(req)
			assert.NoError(t, err)
			defer func() {
				err = resp.Body.Close()
				assert.NoError(t, err)
			}()
			body, err := io.ReadAll(resp.Body)
			assert.NoError(t, err)
			assert.Equal(t, tt.resp.status, resp.StatusCode)
			assert.Contains(t, string(body), tt.resp.body)
		})
	}
}
