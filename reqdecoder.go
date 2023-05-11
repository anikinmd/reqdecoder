// Package reqdecoder implements a request body decoder/decompressor.
// This package contains gzip and  deflate decompression and supports custom decoders.
package reqdecoder

import (
	"io"
	"net/http"
)

var defaultDecoder *requestDecoder

func init() {
	defaultDecoder = newRequestDecoder()
	defaultDecoder.fillDecoders()
}

// decoderFunc should return nil on decode error.
type decoderFunc func(reader io.ReadCloser) io.ReadCloser

type requestDecoder struct {
	decoders map[string]decoderFunc
}

func newRequestDecoder() *requestDecoder {
	d := &requestDecoder{
		decoders: make(map[string]decoderFunc),
	}

	return d
}

// RequestDecoder is decoding middleware function.
func RequestDecoder(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		enc := r.Header.Get("Content-Encoding")
		if enc != "" && enc != "identity" {
			f, ok := defaultDecoder.decoders[enc]
			if !ok {
				http.Error(w, "Unsupported Content-Encoding", http.StatusUnsupportedMediaType)
				return
			}
			decoded := f(r.Body)
			if decoded == nil {
				http.Error(w, "Can't decode body", http.StatusInternalServerError)
				return
			}
			r.Body = decoded
		}
		next.ServeHTTP(w, r)
	})
}

// AddDecoder adds new decoder to the list of decoders.
// contentEncoding should match "Content-Encoding" http header tag.
// decoder should return nil on error.
func AddDecoder(contentEncoding string, decoder decoderFunc) {
	defaultDecoder.addDecoder(contentEncoding, decoder)
}
